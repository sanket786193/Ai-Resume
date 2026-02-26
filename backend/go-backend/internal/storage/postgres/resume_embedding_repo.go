package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// ResumeEmbeddingRepo persists resume embeddings for vector search (resume_embeddings table).
type ResumeEmbeddingRepo struct {
	db *sql.DB
}

// NewResumeEmbeddingRepo returns a new ResumeEmbeddingRepo.
func NewResumeEmbeddingRepo(db *sql.DB) *ResumeEmbeddingRepo {
	return &ResumeEmbeddingRepo{db: db}
}

// formatVector returns a string suitable for PostgreSQL vector type: "[a,b,c,...]"
func formatVector(embedding []float64) string {
	if len(embedding) == 0 {
		return "[]"
	}
	b := strings.Builder{}
	b.WriteString("[")
	for i, v := range embedding {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(fmt.Sprintf("%g", v))
	}
	b.WriteString("]")
	return b.String()
}

// Upsert inserts or updates embedding for (resume_id, job_id); candidate_id and model_version are stored for indexing.
func (r *ResumeEmbeddingRepo) Upsert(ctx context.Context, resumeID, jobID, candidateID string, embedding []float64, modelVersion string) error {
	if len(embedding) == 0 {
		return nil
	}
	id := uuid.New().String()
	vecStr := formatVector(embedding)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO resume_embeddings (id, resume_id, job_id, candidate_id, embedding, model_version, created_at)
		 VALUES ($1, $2, $3, $4, $5::vector, $6, CURRENT_TIMESTAMP)
		 ON CONFLICT (resume_id, job_id) DO UPDATE SET
		   candidate_id = EXCLUDED.candidate_id,
		   embedding = EXCLUDED.embedding,
		   model_version = EXCLUDED.model_version`,
		id, resumeID, jobID, candidateID, vecStr, modelVersion)
	return err
}
