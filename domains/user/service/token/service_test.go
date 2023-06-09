package token_service

import (
	"crypto/ed25519"
	"github.com/a-novel/agora-backend/domains/keys/storage/jwk"
	"github.com/a-novel/agora-backend/framework/test"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	baseTime = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
)

func TestTokenService_EncodeAndDecode(t *testing.T) {
	var (
		token, tokenOld string
		err             error
	)

	service := NewService()

	t.Run("Encode/Success", func(t *testing.T) {
		token, err = service.Encode(
			models.UserTokenPayload{ID: test_utils.NumberUUID(1000)}, time.Hour,
			jwk_storage.MockedKeys[0],
			test_utils.NumberUUID(1),
			baseTime,
		)
		require.NoError(t, err)
	})

	t.Run("Encode/Error/NoSignatureKey", func(t *testing.T) {
		_, err = service.Encode(
			models.UserTokenPayload{ID: test_utils.NumberUUID(1000)},
			time.Hour,
			nil,
			test_utils.NumberUUID(1),
			baseTime,
		)
		require.Error(t, err)
	})

	tokenOld, err = service.Encode(
		models.UserTokenPayload{ID: test_utils.NumberUUID(2000)},
		time.Hour, jwk_storage.MockedKeys[2],
		test_utils.NumberUUID(2),
		baseTime,
	)
	require.NoError(t, err)

	signatureKeys := []ed25519.PublicKey{
		jwk_storage.MockedKeys[0].Public().(ed25519.PublicKey),
		jwk_storage.MockedKeys[1].Public().(ed25519.PublicKey),
		jwk_storage.MockedKeys[2].Public().(ed25519.PublicKey),
	}

	data := []struct {
		name string

		source        string
		signatureKeys []ed25519.PublicKey
		now           time.Time

		expect    *models.UserToken
		expectErr error
	}{
		{
			name:          "Decode/Success",
			source:        token,
			signatureKeys: signatureKeys,
			now:           baseTime,
			expect: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime,
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1000),
				},
			},
		},
		{
			name:          "Decode/Success/OldToken",
			source:        tokenOld,
			signatureKeys: signatureKeys,
			now:           baseTime,
			expect: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime,
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(2),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(2000),
				},
			},
		},
		{
			name:          "Decode/Success/AfterIssuing",
			source:        token,
			signatureKeys: signatureKeys,
			now:           baseTime.Add(30 * time.Minute),
			expect: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime,
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1000),
				},
			},
		},
		{
			name:          "Decode/Error/NoKeyFound",
			source:        tokenOld,
			signatureKeys: signatureKeys[:1],
			now:           baseTime,
			expectErr:     validation.ErrInvalidCredentials,
		},
		{
			name:      "Decode/Error/NoKeysAtAll",
			source:    tokenOld,
			now:       baseTime,
			expectErr: validation.ErrInvalidCredentials,
		},
		{
			name:          "Decode/Error/NoToken",
			signatureKeys: signatureKeys,
			now:           baseTime,
			expectErr:     validation.ErrNil,
		},
		{
			name:          "Decode/Error/NotAToken",
			source:        "not a token",
			signatureKeys: signatureKeys,
			now:           baseTime,
			expectErr:     validation.ErrInvalidEntity,
		},
		{
			name:          "Decode/Error/NotASignedToken",
			source:        "foo.bar",
			signatureKeys: signatureKeys,
			now:           baseTime,
			expectErr:     validation.ErrInvalidEntity,
		},
		{
			name:          "Decode/Error/Expired",
			source:        token,
			signatureKeys: signatureKeys,
			now:           baseTime.Add(2 * time.Hour),
			expectErr:     validation.ErrInvalidCredentials,
		},
		{
			name:          "Decode/Error/NotIssuedYet",
			source:        token,
			signatureKeys: signatureKeys,
			now:           baseTime.Add(-time.Hour),
			expectErr:     validation.ErrInvalidCredentials,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			token, err := service.Decode(d.source, d.signatureKeys, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, token)
		})
	}
}
