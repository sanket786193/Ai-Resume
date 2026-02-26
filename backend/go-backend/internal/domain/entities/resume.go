package entities

import (
	"time"

	"resume/internal/domain/enums"
)

// Resume stores metadata; file stored securely (path or object key).
type Resume struct {
	ID           string
	CandidateID  string
	FileName     string
	StoragePath  string // local path or object key
	FileSize     int64
	MimeType     string
	Status       enums.ResumeStatus // PENDING until AI has processed; PROCESSED when parsed/feedback available
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
