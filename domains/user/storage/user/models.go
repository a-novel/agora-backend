package user_storage

import (
	"github.com/a-novel/agora-backend/domains/user/storage/credentials"
	"github.com/a-novel/agora-backend/domains/user/storage/identity"
	"github.com/a-novel/agora-backend/domains/user/storage/profile"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"time"
)

type Model struct {
	ID        uuid.UUID  `json:"id" bun:"id,pk,type:uuid"`
	CreatedAt time.Time  `json:"created_at,omitempty" bun:"created_at,notnull"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bun:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bun:"deleted_at"`

	Core
}

// Core contains the explicitly editable data of the current model.
type Core struct {
	Credentials credentials_storage.Core
	Identity    identity_storage.Core
	Profile     profile_storage.Core
}

// PublicPreview represents the preview information about a user, visible publicly.
type PublicPreview struct {
	ID        uuid.UUID `json:"id" bun:"id,pk,type:uuid"`
	Slug      string    `json:"slug" bun:"slug"`
	Username  string    `json:"username" bun:"username"`
	FirstName string    `json:"first_name" bun:"first_name"`
	LastName  string    `json:"last_name" bun:"last_name"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
}

type Public struct {
	ID        uuid.UUID  `json:"id" bun:"id,pk,type:uuid"`
	Username  string     `json:"username" bun:"username"`
	FirstName string     `json:"first_name" bun:"first_name"`
	LastName  string     `json:"last_name" bun:"last_name"`
	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	Sex       models.Sex `json:"sex" bun:"sex"`
}

// Preview represents the information about a user, visible privately in the application.
type Preview struct {
	ID        uuid.UUID    `json:"id" bun:"id,pk,type:uuid"`
	Username  string       `json:"username" bun:"username"`
	FirstName string       `json:"first_name" bun:"first_name"`
	LastName  string       `json:"last_name" bun:"last_name"`
	Email     models.Email `json:"email" bun:"email"`
	Slug      string       `json:"slug" bun:"slug"`
	Sex       models.Sex   `json:"sex" bun:"sex"`
}
