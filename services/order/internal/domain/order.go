package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID
	BuyerID   string
	Status    string // pending, paid, cancelled, completed
	Total     float64
	PaymentID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	ProductID string
	Name      string
	Price     float64
	Quantity  int
}
