package enums

// ExperienceLevel represents required experience for a job.
type ExperienceLevel string

const (
	ExperienceFresher    ExperienceLevel = "FRESHER"
	ExperienceExperienced ExperienceLevel = "EXPERIENCED"
	ExperienceAny        ExperienceLevel = "ANY"
)

func (e ExperienceLevel) Valid() bool {
	switch e {
	case ExperienceFresher, ExperienceExperienced, ExperienceAny:
		return true
	default:
		return false
	}
}
