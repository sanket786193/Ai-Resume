package auth

import (
	"net/http"

	"resume/internal/server/response"

	"github.com/gin-gonic/gin"
)

// Handler exposes auth HTTP endpoints.
type Handler struct {
	svc *Service
}

// NewHandler creates an auth handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterCandidateRequest body.
type RegisterCandidateRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	LinkedIn string `json:"linkedin"`
}

// RegisterHRRequest body.
type RegisterHRRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// LoginRequest body.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest body.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest body (refresh_token optional; if sent, server invalidates that session).
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RegisterCandidate godoc
// @Summary Register candidate
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterCandidateRequest true "body"
// @Success 201 {object} auth.AuthResponse
// @Router /auth/register/candidate [post]
func (h *Handler) RegisterCandidate(c *gin.Context) {
	var req RegisterCandidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	out, err := h.svc.RegisterCandidate(c.Request.Context(), RegisterCandidateInput{
		Email: req.Email, Password: req.Password, Name: req.Name,
		Phone: req.Phone, LinkedIn: req.LinkedIn,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, out)
}

// RegisterHR godoc
// @Summary Register HR
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterHRRequest true "body"
// @Success 201 {object} auth.AuthResponse
// @Router /auth/register/hr [post]
func (h *Handler) RegisterHR(c *gin.Context) {
	var req RegisterHRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	out, err := h.svc.RegisterHR(c.Request.Context(), RegisterHRInput{
		Email: req.Email, Password: req.Password, Name: req.Name,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, out)
}

// Login godoc
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "body"
// @Success 200 {object} auth.AuthResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	out, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, out)
}

// Refresh godoc
// @Summary Refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RefreshRequest true "body"
// @Success 200 {object} auth.AuthResponse
// @Router /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	out, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, out)
}

// Logout godoc
// @Summary Logout
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LogoutRequest false "optional refresh_token to invalidate"
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	_ = c.ShouldBindJSON(&req) // optional body
	if err := h.svc.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "logged out"})
}

// Me returns the current user (requires Auth middleware).
func (h *Handler) Me(c *gin.Context) {
	v, ok := c.Get(ContextKeyClaims)
	if !ok {
		return
	}
	claims, _ := v.(*Claims)
	if claims == nil {
		return
	}
	out, err := h.svc.Me(c.Request.Context(), claims.UserID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, out)
}
