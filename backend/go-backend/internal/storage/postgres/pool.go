package postgres

import (
	"database/sql"
	"fmt"

	"resume/internal/config"

	_ "github.com/lib/pq"
)

// DB wraps *sql.DB for dependency injection.
type DB struct {
	*sql.DB
}

// New creates a new database connection with pooling.
func New(cfg *config.DatabaseConfig) (*DB, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("database is disabled")
	}
	db, err := sql.Open("postgres", cfg.GetPostgresDSN())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	return &DB{DB: db}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.DB.Close()
}
