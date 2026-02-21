package entities

import (
	"time"
)

// Interview is tied to an ATS record; supports multiple rounds.
type Interview struct {
	ID         string
	ATSID      string
	ScheduledAt time.Time
	Duration   int     // minutes
	Location   string  // or meeting link
	Round      int
	Notes      string
	Status     string  // e.g. SCHEDULED, COMPLETED, CANCELLED
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
