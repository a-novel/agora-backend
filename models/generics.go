package models

import "fmt"

// Email represents an email address as a structure, rather than a single string. This facilitates indexing:
// for example, when looking for a user, only the User is relevant, so we may only index this field for searching.
// The String method converts the email object back to the standard string representation, in the format [user]@[domain].
type Email struct {
	// Validation is the hashed key used to validate a user email. The raw key is sent to the email address only.
	// When user manages to successfully prove its authenticity, the email is validated and this code is removed.
	// Like a password, the raw key should never be stored or cached.
	Validation string `json:"validationCode" bun:"validation_code"`
	// User of the email. This is the unique name that comes before the provider.
	User string `json:"user" bun:"user"`
	// Domain is the host of the mailing service provider, for example 'gmail.com'.
	Domain string `json:"domain" bun:"domain"`
}

// String converts the email object back to the standard string representation, in the format [user]@[domain].
// This method should never return an invalid format, so if the email is incomplete, it should return an empty string.
func (email Email) String() string {
	// Email is not valid, so it cannot be represented properly.
	if email.User == "" || email.Domain == "" {
		return ""
	}

	return fmt.Sprintf("%s@%s", email.User, email.Domain)
}

// Password represents a hashed password, that can be safely stored in the database.
type Password struct {
	// Validation is used to reset a password, for example when the original one has been forgotten. This field
	// contains the hashed key only. The raw key is sent to the user through a secure channel (an email address),
	// and once the user has managed to prove its identity, it can then create a new password.
	Validation string `json:"validationCode" bun:"validation_code"`
	// Hashed is the hashed password, used to validate user claims when trying to authenticate.
	Hashed string `json:"hashed" bun:"hashed"`
}

// Sex represents the biological gender of a user, either SexMale or SexFemale.
type Sex string

const (
	SexMale   Sex = "male"
	SexFemale Sex = "female"
)
