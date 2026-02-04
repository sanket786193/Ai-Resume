package database

import (
	"database/sql"
	"fmt"
	"log"

	"resume/internal/config"

	_ "github.com/lib/pq"
)

// DB holds the database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB(cfg *config.DatabaseConfig) (*sql.DB, error) {
	if !cfg.Enabled {
		log.Println("Database is disabled")
		return nil, nil
	}

	dsn := cfg.GetPostgresDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connected successfully to %s/%s", cfg.Host, cfg.DBName)
	DB = db

	return db, nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}
