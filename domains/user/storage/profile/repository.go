package profile_storage

import (
	"context"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// Repository of the current layer. You can instantiate a new one with NewRepository.
type Repository interface {
	// Read reads a profile object, based on a user id.
	Read(ctx context.Context, id uuid.UUID) (*Model, error)
	// ReadSlug reads a profile object, based on a slug.
	ReadSlug(ctx context.Context, slug string) (*Model, error)
	// SlugExists looks if a given slug is already used by another profile.
	SlugExists(ctx context.Context, slug string) (bool, error)

	// Update the slug of the targeted user.
	Update(ctx context.Context, data *Core, id uuid.UUID, now time.Time) (*Model, error)
}

func NewRepository(db bun.IDB) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db bun.IDB
}

func (repository *repositoryImpl) Read(ctx context.Context, id uuid.UUID) (*Model, error) {
	model := &Model{ID: id}
	if err := repository.db.NewSelect().Model(model).WherePK().Scan(ctx); err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) ReadSlug(ctx context.Context, slug string) (*Model, error) {
	model := new(Model)
	if err := repository.db.NewSelect().Model(model).Where("slug = ?", slug).Scan(ctx); err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) SlugExists(ctx context.Context, slug string) (bool, error) {
	count, err := repository.db.NewSelect().Model(new(Model)).Where("slug = ?", slug).Count(ctx)
	return count > 0, validation.HandlePGError(err)
}

func (repository *repositoryImpl) Update(ctx context.Context, data *Core, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{ID: id, UpdatedAt: &now, Core: *data}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		Column(
			"slug",
			"username",
			"updated_at",
		).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, validation.HandlePGError(err)
	}
	if err := validation.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}
