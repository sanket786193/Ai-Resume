package entities

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"resume/internal/domain/enums"
)

// RoleVacancy defines a vacancy limit per role (e.g. 4 Golang, 3 Python).
type RoleVacancy struct {
	Role  string `json:"role"`
	Limit int    `json:"limit"`
}

// VacancyLimits is a slice of RoleVacancy stored as JSONB.
type VacancyLimits []RoleVacancy

// Scan implements sql.Scanner for JSONB.
func (v *VacancyLimits) Scan(src interface{}) error {
	if src == nil {
		*v = nil
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return nil
	}
	if len(b) == 0 {
		*v = nil
		return nil
	}
	return json.Unmarshal(b, v)
}

// Value implements driver.Valuer for JSONB.
func (v VacancyLimits) Value() (driver.Value, error) {
	if v == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(v)
}

// Job represents a job posting.
type Job struct {
	ID              string
	Title           string
	Description     string
	Location        string
	Department      string
	Status          enums.JobStatus
	ExperienceLevel enums.ExperienceLevel
	Qualification   string
	Skills          []string
	VacancyLimits   VacancyLimits
	CreatedBy       string // User ID (HR)
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
