package postgres

import (
	"context"
	"database/sql"
	"time"
)

// SessionRepo stores refresh token sessions.
type SessionRepo struct {
	db *sql.DB
}

// NewSessionRepo returns a new SessionRepo.
func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

// Create inserts a session.
func (r *SessionRepo) Create(ctx context.Context, userID, refreshTokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO ats_sessions (user_id, refresh_token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, refreshTokenHash, expiresAt)
	return err
}

// GetUserIDByRefreshHash returns user ID if session exists and not expired.
func (r *SessionRepo) GetUserIDByRefreshHash(ctx context.Context, refreshTokenHash string) (string, error) {
	var userID string
	err := r.db.QueryRowContext(ctx,
		`SELECT user_id FROM ats_sessions WHERE refresh_token_hash = $1 AND expires_at > NOW()`,
		refreshTokenHash).Scan(&userID)
	return userID, err
}

// DeleteByRefreshHash removes a session (e.g. on refresh or logout).
func (r *SessionRepo) DeleteByRefreshHash(ctx context.Context, refreshTokenHash string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ats_sessions WHERE refresh_token_hash = $1`, refreshTokenHash)
	return err
}
