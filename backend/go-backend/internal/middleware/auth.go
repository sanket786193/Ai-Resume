package middleware

import (
	"net/http"
	"strings"

	"resume/internal/service/auth"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware validates JWT tokens and sets user claims in context
func AuthMiddleware(authService *auth.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header is required",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format. Expected: Bearer <token>",
				})
			}

			token := parts[1]
			claims, err := authService.ValidateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
			}

			// Set user claims in context for use in handlers
			c.Set("user", claims)
			return next(c)
		}
	}
}
