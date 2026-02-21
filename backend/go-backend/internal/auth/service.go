package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"resume/internal/domain/entities"
	domainerrors "resume/internal/domain/errors"
	"resume/internal/domain/enums"
	"resume/internal/storage/postgres"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service handles registration, login, and refresh; no HTTP concerns.
type Service struct {
	jwt    *JWTService
	users  *postgres.UserRepo
	sessions *postgres.SessionRepo
	candidates *postgres.CandidateRepo
	cfg    AuthConfig
}

// AuthConfig for JWT and refresh.
type AuthConfig interface {
	GetJWTSecret() string
	GetJWTExpiryHours() int
	GetRefreshExpiryDays() int
}

// AuthResponse returned after login/refresh.
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"` // seconds
	User         UserResponse `json:"user"`
}

// UserResponse for API.
type UserResponse struct {
	ID           string  `json:"id"`
	Email        string  `json:"email"`
	Name         string  `json:"name"`
	Role         string  `json:"role"`
	CandidateID  *string `json:"candidate_id,omitempty"` // set when role is CANDIDATE
}

// RegisterCandidateInput for candidate registration.
type RegisterCandidateInput struct {
	Email    string
	Password string
	Name     string
	Phone    string
	LinkedIn string
}

// RegisterHRInput for HR registration (or admin-created).
type RegisterHRInput struct {
	Email    string
	Password string
	Name     string
}

// NewService constructs the auth service.
func NewService(
	jwt *JWTService,
	users *postgres.UserRepo,
	sessions *postgres.SessionRepo,
	candidates *postgres.CandidateRepo,
	cfg AuthConfig,
) *Service {
	return &Service{
		jwt:        jwt,
		users:     users,
		sessions:  sessions,
		candidates: candidates,
		cfg:       cfg,
	}
}

// RegisterCandidate creates a CANDIDATE user and candidate profile.
func (s *Service) RegisterCandidate(ctx context.Context, in RegisterCandidateInput) (*AuthResponse, error) {
	if in.Email == "" || in.Password == "" || in.Name == "" {
		return nil, &domainerrors.ValidationError{Field: "email/password/name", Message: "required"}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	existing, _ := s.users.GetByEmail(ctx, in.Email)
	if existing != nil {
		return nil, &domainerrors.ConflictError{Message: "email already registered"}
	}
	userID := uuid.New().String()
	now := time.Now()
	user := &entities.User{
		ID:           userID,
		Email:        in.Email,
		PasswordHash: string(hash),
		Name:         in.Name,
		Role:         enums.RoleCandidate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	candidateID := uuid.New().String()
	candidate := &entities.Candidate{
		ID:        candidateID,
		UserID:    userID,
		Phone:     in.Phone,
		LinkedIn:  in.LinkedIn,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.candidates.Create(ctx, candidate); err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, user)
}

// RegisterHR creates an HR user.
func (s *Service) RegisterHR(ctx context.Context, in RegisterHRInput) (*AuthResponse, error) {
	if in.Email == "" || in.Password == "" || in.Name == "" {
		return nil, &domainerrors.ValidationError{Field: "email/password/name", Message: "required"}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	existing, _ := s.users.GetByEmail(ctx, in.Email)
	if existing != nil {
		return nil, &domainerrors.ConflictError{Message: "email already registered"}
	}
	userID := uuid.New().String()
	now := time.Now()
	user := &entities.User{
		ID:           userID,
		Email:        in.Email,
		PasswordHash: string(hash),
		Name:         in.Name,
		Role:         enums.RoleHR,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, user)
}

// Login validates credentials and returns tokens.
func (s *Service) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, &domainerrors.UnauthorizedError{Message: "invalid credentials"}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, &domainerrors.UnauthorizedError{Message: "invalid credentials"}
	}
	return s.issueTokens(ctx, user)
}

// Refresh exchanges a valid refresh token for new access and refresh tokens.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	hash := hashToken(refreshToken)
	userID, err := s.sessions.GetUserIDByRefreshHash(ctx, hash)
	if err != nil || userID == "" {
		return nil, &domainerrors.UnauthorizedError{Message: "invalid or expired refresh token"}
	}
	user, err := s.users.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, &domainerrors.UnauthorizedError{Message: "user not found"}
	}
	_ = s.sessions.DeleteByRefreshHash(ctx, hash)
	return s.issueTokens(ctx, user)
}

// ValidateToken returns claims if the access token is valid (for middleware).
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	return s.jwt.ValidateAccessToken(tokenString)
}

func (s *Service) issueTokens(ctx context.Context, user *entities.User) (*AuthResponse, error) {
	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(time.Duration(s.cfg.GetJWTExpiryHours()) * time.Hour)
	if err := s.sessions.Create(ctx, user.ID, hashToken(refreshToken), expiresAt); err != nil {
		return nil, err
	}
	expiresIn := int(time.Until(expiresAt).Seconds())
	userResp := UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  string(user.Role),
	}
	if user.Role == enums.RoleCandidate {
		if c, err := s.candidates.GetByUserID(ctx, user.ID); err == nil && c != nil {
			userResp.CandidateID = &c.ID
		}
	}
	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User:         userResp,
	}, nil
}

// Logout invalidates the refresh token session if a refresh token is provided.
// Safe to call with empty string (no-op). Returns nil.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	hash := hashToken(refreshToken)
	_ = s.sessions.DeleteByRefreshHash(ctx, hash)
	return nil
}

// Me returns the current user for the given user ID (for GET /auth/me).
func (s *Service) Me(ctx context.Context, userID string) (*UserResponse, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, &domainerrors.NotFoundError{Resource: "user", ID: userID}
	}
	resp := UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  string(user.Role),
	}
	if user.Role == enums.RoleCandidate {
		if c, err := s.candidates.GetByUserID(ctx, user.ID); err == nil && c != nil {
			resp.CandidateID = &c.ID
		}
	}
	return &resp, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
