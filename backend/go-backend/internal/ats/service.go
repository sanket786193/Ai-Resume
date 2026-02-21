package ats

import (
	"context"
	"time"

	"resume/internal/domain/entities"
	domainerrors "resume/internal/domain/errors"
	"resume/internal/domain/enums"
	"resume/internal/storage/postgres"

	"github.com/google/uuid"
)

// AIClient interface for resume screening and parsing (advisory only).
type AIClient interface {
	ScreenResume(ctx context.Context, resumeContentOrPath string, jobDescription string) (*AIScreenResult, error)
	Parse(ctx context.Context, resumePathOrContent string) (rawText string, parsedJSON []byte, cleanedText string, err error)
}

// AIScreenResult from Python AI service (full ATS evaluation).
type AIScreenResult struct {
	SkillMatchScore   float64
	RankingScore      float64
	Qualified         bool
	ATSScore          *int
	SkillMatchPct     *int
	MissingSkills     []string
	ExperienceMatch   *string
	Summary           *string
	ModelVersion      *string
}

// Service contains ATS application and status flow logic.
type Service struct {
	atsRepo       *postgres.ATSRepo
	jobRepo       *postgres.JobRepo
	resumeRepo    *postgres.ResumeRepo
	parsedRepo    *postgres.ResumeParsedRepo
	candidateRepo *postgres.CandidateRepo
	ai            AIClient
	aiEnabled     bool
}

// NewService creates an ATS service.
func NewService(
	atsRepo *postgres.ATSRepo,
	jobRepo *postgres.JobRepo,
	resumeRepo *postgres.ResumeRepo,
	parsedRepo *postgres.ResumeParsedRepo,
	candidateRepo *postgres.CandidateRepo,
	ai AIClient,
	aiEnabled bool,
) *Service {
	return &Service{
		atsRepo:       atsRepo,
		jobRepo:       jobRepo,
		resumeRepo:    resumeRepo,
		parsedRepo:    parsedRepo,
		candidateRepo: candidateRepo,
		ai:            ai,
		aiEnabled:     aiEnabled,
	}
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
		go s.runAIScreening(context.Background(), rec.ID, resume.ID, resume.StoragePath, job.Description)
	}
	return rec, nil
}

// runAIScreening calls Parse (persists to resume_parsed_data), then AI screen; errors are logged, not returned.
func (s *Service) runAIScreening(ctx context.Context, atsID, resumeID, resumePath, jobDescription string) {
	// 1) Parse resume and store in resume_parsed_data
	if s.parsedRepo != nil {
		rawText, parsedJSON, cleanedText, err := s.ai.Parse(ctx, resumePath)
		if err == nil && (rawText != "" || len(parsedJSON) > 0 || cleanedText != "") {
			_ = s.parsedRepo.Upsert(ctx, resumeID, rawText, parsedJSON, cleanedText)
		}
	}
	// 2) Screen resume against job (AI service fetches PDF from URL if needed)
	result, err := s.ai.ScreenResume(ctx, resumePath, jobDescription)
	if err != nil {
		return
	}
	now := time.Now()
	_ = s.atsRepo.UpdateAIScores(ctx, atsID, &result.SkillMatchScore, &result.RankingScore, &result.Qualified, now)
	if result.ATSScore != nil || result.SkillMatchPct != nil || len(result.MissingSkills) > 0 || result.ExperienceMatch != nil || result.Summary != nil || result.ModelVersion != nil {
		_ = s.atsRepo.UpdateAIFeedback(ctx, atsID, result.ATSScore, result.SkillMatchPct, result.MissingSkills, result.ExperienceMatch, result.Summary, result.ModelVersion, now)
	}
	_ = s.atsRepo.UpdateStatus(ctx, atsID, enums.ATSScreening)
}

// GetApplicationStatus returns ATS record for candidate+job (candidate view).
func (s *Service) GetApplicationStatus(ctx context.Context, jobID, candidateID string) (*entities.ATSRecord, error) {
	rec, err := s.atsRepo.GetByJobAndCandidate(ctx, jobID, candidateID)
	if err != nil || rec == nil {
		return nil, &domainerrors.NotFoundError{Resource: "application", ID: jobID + "/" + candidateID}
	}
	return rec, nil
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

// UpdateStatus updates ATS status (HR); validate transitions in service.
func (s *Service) UpdateStatus(ctx context.Context, atsID string, status enums.ATSStatus) error {
	if !status.Valid() {
		return &domainerrors.ValidationError{Field: "status", Message: "invalid ATS status"}
	}
	rec, err := s.atsRepo.GetByID(ctx, atsID)
	if err != nil || rec == nil {
		return &domainerrors.NotFoundError{Resource: "ats_record", ID: atsID}
	}
	return s.atsRepo.UpdateStatus(ctx, atsID, status)
}

// GetByID returns an ATS record by ID.
func (s *Service) GetByID(ctx context.Context, id string) (*entities.ATSRecord, error) {
	rec, err := s.atsRepo.GetByID(ctx, id)
	if err != nil || rec == nil {
		return nil, &domainerrors.NotFoundError{Resource: "ats_record", ID: id}
	}
	return rec, nil
}
