package interviews

import (
	"context"
	"time"

	"resume/internal/domain/entities"
	domainerrors "resume/internal/domain/errors"
	"resume/internal/storage/postgres"

	"github.com/google/uuid"
)

// InterviewNotifier sends candidate notification when interview is scheduled (Phase 3). Optional.
type InterviewNotifier interface {
	SendInterviewScheduled(ctx context.Context, candidateEmail, jobTitle string)
}

// Service contains interview scheduling logic.
type Service struct {
	repo         *postgres.InterviewRepo
	atsRepo      *postgres.ATSRepo
	candidateRepo *postgres.CandidateRepo
	userRepo     *postgres.UserRepo
	jobRepo      *postgres.JobRepo
	notifier     InterviewNotifier
}

// NewService creates an interviews service.
func NewService(repo *postgres.InterviewRepo, atsRepo *postgres.ATSRepo, candidateRepo *postgres.CandidateRepo, userRepo *postgres.UserRepo, jobRepo *postgres.JobRepo, notifier InterviewNotifier) *Service {
	return &Service{repo: repo, atsRepo: atsRepo, candidateRepo: candidateRepo, userRepo: userRepo, jobRepo: jobRepo, notifier: notifier}
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
	if s.notifier != nil && s.candidateRepo != nil && s.userRepo != nil && s.jobRepo != nil {
		rec, _ := s.atsRepo.GetByID(ctx, atsID)
		if rec != nil {
			candidate, _ := s.candidateRepo.GetByID(ctx, rec.CandidateID)
			if candidate != nil {
				user, _ := s.userRepo.GetByID(ctx, candidate.UserID)
				job, _ := s.jobRepo.GetByID(ctx, rec.JobID)
				jobTitle := ""
				if job != nil {
					jobTitle = job.Title
				}
				if user != nil && user.Email != "" {
					s.notifier.SendInterviewScheduled(ctx, user.Email, jobTitle)
				}
			}
		}
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

// ConfirmForCandidate records that the candidate confirmed (or declined) the interview (Phase 3). Candidate must own the ATS record.
func (s *Service) ConfirmForCandidate(ctx context.Context, candidateID, interviewID string) error {
	interview, err := s.repo.GetByID(ctx, interviewID)
	if err != nil || interview == nil {
		return &domainerrors.NotFoundError{Resource: "interview", ID: interviewID}
	}
	rec, err := s.atsRepo.GetByID(ctx, interview.ATSID)
	if err != nil || rec == nil {
		return &domainerrors.NotFoundError{Resource: "application", ID: interview.ATSID}
	}
	if rec.CandidateID != candidateID {
		return &domainerrors.ForbiddenError{Message: "interview does not belong to you"}
	}
	return s.repo.SetCandidateConfirmedAt(ctx, interviewID, time.Now())
}
