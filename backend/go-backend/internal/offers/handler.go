package offers

import (
	"net/http"
	"time"

	"resume/internal/domain/entities"
	"resume/internal/middleware"
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

// List returns all offers for HR (offers for jobs created by the authenticated HR user).
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
		list = []*entities.Offer{}
	}
	response.JSON(c, http.StatusOK, list)
}

// GetByID returns an offer by ID (HR).
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	offer, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, offer)
}

// GetByIDForCandidate returns the offer by ID for the candidate (Phase 3: view offer / download letter).
func (h *Handler) GetByIDForCandidate(c *gin.Context) {
	candidateID := c.Param("candidate_id")
	offerID := c.Param("id")
	offer, err := h.svc.GetByIDForCandidate(c.Request.Context(), offerID, candidateID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, offer)
}
