package security

import "crypto/ed25519"

func JWKKeyGen() (ed25519.PrivateKey, error) {
	_, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
