package jwk_storage

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
)

func writeKeyToOutput(out io.Writer, key ed25519.PrivateKey) error {
	marshalledKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to encode key: %w", err)
	}

	if err := pem.Encode(out, &pem.Block{Type: "PRIVATE KEY", Bytes: marshalledKey}); err != nil {
		return fmt.Errorf("failed to encode key into file: %w", err)
	}

	return nil
}

func unmarshalPrivateKey(data []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("file does not contain a valid ed25519 private key: no block found")
	}
	keyData, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := keyData.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("file does not contain a valid ed25519 private key: unexpected key type")
	}

	return key, nil
}
