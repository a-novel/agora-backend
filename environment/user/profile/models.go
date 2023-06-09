package profile

import (
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"time"
)

type Model struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	CreatedAt time.Time  `json:"createdAt"`
	Sex       models.Sex `json:"sex"`
}

type Preview struct {
	ID        uuid.UUID `json:"id"`
	Slug      string    `json:"slug"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
}
