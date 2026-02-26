package entities

import (
	"time"
)

// Offer is created from an ATS record; ATS → HIRED only after accept.
type Offer struct {
	ID          string
	ATSID       string
	Amount      string
	Currency    string
	StartsAt    *time.Time
	Status      string // PENDING, ACCEPTED, REJECTED
	RespondedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
