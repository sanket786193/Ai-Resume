package enums

// ResumeStatus is the processing state of a resume.
type ResumeStatus string

const (
	ResumePending   ResumeStatus = "PENDING"
	ResumeProcessed ResumeStatus = "PROCESSED"
)

func (s ResumeStatus) Valid() bool {
	switch s {
	case ResumePending, ResumeProcessed:
		return true
	default:
		return false
	}
}
