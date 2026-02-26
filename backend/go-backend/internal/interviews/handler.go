package interviews

import (
	"net/http"
	"time"

	"resume/internal/middleware"
	"resume/internal/server/response"

	"github.com/gin-gonic/gin"
)

// Handler exposes interview HTTP endpoints.
type Handler struct {
	svc *Service
}

// NewHandler creates an interviews handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// ScheduleRequest body.
type ScheduleRequest struct {
	ATSID       string `json:"ats_id" binding:"required"`
	ScheduledAt string `json:"scheduled_at" binding:"required"` // RFC3339
	Duration   int    `json:"duration_minutes"`
	Location   string `json:"location"`
	Round      int    `json:"round"`
	Notes      string `json:"notes"`
}

// Schedule creates an interview (HR).
func (h *Handler) Schedule(c *gin.Context) {
	var req ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		response.ValidationError(c, err)
		return
	}
	interview, err := h.svc.Schedule(c.Request.Context(), req.ATSID, scheduledAt, req.Duration, req.Location, req.Round, req.Notes)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, interview)
}

// List returns all scheduled interviews for HR (jobs created by the authenticated HR user).
func (h *Handler) List(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return
	}
	list, err := h.svc.ListForHR(c.Request.Context(), claims.UserID)
	if err != nil {
		response.Error(c, err)
		return
	}
	if list == nil {
		list = []*InterviewForHR{}
	}
	response.JSON(c, http.StatusOK, list)
}

// GetByID returns an interview by ID.
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	interview, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, interview)
}

// ListByATSID returns interviews for an ATS record.
func (h *Handler) ListByATSID(c *gin.Context) {
	atsID := c.Param("ats_id")
	list, err := h.svc.ListByATSID(c.Request.Context(), atsID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, list)
}

// UpdateRequest body.
type UpdateRequest struct {
	ScheduledAt string `json:"scheduled_at"`
	Duration   int    `json:"duration_minutes"`
	Location   string `json:"location"`
	Round      int    `json:"round"`
	Notes      string `json:"notes"`
	Status     string `json:"status"`
}

// Update updates an interview (HR).
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	interview, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	if req.ScheduledAt != "" {
		t, err := time.Parse(time.RFC3339, req.ScheduledAt)
		if err == nil {
			interview.ScheduledAt = t
		}
	}
	if req.Duration > 0 {
		interview.Duration = req.Duration
	}
	interview.Location = req.Location
	if req.Round > 0 {
		interview.Round = req.Round
	}
	interview.Notes = req.Notes
	if req.Status != "" {
		interview.Status = req.Status
	}
	if err := h.svc.Update(c.Request.Context(), interview); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, interview)
}

// ConfirmByCandidate records that the candidate confirmed the interview (Phase 3). Route: POST /api/candidates/:candidate_id/interviews/:id/confirm
func (h *Handler) ConfirmByCandidate(c *gin.Context) {
	candidateID := c.Param("candidate_id")
	interviewID := c.Param("id")
	if err := h.svc.ConfirmForCandidate(c.Request.Context(), candidateID, interviewID); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "confirmed"})
}
