package entities

import (
	"time"

	"resume/internal/domain/enums"
)

// Job represents a job posting.
type Job struct {
	ID          string
	Title       string
	Description string
	Location    string
	Department  string
	Status      enums.JobStatus
	CreatedBy   string // User ID (HR)
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
