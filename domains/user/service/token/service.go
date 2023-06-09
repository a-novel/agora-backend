package token_service

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"regexp"
	"strings"
	"time"
)

var (
	// Quickly identify if the string is parsable as a signed token.
	signedTokenRegexp = regexp.MustCompile(`^([a-zA-Z\d-_]{2,}).([a-zA-Z\d-_]{2,}).([a-zA-Z\d-_]{2,})$`)
)

type Service interface {
	Encode(data models.UserTokenPayload, ttl time.Duration, signatureKey ed25519.PrivateKey, id uuid.UUID, now time.Time) (string, error)
	Decode(source string, signatureKeys []ed25519.PublicKey, now time.Time) (*models.UserToken, error)
}

type serviceImpl struct{}

func NewService() Service {
	return new(serviceImpl)
}

func (service *serviceImpl) Encode(data models.UserTokenPayload, ttl time.Duration, signatureKey ed25519.PrivateKey, id uuid.UUID, now time.Time) (string, error) {
	if signatureKey == nil {
		return "", fmt.Errorf("no signature key provided")
	}

	source := models.UserToken{
		Header: models.UserTokenHeader{
			IAT: now,
			EXP: now.Add(ttl),
			ID:  id,
		},
		Payload: data,
	}

	mrshHeader, err := json.Marshal(source.Header)
	if err != nil {
		return "", fmt.Errorf("failed to encode token header: %w", err)
	}
	header := base64.RawURLEncoding.EncodeToString(mrshHeader)

	mrshPayload, err := json.Marshal(source.Payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode token payload: %w", err)
	}
	payload := base64.RawURLEncoding.EncodeToString(mrshPayload)

	unsigned := fmt.Sprintf("%s.%s", header, payload)
	signature := base64.RawURLEncoding.EncodeToString(ed25519.Sign(signatureKey, []byte(unsigned)))
	return fmt.Sprintf("%s.%s", unsigned, signature), nil
}

func (service *serviceImpl) Decode(source string, signatureKeys []ed25519.PublicKey, now time.Time) (*models.UserToken, error) {
	if source == "" {
		return nil, validation.NewErrNil("token")
	}
	if err := validation.CheckRegexp("token", source, signedTokenRegexp); err != nil {
		return nil, err
	}

	parts := strings.Split(source, ".")
	header := parts[0]
	payload := parts[1]
	signature := parts[2]

	decodedSignature, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return nil, validation.NewErrInvalidCredentials("the signature is invalid")
	}

	var found bool
	for _, signatureKey := range signatureKeys {
		if ok := ed25519.Verify(signatureKey, []byte(fmt.Sprintf("%s.%s", header, payload)), decodedSignature); ok {
			found = true
			break
		}
	}
	if !found {
		return nil, validation.NewErrInvalidCredentials("no signature keys can decode the current token signature")
	}

	token := new(models.UserToken)
	decodedHeader, err := base64.RawURLEncoding.DecodeString(header)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token header: %w", err)
	}
	if err := json.Unmarshal(decodedHeader, &token.Header); err != nil {
		return nil, fmt.Errorf("unable to unmarshal token header: %w", err)
	}

	decodedPayload, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %w", err)
	}
	if err := json.Unmarshal(decodedPayload, &token.Payload); err != nil {
		return nil, fmt.Errorf("unable to unmarshal token payload: %w", err)
	}

	if token.Header.ID == uuid.Nil {
		return nil, validation.NewErrInvalidEntity("header.id", "ID cannot be empty")
	}

	if token.Header.IAT.After(now) {
		return nil, validation.NewErrInvalidCredentials(fmt.Sprintf("token is not available until %s", token.Header.IAT))
	}
	if token.Header.EXP.Before(now) {
		return nil, validation.NewErrInvalidCredentials(fmt.Sprintf("token has expired since %s", token.Header.EXP))
	}

	return token, nil
}
