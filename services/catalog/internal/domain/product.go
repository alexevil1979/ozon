package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	CategoryID  uuid.UUID
	Name        string
	Slug        string
	Description string
	Price       float64
	Brand       string
	SellerID    string
	Status      string // draft, active
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductAttributeValue struct {
	ProductID   uuid.UUID
	AttributeID uuid.UUID
	Value       string
}
