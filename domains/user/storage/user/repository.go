package user_storage

import (
	"context"
	"fmt"
	"github.com/a-novel/agora-backend/domains/user/storage/credentials"
	"github.com/a-novel/agora-backend/domains/user/storage/identity"
	"github.com/a-novel/agora-backend/domains/user/storage/profile"
	"github.com/a-novel/agora-backend/domains/user/storage/user/queries"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// Repository of the current layer. You can instantiate a new one with NewRepository.
type Repository interface {
	// Create creates a new user. The credentials, identity and profile objects will share the same ID and create time.
	// If any error occurs, no data is created.
	Create(ctx context.Context, data *Core, id uuid.UUID, now time.Time) (*Model, error)
	// Delete deletes a user. Credentials and identity are only soft-deleted, with sensitive data erased, while
	// profile is permanently destroyed.
	// If an error occurs, all data will remain, unaltered.
	Delete(ctx context.Context, id uuid.UUID, now time.Time) (*Model, error)
	// Search performs a cross-table search query over the user repository.
	Search(ctx context.Context, query string, limit, offset int) ([]*PublicPreview, int64, error)
	// GetPreview returns a preview of a user, for private display within the application.
	GetPreview(ctx context.Context, id uuid.UUID) (*Preview, error)
	// GetPublic returns a user information, for public display on the application.
	GetPublic(ctx context.Context, slug string) (*Public, error)
	// GetPublicPreviews returns a list of user previews, for public display on the application.
	GetPublicPreviews(ctx context.Context, ids []uuid.UUID) ([]*PublicPreview, error)
}

func NewRepository(db bun.IDB) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db bun.IDB
}

func (repository *repositoryImpl) Create(ctx context.Context, data *Core, id uuid.UUID, now time.Time) (*Model, error) {
	model := new(Model)

	// Create all in a transaction, to avoid partially created users if any part of the operation fails.
	err := repository.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		credentialsModel := &credentials_storage.Model{ID: id, CreatedAt: now, Core: data.Credentials}
		_, err := repository.db.NewInsert().Model(credentialsModel).Exec(ctx)
		if err != nil {
			return err
		}

		identityModel := &identity_storage.Model{ID: id, CreatedAt: now, Core: data.Identity}
		_, err = repository.db.NewInsert().Model(identityModel).Exec(ctx)
		if err != nil {
			return err
		}

		profileModel := &profile_storage.Model{ID: id, CreatedAt: now, Core: data.Profile}
		_, err = repository.db.NewInsert().Model(profileModel).Exec(ctx)
		if err != nil {
			return err
		}

		// Assign common metadata.
		model.ID = id
		model.CreatedAt = now

		model.Credentials = credentialsModel.Core
		model.Identity = identityModel.Core
		model.Profile = profileModel.Core
		return nil
	})
	if err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) Delete(ctx context.Context, id uuid.UUID, now time.Time) (*Model, error) {
	model := new(Model)

	// Delete all in a transaction, to avoid partially deleted users if any part of the operation fails.
	err := repository.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		identityModel := &identity_storage.Model{ID: id, DeletedAt: &now}
		res, err := repository.db.NewUpdate().Model(identityModel).
			WherePK().
			// Soft delete, remove sensitive information that are not useful for statistics.
			Column("first_name", "last_name", "deleted_at").
			Returning("*").
			Exec(ctx)
		if err != nil {
			return err
		}
		if err := validation.ForceRowsUpdate(res); err != nil {
			return err
		}

		credentialsModel := &credentials_storage.Model{ID: id}
		_, err = repository.db.NewDelete().Model(credentialsModel).WherePK().Exec(ctx)
		if err != nil {
			return err
		}

		profileModel := &profile_storage.Model{ID: id}
		_, err = repository.db.NewDelete().Model(profileModel).WherePK().Exec(ctx)
		if err != nil {
			return err
		}

		// Assign common metadata.
		model.ID = id
		model.CreatedAt = identityModel.CreatedAt
		model.UpdatedAt = identityModel.UpdatedAt
		model.DeletedAt = &now

		model.Identity = identityModel.Core
		return nil
	})
	if err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) Search(ctx context.Context, query string, limit, offset int) ([]*PublicPreview, int64, error) {
	var (
		count   int64
		results []*PublicPreview
	)

	rows, err := repository.db.QueryContext(ctx, user_queries.SearchUserQuery, query, limit, offset)
	if err != nil {
		return nil, 0, validation.HandlePGError(err)
	}

	defer rows.Close()

	for rows.Next() {
		var result PublicPreview

		if err := rows.Scan(
			&result.ID,
			&result.Username,
			&result.Slug,
			&result.FirstName,
			&result.LastName,
			&result.CreatedAt,
			&count,
		); err != nil {
			return nil, 0, err
		}

		results = append(results, &result)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, validation.HandlePGError(err)
	}

	return results, count, nil
}

