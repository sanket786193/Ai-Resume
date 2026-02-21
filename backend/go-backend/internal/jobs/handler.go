package jobs

import (
	"net/http"
	"strconv"

	"resume/internal/domain/enums"
	"resume/internal/middleware"
	"resume/internal/server/response"

	"github.com/gin-gonic/gin"
)

// Handler exposes job HTTP endpoints.
type Handler struct {
	svc *Service
}

// NewHandler creates a jobs handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// CreateJobRequest body.
type CreateJobRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Location    string `json:"location"`
	Department  string `json:"department"`
}

// UpdateJobRequest body.
type UpdateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Department  string `json:"department"`
}

// Create creates a job (HR).
func (h *Handler) Create(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return
	}
	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	job, err := h.svc.Create(c.Request.Context(), req.Title, req.Description, req.Location, req.Department, claims.UserID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, job)
}

// GetByID returns a job by ID (HR sees any; public sees only PUBLISHED).
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	job, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, job)
}

// List returns jobs for public listing. Defaults to PUBLISHED only when no status query.
func (h *Handler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	statusStr := c.Query("status")
	var status *enums.JobStatus
	if statusStr != "" {
		s := enums.JobStatus(statusStr)
		if s.Valid() {
			status = &s
		}
	}
	if status == nil {
		status = ptr(enums.JobPublished) // public list: only published by default
	}
	list, err := h.svc.List(c.Request.Context(), status, limit, offset)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, list)
}

// ListForHR returns all jobs for HR (any status). Requires HR role.
func (h *Handler) ListForHR(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	list, err := h.svc.List(c.Request.Context(), nil, limit, offset)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, list)
}

// Publish publishes a job (HR).
func (h *Handler) Publish(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Publish(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "published"})
}

// Update updates a job (HR, DRAFT only).
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	job, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	if req.Title != "" {
		job.Title = req.Title
	}
	if req.Description != "" {
		job.Description = req.Description
	}
	job.Location = req.Location
	job.Department = req.Department
	if err := h.svc.Update(c.Request.Context(), job); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, job)
}

// Close closes a job (HR).
func (h *Handler) Close(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Close(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "closed"})
}

// Delete soft-deletes a job (HR).
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.SoftDelete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "deleted"})
}

func ptr(s enums.JobStatus) *enums.JobStatus { return &s }
