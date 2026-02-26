package candidates

import (
	"context"
	"time"

	"resume/internal/domain/entities"
	domainerrors "resume/internal/domain/errors"
	"resume/internal/domain/enums"
	"resume/internal/storage/postgres"

	"github.com/google/uuid"
)

// Service contains candidate application and resume logic.
type Service struct {
	candidateRepo *postgres.CandidateRepo
	resumeRepo    *postgres.ResumeRepo
	atsRepo       *postgres.ATSRepo
	jobRepo       *postgres.JobRepo
}

// NewService creates a candidates service.
func NewService(
	candidateRepo *postgres.CandidateRepo,
	resumeRepo *postgres.ResumeRepo,
	atsRepo *postgres.ATSRepo,
	jobRepo *postgres.JobRepo,
) *Service {
	return &Service{
		candidateRepo: candidateRepo,
		resumeRepo:    resumeRepo,
		atsRepo:       atsRepo,
		jobRepo:       jobRepo,
	}
}

// GetCandidateByUserID returns the candidate profile for a user (CANDIDATE role).
func (s *Service) GetCandidateByUserID(ctx context.Context, userID string) (*entities.Candidate, error) {
	c, err := s.candidateRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, &domainerrors.NotFoundError{Resource: "candidate", ID: userID}
	}
	return c, nil
}

// AddResume stores resume metadata; caller stores file and provides storagePath.
func (s *Service) AddResume(ctx context.Context, candidateID, fileName, storagePath, mimeType string, fileSize int64) (*entities.Resume, error) {
	c, err := s.candidateRepo.GetByID(ctx, candidateID)
	if err != nil || c == nil {
		return nil, &domainerrors.NotFoundError{Resource: "candidate", ID: candidateID}
	}
	now := time.Now()
	res := &entities.Resume{
		ID:          uuid.New().String(),
		CandidateID: candidateID,
		FileName:    fileName,
		StoragePath: storagePath,
		FileSize:    fileSize,
		MimeType:    mimeType,
		Status:      enums.ResumePending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.resumeRepo.Create(ctx, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListResumes returns resumes for a candidate.
func (s *Service) ListResumes(ctx context.Context, candidateID string) ([]*entities.Resume, error) {
	return s.resumeRepo.ListByCandidateID(ctx, candidateID)
}

// GetResumeByID returns a resume by ID if it belongs to the candidate.
func (s *Service) GetResumeByID(ctx context.Context, resumeID, candidateID string) (*entities.Resume, error) {
	res, err := s.resumeRepo.GetByID(ctx, resumeID)
	if err != nil || res == nil {
		return nil, &domainerrors.NotFoundError{Resource: "resume", ID: resumeID}
	}
	if res.CandidateID != candidateID {
		return nil, &domainerrors.ForbiddenError{Message: "resume not found"}
	}
	return res, nil
}
