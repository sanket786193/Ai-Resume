package postgres

import (
	"context"
	"database/sql"

	"resume/internal/domain/entities"
	"resume/internal/domain/enums"

	"github.com/lib/pq"
)

// JobRepo persists jobs.
type JobRepo struct {
	db *sql.DB
}

// NewJobRepo returns a new JobRepo.
func NewJobRepo(db *sql.DB) *JobRepo {
	return &JobRepo{db: db}
}

// Create inserts a job.
func (r *JobRepo) Create(ctx context.Context, j *entities.Job) error {
	query := `INSERT INTO jobs (id, title, description, location, department, status, experience_level, qualification, skills, vacancy_limits, created_by, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	expLevel := string(j.ExperienceLevel)
	if expLevel == "" {
		expLevel = string(enums.ExperienceAny)
	}
	_, err := r.db.ExecContext(ctx, query,
		j.ID, j.Title, j.Description, j.Location, j.Department,
		string(j.Status), expLevel, j.Qualification, pq.Array(j.Skills), j.VacancyLimits,
		j.CreatedBy, j.CreatedAt, j.UpdatedAt)
	return err
}

// GetByID returns a job by ID (excluding soft-deleted).
func (r *JobRepo) GetByID(ctx context.Context, id string) (*entities.Job, error) {
	query := `SELECT id, title, description, location, department, status, experience_level, qualification, skills, vacancy_limits, created_by, created_at, updated_at
	          FROM jobs WHERE id = $1 AND deleted_at IS NULL`
	var j entities.Job
	var status, expLevel string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&j.ID, &j.Title, &j.Description, &j.Location, &j.Department,
		&status, &expLevel, &j.Qualification, pq.Array(&j.Skills), &j.VacancyLimits,
		&j.CreatedBy, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return nil, err
	}
	j.Status = enums.JobStatus(status)
	if expLevel != "" {
		j.ExperienceLevel = enums.ExperienceLevel(expLevel)
	} else {
		j.ExperienceLevel = enums.ExperienceAny
	}
	return &j, nil
}

// ListIDsByCreatedBy returns all job IDs created by the given user (for HR's applications list).
func (r *JobRepo) ListIDsByCreatedBy(ctx context.Context, createdBy string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id FROM jobs WHERE deleted_at IS NULL AND created_by = $1 ORDER BY created_at DESC`, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// List returns jobs with optional status and pagination.
func (r *JobRepo) List(ctx context.Context, status *enums.JobStatus, limit, offset int) ([]*entities.Job, error) {
	var query string
	var args []interface{}
	if status != nil {
		query = `SELECT id, title, description, location, department, status, experience_level, qualification, skills, vacancy_limits, created_by, created_at, updated_at
	          FROM jobs WHERE deleted_at IS NULL AND status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{string(*status), limit, offset}
	} else {
		query = `SELECT id, title, description, location, department, status, experience_level, qualification, skills, vacancy_limits, created_by, created_at, updated_at
	          FROM jobs WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		args = []interface{}{limit, offset}
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*entities.Job
	for rows.Next() {
		var j entities.Job
		var statusStr, expLevel string
		if err := rows.Scan(&j.ID, &j.Title, &j.Description, &j.Location, &j.Department,
			&statusStr, &expLevel, &j.Qualification, pq.Array(&j.Skills), &j.VacancyLimits,
			&j.CreatedBy, &j.CreatedAt, &j.UpdatedAt); err != nil {
			return nil, err
		}
		j.Status = enums.JobStatus(statusStr)
		if expLevel != "" {
			j.ExperienceLevel = enums.ExperienceLevel(expLevel)
		} else {
			j.ExperienceLevel = enums.ExperienceAny
		}
		list = append(list, &j)
	}
	return list, rows.Err()
}

// UpdateStatus updates job status.
func (r *JobRepo) UpdateStatus(ctx context.Context, id string, status enums.JobStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE jobs SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND deleted_at IS NULL`,
		string(status), id)
	return err
}

// Update updates job fields (title, description, location, department, status, requirements, vacancy_limits).
func (r *JobRepo) Update(ctx context.Context, j *entities.Job) error {
	expLevel := string(j.ExperienceLevel)
	if expLevel == "" {
		expLevel = string(enums.ExperienceAny)
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE jobs SET title = $1, description = $2, location = $3, department = $4, status = $5,
		 experience_level = $6, qualification = $7, skills = $8, vacancy_limits = $9, updated_at = $10
		 WHERE id = $11 AND deleted_at IS NULL`,
		j.Title, j.Description, j.Location, j.Department, string(j.Status),
		expLevel, j.Qualification, pq.Array(j.Skills), j.VacancyLimits, j.UpdatedAt, j.ID)
	return err
}

// SoftDelete sets deleted_at.
func (r *JobRepo) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE jobs SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, id)
	return err
}
