package entities

import (
	"time"
)

// Candidate profile (linked to User with Role=CANDIDATE).
type Candidate struct {
	ID        string
	UserID    string
	Phone     string
	LinkedIn  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
