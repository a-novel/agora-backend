package jwk_storage

import (
	"crypto/ed25519"
	"time"
)

type Model struct {
	// Key returns the decoded key for the current entry.
	Key ed25519.PrivateKey
	// Date returns the date when the key was created.
	Date time.Time
	// Name of the record (file) that stores the entry.
	Name string
}
