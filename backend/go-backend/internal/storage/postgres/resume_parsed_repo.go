package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// ResumeParsedRepo persists parsed resume data (raw_text, parsed_json, cleaned_text).
type ResumeParsedRepo struct {
	db *sql.DB
}

// NewResumeParsedRepo returns a new ResumeParsedRepo.
func NewResumeParsedRepo(db *sql.DB) *ResumeParsedRepo {
	return &ResumeParsedRepo{db: db}
}

// Upsert inserts or updates parsed data for a resume (unique on resume_id).
func (r *ResumeParsedRepo) Upsert(ctx context.Context, resumeID, rawText string, parsedJSON []byte, cleanedText string) error {
	id := uuid.New().String()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO resume_parsed_data (id, resume_id, raw_text, parsed_json, cleaned_text, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		 ON CONFLICT (resume_id) DO UPDATE SET
		   raw_text = EXCLUDED.raw_text,
		   parsed_json = EXCLUDED.parsed_json,
		   cleaned_text = EXCLUDED.cleaned_text,
		   updated_at = CURRENT_TIMESTAMP`,
		id, resumeID, rawText, parsedJSON, cleanedText)
	return err
}
