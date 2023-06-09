package profile_storage

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Model struct {
	bun.BaseModel `bun:"table:profiles"`

	ID        uuid.UUID  `json:"id" bun:"id,pk,type:uuid"`
	CreatedAt time.Time  `json:"created_at,omitempty" bun:"created_at,notnull"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bun:"updated_at"`

	Core
}

// Core contains the explicitly editable data of the current model.
type Core struct {
	// Username is a fake name displayed for a user.
	Username string `json:"username" bun:"username"`
	// Slug is the unique url suffix used to access the current profile.
	Slug string `json:"slug" bun:"slug"`
}
