package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck handles GET /health - health check endpoint
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "UP",
		"service": "resume-backend",
	})
}
