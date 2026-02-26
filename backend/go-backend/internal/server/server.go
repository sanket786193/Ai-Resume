package server

import (
	"net/http"
	"time"

	"resume/internal/ats"
	"resume/internal/auth"
	"resume/internal/candidates"
	"resume/internal/interviews"
	"resume/internal/jobs"
	"resume/internal/middleware"
	"resume/internal/offers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Config for the HTTP server.
type Config struct {
	AuthHandler        *auth.Handler
	JobsHandler        *jobs.Handler
	CandidatesHandler  *candidates.Handler
	ATSHandler         *ats.Handler
	InterviewsHandler  *interviews.Handler
	OffersHandler      *offers.Handler
	AuthService        *auth.Service
	Port               string
}

// NewRouter returns a configured Gin engine for use with http.Server (e.g. graceful shutdown).
func NewRouter(cfg Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Public auth
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register/candidate", cfg.AuthHandler.RegisterCandidate)
		authGroup.POST("/register/hr", cfg.AuthHandler.RegisterHR)
		authGroup.POST("/login", cfg.AuthHandler.Login)
		authGroup.POST("/refresh", cfg.AuthHandler.Refresh)
		authGroup.POST("/logout", cfg.AuthHandler.Logout)
	}
	// Protected auth (current user)
	r.GET("/auth/me", middleware.Auth(cfg.AuthService), cfg.AuthHandler.Me)

	// Public job listing (published only can be exposed via query)
	r.GET("/jobs", cfg.JobsHandler.List)
	r.GET("/jobs/:id", cfg.JobsHandler.GetByID)

	// Protected: HR
	hr := r.Group("/api/hr")
	hr.Use(middleware.Auth(cfg.AuthService))
	hr.Use(middleware.RequireRole("HR"))
	{
		hr.GET("/jobs", cfg.JobsHandler.ListForHR)
		hr.POST("/jobs", cfg.JobsHandler.Create)
		hr.PUT("/jobs/:id", cfg.JobsHandler.Update)
		hr.POST("/jobs/:id/publish", cfg.JobsHandler.Publish)
		hr.POST("/jobs/:id/close", cfg.JobsHandler.Close)
		hr.DELETE("/jobs/:id", cfg.JobsHandler.Delete)
		hr.GET("/applications", cfg.ATSHandler.ListForHR)         // all HR's applications (optional ?job_id=)
		hr.GET("/applications/:id", cfg.ATSHandler.GetApplicationByID) // single application detail
		hr.GET("/applications/:id/resume", cfg.ATSHandler.GetApplicationResume) // resume PDF (proxied for iframe)
		hr.GET("/jobs/:id/applications", cfg.ATSHandler.ListByJob) // id = job_id
		hr.POST("/jobs/:id/bulk-apply", cfg.ATSHandler.BulkApply)  // Phase 4: bulk add applicants
		hr.PUT("/applications/:id/status", cfg.ATSHandler.UpdateStatus)
		hr.GET("/interviews", cfg.InterviewsHandler.List)
		hr.POST("/interviews", cfg.InterviewsHandler.Schedule)
		hr.GET("/interviews/:id", cfg.InterviewsHandler.GetByID)
		hr.PUT("/interviews/:id", cfg.InterviewsHandler.Update)
		hr.GET("/ats/:ats_id/interviews", cfg.InterviewsHandler.ListByATSID)
		hr.GET("/offers", cfg.OffersHandler.List)
		hr.POST("/offers", cfg.OffersHandler.Initiate)
		hr.GET("/offers/:id", cfg.OffersHandler.GetByID)
	}

	// Protected: Candidate
	candidate := r.Group("/api/candidates/:candidate_id")
	candidate.Use(middleware.Auth(cfg.AuthService))
	candidate.Use(middleware.RequireRole("CANDIDATE"))
	{
		candidate.GET("/profile", cfg.CandidatesHandler.GetProfile)
		candidate.POST("/resumes", cfg.CandidatesHandler.UploadResume)
		candidate.POST("/resumes/upload", cfg.CandidatesHandler.UploadResumePDF)
		candidate.GET("/resumes", cfg.CandidatesHandler.ListResumes)
		candidate.GET("/resumes/:resume_id", cfg.CandidatesHandler.GetResume)
		candidate.POST("/applications", cfg.ATSHandler.Apply)
		candidate.GET("/applications", cfg.ATSHandler.ListMyApplications)
		candidate.GET("/applications/:job_id/status", cfg.ATSHandler.GetApplicationStatus)   // candidate_id from path
		candidate.GET("/applications/:job_id/feedback", cfg.ATSHandler.GetApplicationFeedback) // AI feedback (safe subset)
		candidate.POST("/interviews/:id/confirm", cfg.InterviewsHandler.ConfirmByCandidate)   // confirm interview (Phase 3)
		candidate.GET("/offers/:id", cfg.OffersHandler.GetByIDForCandidate)                  // view offer / download letter (Phase 3)
		candidate.POST("/offers/:id/accept", cfg.OffersHandler.Accept)
		candidate.POST("/offers/:id/reject", cfg.OffersHandler.Reject)
	}

	return r
}

// Run starts the Gin server (blocking). Prefer main with NewRouter + http.Server for graceful shutdown.
func Run(cfg Config) error {
	r := NewRouter(cfg)
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return srv.ListenAndServe()
}
