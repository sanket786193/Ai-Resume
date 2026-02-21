package jobs

import (
	"context"
	"time"

	"resume/internal/domain/entities"
	domainerrors "resume/internal/domain/errors"
	"resume/internal/domain/enums"
	"resume/internal/storage/postgres"

	"github.com/google/uuid"
)

// Service contains job business logic; no HTTP.
type Service struct {
	repo *postgres.JobRepo
}

// NewService creates a jobs service.
func NewService(repo *postgres.JobRepo) *Service {
	return &Service{repo: repo}
}

// Create creates a job in DRAFT; createdBy must be HR user ID.
func (s *Service) Create(ctx context.Context, title, description, location, department, createdBy string) (*entities.Job, error) {
	if title == "" || description == "" {
		return nil, &domainerrors.ValidationError{Field: "title/description", Message: "required"}
	}
	now := time.Now()
	j := &entities.Job{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Location:    location,
		Department:  department,
		Status:      enums.JobDraft,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.Create(ctx, j); err != nil {
		return nil, err
	}
	return j, nil
}

// GetByID returns a job by ID.
func (s *Service) GetByID(ctx context.Context, id string) (*entities.Job, error) {
	j, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, err
		}
		return nil, &domainerrors.NotFoundError{Resource: "job", ID: id}
	}
	return j, nil
}

// List returns jobs with optional status filter and pagination.
func (s *Service) List(ctx context.Context, status *enums.JobStatus, limit, offset int) ([]*entities.Job, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.List(ctx, status, limit, offset)
}

// Publish transitions job from DRAFT to PUBLISHED.
func (s *Service) Publish(ctx context.Context, id string) error {
	j, err := s.repo.GetByID(ctx, id)
	if err != nil || j == nil {
		return &domainerrors.NotFoundError{Resource: "job", ID: id}
	}
	if j.Status != enums.JobDraft {
		return &domainerrors.ValidationError{Field: "status", Message: "only DRAFT jobs can be published"}
	}
	return s.repo.UpdateStatus(ctx, id, enums.JobPublished)
}

// Update updates job fields; only DRAFT can be fully updated.
func (s *Service) Update(ctx context.Context, j *entities.Job) error {
	existing, err := s.repo.GetByID(ctx, j.ID)
	if err != nil || existing == nil {
		return &domainerrors.NotFoundError{Resource: "job", ID: j.ID}
	}
	if existing.Status != enums.JobDraft {
		return &domainerrors.ValidationError{Field: "status", Message: "only DRAFT jobs can be updated"}
	}
	return s.repo.Update(ctx, j)
}

// Close transitions job to CLOSED.
func (s *Service) Close(ctx context.Context, id string) error {
	j, err := s.repo.GetByID(ctx, id)
	if err != nil || j == nil {
		return &domainerrors.NotFoundError{Resource: "job", ID: id}
	}
	if j.Status == enums.JobClosed {
		return &domainerrors.ConflictError{Message: "job already closed"}
	}
	return s.repo.UpdateStatus(ctx, id, enums.JobClosed)
}

// SoftDelete soft-deletes a job (HR only; enforce in handler).
func (s *Service) SoftDelete(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return &domainerrors.NotFoundError{Resource: "job", ID: id}
	}
	return s.repo.SoftDelete(ctx, id)
}
