package middleware

import (
	"log"
	"net/http"

	"resume/internal/server/response"

	"github.com/gin-gonic/gin"
)

// Recovery recovers from panics and returns 500.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, response.Body{
					Success: false,
					Error:   &response.ErrorInfo{Code: "INTERNAL_ERROR", Message: "internal server error"},
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
