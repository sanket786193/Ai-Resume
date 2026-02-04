package main

import (
	"database/sql"
	"flag"
	"log"

	"resume/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	var (
		command = flag.String("command", "up", "Migration command: up, down, status, or create")
		dir     = flag.String("dir", "migrations", "Directory with migration files")
	)
	flag.Parse()

	// Load configuration
	cfg := config.LoadConfig()

	if !cfg.Database.Enabled {
		log.Fatal("Database is disabled. Set DB_ENABLED=true to enable migrations.")
	}

	// Get DSN for goose
	dsn := cfg.Database.GetDSN()

	// Set migration directory
	migrationDir := *dir
	if migrationDir == "" {
		migrationDir = "migrations"
	}

	// Execute migration command
	switch *command {
	case "up":
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		defer db.Close()

		if err := goose.Up(db, migrationDir); err != nil {
			log.Fatalf("Failed to migrate up: %v", err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		defer db.Close()

		if err := goose.Down(db, migrationDir); err != nil {
			log.Fatalf("Failed to migrate down: %v", err)
		}
		log.Println("Migrations rolled back successfully")
	case "status":
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		defer db.Close()

		if err := goose.Status(db, migrationDir); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
	case "create":
		name := flag.Arg(0)
		if name == "" {
			log.Fatal("Migration name is required. Usage: migrate -command=create <name>")
		}
		if err := goose.Create(nil, migrationDir, name, "sql"); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}
		log.Printf("Created migration: %s", name)
	default:
		log.Fatalf("Unknown command: %s. Use: up, down, status, or create", *command)
	}
}
