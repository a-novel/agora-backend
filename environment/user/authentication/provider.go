package authentication

import (
	"context"
	"fmt"
	jwk_service "github.com/a-novel/agora-backend/domains/keys/service/jwk"
	"github.com/a-novel/agora-backend/domains/user/service/credentials"
	"github.com/a-novel/agora-backend/domains/user/service/token"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"time"
)

type Provider interface {
	// Authenticate checks whether the provided token is valid or not. It returns an error if the token is incorrect.
	// Otherwise, it returns the token itself.
	//
	// If the autoRenew option is set to true, the token will be renewed when close to expiration, based on the
	// provider's configuration; in this case, the new token will be returned (otherwise, the current token is
	// returned).
	//
	// It is generally safer to use the returned token no matter what.
	//
	//  var err error
	//  // Replace the current token value with the one returned by the function.
	//  token, err = provider.Authenticate(ctx, token, autoRenew)
	Authenticate(ctx context.Context, token string, autoRenew bool) (string, error)
	// Login the user. On success, it returns the user's token.
	Login(ctx context.Context, form models.UserCredentialsLoginForm) (string, error)
}

type Config struct {
	CredentialsService credentials_service.Service
	TokenService       token_service.Service
	KeysService        jwk_service.ServiceCached

	Time func() time.Time
	ID   func() uuid.UUID

	TokenTTL        time.Duration
	TokenRenewDelta time.Duration
}

type providerImpl struct {
	credentialsService credentials_service.Service
	tokenService       token_service.Service
	keysService        jwk_service.ServiceCached
	time               func() time.Time
	id                 func() uuid.UUID

	tokenTTL        time.Duration
	tokenRenewDelta time.Duration
}

func NewProvider(cfg Config) Provider {
	return &providerImpl{
		credentialsService: cfg.CredentialsService,
		tokenService:       cfg.TokenService,
		keysService:        cfg.KeysService,
		time:               cfg.Time,
		id:                 cfg.ID,

		tokenTTL:        cfg.TokenTTL,
		tokenRenewDelta: cfg.TokenRenewDelta,
	}
}

func (provider *providerImpl) Authenticate(_ context.Context, token string, autoRenew bool) (string, error) {
	now := provider.time()
	claims, err := ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return "", err
	}

	if autoRenew && claims.Header.EXP.Sub(now) < provider.tokenRenewDelta {
		newToken, err := provider.tokenService.Encode(
			claims.Payload,
			provider.tokenTTL,
			provider.keysService.GetPrivate(),
			provider.id(),
			now,
		)
		// Allow normal processing if token renewal fails. We let the API handle the error softly.
		if err != nil {
			return token, fmt.Errorf(
				"failed to generate new token for user %q: %w", claims.Payload.ID.String(), err,
			)
		}

		token = newToken
	}

	return token, nil
}

func (provider *providerImpl) Login(ctx context.Context, form models.UserCredentialsLoginForm) (string, error) {
	credentials, err := provider.credentialsService.Authenticate(ctx, &form)
	if err != nil {
		return "", fmt.Errorf("failed to login user with email %q: %w", form.Email, err)
	}

	token, err := provider.tokenService.Encode(
		models.UserTokenPayload{ID: credentials.ID},
		provider.tokenTTL,
		provider.keysService.GetPrivate(),
		provider.id(),
		provider.time(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate token for user %q: %w", form.Email, err)
	}

	return token, nil
}
