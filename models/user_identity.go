package models

import (
	"github.com/google/uuid"
	"time"
)

type UserIdentity struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Birthday  time.Time `json:"birthday"`
	Sex       Sex       `json:"sex"`
}

type UserIdentityRegistrationForm struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Birthday  time.Time `json:"birthday"`
	Sex       Sex       `json:"sex"`
}

type UserIdentityUpdateForm struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Birthday  time.Time `json:"birthday"`
	Sex       Sex       `json:"sex"`
}
