package ats

import (
	"net/http"
	"strconv"

	"resume/internal/domain/enums"
	"resume/internal/middleware"
	"resume/internal/server/response"

	"github.com/gin-gonic/gin"
)

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

// ListByJob returns ATS records for a job (HR).
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
	list, err := h.svc.ListByJob(c.Request.Context(), jobID, status, limit, offset)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, list)
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
	list, err := h.svc.ListForHR(c.Request.Context(), claims.UserID, jobID, status, limit, offset)
	if err != nil {
		response.Error(c, err)
		return
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
