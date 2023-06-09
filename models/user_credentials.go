package models

import (
	"github.com/google/uuid"
	"time"
)

// UserCredentials represents the credentials of a user. Credentials are unique combination of data used to
// authenticate a user.
type UserCredentials struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	// Email is the main, active email of the user, used for communications and login.
	Email string `json:"email"`
	// NewEmail is set during the process of changing the main email. It remains here until validation.
	NewEmail string `json:"newEmail"`

	// Validated indicates whether the main email (Email) is validated or not, for the current user.
	Validated bool `json:"validated"`
}

// UserCredentialsRegistrationForm represents the parsed data struct, containing useful post-registration data.
// It is computed from UserCredentialsLoginForm.
type UserCredentialsRegistrationForm struct {
	// Email is the parsed email for the new user.
	Email Email `json:"email"`
	// Password is the hashed password for the new user.
	Password Password `json:"password"`

	// EmailPublicValidationCode is the raw, un-hashed key, that should be used to validate the new user email.
	// It must be immediately sent to the user address, and not persisted to the local disk once the request is
	// closed.
	EmailPublicValidationCode string `json:"emailPublicValidationCode"`
}

// UserCredentialsLoginForm is the form used to both log in and register.
type UserCredentialsLoginForm struct {
	// Email of the user. In case of a login, it MUST exist in the database. Otherwise, it must not.
	Email string `json:"email"`
	// Password is the raw, un-hashed password of the current user. It MUST not be persisted.
	Password string `json:"password"`
}
