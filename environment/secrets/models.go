package secrets

import (
	"crypto/ed25519"
	"github.com/google/uuid"
	"time"
)

type Config struct {
	KeysService jwk_service.Service
	KeyGen      func() (ed25519.PrivateKey, error)

	Now func() time.Time
	ID  func() uuid.UUID

	MaxBackups     int
	UpdateInterval time.Duration
}