func (repository *repositoryImpl) GetPreview(ctx context.Context, id uuid.UUID) (*Preview, error) {
	var (
		credentialsModel credentials_storage.Model
		profileModel     profile_storage.Model
		identityModel    identity_storage.Model
	)

	err := repository.db.NewSelect().Model(&credentialsModel).Column("email_user", "email_domain", "id").Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", validation.HandlePGError(err))
	}
	err = repository.db.NewSelect().Model(&identityModel).Column("first_name", "last_name", "sex").Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity: %w", validation.HandlePGError(err))
	}
	err = repository.db.NewSelect().Model(&profileModel).Column("username", "slug").Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", validation.HandlePGError(err))
	}

	return &Preview{
		ID:        credentialsModel.ID,
		Username:  profileModel.Username,
		FirstName: identityModel.FirstName,
		LastName:  identityModel.LastName,
		Email:     credentialsModel.Email,
		Slug:      profileModel.Slug,
		Sex:       identityModel.Sex,
	}, nil
}

func (repository *repositoryImpl) GetPublicPreviews(ctx context.Context, ids []uuid.UUID) ([]*PublicPreview, error) {
	var (
		profileModels  []profile_storage.Model
		identityModels []identity_storage.Model
	)

	err := repository.db.NewSelect().Model(&identityModels).
		Column("created_at", "first_name", "last_name", "id").
		Where("id IN (?)", bun.In(ids)).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity: %w", validation.HandlePGError(err))
	}

	err = repository.db.NewSelect().Model(&profileModels).
		Column("username", "slug", "id").
		Where("id IN (?)", bun.In(ids)).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", validation.HandlePGError(err))
	}

	results := make([]*PublicPreview, len(profileModels))
	for i, profileModel := range profileModels {
		// Stay stable in case order is not consistent.
		var identityModel identity_storage.Model
		for j := range identityModels {
			if identityModels[j].ID == profileModel.ID {
				identityModel = identityModels[j]
				break
			}
		}

		results[i] = &PublicPreview{
			ID:        profileModel.ID,
			Username:  profileModel.Username,
			Slug:      profileModel.Slug,
			FirstName: identityModel.FirstName,
			LastName:  identityModel.LastName,
			CreatedAt: identityModel.CreatedAt,
		}
	}

	return results, nil
}

func (repository *repositoryImpl) GetPublic(ctx context.Context, slug string) (*Public, error) {
	var (
		profileModel  profile_storage.Model
		identityModel identity_storage.Model
	)

	err := repository.db.NewSelect().Model(&profileModel).
		Column("username", "id").
		Where("slug = ?", slug).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", validation.HandlePGError(err))
	}

	err = repository.db.NewSelect().Model(&identityModel).
		Column("first_name", "last_name", "created_at", "sex").
		Where("id = ?", profileModel.ID).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity: %w", validation.HandlePGError(err))
	}

	return &Public{
		ID:        profileModel.ID,
		Username:  profileModel.Username,
		FirstName: identityModel.FirstName,
		LastName:  identityModel.LastName,
		CreatedAt: identityModel.CreatedAt,
		Sex:       identityModel.Sex,
	}, nil
}
