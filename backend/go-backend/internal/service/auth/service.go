package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"resume/internal/config"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type Service struct {
	jwt    *JWTService
	oauth  *OAuthService
	db     *sql.DB
	config *config.AuthConfig
}

type User struct {
	ID         string
	Email      string
	Name       string
	Provider   string
	ProviderID string
	AvatarURL  string
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

func NewService(cfg *config.AuthConfig, db *sql.DB) *Service {
	return &Service{
		jwt:    NewJWTService(cfg.JWTSecret, cfg.JWTExpiryHours, cfg.RefreshExpiryDays),
		oauth:  NewOAuthService(cfg),
		db:     db,
		config: cfg,
	}
}

func (s *Service) GetOAuthURL(provider string) (string, error) {
	var config *oauth2.Config
	switch provider {
	case "google":
		if !s.config.Google.Enabled {
			return "", fmt.Errorf("Google OAuth is not enabled")
		}
		config = s.oauth.GetGoogleConfig()
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}

	state := uuid.New().String()
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (s *Service) HandleOAuthCallback(provider, code string) (*AuthResponse, error) {
	ctx := context.Background()
	var config *oauth2.Config
	var userInfo interface{}

	switch provider {
	case "google":
		if !s.config.Google.Enabled {
			return nil, fmt.Errorf("Google OAuth is not enabled")
		}
		config = s.oauth.GetGoogleConfig()
		token, err := config.Exchange(ctx, code)
		if err != nil {
			return nil, err
		}
		info, err := s.oauth.GetGoogleUserInfo(ctx, token)
		if err != nil {
			return nil, err
		}
		userInfo = info

	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Create or get user
	user, err := s.createOrGetUser(provider, userInfo)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Store session
	if err := s.storeSession(user.ID, accessToken, refreshToken); err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *Service) createOrGetUser(provider string, userInfo interface{}) (*User, error) {
	var email, name, providerID, avatarURL string

	switch v := userInfo.(type) {
	case *GoogleUserInfo:
		email = v.Email
		name = v.Name
		providerID = v.ID
		avatarURL = v.Picture

	}

	// Check if user exists
	var userID string
	err := s.db.QueryRow(
		"SELECT id FROM users WHERE provider = $1 AND provider_id = $2",
		provider, providerID,
	).Scan(&userID)

	if err == sql.ErrNoRows {
		// Create new user
		userID = uuid.New().String()
		_, err = s.db.Exec(
			`INSERT INTO users (id, email, name, provider, provider_id, avatar_url, last_login)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			userID, email, name, provider, providerID, avatarURL, time.Now(),
		)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// Update last login
		if _, err := s.db.Exec("UPDATE users SET last_login = $1 WHERE id = $2", time.Now(), userID); err != nil {
			return nil, err
		}
		// Update avatar if provided
		if avatarURL != "" {
			if _, err := s.db.Exec("UPDATE users SET avatar_url = $1 WHERE id = $2", avatarURL, userID); err != nil {
				return nil, err
			}
		}
	}

	return &User{
		ID:         userID,
		Email:      email,
		Name:       name,
		Provider:   provider,
		ProviderID: providerID,
		AvatarURL:  avatarURL,
	}, nil
}

func (s *Service) storeSession(userID, accessToken, refreshToken string) error {
	accessHash := hashToken(accessToken)
	refreshHash := hashToken(refreshToken)
	expiresAt := time.Now().Add(time.Duration(s.config.JWTExpiryHours) * time.Hour)

	_, err := s.db.Exec(
		`INSERT INTO user_sessions (user_id, token_hash, refresh_token_hash, expires_at)
		 VALUES ($1, $2, $3, $4)`,
		userID, accessHash, refreshHash, expiresAt,
	)
	return err
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	return s.jwt.ValidateToken(tokenString)
}

func (s *Service) RefreshToken(refreshToken string) (*AuthResponse, error) {
	// Look up user from refresh token hash in database
	refreshHash := hashToken(refreshToken)

	var userID string
	err := s.db.QueryRow(
		"SELECT user_id FROM user_sessions WHERE refresh_token_hash = $1 AND expires_at > NOW()",
		refreshHash,
	).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid or expired refresh token")
		}
		return nil, err
	}

	// Get user
	var user User
	err = s.db.QueryRow(
		"SELECT id, email, name, provider, provider_id, avatar_url FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Provider, &user.ProviderID, &user.AvatarURL)

	if err != nil {
		return nil, err
	}

	// Generate new tokens
	accessToken, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Update session
	if err := s.storeSession(user.ID, accessToken, newRefreshToken); err != nil {
		return nil, err
	}

	// Delete old session
	_, err = s.db.Exec("DELETE FROM user_sessions WHERE refresh_token_hash = $1", refreshHash)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}
