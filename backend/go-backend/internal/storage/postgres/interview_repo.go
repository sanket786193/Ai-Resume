package postgres

import (
	"context"
	"database/sql"

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
		`INSERT INTO interviews (id, ats_id, scheduled_at, duration_minutes, location, round, notes, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		i.ID, i.ATSID, i.ScheduledAt, i.Duration, i.Location, i.Round, i.Notes, i.Status, i.CreatedAt, i.UpdatedAt)
	return err
}

// GetByID returns an interview by ID.
func (r *InterviewRepo) GetByID(ctx context.Context, id string) (*entities.Interview, error) {
	var i entities.Interview
	err := r.db.QueryRowContext(ctx,
		`SELECT id, ats_id, scheduled_at, duration_minutes, location, round, notes, status, created_at, updated_at
		 FROM interviews WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&i.ID, &i.ATSID, &i.ScheduledAt, &i.Duration, &i.Location, &i.Round, &i.Notes, &i.Status, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

// ListByATSID returns interviews for an ATS record.
func (r *InterviewRepo) ListByATSID(ctx context.Context, atsID string) ([]*entities.Interview, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, ats_id, scheduled_at, duration_minutes, location, round, notes, status, created_at, updated_at
		 FROM interviews WHERE ats_id = $1 AND deleted_at IS NULL ORDER BY round, scheduled_at`,
		atsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*entities.Interview
	for rows.Next() {
		var i entities.Interview
		if err := rows.Scan(&i.ID, &i.ATSID, &i.ScheduledAt, &i.Duration, &i.Location, &i.Round, &i.Notes, &i.Status, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &i)
	}
	return list, rows.Err()
}

// Update updates an interview.
func (r *InterviewRepo) Update(ctx context.Context, i *entities.Interview) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE interviews SET scheduled_at = $1, duration_minutes = $2, location = $3, round = $4, notes = $5, status = $6, updated_at = $7
		 WHERE id = $8 AND deleted_at IS NULL`,
		i.ScheduledAt, i.Duration, i.Location, i.Round, i.Notes, i.Status, i.UpdatedAt, i.ID)
	return err
}
