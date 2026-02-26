package interviews

import (
	"context"
	"sort"
	"time"

	"resume/internal/domain/entities"
	domainerrors "resume/internal/domain/errors"
	"resume/internal/storage/postgres"

	"github.com/google/uuid"
)

// InterviewForHR is an interview with candidate and job display info for HR list.
type InterviewForHR struct {
	ID             string    `json:"id"`
	ATSID          string    `json:"ats_id"`
	JobID          string    `json:"job_id"`
	ScheduledAt    time.Time `json:"scheduled_at"`
	DurationMin    int       `json:"duration_minutes"`
	Location       string    `json:"location"`
	Round          int       `json:"round"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes,omitempty"`
	CandidateName  string    `json:"candidate_name"`
	JobTitle       string    `json:"job_title"`
}

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

// ListForHR returns all scheduled interviews for jobs created by the HR user, with candidate name and job title.
func (s *Service) ListForHR(ctx context.Context, hrUserID string) ([]*InterviewForHR, error) {
	jobIDs, err := s.jobRepo.ListIDsByCreatedBy(ctx, hrUserID)
	if err != nil || len(jobIDs) == 0 {
		return nil, err
	}
	recs, err := s.atsRepo.ListByJobIDs(ctx, jobIDs, nil, 500, 0)
	if err != nil || len(recs) == 0 {
		return nil, err
	}
	var out []*InterviewForHR
	for _, rec := range recs {
		interviews, err := s.repo.ListByATSID(ctx, rec.ID)
		if err != nil || len(interviews) == 0 {
			continue
		}
		candidateName := ""
		if rec.CandidateID != "" {
			if c, _ := s.candidateRepo.GetByID(ctx, rec.CandidateID); c != nil {
				if u, _ := s.userRepo.GetByID(ctx, c.UserID); u != nil {
					candidateName = u.Name
				}
			}
		}
		jobTitle := ""
		if rec.JobID != "" {
			if j, _ := s.jobRepo.GetByID(ctx, rec.JobID); j != nil {
				jobTitle = j.Title
			}
		}
		for _, i := range interviews {
			out = append(out, &InterviewForHR{
				ID:            i.ID,
				ATSID:         i.ATSID,
				JobID:         rec.JobID,
				ScheduledAt:   i.ScheduledAt,
				DurationMin:   i.Duration,
				Location:      i.Location,
				Round:         i.Round,
				Status:        i.Status,
				Notes:         i.Notes,
				CandidateName: candidateName,
				JobTitle:      jobTitle,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ScheduledAt.After(out[j].ScheduledAt) })
	return out, nil
}
