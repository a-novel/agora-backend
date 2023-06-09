package credentials_storage

import (
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Model struct {
	bun.BaseModel `bun:"table:credentials"`

	ID        uuid.UUID  `json:"id" bun:"id,pk,type:uuid"`
	CreatedAt time.Time  `json:"created_at,omitempty" bun:"created_at,notnull"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bun:"updated_at"`

	Core
}

// Core contains the explicitly editable data of the current model.
type Core struct {
	// Email is the main email address of a user, used to authenticate it and to communicate with.
	Email models.Email `json:"email" bun:"embed:email_"`
	// NewEmail is set when user wants to change its email address. Because email is the primary way to
	// authenticate a user, the Email value is not directly updated, but saved here in a pending state, with
	// Email.Validation set.
	// Once this email is validated, the Email field is updated, and this one is nullified.
	NewEmail models.Email `json:"new_email" bun:"embed:new_email_"`
	// Password used to authenticate the user.
	Password models.Password `json:"password" bun:"embed:password_"`
}
