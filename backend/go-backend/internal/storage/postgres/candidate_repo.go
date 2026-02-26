package postgres

import (
	"context"
	"database/sql"

	"resume/internal/domain/entities"
)

// CandidateRepo persists candidates.
type CandidateRepo struct {
	db *sql.DB
}

// NewCandidateRepo returns a new CandidateRepo.
func NewCandidateRepo(db *sql.DB) *CandidateRepo {
	return &CandidateRepo{db: db}
}

// Create inserts a candidate.
func (r *CandidateRepo) Create(ctx context.Context, c *entities.Candidate) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO candidates (id, user_id, phone, linkedin, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		c.ID, c.UserID, c.Phone, c.LinkedIn, c.CreatedAt, c.UpdatedAt)
	return err
}

// GetByID returns a candidate by ID.
func (r *CandidateRepo) GetByID(ctx context.Context, id string) (*entities.Candidate, error) {
	var c entities.Candidate
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, phone, linkedin, created_at, updated_at
		 FROM candidates WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&c.ID, &c.UserID, &c.Phone, &c.LinkedIn, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetByUserID returns a candidate by user ID.
func (r *CandidateRepo) GetByUserID(ctx context.Context, userID string) (*entities.Candidate, error) {
	var c entities.Candidate
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, phone, linkedin, created_at, updated_at
		 FROM candidates WHERE user_id = $1 AND deleted_at IS NULL`, userID).
		Scan(&c.ID, &c.UserID, &c.Phone, &c.LinkedIn, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
