// Package jwt — JWT issue and validate. Access tokens only here; refresh in next phase.
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

// Claims standard + subject (user id).
type Claims struct {
	jwt.RegisteredClaims
	Subject string `json:"sub"`
}

// JWT helper with secret and expiry.
type JWT struct {
	secret []byte
	expiry time.Duration
}

// New builds JWT util. expiryStr e.g. "24h".
func New(secret string, expiryStr string) *JWT {
	expiry, _ := time.ParseDuration(expiryStr)
	if expiry == 0 {
		expiry = 24 * time.Hour
	}
	return &JWT{secret: []byte(secret), expiry: expiry}
}

// Issue creates a signed token for userID.
func (j *JWT) Issue(userID string) (string, error) {
	now := time.Now()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		Subject: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// Validate parses and validates token; returns subject (user id) or error.
func (j *JWT) Validate(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return "", ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", ErrInvalidToken
	}
	return claims.Subject, nil
}

// ExpirySeconds returns access token TTL in seconds (for API response).
func (j *JWT) ExpirySeconds() time.Duration {
	return j.expiry
}
