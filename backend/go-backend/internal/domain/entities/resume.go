package entities

import (
	"time"
)

// Resume stores metadata; file stored securely (path or object key).
type Resume struct {
	ID         string
	CandidateID string
	FileName   string
	StoragePath string // local path or object key
	FileSize   int64
	MimeType   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
