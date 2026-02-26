package enums

// JobStatus represents job lifecycle state.
type JobStatus string

const (
	JobDraft    JobStatus = "DRAFT"
	JobPublished JobStatus = "PUBLISHED"
	JobClosed   JobStatus = "CLOSED"
)

func (s JobStatus) Valid() bool {
	switch s {
	case JobDraft, JobPublished, JobClosed:
		return true
	default:
		return false
	}
}
