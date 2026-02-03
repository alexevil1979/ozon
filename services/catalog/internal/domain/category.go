package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	ParentID  *uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
