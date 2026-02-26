package config

import (
	"fmt"
	"os"
)

// Config holds application configuration.
type Config struct {
	Port      string
	UploadDir string

	Database         DatabaseConfig
	Auth             AuthConfig
	AI               AIConfig
	CloudinaryStorage CloudinaryStorageConfig
	SMTP             SMTPConfig
}

// SMTPConfig for email notifications (e.g. Gmail).
type SMTPConfig struct {
	Host     string
	Port     string
	Email    string
	Password string
	Enabled  bool
}

// CloudinaryStorageConfig for resume/file uploads via Cloudinary.
type CloudinaryStorageConfig struct {
	Enabled    bool
	CloudName  string
	APIKey     string
	APISecret  string
	Folder     string // e.g. resumes
}

// DatabaseConfig holds database connection parameters.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Enabled  bool
}

// AuthConfig holds JWT and refresh token configuration.
type AuthConfig struct {
	JWTSecret         string
	JWTExpiryHours    int
	RefreshExpiryDays int
}

// AIConfig for Python AI service (Ollama/ADK); supports HTTP or gRPC.
type AIConfig struct {
	BaseURL    string
	Enabled    bool
	TimeoutSec int
	UseGRPC    bool   // if true, use gRPC instead of HTTP
	GRPCTarget string // e.g. localhost:50051
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		UploadDir: getEnv("UPLOAD_DIR", "uploads"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "resume_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Enabled:  getEnv("DB_ENABLED", "true") == "true",
		},
		Auth: AuthConfig{
			JWTSecret:         getEnv("JWT_SECRET", "change-me-in-production"),
			JWTExpiryHours:    getEnvInt("JWT_EXPIRY_HOURS", 24),
			RefreshExpiryDays: getEnvInt("REFRESH_EXPIRY_DAYS", 30),
		},
		AI: AIConfig{
			BaseURL:    getEnv("AI_SERVICE_URL", "http://localhost:8000"),
			Enabled:    getEnv("AI_ENABLED", "true") == "true",
			TimeoutSec: getEnvInt("AI_TIMEOUT_SEC", 60),
			UseGRPC:    getEnv("AI_USE_GRPC", "false") == "true",
			GRPCTarget: getEnv("AI_GRPC_TARGET", "localhost:50051"),
		},
		CloudinaryStorage: CloudinaryStorageConfig{
			Enabled:   getEnv("CLOUDINARY_ENABLED", "false") == "true",
			CloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
			APIKey:    getEnv("CLOUDINARY_API_KEY", ""),
			APISecret: getEnv("CLOUDINARY_API_SECRET", ""),
			Folder:    getEnv("CLOUDINARY_FOLDER", "resumes"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnv("SMTP_PORT", "587"),
			Email:    getEnv("SMTP_EMAIL", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			Enabled:  getEnv("AUTO_EMAIL_ENABLED", "false") == "true",
		},
	}
}

// GetDSN returns the database connection string for Goose.
func (c *DatabaseConfig) GetDSN() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.DBName + "?sslmode=" + c.SSLMode
}

// GetPostgresDSN returns the database connection string for lib/pq.
func (c *DatabaseConfig) GetPostgresDSN() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}
