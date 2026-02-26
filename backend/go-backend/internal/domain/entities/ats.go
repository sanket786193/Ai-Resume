package entities

import (
	"time"

	"resume/internal/domain/enums"
)

// ATSRecord is the single application record per candidate per job.
type ATSRecord struct {
	ID                string
	JobID             string
	CandidateID       string
	ResumeID          string
	Status            enums.ATSStatus
	SkillMatchScore   *float64
	RankingScore      *float64
	Qualified         *bool   // AI qualification decision
	AIProcessedAt     *time.Time
	ATSScore            *int      // 0-100 from LLM
	SkillMatchPct       *int      // skill match percentage
	MissingSkills       []string
	ExperienceMatch     *string
	ExperienceWarnings  []string  // specific experience mismatch warnings
	KeywordMatches      []string  // JD terms found verbatim in resume
	SemanticMatches     []string  // JD requirements satisfied by meaning
	AISummary           *string
	ModelVersion      *string // e.g. llama3:8b
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}
