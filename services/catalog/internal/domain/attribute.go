package domain

import (
	"time"

	"github.com/google/uuid"
)

type Attribute struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	ValueType string // text, number, boolean
	CreatedAt time.Time
}
