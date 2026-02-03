// Package auth — application service: Register, Login, Refresh.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"marketplace/auth/internal/domain"
	"marketplace/auth/internal/jwt"
	"marketplace/auth/internal/repository"
)

// TokenPair returned on login/register/refresh.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
}

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidRefresh     = errors.New("invalid or expired refresh token")
)

// Service implements auth use cases.
type Service struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	jwtUtil          *jwt.JWT
	refreshExpiry    time.Duration
}

// NewService creates auth service. refreshExpiry e.g. 7*24*time.Hour.
func NewService(
	userRepo *repository.UserRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	jwtUtil *jwt.JWT,
	refreshExpiry time.Duration,
) *Service {
	if refreshExpiry == 0 {
		refreshExpiry = 7 * 24 * time.Hour
	}
	return &Service{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtUtil:          jwtUtil,
		refreshExpiry:    refreshExpiry,
	}
}

// Register creates user and returns tokens.
func (s *Service) Register(ctx context.Context, email string, password string) (*domain.User, *TokenPair, error) {
	email = trimLower(email)
	if !ValidateEmail(email) {
		return nil, nil, ErrInvalidCredentials
	}
	if err := ValidatePassword(password); err != nil {
		return nil, nil, err
	}
	hash, err := HashPassword(password)
	if err != nil {
		return nil, nil, err
	}
	u := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
	}
	if err := s.userRepo.Create(ctx, u); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, nil, ErrEmailTaken
		}
		return nil, nil, err
	}
	tp, err := s.issueTokenPair(ctx, u.ID)
	if err != nil {
		return u, nil, err
	}
	return u, tp, nil
}

// Login validates credentials and returns tokens.
func (s *Service) Login(ctx context.Context, email string, password string) (*TokenPair, error) {
	email = trimLower(email)
	if !ValidateEmail(email) {
		return nil, ErrInvalidCredentials
	}
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || u == nil {
		return nil, ErrInvalidCredentials
	}
	if !CheckPassword(u.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}
	return s.issueTokenPair(ctx, u.ID)
}

// Refresh validates refresh token and returns new token pair (rotation).
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	if refreshToken == "" {
		return nil, ErrInvalidRefresh
	}
	hash := hashRefreshToken(refreshToken)
	row, err := s.refreshTokenRepo.GetByTokenHash(ctx, hash)
	if err != nil || row == nil {
		return nil, ErrInvalidRefresh
	}
	// Rotate: delete old refresh token
	_ = s.refreshTokenRepo.DeleteByID(ctx, row.ID)
	return s.issueTokenPair(ctx, row.UserID)
}

func (s *Service) issueTokenPair(ctx context.Context, userID uuid.UUID) (*TokenPair, error) {
	access, err := s.jwtUtil.Issue(userID.String())
	if err != nil {
		return nil, err
	}
	refresh, hash, expiresAt, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}
	if err := s.refreshTokenRepo.Create(ctx, userID, hash, expiresAt); err != nil {
		return nil, err
	}
	expSec := int64(s.jwtUtil.ExpirySeconds() / time.Second)
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expSec,
	}, nil
}

func (s *Service) generateRefreshToken() (plain, hash string, expiresAt time.Time, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", time.Time{}, err
	}
	plain = hex.EncodeToString(b)
	hash = hashRefreshToken(plain)
	expiresAt = time.Now().UTC().Add(s.refreshExpiry)
	return plain, hash, expiresAt, nil
}

func hashRefreshToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func trimLower(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
