package postgres

import (
	"context"
	"database/sql"

	"resume/internal/domain/entities"
	"resume/internal/domain/enums"
)

// UserRepo persists users.
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo returns a new UserRepo.
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create inserts a user.
func (r *UserRepo) Create(ctx context.Context, u *entities.User) error {
	query := `INSERT INTO ats_users (id, email, password_hash, name, role, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.Email, u.PasswordHash, u.Name, string(u.Role), u.CreatedAt, u.UpdatedAt)
	return err
}

// GetByID returns a user by ID (excluding soft-deleted).
func (r *UserRepo) GetByID(ctx context.Context, id string) (*entities.User, error) {
	query := `SELECT id, email, password_hash, name, role, created_at, updated_at
	          FROM ats_users WHERE id = $1 AND deleted_at IS NULL`
	var u entities.User
	var role string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &role,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	u.Role = enums.Role(role)
	return &u, nil
}

// GetByEmail returns a user by email (excluding soft-deleted).
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `SELECT id, email, password_hash, name, role, created_at, updated_at
	          FROM ats_users WHERE email = $1 AND deleted_at IS NULL`
	var u entities.User
	var role string
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &role,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	u.Role = enums.Role(role)
	return &u, nil
}
