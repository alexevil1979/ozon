// Package domain — Auth bounded context: User (identity only).
package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is the identity entity (email + password hash). Profile lives in User service.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
