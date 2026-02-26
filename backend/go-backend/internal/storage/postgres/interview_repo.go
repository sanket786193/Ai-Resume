package postgres

import (
	"context"
	"database/sql"
	"time"

	"resume/internal/domain/entities"
)

// InterviewRepo persists interviews.
type InterviewRepo struct {
	db *sql.DB
}

// NewInterviewRepo returns a new InterviewRepo.
func NewInterviewRepo(db *sql.DB) *InterviewRepo {
	return &InterviewRepo{db: db}
}

// Create inserts an interview.
func (r *InterviewRepo) Create(ctx context.Context, i *entities.Interview) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO interviews (id, ats_id, scheduled_at, duration_minutes, location, round, notes, status, candidate_confirmed_at, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		i.ID, i.ATSID, i.ScheduledAt, i.Duration, i.Location, i.Round, i.Notes, i.Status, i.CandidateConfirmedAt, i.CreatedAt, i.UpdatedAt)
	return err
}

// GetByID returns an interview by ID.
func (r *InterviewRepo) GetByID(ctx context.Context, id string) (*entities.Interview, error) {
	var i entities.Interview
	var confirmedAt *time.Time
	err := r.db.QueryRowContext(ctx,
		`SELECT id, ats_id, scheduled_at, duration_minutes, location, round, notes, status, candidate_confirmed_at, created_at, updated_at
		 FROM interviews WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&i.ID, &i.ATSID, &i.ScheduledAt, &i.Duration, &i.Location, &i.Round, &i.Notes, &i.Status, &confirmedAt, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, err
	}
	i.CandidateConfirmedAt = confirmedAt
	return &i, nil
}

// ListByATSID returns interviews for an ATS record.
func (r *InterviewRepo) ListByATSID(ctx context.Context, atsID string) ([]*entities.Interview, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, ats_id, scheduled_at, duration_minutes, location, round, notes, status, candidate_confirmed_at, created_at, updated_at
		 FROM interviews WHERE ats_id = $1 AND deleted_at IS NULL ORDER BY round, scheduled_at`,
		atsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*entities.Interview
	for rows.Next() {
		var i entities.Interview
		var confirmedAt *time.Time
		if err := rows.Scan(&i.ID, &i.ATSID, &i.ScheduledAt, &i.Duration, &i.Location, &i.Round, &i.Notes, &i.Status, &confirmedAt, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		i.CandidateConfirmedAt = confirmedAt
		list = append(list, &i)
	}
	return list, rows.Err()
}

// Update updates an interview.
func (r *InterviewRepo) Update(ctx context.Context, i *entities.Interview) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE interviews SET scheduled_at = $1, duration_minutes = $2, location = $3, round = $4, notes = $5, status = $6, candidate_confirmed_at = $7, updated_at = $8
		 WHERE id = $9 AND deleted_at IS NULL`,
		i.ScheduledAt, i.Duration, i.Location, i.Round, i.Notes, i.Status, i.CandidateConfirmedAt, i.UpdatedAt, i.ID)
	return err
}

// SetCandidateConfirmedAt sets when the candidate confirmed/declined (Phase 3).
func (r *InterviewRepo) SetCandidateConfirmedAt(ctx context.Context, id string, at time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE interviews SET candidate_confirmed_at = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND deleted_at IS NULL`,
		at, id)
	return err
}
