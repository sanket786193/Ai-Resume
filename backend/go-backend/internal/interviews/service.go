package interviews

import (
	"context"
	"time"

	"resume/internal/domain/entities"
	domainerrors "resume/internal/domain/errors"
	"resume/internal/storage/postgres"

	"github.com/google/uuid"
)

// Service contains interview scheduling logic.
type Service struct {
	repo   *postgres.InterviewRepo
	atsRepo *postgres.ATSRepo
}

// NewService creates an interviews service.
func NewService(repo *postgres.InterviewRepo, atsRepo *postgres.ATSRepo) *Service {
	return &Service{repo: repo, atsRepo: atsRepo}
}

// Schedule creates an interview tied to an ATS record.
func (s *Service) Schedule(ctx context.Context, atsID string, scheduledAt time.Time, durationMin int, location string, round int, notes string) (*entities.Interview, error) {
	_, err := s.atsRepo.GetByID(ctx, atsID)
	if err != nil {
		return nil, &domainerrors.NotFoundError{Resource: "ats_record", ID: atsID}
	}
	now := time.Now()
	i := &entities.Interview{
		ID:           uuid.New().String(),
		ATSID:        atsID,
		ScheduledAt:  scheduledAt,
		Duration:     durationMin,
		Location:     location,
		Round:        round,
		Notes:        notes,
		Status:       "SCHEDULED",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if i.Duration <= 0 {
		i.Duration = 60
	}
	if i.Round <= 0 {
		i.Round = 1
	}
	if err := s.repo.Create(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}

// GetByID returns an interview by ID.
func (s *Service) GetByID(ctx context.Context, id string) (*entities.Interview, error) {
	i, err := s.repo.GetByID(ctx, id)
	if err != nil || i == nil {
		return nil, &domainerrors.NotFoundError{Resource: "interview", ID: id}
	}
	return i, nil
}

// ListByATSID returns interviews for an ATS record.
func (s *Service) ListByATSID(ctx context.Context, atsID string) ([]*entities.Interview, error) {
	return s.repo.ListByATSID(ctx, atsID)
}

// Update updates an interview.
func (s *Service) Update(ctx context.Context, i *entities.Interview) error {
	existing, err := s.repo.GetByID(ctx, i.ID)
	if err != nil || existing == nil {
		return &domainerrors.NotFoundError{Resource: "interview", ID: i.ID}
	}
	i.UpdatedAt = time.Now()
	return s.repo.Update(ctx, i)
}
