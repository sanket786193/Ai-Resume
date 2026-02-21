package enums

// Role for RBAC.
type Role string

const (
	RoleHR       Role = "HR"
	RoleCandidate Role = "CANDIDATE"
)

func (r Role) Valid() bool {
	switch r {
	case RoleHR, RoleCandidate:
		return true
	default:
		return false
	}
}
