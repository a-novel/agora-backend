package identity_storage

import (
	"context"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// Repository of the current layer. You can instantiate a new one with NewRepository.
type Repository interface {
	// Read reads an identity object, based on a user id.
	Read(ctx context.Context, id uuid.UUID) (*Model, error)
	// Update the identity of the targeted user.
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

func (repository *repositoryImpl) Update(ctx context.Context, data *Core, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{ID: id, UpdatedAt: &now, Core: *data}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		Column(
			"first_name",
			"last_name",
			"birthday",
			"sex",
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
