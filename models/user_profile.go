package models

import (
	"github.com/google/uuid"
	"time"
)

type UserProfile struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	Username string `json:"username"`
	Slug     string `json:"slug"`
}

type UserProfileRegistrationForm struct {
	Username string `json:"username"`
	Slug     string `json:"slug"`
}

type UserProfileUpdateForm struct {
	Username string `json:"username"`
	Slug     string `json:"slug"`
}
