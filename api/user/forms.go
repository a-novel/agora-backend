package userapi

import (
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"time"
)

type RegisterForm struct {
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Birthday  time.Time  `json:"birthday"`
	Sex       models.Sex `json:"sex"`
	Username  string     `json:"username"`
	Slug      string     `json:"slug"`
}

type IdentityUpdateForm struct {
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Birthday  time.Time  `json:"birthday"`
	Sex       models.Sex `json:"sex"`
}

type ProfileUpdateForm struct {
	Username string `json:"username"`
	Slug     string `json:"slug"`
}

type PasswordUpdateForm struct {
	ID          uuid.UUID `json:"id"`
	Password    string    `json:"password"`
	OldPassword string    `json:"oldPassword"`
}

type EmailUpdateForm struct {
	Email string `json:"email"`
}

type PasswordResetForm struct {
	Email string `json:"email"`
}

type ValidateEmailForm struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
}

type EmailExistsForm struct {
	Email string `json:"email"`
}

type SlugExistsForm struct {
	Slug string `json:"slug"`
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReadProfileForm struct {
	Slug string `uri:"slug"`
}

type SearchProfileForm struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type PreviewProfilesForm struct {
	IDs []uuid.UUID `json:"ids"`
}
