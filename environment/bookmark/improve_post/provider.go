package improve_post

import (
	"context"
	"fmt"
	user_service "github.com/a-novel/agora-backend/domains/user/service/user"
	"github.com/a-novel/agora-backend/environment/user/authentication"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"time"
)

type Provider interface {
	Bookmark(ctx context.Context, token string, requestID uuid.UUID, target models.BookmarkTarget, level models.BookmarkLevel) (*models.Bookmark, error)
	UnBookmark(ctx context.Context, token string, requestID uuid.UUID, target models.BookmarkTarget) error
	IsBookmarked(ctx context.Context, userID, requestID uuid.UUID, target models.BookmarkTarget) (*models.BookmarkLevel, error)
	List(ctx context.Context, userID uuid.UUID, level models.BookmarkLevel, target models.BookmarkTarget, limit, offset int) ([]*models.Bookmark, int64, error)
}

type providerImpl struct {
	bookmarkService improve_post_service.Service
	tokenService    token_service.Service
	keysService     jwk_service.ServiceCached
	userService     user_service.Service

	time func() time.Time
}

type Config struct {
	BookmarkService improve_post_service.Service
	TokenService    token_service.Service
	KeysService     jwk_service.ServiceCached
	UserService     user_service.Service

	Time func() time.Time
}

func NewProvider(config Config) Provider {
	return &providerImpl{
		bookmarkService: config.BookmarkService,
		tokenService:    config.TokenService,
		keysService:     config.KeysService,
		userService:     config.UserService,

		time: config.Time,
	}
}

func (provider *providerImpl) Bookmark(ctx context.Context, token string, requestID uuid.UUID, target models.BookmarkTarget, level models.BookmarkLevel) (*models.Bookmark, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return nil, err
	}

	ok, err := provider.userService.HasAuthorizations(ctx, claims.Payload.ID, models.UserAuthorizations{
		{models.UserAuthorizationsAccountValidated},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to check user authorizations: %w", err)
	}
	if !ok {
		return nil, validation.NewErrUnauthorized("user email is not validated")
	}

	res, err := provider.bookmarkService.Bookmark(ctx, claims.Payload.ID, requestID, target, level, now)
	if err != nil {
		return nil, fmt.Errorf("unable to bookmark %q %q: %w", target, requestID, err)
	}

	return res, nil
}

func (provider *providerImpl) UnBookmark(ctx context.Context, token string, requestID uuid.UUID, target models.BookmarkTarget) error {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return err
	}

	if err = provider.bookmarkService.UnBookmark(ctx, claims.Payload.ID, requestID, target); err != nil {
		return fmt.Errorf("unable to unbookmark %q %q: %w", target, requestID, err)
	}

	return nil
}

func (provider *providerImpl) IsBookmarked(ctx context.Context, userID, requestID uuid.UUID, target models.BookmarkTarget) (*models.BookmarkLevel, error) {
	res, err := provider.bookmarkService.IsBookmarked(ctx, userID, requestID, target)
	if err != nil {
		return nil, fmt.Errorf("unable to check if %q %q is bookmarked: %w", target, requestID, err)
	}

	return res, nil
}

func (provider *providerImpl) List(ctx context.Context, userID uuid.UUID, level models.BookmarkLevel, target models.BookmarkTarget, limit, offset int) ([]*models.Bookmark, int64, error) {
	res, total, err := provider.bookmarkService.List(ctx, userID, level, target, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to list bookmarks: %w", err)
	}

	return res, total, nil
}
