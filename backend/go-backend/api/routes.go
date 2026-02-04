package api

import (
	"resume/api/handlers"
	"resume/internal/middleware"
	"resume/internal/service/auth"
	"resume/internal/service/ocr"

	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all API routes
func SetupRoutes(e *echo.Echo, ocrHandler *ocr.Handler, authHandler *auth.Handler, authService *auth.Service) {
	// Health check
	e.GET("/health", handlers.HealthCheck)

	// Auth endpoints (public)
	authGroup := e.Group("/auth")
	{
		authGroup.GET("/login/:provider", authHandler.OAuthLogin)
		authGroup.GET("/callback/:provider", authHandler.OAuthCallback)
		authGroup.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	protected := e.Group("/api")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.GET("/me", authHandler.Me)
		protected.POST("/logout", authHandler.Logout)
	}

	// OCR endpoints
	ocrGroup := e.Group("/ocr")
	{
		// Upload and process image or PDF
		ocrGroup.POST("/image", ocrHandler.ImageOCR)
		ocrGroup.POST("/file", ocrHandler.ImageOCR) // Alias for image endpoint

		// Test endpoint - process files from uploads folder
		ocrGroup.GET("/test/:filename", ocrHandler.TestOCR)
		ocrGroup.GET("/test", ocrHandler.TestOCR)
	}
}
