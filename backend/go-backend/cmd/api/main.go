package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"resume/internal/ai/client"
	"resume/internal/ats"
	"resume/internal/auth"
	"resume/internal/config"
	"resume/internal/candidates"
	"resume/internal/interviews"
	"resume/internal/jobs"
	"resume/internal/notifications"
	"resume/internal/offers"
	"resume/internal/server"
	"resume/internal/cloudinary"
	"resume/internal/storage/postgres"
	"resume/internal/storage/tx"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // load .env from cwd so os.Getenv sees CLOUDINARY_*, DB_*, etc.
	cfg := config.LoadConfig()

	var db *postgres.DB
	if cfg.Database.Enabled {
		var err error
		db, err = postgres.New(&cfg.Database)
		if err != nil {
			log.Fatalf("database: %v", err)
		}
		defer db.Close()
	} else {
		log.Fatal("database is required; set DB_ENABLED=true")
	}

	// Repositories
	userRepo := postgres.NewUserRepo(db.DB)
	sessionRepo := postgres.NewSessionRepo(db.DB)
	jobRepo := postgres.NewJobRepo(db.DB)
	candidateRepo := postgres.NewCandidateRepo(db.DB)
	resumeRepo := postgres.NewResumeRepo(db.DB)
	parsedRepo := postgres.NewResumeParsedRepo(db.DB)
	atsRepo := postgres.NewATSRepo(db.DB)
	interviewRepo := postgres.NewInterviewRepo(db.DB)
	offerRepo := postgres.NewOfferRepo(db.DB)

	txRunner := &tx.Runner{DB: db.DB}

	// Auth
	jwtSvc := auth.NewJWTService(cfg.Auth.JWTSecret, cfg.Auth.JWTExpiryHours, cfg.Auth.RefreshExpiryDays)
	authConfig := auth.AuthConfigAdapter(cfg)
	authSvc := auth.NewService(jwtSvc, userRepo, sessionRepo, candidateRepo, authConfig)
	authHandler := auth.NewHandler(authSvc)

	// AI client: gRPC or HTTP
	var atsAIClient ats.AIClient
	if cfg.AI.UseGRPC {
		var err error
		atsAIClient, err = ats.NewGRPCAIClient(cfg.AI.GRPCTarget, cfg.AI.TimeoutSec, cfg.AI.Enabled)
		if err != nil {
			log.Fatalf("AI gRPC client: %v", err)
		}
	} else {
		httpClient := client.New(cfg.AI.BaseURL, cfg.AI.TimeoutSec, cfg.AI.Enabled)
		atsAIClient = ats.NewAIAdapter(httpClient)
	}

	// Services & handlers
	jobSvc := jobs.NewService(jobRepo)
	jobHandler := jobs.NewHandler(jobSvc)

	candidateSvc := candidates.NewService(candidateRepo, resumeRepo, atsRepo, jobRepo)
	var resumeUploader candidates.ResumeUploader
	var atsResumeURLResolver ats.ResumeURLResolver
	if cfg.CloudinaryStorage.Enabled && cfg.CloudinaryStorage.CloudName != "" && cfg.CloudinaryStorage.APIKey != "" && cfg.CloudinaryStorage.APISecret != "" {
		cl := cloudinary.NewClient(cloudinary.Config{
			CloudName:     cfg.CloudinaryStorage.CloudName,
			APIKey:        cfg.CloudinaryStorage.APIKey,
			APISecret:     cfg.CloudinaryStorage.APISecret,
			UploadFolder:  cfg.CloudinaryStorage.Folder,
		})
		resumeUploader = &cloudinaryResumeUploader{client: cl}
		atsResumeURLResolver = func(storagePath string) string {
			if storagePath == "" {
				return ""
			}
			if strings.HasPrefix(storagePath, "http://") || strings.HasPrefix(storagePath, "https://") {
				return storagePath
			}
			return cl.PublicURL(storagePath)
		}
	}
	candidateHandler := candidates.NewHandler(candidateSvc, resumeUploader)

	notifSvc := notifications.NewServiceWithSMTP(&cfg.SMTP)
	atsSvc := ats.NewServiceWithResolver(atsRepo, jobRepo, resumeRepo, parsedRepo, candidateRepo, userRepo, atsAIClient, cfg.AI.Enabled, atsResumeURLResolver)
	atsSvc.SetNotifier(notifSvc)
	atsHandler := ats.NewHandler(atsSvc)

	interviewSvc := interviews.NewService(interviewRepo, atsRepo, candidateRepo, userRepo, jobRepo, notifSvc)
	interviewHandler := interviews.NewHandler(interviewSvc)

	offerSvc := offers.NewService(offerRepo, atsRepo, txRunner, candidateRepo, userRepo, jobRepo, notifSvc)
	offerHandler := offers.NewHandler(offerSvc)

	// Server
	srvCfg := server.Config{
		AuthHandler:        authHandler,
		JobsHandler:        jobHandler,
		CandidatesHandler:  candidateHandler,
		ATSHandler:         atsHandler,
		InterviewsHandler:  interviewHandler,
		OffersHandler:      offerHandler,
		AuthService:        authSvc,
		Port:               cfg.Port,
	}

	router := setupRouter(srvCfg)
	httpSrv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
	log.Println("server stopped")
}

// cloudinaryResumeUploader adapts cloudinary.Client to candidates.ResumeUploader.
type cloudinaryResumeUploader struct {
	client *cloudinary.Client
}

func (u *cloudinaryResumeUploader) Upload(ctx context.Context, fileContent []byte, fileName string) (string, error) {
	res, err := u.client.Upload(ctx, fileContent, fileName)
	if err != nil {
		return "", err
	}
	return res.SecureURL, nil
}

// setupRouter builds the Gin router (extracted for testing).
func setupRouter(cfg server.Config) *gin.Engine {
	// Use server.Run logic inline so we can return *gin.Engine for shutdown
	// For simplicity we duplicate route setup here and start server in main with http.Server
	// Alternatively server.Run could return (http.Handler, error). We'll use server's handler setup.
	return server.NewRouter(cfg)
}
