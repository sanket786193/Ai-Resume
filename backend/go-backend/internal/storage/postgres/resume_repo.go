package postgres

import (
	"context"
	"database/sql"

	"resume/internal/domain/entities"
)

// ResumeRepo persists resume metadata.
type ResumeRepo struct {
	db *sql.DB
}

// NewResumeRepo returns a new ResumeRepo.
func NewResumeRepo(db *sql.DB) *ResumeRepo {
	return &ResumeRepo{db: db}
}

// Create inserts a resume.
func (r *ResumeRepo) Create(ctx context.Context, res *entities.Resume) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO resumes (id, candidate_id, file_name, storage_path, file_size, mime_type, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		res.ID, res.CandidateID, res.FileName, res.StoragePath, res.FileSize, res.MimeType, res.CreatedAt, res.UpdatedAt)
	return err
}

// GetByID returns a resume by ID.
func (r *ResumeRepo) GetByID(ctx context.Context, id string) (*entities.Resume, error) {
	var res entities.Resume
	err := r.db.QueryRowContext(ctx,
		`SELECT id, candidate_id, file_name, storage_path, file_size, mime_type, created_at, updated_at
		 FROM resumes WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&res.ID, &res.CandidateID, &res.FileName, &res.StoragePath, &res.FileSize, &res.MimeType, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ListByCandidateID returns resumes for a candidate.
func (r *ResumeRepo) ListByCandidateID(ctx context.Context, candidateID string) ([]*entities.Resume, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, candidate_id, file_name, storage_path, file_size, mime_type, created_at, updated_at
		 FROM resumes WHERE candidate_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`,
		candidateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*entities.Resume
	for rows.Next() {
		var res entities.Resume
		if err := rows.Scan(&res.ID, &res.CandidateID, &res.FileName, &res.StoragePath, &res.FileSize, &res.MimeType, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &res)
	}
	return list, rows.Err()
}
