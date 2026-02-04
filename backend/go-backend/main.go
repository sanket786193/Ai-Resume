package main

import (
	"database/sql"
	"log"

	"resume/api"
	"resume/internal/config"
	"resume/internal/database"
	"resume/internal/service/auth"
	"resume/internal/service/ocr"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	var db *sql.DB
	if cfg.Database.Enabled {
		dbConn, err := database.InitDB(&cfg.Database)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		defer database.CloseDB()
		db = dbConn
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize services
	ocrService := ocr.NewService(cfg.UploadDir)

	// Initialize storage (Cloudinary or local)
	var ocrStorage ocr.Storage
	var err error

	if cfg.Cloudinary.Enabled && cfg.Cloudinary.CloudName != "" {
		// Use Cloudinary storage
		ocrStorage, err = ocr.NewCloudinaryStorage(
			cfg.Cloudinary.CloudName,
			cfg.Cloudinary.APIKey,
			cfg.Cloudinary.APISecret,
			cfg.Cloudinary.Folder,
		)
		if err != nil {
			log.Fatalf("Failed to initialize Cloudinary storage: %v", err)
		}
		log.Printf("Using Cloudinary storage (folder: %s)", cfg.Cloudinary.Folder)
	} else {
		// Use local file storage
		ocrStorage = ocr.NewStorage(cfg.UploadDir)
		log.Printf("Using local file storage (directory: %s)", cfg.UploadDir)
	}

	ocrHandler := ocr.NewHandler(ocrService, ocrStorage)

	// Initialize auth service and handler
	var authService *auth.Service
	var authHandler *auth.Handler
	if cfg.Database.Enabled && db != nil {
		authService = auth.NewService(&cfg.Auth, db)
		authHandler = auth.NewHandler(authService, db)
		log.Println("Auth service initialized")
	}

	// Setup routes
	api.SetupRoutes(e, ocrHandler, authHandler, authService)

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)

	if err := e.Start(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
