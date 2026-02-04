package auth

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
	db      *sql.DB
}

func NewHandler(service *Service, db *sql.DB) *Handler {
	return &Handler{
		service: service,
		db:      db,
	}
}

func (h *Handler) OAuthLogin(c echo.Context) error {
	provider := c.Param("provider")

	url, err := h.service.GetOAuthURL(provider)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"auth_url": url,
	})
}

func (h *Handler) OAuthCallback(c echo.Context) error {
	provider := c.Param("provider")
	code := c.QueryParam("code")

	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Authorization code is required",
		})
	}

	response, err := h.service.HandleOAuthCallback(provider, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) RefreshToken(c echo.Context) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	if req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Refresh token is required",
		})
	}

	response, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid refresh token",
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) Me(c echo.Context) error {
	claims := c.Get("user").(*Claims)

	var dbUser struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Provider  string `json:"provider"`
	}

	err := h.db.QueryRow(
		"SELECT id, email, name, avatar_url, provider FROM users WHERE id = $1",
		claims.UserID,
	).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Name, &dbUser.AvatarURL, &dbUser.Provider)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(http.StatusOK, dbUser)
}

func (h *Handler) Logout(c echo.Context) error {
	claims := c.Get("user").(*Claims)

	// Delete all sessions for the user
	_, err := h.db.Exec("DELETE FROM user_sessions WHERE user_id = $1", claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to logout",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
