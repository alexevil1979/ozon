package domain

import "time"

type Profile struct {
	UserID      string
	DisplayName string
	AvatarURL   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
