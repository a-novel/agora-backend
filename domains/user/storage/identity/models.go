package identity_storage

import (
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Model struct {
	bun.BaseModel `bun:"table:identities"`

	ID        uuid.UUID  `json:"id" bun:"id,pk,type:uuid"`
	CreatedAt time.Time  `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bun:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bun:"deleted_at"`

	Core
}

// Core contains the explicitly editable data of the current model.
type Core struct {
	FirstName string     `json:"first_name" bun:"first_name"`
	LastName  string     `json:"last_name" bun:"last_name"`
	Birthday  time.Time  `json:"birthday" bun:"birthday"`
	Sex       models.Sex `json:"sex" bun:"sex,notnull"`
}
