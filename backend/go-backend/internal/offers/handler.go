package offers

import (
	"net/http"
	"time"

	"resume/internal/server/response"

	"github.com/gin-gonic/gin"
)

// Handler exposes offer HTTP endpoints.
type Handler struct {
	svc *Service
}

// NewHandler creates an offers handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// InitiateRequest body.
type InitiateRequest struct {
	ATSID    string  `json:"ats_id" binding:"required"`
	Amount   string  `json:"amount"`
	Currency string  `json:"currency"`
	StartsAt *string `json:"starts_at"` // RFC3339
}

// Initiate creates an offer (HR).
func (h *Handler) Initiate(c *gin.Context) {
	var req InitiateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	var startsAt *time.Time
	if req.StartsAt != nil && *req.StartsAt != "" {
		t, err := time.Parse(time.RFC3339, *req.StartsAt)
		if err == nil {
			startsAt = &t
		}
	}
	offer, err := h.svc.Initiate(c.Request.Context(), req.ATSID, req.Amount, req.Currency, startsAt)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, offer)
}

// Accept accepts an offer (candidate); ATS → HIRED atomically.
func (h *Handler) Accept(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Accept(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "accepted"})
}

// Reject rejects an offer (candidate).
func (h *Handler) Reject(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Reject(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "rejected"})
}

// List returns all offers for HR. Placeholder: returns empty list until ListForHR is implemented.
func (h *Handler) List(c *gin.Context) {
	response.JSON(c, http.StatusOK, []interface{}{})
}

// GetByID returns an offer by ID.
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	offer, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, offer)
}
