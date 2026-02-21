package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"resume/internal/domain/entities"
	"resume/internal/domain/enums"
)

// ATSRepo persists ATS records.
type ATSRepo struct {
	db *sql.DB
}

// NewATSRepo returns a new ATSRepo.
func NewATSRepo(db *sql.DB) *ATSRepo {
	return &ATSRepo{db: db}
}

// Create inserts an ATS record.
func (r *ATSRepo) Create(ctx context.Context, a *entities.ATSRecord) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO ats_records (id, job_id, candidate_id, resume_id, status, skill_match_score, ranking_score, qualified, ai_processed_at, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		a.ID, a.JobID, a.CandidateID, a.ResumeID, string(a.Status),
		a.SkillMatchScore, a.RankingScore, a.Qualified, a.AIProcessedAt, a.CreatedAt, a.UpdatedAt)
	return err
}

// atsSelectCols is the full column list for ats_records (including AI feedback).
const atsSelectCols = `id, job_id, candidate_id, resume_id, status, skill_match_score, ranking_score, qualified, ai_processed_at,
	ats_score, skill_match_pct, missing_skills, experience_match, experience_warnings, keyword_matches, semantic_matches, ai_summary, model_version, created_at, updated_at`

// GetByID returns an ATS record by ID.
func (r *ATSRepo) GetByID(ctx context.Context, id string) (*entities.ATSRecord, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+atsSelectCols+` FROM ats_records WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanATSRecord(row)
}

// GetByJobAndCandidate returns the ATS record for a job and candidate.
func (r *ATSRepo) GetByJobAndCandidate(ctx context.Context, jobID, candidateID string) (*entities.ATSRecord, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+atsSelectCols+` FROM ats_records WHERE job_id = $1 AND candidate_id = $2 AND deleted_at IS NULL`, jobID, candidateID)
	return scanATSRecord(row)
}

// UpdateStatus updates ATS status.
func (r *ATSRepo) UpdateStatus(ctx context.Context, id string, status enums.ATSStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE ats_records SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND deleted_at IS NULL`,
		string(status), id)
	return err
}

// UpdateStatusTx updates ATS status within a transaction.
func (r *ATSRepo) UpdateStatusTx(tx *sql.Tx, ctx context.Context, id string, status enums.ATSStatus) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE ats_records SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND deleted_at IS NULL`,
		string(status), id)
	return err
}

// UpdateAIScores updates AI-derived fields (advisory; backend validates).
func (r *ATSRepo) UpdateAIScores(ctx context.Context, id string, skillMatch, ranking *float64, qualified *bool, processedAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE ats_records SET skill_match_score = $1, ranking_score = $2, qualified = $3, ai_processed_at = $4, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $5 AND deleted_at IS NULL`,
		skillMatch, ranking, qualified, processedAt, id)
	return err
}

