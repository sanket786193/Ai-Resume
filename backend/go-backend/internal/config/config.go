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
	SupabaseStorage  SupabaseStorageConfig
}

// SupabaseStorageConfig for resume/file uploads via Supabase Storage.
type SupabaseStorageConfig struct {
	Enabled        bool
	URL            string // e.g. https://xxx.supabase.co
	ServiceRoleKey string
	Bucket         string // e.g. resumes
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
		SupabaseStorage: SupabaseStorageConfig{
			Enabled:        getEnv("SUPABASE_STORAGE_ENABLED", "false") == "true",
			URL:            getEnv("SUPABASE_URL", ""),
			ServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
			Bucket:         getEnv("SUPABASE_STORAGE_BUCKET", "resumes"),
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
