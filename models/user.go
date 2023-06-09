package models

import (
	"github.com/google/uuid"
	"time"
)

// User is a registered user of the application.
type User struct {
	// ID of the user, used to link other objects to it.
	ID uuid.UUID `json:"id"`
	// CreatedAt stores the time at which the user was created.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt stores the time at which the user was last updated.
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	// DeletedAt stores the time at which the user was deleted. It is nil if the user is not deleted.
	DeletedAt *time.Time `json:"deletedAt,omitempty"`

	Credentials UserCredentials `json:"credentials"`
	Identity    UserIdentity    `json:"identity"`
	Profile     UserProfile     `json:"profile"`
}

// ToPayload generates a JWT payload for the current user.
func (model *User) ToPayload() UserTokenPayload {
	return UserTokenPayload{ID: model.ID}
}

// UserPostRegistration is returned after a user first registers. It contains information for further steps.
type UserPostRegistration struct {
	// EmailValidationCode is the code that should be sent to the user email address for validation.
	EmailValidationCode string `json:"email_validation_link"`
}

// UserCreateForm is the form sent by the user for registration.
type UserCreateForm struct {
	Credentials UserCredentialsLoginForm `json:"credentials"`
	Identity    UserIdentityUpdateForm   `json:"identity"`
	Profile     UserProfileUpdateForm    `json:"profile"`
}

// UserPublic is the user data publicly available to other users, once created.
type UserPublic struct {
	// ID of the user, used to link other objects to it.
	ID uuid.UUID `json:"id"`
	// Slug of the user, used to retrieve the user profile.
	Username string `json:"username"`
	// Username of the user. Hides FirstName and LastName if set.
	FirstName string `json:"firstName"`
	// LastName of the user. This value is hidden if Username is set.
	LastName string `json:"lastName"`
	// CreatedAt stores the time at which the user was created.
	CreatedAt time.Time `json:"createdAt"`
	// Sex of the user.
	Sex Sex `json:"sex"`
}

// UserPublicPreview is the user data publicly available to other users, once created.
// This is an alternative version of UserPublic used for display in lists.
type UserPublicPreview struct {
	// ID of the user, used to link other objects to it.
	ID uuid.UUID `json:"id"`
	// Slug of the user, used to retrieve the user profile.
	Slug string `json:"slug"`
	// Username of the user. Hides FirstName and LastName if set.
	Username string `json:"username"`
	// FirstName of the user. This value is hidden if Username is set.
	FirstName string `json:"firstName"`
	// LastName of the user. This value is hidden if Username is set.
	LastName string `json:"lastName"`
	// CreatedAt stores the time at which the user was created.
	CreatedAt time.Time `json:"createdAt"`
}

// UserPreview is the user data available to the user itself.
// It is used for generic previews across the application.
type UserPreview struct {
	// ID of the user, used to link other objects to it.
	ID uuid.UUID `json:"id"`
	// Username of the user. Hides FirstName and LastName if set.
	Username string `json:"username"`
	// FirstName of the user. This value is hidden if Username is set.
	FirstName string `json:"firstName"`
	// LastName of the user. This value is hidden if Username is set.
	LastName string `json:"lastName"`
	// Email is the main, active email of the user, used for communications and login.
	Email string `json:"email"`
	// Slug of the user, used to retrieve the user profile.
	Slug string `json:"slug"`
	// Sex of the user.
	Sex Sex `json:"sex"`
}

// UserInfo is the user data available to the user itself, on its private profile page.
type UserInfo struct {
	// ID of the user, used to link other objects to it.
	ID uuid.UUID `json:"id"`
	// CreatedAt stores the time at which the user was created.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt stores the time at which the user was last updated.
	UpdatedAt *time.Time `json:"updatedAt"`
	// Email is the main, active email of the user, used for communications and login.
	Email string `json:"email"`
	// NewEmail is set during the process of changing the main email. It remains here until validation.
	NewEmail string           `json:"newEmail"`
	Identity UserInfoIdentity `json:"identity"`
	Profile  UserInfoProfile  `json:"profile"`
}

// UserFlat is a flattened version of UserInfo, mainly used for frontend display.
// TODO: we really don't need this, should be removed to stick with UserInfo and other variations.
type UserFlat struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	Email     string `json:"email"`
	NewEmail  string `json:"newEmail"`
	Validated bool   `json:"validated"`

	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Birthday  time.Time `json:"birthday"`
	Sex       Sex       `json:"sex"`

	Username string `json:"username"`
	Slug     string `json:"slug"`
}

// UserInfoIdentity contains the identity information of a user.
type UserInfoIdentity struct {
	// FirstName of the user.
	FirstName string `json:"firstName"`
	// LastName of the user.
	LastName string `json:"lastName"`
	// Birthday of the user. Hour should be ignored.
	Birthday time.Time `json:"birthday"`
	// Sex of the user.
	Sex Sex `json:"sex"`
}

// UserInfoProfile contains the profile information of a user.
type UserInfoProfile struct {
	// Username of the user. Is used as an alternate, anonymous display name.
	Username string `json:"username"`
	// Slug of the user. Is used as a unique identifier in URLs.
	Slug string `json:"slug"`
}

// UserEmailValidationStatus contains the status of the email validation process.
type UserEmailValidationStatus struct {
	// Email is the main, active email of the user, used for communications and login.
	Email string `json:"email"`
	// NewEmail is set during the process of changing the main email. It remains here until validation.
	NewEmail string `json:"newEmail"`
	// Validated is true if the main Email is validated.
	Validated bool `json:"validated"`
}

// UserValidateEmailForm is the form sent by the frontend page linked in the validation email, to validate the email.
type UserValidateEmailForm struct {
	// ID of the user to validate.
	ID uuid.UUID `json:"id"`
	// Code is the validation code sent by email, after user registered/asked for an email update.
	Code string `json:"code"`
}

// UserEmailUpdateForm is the form sent by a user to update its email.
type UserEmailUpdateForm struct {
	// Email is the new email the user wants to use.
	Email string `json:"email"`
}

// UserPasswordResetForm is the form sent by a user who forgot its password.
type UserPasswordResetForm struct {
	// Email is the currently active email of the user. A password reset link will be sent here, if a user exists.
	Email string `json:"email"`
}

// UserPasswordUpdateForm is the form sent by a user to update its password.
type UserPasswordUpdateForm struct {
	// ID of the user to update.
	ID uuid.UUID `json:"id"`
	// Password is the new password for the user account.
	Password string `json:"password"`
	// OldPassword is the current password of the user. If password reset was called, it must be the
	// validation code instead.
	OldPassword string `json:"oldPassword"`
}

// UserAuthorizations is a list of authorizations the current user has.
type UserAuthorizations [][]string

const (
	// UserAuthorizationsAccountValidated is a special authorization, set once the user validated its account at least
	// once.
	UserAuthorizationsAccountValidated = "account-validated"
)