// UpdateAIFeedback updates full AI evaluation (ATS score, missing skills, experience warnings, keyword/semantic matches, summary, model version).
func (r *ATSRepo) UpdateAIFeedback(ctx context.Context, id string, atsScore, skillMatchPct *int, missingSkills []string, experienceMatch, aiSummary, modelVersion *string, experienceWarnings, keywordMatches, semanticMatches []string, processedAt time.Time) error {
	var missingJSON []byte
	if len(missingSkills) > 0 {
		missingJSON, _ = json.Marshal(missingSkills)
	}
	var expWarnJSON, keywordJSON, semanticJSON []byte
	if len(experienceWarnings) > 0 {
		expWarnJSON, _ = json.Marshal(experienceWarnings)
	}
	if len(keywordMatches) > 0 {
		keywordJSON, _ = json.Marshal(keywordMatches)
	}
	if len(semanticMatches) > 0 {
		semanticJSON, _ = json.Marshal(semanticMatches)
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE ats_records SET ats_score = $1, skill_match_pct = $2, missing_skills = $3, experience_match = $4, experience_warnings = $5, keyword_matches = $6, semantic_matches = $7, ai_summary = $8, model_version = $9, ai_processed_at = $10, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $11 AND deleted_at IS NULL`,
		atsScore, skillMatchPct, missingJSON, experienceMatch, expWarnJSON, keywordJSON, semanticJSON, aiSummary, modelVersion, processedAt, id)
	return err
}

// ListByJobID returns ATS records for a job with optional status filter.
func (r *ATSRepo) ListByJobID(ctx context.Context, jobID string, status *enums.ATSStatus, limit, offset int) ([]*entities.ATSRecord, error) {
	var query string
	var args []interface{}
	if status != nil {
		query = `SELECT ` + atsSelectCols + ` FROM ats_records WHERE job_id = $1 AND status = $2 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		args = []interface{}{jobID, string(*status), limit, offset}
	} else {
		query = `SELECT ` + atsSelectCols + ` FROM ats_records WHERE job_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{jobID, limit, offset}
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanATSRecords(rows)
}

// ListByJobIDs returns ATS records for any of the given job IDs (for HR pipeline).
func (r *ATSRepo) ListByJobIDs(ctx context.Context, jobIDs []string, status *enums.ATSStatus, limit, offset int) ([]*entities.ATSRecord, error) {
	if len(jobIDs) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(jobIDs))
	args := make([]interface{}, 0, len(jobIDs)+3)
	for i := range jobIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args = append(args, jobIDs[i])
	}
	query := `SELECT ` + atsSelectCols + ` FROM ats_records WHERE job_id IN (` + strings.Join(placeholders, ",") + `) AND deleted_at IS NULL`
	pos := len(args) + 1
	if status != nil {
		args = append(args, string(*status))
		query += ` AND status = $` + strconv.Itoa(pos)
		pos++
	}
	args = append(args, limit, offset)
	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(pos) + ` OFFSET $` + strconv.Itoa(pos+1)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanATSRecords(rows)
}

// ListByCandidateID returns ATS records for a candidate.
func (r *ATSRepo) ListByCandidateID(ctx context.Context, candidateID string, limit, offset int) ([]*entities.ATSRecord, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+atsSelectCols+` FROM ats_records WHERE candidate_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		candidateID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanATSRecords(rows)
}

func scanATSRecord(row *sql.Row) (*entities.ATSRecord, error) {
	var a entities.ATSRecord
	var status string
	var atsScore, skillMatchPct sql.NullInt64
	var experienceMatch, aiSummary, modelVersion sql.NullString
	var missingSkillsJSON, expWarningsJSON, keywordMatchesJSON, semanticMatchesJSON []byte
	err := row.Scan(&a.ID, &a.JobID, &a.CandidateID, &a.ResumeID, &status,
		&a.SkillMatchScore, &a.RankingScore, &a.Qualified, &a.AIProcessedAt,
		&atsScore, &skillMatchPct, &missingSkillsJSON, &experienceMatch, &expWarningsJSON, &keywordMatchesJSON, &semanticMatchesJSON, &aiSummary, &modelVersion,
		&a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	a.Status = enums.ATSStatus(status)
	if atsScore.Valid {
		v := int(atsScore.Int64)
		a.ATSScore = &v
	}
	if skillMatchPct.Valid {
		v := int(skillMatchPct.Int64)
		a.SkillMatchPct = &v
	}
	if len(missingSkillsJSON) > 0 {
		_ = json.Unmarshal(missingSkillsJSON, &a.MissingSkills)
	}
	if experienceMatch.Valid {
		a.ExperienceMatch = &experienceMatch.String
	}
	if len(expWarningsJSON) > 0 {
		_ = json.Unmarshal(expWarningsJSON, &a.ExperienceWarnings)
	}
	if len(keywordMatchesJSON) > 0 {
		_ = json.Unmarshal(keywordMatchesJSON, &a.KeywordMatches)
	}
	if len(semanticMatchesJSON) > 0 {
		_ = json.Unmarshal(semanticMatchesJSON, &a.SemanticMatches)
	}
	if aiSummary.Valid {
		a.AISummary = &aiSummary.String
	}
	if modelVersion.Valid {
		a.ModelVersion = &modelVersion.String
	}
	return &a, nil
}

func scanATSRecords(rows *sql.Rows) ([]*entities.ATSRecord, error) {
	var list []*entities.ATSRecord
	for rows.Next() {
		var a entities.ATSRecord
		var status string
		var atsScore, skillMatchPct sql.NullInt64
		var experienceMatch, aiSummary, modelVersion sql.NullString
		var missingSkillsJSON, expWarningsJSON, keywordMatchesJSON, semanticMatchesJSON []byte
		if err := rows.Scan(&a.ID, &a.JobID, &a.CandidateID, &a.ResumeID, &status,
			&a.SkillMatchScore, &a.RankingScore, &a.Qualified, &a.AIProcessedAt,
			&atsScore, &skillMatchPct, &missingSkillsJSON, &experienceMatch, &expWarningsJSON, &keywordMatchesJSON, &semanticMatchesJSON, &aiSummary, &modelVersion,
			&a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		a.Status = enums.ATSStatus(status)
		if atsScore.Valid {
			v := int(atsScore.Int64)
			a.ATSScore = &v
		}
		if skillMatchPct.Valid {
			v := int(skillMatchPct.Int64)
			a.SkillMatchPct = &v
		}
		if len(missingSkillsJSON) > 0 {
			_ = json.Unmarshal(missingSkillsJSON, &a.MissingSkills)
		}
		if experienceMatch.Valid {
			a.ExperienceMatch = &experienceMatch.String
		}
		if len(expWarningsJSON) > 0 {
			_ = json.Unmarshal(expWarningsJSON, &a.ExperienceWarnings)
		}
		if len(keywordMatchesJSON) > 0 {
			_ = json.Unmarshal(keywordMatchesJSON, &a.KeywordMatches)
		}
		if len(semanticMatchesJSON) > 0 {
			_ = json.Unmarshal(semanticMatchesJSON, &a.SemanticMatches)
		}
		if aiSummary.Valid {
			a.AISummary = &aiSummary.String
		}
		if modelVersion.Valid {
			a.ModelVersion = &modelVersion.String
		}
		list = append(list, &a)
	}
	return list, rows.Err()
}
