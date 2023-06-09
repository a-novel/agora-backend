package authentication

import (
	"context"
	"fmt"
	jwk_service "github.com/a-novel/agora-backend/domains/keys/service/jwk"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"google.golang.org/api/oauth2/v2"
	"strings"
	"time"
)

// ForceAuthentication verifies the given token, and return an error if empty or not valid.
func ForceAuthentication(token string, service token_service.Service, keys jwk_service.ServiceCached, now time.Time) (*models.UserToken, error) {
	if token == "" {
		return nil, fmt.Errorf("%w: no token found", validation.ErrInvalidCredentials)
	}

	claims, err := service.Decode(token, keys.ListPublic(), now)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	return claims, nil
}

type BackendServiceAuth struct {
	UserAgent     string
	Authorization string
	AllowedUsers  []string
}

// ForceBackendService verifies the given token authenticates a GCP backend service, and return an error if empty or not
// valid.
// If the provided auth config is nil, this method does nothing.
func ForceBackendService(ctx context.Context, auth *BackendServiceAuth) error {
	if auth == nil {
		return nil
	}

	// https://jackcuthbert.dev/blog/verifying-google-cloud-scheduler-requests-in-cloud-run-with-typescript
	if auth.UserAgent != "Google-Cloud-Scheduler" {
		return fmt.Errorf(
			"%w: bad user agent: expected %q, got %q",
			validation.ErrInvalidCredentials, "Google-Cloud-Scheduler", auth.UserAgent,
		)
	}

	// https://stackoverflow.com/questions/53181297/verify-http-request-from-google-cloud-scheduler
	if auth.Authorization == "" {
		return fmt.Errorf("%w: missing authorization header", validation.ErrInvalidCredentials)
	}

	idToken := strings.Split(auth.Authorization, "Bearer ")[0]

	authenticator, err := oauth2.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire authenticator: %w", err)
	}

	info, err := authenticator.Tokeninfo().IdToken(idToken).Do()
	if err != nil {
		return fmt.Errorf("failed to retrieve token information: %w", err)
	}

	for _, allowedUser := range auth.AllowedUsers {
		if info.Email == allowedUser {
			return nil
		}
	}

	return fmt.Errorf("unexpected user %q", info.Email)
}
