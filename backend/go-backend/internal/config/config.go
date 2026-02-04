package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	Port      string
	UploadDir string

	// Database configuration
	Database DatabaseConfig

	// Cloudinary configuration
	Cloudinary CloudinaryConfig

	// Auth configuration
	Auth AuthConfig
}

// DatabaseConfig holds database connection parameters
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Enabled  bool
}

// CloudinaryConfig holds Cloudinary credentials
type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
	Folder    string
	Enabled   bool
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret         string
	JWTExpiryHours    int
	RefreshExpiryDays int
	RedirectURL       string
	Google            OAuthProvider
	GitHub            OAuthProvider
	Microsoft         OAuthProvider
}

// OAuthProvider holds OAuth provider configuration
type OAuthProvider struct {
	ClientID     string
	ClientSecret string
	Enabled      bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		UploadDir: getEnv("UPLOAD_DIR", "uploads"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "resume_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Enabled:  getEnv("DB_ENABLED", "true") == "true",
		},
		Cloudinary: CloudinaryConfig{
			CloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
			APIKey:    getEnv("CLOUDINARY_API_KEY", ""),
			APISecret: getEnv("CLOUDINARY_API_SECRET", ""),
			Folder:    getEnv("CLOUDINARY_FOLDER", "resume-ocr"),
			Enabled:   getEnv("CLOUDINARY_ENABLED", "true") == "true",
		},
		Auth: AuthConfig{
			JWTSecret:         getEnv("JWT_SECRET", "change-me-in-production"),
			JWTExpiryHours:    getEnvInt("JWT_EXPIRY_HOURS", 24),
			RefreshExpiryDays: getEnvInt("REFRESH_EXPIRY_DAYS", 30),
			RedirectURL:       getEnv("AUTH_REDIRECT_URL", "http://localhost:3000/auth/callback"),
			Google: OAuthProvider{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				Enabled:      getEnv("GOOGLE_ENABLED", "false") == "true",
			},
			GitHub: OAuthProvider{
				ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
				ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
				Enabled:      getEnv("GITHUB_ENABLED", "false") == "true",
			},
			Microsoft: OAuthProvider{
				ClientID:     getEnv("MICROSOFT_CLIENT_ID", ""),
				ClientSecret: getEnv("MICROSOFT_CLIENT_SECRET", ""),
				Enabled:      getEnv("MICROSOFT_ENABLED", "false") == "true",
			},
		},
	}
}

// GetDSN returns the database connection string for goose
func (c *DatabaseConfig) GetDSN() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.DBName + "?sslmode=" + c.SSLMode
}

// GetPostgresDSN returns the database connection string for lib/pq
func (c *DatabaseConfig) GetPostgresDSN() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as int or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}
