package enums

// ATSStatus is the single source of truth for applicant tracking state.
type ATSStatus string

const (
	ATSApplied    ATSStatus = "APPLIED"
	ATSScreening  ATSStatus = "SCREENING"
	ATSShortlisted ATSStatus = "SHORTLISTED"
	ATSInterview  ATSStatus = "INTERVIEW"
	ATSRejected   ATSStatus = "REJECTED"
	ATSHired      ATSStatus = "HIRED"
)

// Valid returns true if s is a known status.
func (s ATSStatus) Valid() bool {
	switch s {
	case ATSApplied, ATSScreening, ATSShortlisted, ATSInterview, ATSRejected, ATSHired:
		return true
	default:
		return false
	}
}

// All returns all valid statuses in display order.
func AllATSStatuses() []ATSStatus {
	return []ATSStatus{ATSApplied, ATSScreening, ATSShortlisted, ATSInterview, ATSRejected, ATSHired}
}
