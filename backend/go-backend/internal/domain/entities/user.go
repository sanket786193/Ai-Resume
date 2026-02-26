package entities

import (
	"time"

	"resume/internal/domain/enums"
)

// User is the core identity (HR or Candidate).
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	Role         enums.Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
