package ats

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"resume/internal/domain/enums"
	domainerrors "resume/internal/domain/errors"
	"resume/internal/middleware"
	"resume/internal/server/response"

	"github.com/gin-gonic/gin"
)

var resumeProxyClient = &http.Client{Timeout: 60 * time.Second}

// Handler exposes ATS HTTP endpoints.
type Handler struct {
	svc *Service
}

// NewHandler creates an ATS handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// ApplyRequest body.
type ApplyRequest struct {
	JobID      string `json:"job_id" binding:"required"`
	ResumeID   string `json:"resume_id" binding:"required"`
}

// Apply creates an application (candidate).
func (h *Handler) Apply(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return
	}
	// Resolve candidate ID from user ID (middleware ensures auth)
	candidateID := c.Param("candidate_id")
	if candidateID == "" {
		response.ValidationError(c, nil)
		return
	}
	var req ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	rec, err := h.svc.Apply(c.Request.Context(), req.JobID, candidateID, req.ResumeID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, rec)
}

// GetApplicationStatus returns status for candidate's application to a job.
func (h *Handler) GetApplicationStatus(c *gin.Context) {
	jobID := c.Param("job_id")
	candidateID := c.Param("candidate_id")
	rec, err := h.svc.GetApplicationStatus(c.Request.Context(), jobID, candidateID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, rec)
}

// GetApplicationFeedback returns AI feedback (safe subset) for candidate's application (Phase 3).
func (h *Handler) GetApplicationFeedback(c *gin.Context) {
	jobID := c.Param("job_id")
	candidateID := c.Param("candidate_id")
	feedback, err := h.svc.GetApplicationFeedbackForCandidate(c.Request.Context(), jobID, candidateID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, feedback)
}

// ListMyApplications returns ATS records for the candidate with job_title and job_status for UI.
func (h *Handler) ListMyApplications(c *gin.Context) {
	candidateID := c.Param("candidate_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	list, err := h.svc.ListByCandidateWithJobs(c.Request.Context(), candidateID, limit, offset)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, list)
}

// ListByJob returns enriched applications for a job (HR): candidate name, email, job title.
func (h *Handler) ListByJob(c *gin.Context) {
	jobID := c.Param("id") // from route /jobs/:id/applications
	statusStr := c.Query("status")
	var status *enums.ATSStatus
	if statusStr != "" {
		s := enums.ATSStatus(statusStr)
		if s.Valid() {
			status = &s
		}
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	list, err := h.svc.ListByJobEnriched(c.Request.Context(), jobID, status, limit, offset)
	if err != nil {
		response.Error(c, err)
		return
	}
	if list == nil {
		list = []*ApplicationForHR{}
	}
	response.JSON(c, http.StatusOK, list)
}

// GetApplicationByID returns one application with full detail for HR (candidate contact, AI feedback, resume).
func (h *Handler) GetApplicationByID(c *gin.Context) {
	id := c.Param("id")
	detail, err := h.svc.GetApplicationByIDEnriched(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, detail)
}

// GetApplicationResume streams the application's resume PDF (proxied from storage) so the iframe can load it same-origin.
func (h *Handler) GetApplicationResume(c *gin.Context) {
	id := c.Param("id")
	resumeURL, err := h.svc.GetApplicationResumeURL(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	if resumeURL == "" {
		response.Error(c, &domainerrors.NotFoundError{Resource: "resume", ID: id})
		return
	}
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, resumeURL, nil)
	if err != nil {
		response.Error(c, err)
		return
	}
	// Cloudinary and many CDNs return 401 for requests without a browser-like User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/115.0")
	resp, err := resumeProxyClient.Do(req)
	if err != nil {
		response.Error(c, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg := "resume fetch failed: " + resp.Status
		if resp.StatusCode == http.StatusUnauthorized {
			msg = "resume fetch failed (401). If using Cloudinary, enable 'Allow delivery of PDF and ZIP files' in Dashboard → Settings → Security."
		}
		response.Error(c, fmt.Errorf("%s", msg))
		return
	}
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline; filename=\"resume.pdf\"")
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		_ = c.Error(err)
	}
}

// ListForHR returns ATS records for all of HR's jobs (or one job if job_id query is set).
func (h *Handler) ListForHR(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return
	}
	jobIDStr := c.Query("job_id")
	var jobID *string
	if jobIDStr != "" {
		jobID = &jobIDStr
	}
	statusStr := c.Query("status")
	var status *enums.ATSStatus
	if statusStr != "" {
		s := enums.ATSStatus(statusStr)
		if s.Valid() {
			status = &s
		}
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	list, err := h.svc.ListForHREnriched(c.Request.Context(), claims.UserID, jobID, status, limit, offset)
	if err != nil {
		response.Error(c, err)
		return
	}
	if list == nil {
		list = []*ApplicationForHR{}
	}
	response.JSON(c, http.StatusOK, list)
}

// UpdateStatus updates ATS status (HR).
func (h *Handler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	st := enums.ATSStatus(req.Status)
	if !st.Valid() {
		response.ValidationError(c, nil)
		return
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), id, st); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "updated"})
}

// BulkApplyRequest body (Phase 4).
type BulkApplyRequest struct {
	Applications []struct {
		CandidateID string `json:"candidate_id"`
		ResumeID    string `json:"resume_id"`
	} `json:"applications" binding:"required"`
}

// BulkApply creates multiple applications for a job (HR). Route: POST /api/hr/jobs/:id/bulk-apply.
func (h *Handler) BulkApply(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return
	}
	jobID := c.Param("id")
	if jobID == "" {
		response.ValidationError(c, nil)
		return
	}
	var req BulkApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	items := make([]BulkApplyItem, 0, len(req.Applications))
	for _, a := range req.Applications {
		items = append(items, BulkApplyItem{CandidateID: a.CandidateID, ResumeID: a.ResumeID})
	}
	result, err := h.svc.BulkApply(c.Request.Context(), jobID, claims.UserID, items)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
}

