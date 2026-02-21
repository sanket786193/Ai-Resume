package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger logs method, path, status, latency, and request ID.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		rid, _ := c.Get("request_id")
		requestID, _ := rid.(string)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		log.Printf("[%s] %s %s %d %v %s", requestID, method, path, status, latency, clientIP)
	}
}
