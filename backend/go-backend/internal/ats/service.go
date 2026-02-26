package ats

import (
	"context"
	"log"
	"strings"
	"time"

	"resume/internal/ai/client"
	"resume/internal/domain/entities"
	domainerrors "resume/internal/domain/errors"
	"resume/internal/domain/enums"
	"resume/internal/storage/postgres"

	"github.com/google/uuid"
)

// AIClient interface for resume screening and parsing (advisory only).
type AIClient interface {
	ScreenResume(ctx context.Context, resumeContentOrPath, jobDescription string, jobRequirements *client.JobRequirements) (*AIScreenResult, error)
	Parse(ctx context.Context, resumePathOrContent string) (rawText string, parsedJSON []byte, cleanedText string, err error)
}

// AIScreenResult from Python AI service (full ATS evaluation).
type AIScreenResult struct {
	SkillMatchScore    float64
	RankingScore       float64
	Qualified          bool
	ATSScore           *int
	SkillMatchPct      *int
	MissingSkills      []string
	ExperienceMatch    *string
	ExperienceWarnings []string
	KeywordMatches     []string
	SemanticMatches    []string
	Summary            *string
	ModelVersion       *string
}

// ApplicationForHR is an ATS record enriched with candidate and job display info for the HR pipeline.
type ApplicationForHR struct {
	ID              string    `json:"id"`
	JobID           string    `json:"job_id"`
	CandidateID     string    `json:"candidate_id"`
	ResumeID        string    `json:"resume_id"`
	Status          string    `json:"status"`
	SkillMatchScore *float64  `json:"skill_match_score,omitempty"`
	ATSScore        *int      `json:"ats_score,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CandidateName   string    `json:"candidate_name"`
	CandidateEmail  string    `json:"candidate_email"`
	JobTitle        string    `json:"job_title"`
}

// ApplicationDetailForHR is a single application with full candidate and AI details for HR detail view.
type ApplicationDetailForHR struct {
	ApplicationForHR
	CandidatePhone     string    `json:"candidate_phone,omitempty"`
	CandidateLinkedIn  string    `json:"candidate_linkedin,omitempty"`
	ResumeFileName     string    `json:"resume_file_name,omitempty"`
	ResumeURL          string    `json:"resume_url,omitempty"` // public URL for PDF viewer
	AIProcessedAt     *time.Time `json:"ai_processed_at,omitempty"`
	SkillMatchPct     *int       `json:"skill_match_pct,omitempty"`
	AISummary         string     `json:"ai_summary,omitempty"`
	ExperienceMatch   string     `json:"experience_match,omitempty"`
	ExperienceWarnings []string  `json:"experience_warnings,omitempty"`
	MissingSkills     []string   `json:"missing_skills,omitempty"`
	KeywordMatches    []string   `json:"keyword_matches,omitempty"`
	SemanticMatches   []string   `json:"semantic_matches,omitempty"`
	Qualified         *bool      `json:"qualified,omitempty"`
}

// ResumeURLResolver converts a storage path or key into a fetchable public URL for the AI service.
// If nil, the raw path is passed through (must be a full URL for screening to work).
type ResumeURLResolver func(storagePath string) string

// ATSNotifier sends candidate notifications on status change (Phase 3). Optional; if nil, no emails sent.
type ATSNotifier interface {
	SendShortlisted(ctx context.Context, candidateEmail, jobTitle string)
	SendRejected(ctx context.Context, candidateEmail, jobTitle string)
}

// Service contains ATS application and status flow logic.
type Service struct {
	atsRepo           *postgres.ATSRepo
	jobRepo           *postgres.JobRepo
	resumeRepo        *postgres.ResumeRepo
	parsedRepo        *postgres.ResumeParsedRepo
	candidateRepo     *postgres.CandidateRepo
	userRepo          *postgres.UserRepo
	ai                AIClient
	aiEnabled         bool
	resumeURLResolver ResumeURLResolver
	notifier          ATSNotifier
}

// NewService creates an ATS service.
func NewService(
	atsRepo *postgres.ATSRepo,
	jobRepo *postgres.JobRepo,
	resumeRepo *postgres.ResumeRepo,
	parsedRepo *postgres.ResumeParsedRepo,
	candidateRepo *postgres.CandidateRepo,
	userRepo *postgres.UserRepo,
	ai AIClient,
	aiEnabled bool,
) *Service {
	return &Service{
		atsRepo:       atsRepo,
		jobRepo:       jobRepo,
		resumeRepo:    resumeRepo,
		parsedRepo:    parsedRepo,
		candidateRepo: candidateRepo,
		userRepo:      userRepo,
		ai:            ai,
		aiEnabled:     aiEnabled,
	}
}

// NewServiceWithResolver creates an ATS service with an optional resolver to turn storage paths into public URLs.
func NewServiceWithResolver(
	atsRepo *postgres.ATSRepo,
	jobRepo *postgres.JobRepo,
	resumeRepo *postgres.ResumeRepo,
	parsedRepo *postgres.ResumeParsedRepo,
	candidateRepo *postgres.CandidateRepo,
	userRepo *postgres.UserRepo,
	ai AIClient,
	aiEnabled bool,
	resumeURLResolver ResumeURLResolver,
) *Service {
	s := NewService(atsRepo, jobRepo, resumeRepo, parsedRepo, candidateRepo, userRepo, ai, aiEnabled)
	s.resumeURLResolver = resumeURLResolver
	return s
}

// SetNotifier sets the optional notifier for shortlist/rejection emails (Phase 3).
func (s *Service) SetNotifier(n ATSNotifier) {
	s.notifier = n
}

// Apply creates an ATS record for candidate+job+resume; triggers AI screening async (non-blocking).
func (s *Service) Apply(ctx context.Context, jobID, candidateID, resumeID string) (*entities.ATSRecord, error) {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil || job == nil {
		return nil, &domainerrors.NotFoundError{Resource: "job", ID: jobID}
	}
	if job.Status != enums.JobPublished {
		return nil, &domainerrors.ValidationError{Field: "job", Message: "job is not published"}
	}
	candidate, err := s.candidateRepo.GetByID(ctx, candidateID)
	if err != nil || candidate == nil {
		return nil, &domainerrors.NotFoundError{Resource: "candidate", ID: candidateID}
	}
	resume, err := s.resumeRepo.GetByID(ctx, resumeID)
	if err != nil || resume == nil {
		return nil, &domainerrors.NotFoundError{Resource: "resume", ID: resumeID}
	}
	if resume.CandidateID != candidateID {
		return nil, &domainerrors.ValidationError{Field: "resume_id", Message: "resume does not belong to candidate"}
	}
	existing, _ := s.atsRepo.GetByJobAndCandidate(ctx, jobID, candidateID)
	if existing != nil {
		return nil, &domainerrors.ConflictError{Message: "already applied to this job"}
	}
	now := time.Now()
	rec := &entities.ATSRecord{
		ID:          uuid.New().String(),
		JobID:       jobID,
		CandidateID: candidateID,
		ResumeID:    resumeID,
		Status:      enums.ATSApplied,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.atsRepo.Create(ctx, rec); err != nil {
		return nil, err
	}
	// Trigger AI screening in background: parse resume (store in resume_parsed_data), then screen.
	if s.aiEnabled && s.ai != nil {
		go s.runAIScreening(context.Background(), rec.ID, resume.ID, resume.StoragePath, job)
	}
	return rec, nil
}

// BulkApplyItem is one (candidate_id, resume_id) for bulk apply (Phase 4).
type BulkApplyItem struct {
	CandidateID string `json:"candidate_id"`
	ResumeID    string `json:"resume_id"`
}

// BulkApplyResult is the result of bulk apply (Phase 4).
type BulkApplyResult struct {
	Created int      `json:"created"`
	Skipped int      `json:"skipped"` // already applied
	Errors  []string `json:"errors,omitempty"`
}

// BulkApply creates ATS records for multiple candidate+resume pairs for a job (HR). hrUserID must own the job. Skips already-applied; triggers AI for each new record (Phase 4).
func (s *Service) BulkApply(ctx context.Context, jobID, hrUserID string, items []BulkApplyItem) (*BulkApplyResult, error) {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil || job == nil {
		return nil, &domainerrors.NotFoundError{Resource: "job", ID: jobID}
	}
	if job.CreatedBy != hrUserID {
		return nil, &domainerrors.ForbiddenError{Message: "job not owned by you"}
	}
	if job.Status != enums.JobPublished {
		return nil, &domainerrors.ValidationError{Field: "job", Message: "job is not published"}
	}
	result := &BulkApplyResult{}
	for _, item := range items {
		if item.CandidateID == "" || item.ResumeID == "" {
			result.Errors = append(result.Errors, "missing candidate_id or resume_id")
			continue
		}
		_, err := s.Apply(ctx, jobID, item.CandidateID, item.ResumeID)
		if err != nil {
			if _, ok := err.(*domainerrors.ConflictError); ok {
				result.Skipped++
			} else {
				result.Errors = append(result.Errors, err.Error())
			}
			continue
		}
		result.Created++
	}
	return result, nil
}

// runAIScreening calls Parse (persists to resume_parsed_data), then AI screen; errors are logged, not returned.
func (s *Service) runAIScreening(ctx context.Context, atsID, resumeID, resumePath string, job *entities.Job) {
	resumeInput := resumePath
	if s.resumeURLResolver != nil && resumePath != "" && !strings.HasPrefix(strings.TrimSpace(resumePath), "http://") && !strings.HasPrefix(strings.TrimSpace(resumePath), "https://") {
		if u := s.resumeURLResolver(resumePath); u != "" {
			resumeInput = u
		}
	}

	// 1) Parse resume and store in resume_parsed_data; mark resume as PROCESSED when done (Phase 3).
	if s.parsedRepo != nil {
		rawText, parsedJSON, cleanedText, err := s.ai.Parse(ctx, resumeInput)
		if err != nil {
			log.Printf("[ATS] AI parse failed for atsID=%s resumeID=%s: %v", atsID, resumeID, err)
		} else if rawText != "" || len(parsedJSON) > 0 || cleanedText != "" {
			if err := s.parsedRepo.Upsert(ctx, resumeID, rawText, parsedJSON, cleanedText); err != nil {
				log.Printf("[ATS] failed to persist parsed resume resumeID=%s: %v", resumeID, err)
			} else {
				_ = s.resumeRepo.UpdateStatus(ctx, resumeID, enums.ResumeProcessed)
			}
		}
	}

	// 2) Build optional job requirements for matching (Phase 2)
	var req *client.JobRequirements
	if len(job.Skills) > 0 || job.Qualification != "" || (job.ExperienceLevel != "" && job.ExperienceLevel != enums.ExperienceAny) {
		req = &client.JobRequirements{
			Skills:          job.Skills,
			ExperienceLevel: string(job.ExperienceLevel),
			Qualification:   job.Qualification,
		}
	}

	// 3) Screen resume against job (AI service uses requirements when provided)
	result, err := s.ai.ScreenResume(ctx, resumeInput, job.Description, req)
	if err != nil {
		log.Printf("[ATS] AI screen failed for atsID=%s: %v", atsID, err)
		return
	}

	now := time.Now()
	_ = s.atsRepo.UpdateAIScores(ctx, atsID, &result.SkillMatchScore, &result.RankingScore, &result.Qualified, now)

	// Persist full feedback when present; otherwise derive from scores so DB columns are populated (e.g. heuristic path)
	atsScore := result.ATSScore
	skillMatchPct := result.SkillMatchPct
	if atsScore == nil && (result.SkillMatchScore >= 0 && result.SkillMatchScore <= 1) {
		v := int(result.SkillMatchScore * 100)
		if v > 100 {
			v = 100
		}
		atsScore = &v
	}
	if skillMatchPct == nil && (result.SkillMatchScore >= 0 && result.SkillMatchScore <= 1) {
		v := int(result.SkillMatchScore * 100)
		if v > 100 {
			v = 100
		}
		skillMatchPct = &v
	}
	hasFeedback := atsScore != nil || skillMatchPct != nil || len(result.MissingSkills) > 0 ||
		result.ExperienceMatch != nil || result.Summary != nil || result.ModelVersion != nil ||
		len(result.ExperienceWarnings) > 0 || len(result.KeywordMatches) > 0 || len(result.SemanticMatches) > 0
	if hasFeedback {
		_ = s.atsRepo.UpdateAIFeedback(ctx, atsID, atsScore, skillMatchPct, result.MissingSkills, result.ExperienceMatch, result.Summary, result.ModelVersion, result.ExperienceWarnings, result.KeywordMatches, result.SemanticMatches, now)
	}

	// 4) Auto-shortlist when qualified and under vacancy limit (Phase 2); otherwise set to SCREENING
	nextStatus := enums.ATSScreening
	if result.Qualified && len(job.VacancyLimits) > 0 {
		totalSlots := 0
		for _, v := range job.VacancyLimits {
			totalSlots += v.Limit
		}
		shortlistedCount, _ := s.atsRepo.CountByJobIDAndStatus(ctx, job.ID, enums.ATSShortlisted)
		if shortlistedCount < totalSlots {
			nextStatus = enums.ATSShortlisted
		}
	}
	_ = s.atsRepo.UpdateStatus(ctx, atsID, nextStatus)
}

// GetApplicationStatus returns ATS record for candidate+job (candidate view).
func (s *Service) GetApplicationStatus(ctx context.Context, jobID, candidateID string) (*entities.ATSRecord, error) {
	rec, err := s.atsRepo.GetByJobAndCandidate(ctx, jobID, candidateID)
	if err != nil || rec == nil {
		return nil, &domainerrors.NotFoundError{Resource: "application", ID: jobID + "/" + candidateID}
	}
	return rec, nil
}

// ApplicationFeedbackForCandidate is the safe subset of AI feedback shown to candidates (Phase 3).
type ApplicationFeedbackForCandidate struct {
	SkillMatchPct   *int     `json:"skill_match_pct,omitempty"`
	Summary         string   `json:"summary,omitempty"`
	ImprovementTips []string `json:"improvement_tips,omitempty"` // from missing_skills as suggestions
}

// GetApplicationFeedbackForCandidate returns AI feedback for the candidate's own application (safe subset only).
func (s *Service) GetApplicationFeedbackForCandidate(ctx context.Context, jobID, candidateID string) (*ApplicationFeedbackForCandidate, error) {
	rec, err := s.atsRepo.GetByJobAndCandidate(ctx, jobID, candidateID)
	if err != nil || rec == nil {
		return nil, &domainerrors.NotFoundError{Resource: "application", ID: jobID + "/" + candidateID}
	}
	out := &ApplicationFeedbackForCandidate{}
	if rec.SkillMatchPct != nil {
		out.SkillMatchPct = rec.SkillMatchPct
	}
	if rec.AISummary != nil && *rec.AISummary != "" {
		out.Summary = *rec.AISummary
	}
	if len(rec.MissingSkills) > 0 {
		out.ImprovementTips = make([]string, 0, len(rec.MissingSkills))
		for _, skill := range rec.MissingSkills {
			if skill != "" {
				out.ImprovementTips = append(out.ImprovementTips, "Consider highlighting: "+skill)
			}
		}
	}
	return out, nil
}

// ListByCandidate returns ATS records for a candidate.
func (s *Service) ListByCandidate(ctx context.Context, candidateID string, limit, offset int) ([]*entities.ATSRecord, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.atsRepo.ListByCandidateID(ctx, candidateID, limit, offset)
}

// ApplicationWithJob is an ATS record enriched with job title and status for candidate UI.
type ApplicationWithJob struct {
	*entities.ATSRecord
	JobTitle  string `json:"job_title"`
	JobStatus string `json:"job_status"`
}

// ListByCandidateWithJobs returns applications with job title and job status for the candidate applications page.
func (s *Service) ListByCandidateWithJobs(ctx context.Context, candidateID string, limit, offset int) ([]*ApplicationWithJob, error) {
	list, err := s.ListByCandidate(ctx, candidateID, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]*ApplicationWithJob, 0, len(list))
	for _, rec := range list {
		job, _ := s.jobRepo.GetByID(ctx, rec.JobID)
		title := ""
		status := ""
		if job != nil {
			title = job.Title
			status = string(job.Status)
		}
		out = append(out, &ApplicationWithJob{ATSRecord: rec, JobTitle: title, JobStatus: status})
	}
	return out, nil
}

// ListByJob returns ATS records for a job (HR).
func (s *Service) ListByJob(ctx context.Context, jobID string, status *enums.ATSStatus, limit, offset int) ([]*entities.ATSRecord, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.atsRepo.ListByJobID(ctx, jobID, status, limit, offset)
}

// ListByJobEnriched returns applications for a job with candidate name, email, and job title (HR applicants list).
func (s *Service) ListByJobEnriched(ctx context.Context, jobID string, status *enums.ATSStatus, limit, offset int) ([]*ApplicationForHR, error) {
	list, err := s.ListByJob(ctx, jobID, status, limit, offset)
	if err != nil || len(list) == 0 {
		return nil, err
	}
	return s.enrichApplications(ctx, list)
}

// ListForHR returns ATS records for all jobs created by the HR user, or for a single job if jobID is set.
func (s *Service) ListForHR(ctx context.Context, hrUserID string, jobID *string, status *enums.ATSStatus, limit, offset int) ([]*entities.ATSRecord, error) {
	if limit <= 0 {
		limit = 20
	}
	if jobID != nil && *jobID != "" {
		job, err := s.jobRepo.GetByID(ctx, *jobID)
		if err != nil || job == nil {
			return nil, &domainerrors.NotFoundError{Resource: "job", ID: *jobID}
		}
		if job.CreatedBy != hrUserID {
			return nil, &domainerrors.ForbiddenError{Message: "job not owned by you"}
		}
		return s.atsRepo.ListByJobID(ctx, *jobID, status, limit, offset)
	}
	jobIDs, err := s.jobRepo.ListIDsByCreatedBy(ctx, hrUserID)
	if err != nil {
		return nil, err
	}
	return s.atsRepo.ListByJobIDs(ctx, jobIDs, status, limit, offset)
}

// ListForHREnriched returns applications with candidate name, email, and job title for the HR pipeline.
func (s *Service) ListForHREnriched(ctx context.Context, hrUserID string, jobID *string, status *enums.ATSStatus, limit, offset int) ([]*ApplicationForHR, error) {
	list, err := s.ListForHR(ctx, hrUserID, jobID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return s.enrichApplications(ctx, list)
}

// enrichApplications converts ATS records into ApplicationForHR with candidate and job details.
func (s *Service) enrichApplications(ctx context.Context, list []*entities.ATSRecord) ([]*ApplicationForHR, error) {
	out := make([]*ApplicationForHR, 0, len(list))
	for _, rec := range list {
		app := &ApplicationForHR{
			ID:           rec.ID,
			JobID:        rec.JobID,
			CandidateID:  rec.CandidateID,
			ResumeID:     rec.ResumeID,
			Status:       string(rec.Status),
			SkillMatchScore: rec.SkillMatchScore,
			ATSScore:     rec.ATSScore,
			CreatedAt:    rec.CreatedAt,
			UpdatedAt:    rec.UpdatedAt,
		}
		candidate, _ := s.candidateRepo.GetByID(ctx, rec.CandidateID)
		if candidate != nil && s.userRepo != nil {
			user, _ := s.userRepo.GetByID(ctx, candidate.UserID)
			if user != nil {
				app.CandidateName = user.Name
				app.CandidateEmail = user.Email
			}
		}
		job, _ := s.jobRepo.GetByID(ctx, rec.JobID)
		if job != nil {
			app.JobTitle = job.Title
		}
		out = append(out, app)
	}
	return out, nil
}

// UpdateStatus updates ATS status (HR); validate transitions in service. Sends notification when shortlisted/rejected (Phase 3).
func (s *Service) UpdateStatus(ctx context.Context, atsID string, status enums.ATSStatus) error {
	if !status.Valid() {
		return &domainerrors.ValidationError{Field: "status", Message: "invalid ATS status"}
	}
	rec, err := s.atsRepo.GetByID(ctx, atsID)
	if err != nil || rec == nil {
		return &domainerrors.NotFoundError{Resource: "ats_record", ID: atsID}
	}
	if err := s.atsRepo.UpdateStatus(ctx, atsID, status); err != nil {
		return err
	}
	if s.notifier != nil && (status == enums.ATSShortlisted || status == enums.ATSRejected) {
		candidate, _ := s.candidateRepo.GetByID(ctx, rec.CandidateID)
		if candidate != nil {
			user, _ := s.userRepo.GetByID(ctx, candidate.UserID)
			job, _ := s.jobRepo.GetByID(ctx, rec.JobID)
			jobTitle := ""
			if job != nil {
				jobTitle = job.Title
			}
			if user != nil && user.Email != "" {
				if status == enums.ATSShortlisted {
					s.notifier.SendShortlisted(ctx, user.Email, jobTitle)
				} else {
					s.notifier.SendRejected(ctx, user.Email, jobTitle)
				}
			}
		}
	}
	return nil
}

// GetByID returns an ATS record by ID.
func (s *Service) GetByID(ctx context.Context, id string) (*entities.ATSRecord, error) {
	rec, err := s.atsRepo.GetByID(ctx, id)
	if err != nil || rec == nil {
		return nil, &domainerrors.NotFoundError{Resource: "ats_record", ID: id}
	}
	return rec, nil
}

// GetApplicationByIDEnriched returns a single application with full candidate and AI details for HR.
func (s *Service) GetApplicationByIDEnriched(ctx context.Context, id string) (*ApplicationDetailForHR, error) {
	rec, err := s.atsRepo.GetByID(ctx, id)
	if err != nil || rec == nil {
		return nil, &domainerrors.NotFoundError{Resource: "application", ID: id}
	}
	list, err := s.enrichApplications(ctx, []*entities.ATSRecord{rec})
	if err != nil || len(list) == 0 {
		return nil, err
	}
	detail := &ApplicationDetailForHR{ApplicationForHR: *list[0]}
	candidate, _ := s.candidateRepo.GetByID(ctx, rec.CandidateID)
	if candidate != nil {
		detail.CandidatePhone = candidate.Phone
		detail.CandidateLinkedIn = candidate.LinkedIn
	}
	if rec.ResumeID != "" {
		resume, _ := s.resumeRepo.GetByID(ctx, rec.ResumeID)
		if resume != nil {
			detail.ResumeFileName = resume.FileName
			if s.resumeURLResolver != nil {
				detail.ResumeURL = s.resumeURLResolver(resume.StoragePath)
			}
		}
	}
	detail.AIProcessedAt = rec.AIProcessedAt
	detail.SkillMatchPct = rec.SkillMatchPct
	if rec.AISummary != nil {
		detail.AISummary = *rec.AISummary
	}
	if rec.ExperienceMatch != nil {
		detail.ExperienceMatch = *rec.ExperienceMatch
	}
	if len(rec.MissingSkills) > 0 {
		detail.MissingSkills = rec.MissingSkills
	}
	if len(rec.ExperienceWarnings) > 0 {
		detail.ExperienceWarnings = rec.ExperienceWarnings
	}
	if len(rec.KeywordMatches) > 0 {
		detail.KeywordMatches = rec.KeywordMatches
	}
	if len(rec.SemanticMatches) > 0 {
		detail.SemanticMatches = rec.SemanticMatches
	}
	detail.Qualified = rec.Qualified
	return detail, nil
}

// GetApplicationResumeURL returns the public URL for an application's resume PDF (for proxy/viewer). Returns empty string if no resume or no resolver.
func (s *Service) GetApplicationResumeURL(ctx context.Context, applicationID string) (string, error) {
	rec, err := s.atsRepo.GetByID(ctx, applicationID)
	if err != nil || rec == nil {
		return "", &domainerrors.NotFoundError{Resource: "application", ID: applicationID}
	}
	if rec.ResumeID == "" || s.resumeURLResolver == nil {
		return "", nil
	}
	resume, err := s.resumeRepo.GetByID(ctx, rec.ResumeID)
	if err != nil || resume == nil {
		return "", nil
	}
	return s.resumeURLResolver(resume.StoragePath), nil
}
