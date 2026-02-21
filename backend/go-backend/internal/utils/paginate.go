package utils

// Pagination defaults.
const (
	DefaultLimit  = 20
	DefaultOffset = 0
	MaxLimit      = 100
)

// LimitOffset returns limit and offset clamped to defaults and max.
func LimitOffset(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	if offset < 0 {
		offset = DefaultOffset
	}
	return limit, offset
}
