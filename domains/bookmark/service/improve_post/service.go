package improve_post_service

import (
	"context"
	"fmt"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"time"
)

var (
	bookmarkTargetValues = []models.BookmarkTarget{models.BookmarkTargetImproveSuggestion, models.BookmarkTargetImproveRequest}
	levelValues          = []models.BookmarkLevel{models.BookmarkLevelBookmark, models.BookmarkLevelFavorite}
)

// Service of the current layer. You can instantiate a new one with NewService.
type Service interface {
	// Bookmark a post. If the post is already bookmarked by the user, the target bookmark will be updated.
	Bookmark(ctx context.Context, userID, requestID uuid.UUID, target models.BookmarkTarget, level models.BookmarkLevel, now time.Time) (*models.Bookmark, error)
	// UnBookmark a post. If the post is not bookmarked by the user, an error will be returned.
	UnBookmark(ctx context.Context, userID, requestID uuid.UUID, target models.BookmarkTarget) error
	// IsBookmarked checks whether a post is bookmarked by the user or not. If the post is bookmarked, the BookmarkLevel
	// is returned. Otherwise, nil is returned.
	// The error is only returned when something unexpected happens.
	IsBookmarked(ctx context.Context, userID, requestID uuid.UUID, target models.BookmarkTarget) (*models.BookmarkLevel, error)
	// List returns all the bookmarked post for a given user. Only one type of bookmark can be retrieved at time.
	// Results must be paginated using the limit and offset parameters.
	List(ctx context.Context, userID uuid.UUID, level models.BookmarkLevel, target models.BookmarkTarget, limit, offset int) ([]*models.Bookmark, int64, error)

	// StorageToModel converts a storage model to a service model.
	StorageToModel(source *improve_post_storage.Model) *models.Bookmark
}

type serviceImpl struct {
	repository improve_post_storage.Repository
}

// NewService returns a new Service instance.
// To use a mocked one, call NewMockService.
func NewService(repository improve_post_storage.Repository) Service {
	return &serviceImpl{repository: repository}
}

func (service *serviceImpl) Bookmark(ctx context.Context, userID, requestID uuid.UUID, target models.BookmarkTarget, level models.BookmarkLevel, now time.Time) (*models.Bookmark, error) {
	if err := validation.CheckRestricted("target", target, bookmarkTargetValues...); err != nil {
		return nil, err
	}
	if err := validation.CheckRestricted("level", level, levelValues...); err != nil {
		return nil, err
	}

	storageModel, err := service.repository.Bookmark(
		ctx, userID, requestID,
		improve_post_storage.BookmarkTarget(target),
		bookmark_storage.Level(level),
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bookmark post: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) UnBookmark(ctx context.Context, userID, requestID uuid.UUID, target models.BookmarkTarget) error {
	if err := validation.CheckRestricted("target", target, bookmarkTargetValues...); err != nil {
		return err
	}

	if err := service.repository.UnBookmark(ctx, userID, requestID, improve_post_storage.BookmarkTarget(target)); err != nil {
		return fmt.Errorf("failed to unbookmark post: %w", err)
	}

	return nil
}

func (service *serviceImpl) IsBookmarked(ctx context.Context, userID, requestID uuid.UUID, target models.BookmarkTarget) (*models.BookmarkLevel, error) {
	if err := validation.CheckRestricted("target", target, bookmarkTargetValues...); err != nil {
		return nil, err
	}

	level, err := service.repository.IsBookmarked(ctx, userID, requestID, improve_post_storage.BookmarkTarget(target))
	if err != nil {
		return nil, fmt.Errorf("failed to check if post is bookmarked: %w", err)
	}

	if level == nil {
		return nil, nil
	}

	levelValue := models.BookmarkLevel(*level)
	return &levelValue, nil
}

func (service *serviceImpl) List(ctx context.Context, userID uuid.UUID, level models.BookmarkLevel, target models.BookmarkTarget, limit, offset int) ([]*models.Bookmark, int64, error) {
	if err := validation.CheckRestricted("level", level, levelValues...); err != nil {
		return nil, 0, err
	}
	if err := validation.CheckRestricted("target", target, bookmarkTargetValues...); err != nil {
		return nil, 0, err
	}

	storageModels, total, err := service.repository.List(
		ctx, userID, bookmark_storage.Level(level), improve_post_storage.BookmarkTarget(target), limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list bookmarks: %w", err)
	}

	serviceModels := make([]*models.Bookmark, len(storageModels))
	for i, storageModel := range storageModels {
		serviceModels[i] = service.StorageToModel(storageModel)
	}

	return serviceModels, total, nil
}

func (service *serviceImpl) StorageToModel(source *improve_post_storage.Model) *models.Bookmark {
	if source == nil {
		return nil
	}

	return &models.Bookmark{
		UserID:    source.UserID,
		RequestID: source.RequestID,
		Target:    models.BookmarkTarget(source.Target),
		Level:     models.BookmarkLevel(source.Level),
		CreatedAt: source.CreatedAt,
	}
}
