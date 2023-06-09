package improve_post_storage

import (
	"context"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// Repository of the current layer. You can instantiate a new one with NewRepository.
type Repository interface {
	// Bookmark a post. If the post is already bookmarked by the user, the target bookmark will be updated.
	Bookmark(ctx context.Context, userID, requestID uuid.UUID, target BookmarkTarget, level bookmark_storage.Level, now time.Time) (*Model, error)
	// UnBookmark a post. If the post is not bookmarked by the user, an error will be returned.
	UnBookmark(ctx context.Context, userID, requestID uuid.UUID, target BookmarkTarget) error
	// IsBookmarked checks whether a post is bookmarked by the user or not. If the post is bookmarked, the Level
	// is returned. Otherwise, nil is returned.
	// The error is only returned when something unexpected happens.
	IsBookmarked(ctx context.Context, userID, requestID uuid.UUID, target BookmarkTarget) (*bookmark_storage.Level, error)
	// List returns all the bookmarked post for a given user. Only one type of bookmark can be retrieved at time.
	// Results must be paginated using the limit and offset parameters.
	// It also returns the total number of available results, to help with pagination.
	List(ctx context.Context, userID uuid.UUID, level bookmark_storage.Level, target BookmarkTarget, limit, offset int) ([]*Model, int64, error)
}

// NewRepository returns a new Repository instance.
// To use a mocked one, call NewMockRepository.
func NewRepository(db bun.IDB) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db bun.IDB
}

func (repository *repositoryImpl) Bookmark(ctx context.Context, userID, requestID uuid.UUID, target BookmarkTarget, level bookmark_storage.Level, now time.Time) (*Model, error) {
	model := &Model{
		UserID:    userID,
		RequestID: requestID,
		Target:    target,
		Level:     level,
		CreatedAt: now,
	}

	if err := repository.db.NewInsert().
		Model(model).
		On("conflict (user_id, request_id, target) do update").
		Set("created_at = ?created_at").
		Set("level = ?level").
		Returning("*").
		Scan(ctx); err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) UnBookmark(ctx context.Context, userID, requestID uuid.UUID, target BookmarkTarget) error {
	model := &Model{
		UserID:    userID,
		RequestID: requestID,
		Target:    target,
	}

	if res, err := repository.db.NewDelete().Model(model).WherePK().Exec(ctx); err != nil {
		return validation.HandlePGError(err)
	} else if err = validation.ForceRowsUpdate(res); err != nil {
		return err
	}

	return nil
}

func (repository *repositoryImpl) IsBookmarked(ctx context.Context, userID, requestID uuid.UUID, target BookmarkTarget) (*bookmark_storage.Level, error) {
	model := &Model{
		UserID:    userID,
		RequestID: requestID,
		Target:    target,
	}

	if err := repository.db.NewSelect().Model(model).Column("level").WherePK().Scan(ctx); err != nil {
		err = validation.HandlePGError(err)
		if err == validation.ErrNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &model.Level, nil
}

func (repository *repositoryImpl) List(ctx context.Context, userID uuid.UUID, level bookmark_storage.Level, target BookmarkTarget, limit, offset int) ([]*Model, int64, error) {
	var models []*Model

	count, err := repository.db.NewSelect().
		Model(&models).
		Where("user_id = ? AND level = ? AND target = ?", userID, level, target).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, validation.HandlePGError(err)
	}

	return models, int64(count), nil
}
