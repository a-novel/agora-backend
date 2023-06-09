package security

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const (
	// Default length for a code.
	defaultLength = 32
)

// GenerateCode generates a random and secure hash. It returns (in order), the raw code and a bcrypt hashed
// version of it (both encoded in base64 URL), along with an error if anything unexpected occurs.
//
//	// Save hashed in a secure location on your server,
//	// and send the code to the user.
//	code, encrypted, err := validator.GenerateCode()
func GenerateCode() (string, string, error) {
	code := []byte(Random(defaultLength))

	encrypted, err := bcrypt.GenerateFromPassword(code, bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate secure hash: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(code), base64.RawURLEncoding.EncodeToString(encrypted), nil
}

// VerifyCode compare a code and its encrypted version. It returns true if the code matches the encrypted version,
// false otherwise.
//
// The error indicates an unexpected error, meaning the code cannot be validated.
//
//	ok, err := validator.VerifyCode(code, encrypted)
func VerifyCode(code string, encrypted string) (bool, error) {
	codeParsed, err := base64.RawURLEncoding.DecodeString(code)
	if err != nil {
		return false, fmt.Errorf("failed to decode code: %w", err)
	}
	encryptedParsed, err := base64.RawURLEncoding.DecodeString(encrypted)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(encryptedParsed, codeParsed)
	if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		return false, fmt.Errorf("failed to verify code against hash: %w", err)
	}

	return err == nil, nil
}
