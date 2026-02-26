package middleware

import (
	"strings"

	"resume/internal/auth"
	domainerrors "resume/internal/domain/errors"
	"resume/internal/server/response"

	"github.com/gin-gonic/gin"
)

// Auth validates JWT and sets claims in context.
// Token can be in Authorization header (Bearer) or in query param "access_token" (e.g. for iframe PDF viewer).
func Auth(svc *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ""
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
		if token == "" {
			token = c.Query("access_token")
		}
		if token == "" {
			response.Error(c, &domainerrors.UnauthorizedError{Message: "Authorization header or access_token required"})
			c.Abort()
			return
		}
		claims, err := svc.ValidateToken(token)
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
			response.Error(c, &domainerrors.UnauthorizedError{Message: "unauthorized"})
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
		response.Error(c, &domainerrors.ForbiddenError{Message: "HR role required"})
		c.Abort()
	}
}
