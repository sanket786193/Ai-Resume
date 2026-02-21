package middleware

import (
	"net/http"
	"strings"

	"resume/internal/auth"
	domainerrors "resume/internal/domain/errors"
	"resume/internal/server/response"

	"github.com/gin-gonic/gin"
)

// Auth validates JWT and sets claims in context.
func Auth(svc *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, &domainerrors.UnauthorizedError{Message: "Authorization header required"})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, &domainerrors.UnauthorizedError{Message: "Invalid Authorization format"})
			c.Abort()
			return
		}
		claims, err := svc.ValidateToken(parts[1])
		if err != nil {
			response.Error(c, &domainerrors.UnauthorizedError{Message: "Invalid or expired token"})
			c.Abort()
			return
		}
		c.Set(auth.ContextKeyClaims, claims)
		c.Next()
	}
}

// GetClaims returns JWT claims from context (nil if not set).
func GetClaims(c *gin.Context) *auth.Claims {
	v, ok := c.Get(auth.ContextKeyClaims)
	if !ok {
		return nil
	}
	claims, _ := v.(*auth.Claims)
	return claims
}

// RequireRole allows only the given roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get(auth.ContextKeyClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		claims := v.(*auth.Claims)
		for _, r := range roles {
			if claims.Role == r {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		c.Abort()
	}
}
