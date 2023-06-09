package improve_suggestion_storage

import (
	"context"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// Repository of the current layer. You can instantiate a new one with NewRepository.
type Repository interface {
	// Read returns the improvement suggestion with the given ID.
	Read(ctx context.Context, id uuid.UUID) (*Model, error)
	// Create creates a new improvement suggestion for a given improvement request revision.
	Create(ctx context.Context, data *Core, userID, sourceID, id uuid.UUID, now time.Time) (*Model, error)
	// Update updates an existing improvement suggestion.
	Update(ctx context.Context, data *Core, id uuid.UUID, now time.Time) (*Model, error)
	// Delete deletes an existing improvement suggestion.
	Delete(ctx context.Context, id uuid.UUID) error

	// Validate validates an existing improvement suggestion.
	Validate(ctx context.Context, validated bool, id uuid.UUID) (*Model, error)

	// List returns a list of improvement suggestions, matching the provided query. Results must be paginated using
	// the limit and offset parameters.
	// It also returns the total number of available results, to help with pagination.
	List(ctx context.Context, query ListQuery, limit, offset int) ([]*Model, int64, error)

	// IsCreator returns whether the user is the creator of the improvement suggestion.
	IsCreator(ctx context.Context, userID, postID uuid.UUID) (bool, error)

	GetPreviews(ctx context.Context, ids []uuid.UUID) ([]*Model, error)
}

// NewRepository returns a new Repository instance.
// To use a mocked one, call NewMockRepository.
func NewRepository(db bun.IDB, cropPreviewContent int) Repository {
	return &repositoryImpl{db: db, cropPreviewContent: cropPreviewContent}
}

type repositoryImpl struct {
	db                 bun.IDB
	cropPreviewContent int
}

func (repository *repositoryImpl) selectPreview() *bun.SelectQuery {
	return repository.db.NewSelect().
		Model((*models.ImproveSuggestion)(nil)).
		Column(
			"id",
			"created_at",
			"updated_at",
			"user_id",
			"source_id",
			"request_id",
			"validated",
			"title",
			"up_votes",
			"down_votes",
		).
		ColumnExpr("content::VARCHAR(?)", repository.cropPreviewContent)
}

func (repository *repositoryImpl) validateSource(ctx context.Context, sourceID, requestID uuid.UUID) error {
	// Source must exist.
	count, err := repository.db.NewSelect().Table("improve_requests").Where("id = ?", sourceID).Count(ctx)
	if err != nil {
		return validation.HandlePGError(err)
	}
	if count == 0 {
		return validation.ErrMissingRelation
	}

	// Revision with given source must exist.
	count, err = repository.db.NewSelect().Table("improve_requests").
		Where("id = ?", requestID).
		Where("source = ?", sourceID).
		Count(ctx)
	if err != nil {
		return validation.HandlePGError(err)
	}
	if count == 0 {
		return validation.ErrMissingRelation
	}

	return nil
}

func (repository *repositoryImpl) Read(ctx context.Context, id uuid.UUID) (*Model, error) {
	model := &Model{ID: id}
	if err := repository.db.NewSelect().Model(model).WherePK().Scan(ctx); err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) Create(ctx context.Context, data *Core, userID, sourceID, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{
		ID:        id,
		Core:      *data,
		SourceID:  sourceID,
		UserID:    userID,
		CreatedAt: now,
	}

	if err := repository.validateSource(ctx, sourceID, data.RequestID); err != nil {
		return nil, err
	}

	if err := repository.db.NewInsert().Model(model).Scan(ctx); err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) Update(ctx context.Context, data *Core, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{
		ID:        id,
		Core:      *data,
		UpdatedAt: &now,
	}

	err := repository.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := repository.db.NewUpdate().
			Model(model).
			WherePK().
			Column("id", "updated_at", "request_id", "title", "content").
			Returning("*").
			Scan(ctx); err != nil {
			return validation.HandlePGError(err)
		}

		if err := repository.validateSource(ctx, model.SourceID, data.RequestID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *repositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	model := &Model{ID: id}

	if res, err := repository.db.NewDelete().Model(model).WherePK().Exec(ctx); err != nil {
		return validation.HandlePGError(err)
	} else if err = validation.ForceRowsUpdate(res); err != nil {
		return err
	}

	return nil
}

func (repository *repositoryImpl) Validate(ctx context.Context, validated bool, id uuid.UUID) (*Model, error) {
	model := &Model{ID: id, Validated: validated}

	if err := repository.db.NewUpdate().
		Model(model).
		Column("validated").
		WherePK().
		Returning("*").
		Scan(ctx); err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) List(ctx context.Context, query ListQuery, limit, offset int) ([]*Model, int64, error) {
	var results []*Model

	dbQuery := repository.selectPreview().Limit(limit).Offset(offset)

	if query.UserID != nil {
		dbQuery = dbQuery.Where("user_id = ?", *query.UserID)
	}
	if query.SourceID != nil {
		dbQuery = dbQuery.Where("source_id = ?", *query.SourceID)
	}
	if query.RequestID != nil {
		dbQuery = dbQuery.Where("request_id = ?", *query.RequestID)
	}
	if query.Validated != nil {
		dbQuery = dbQuery.Where("validated = ?", *query.Validated)
	}

	if query.Order != nil {
		if query.Order.Score {
			dbQuery = dbQuery.OrderExpr("up_votes - down_votes DESC")
		}
	}

	// Order by date by default.
	dbQuery = dbQuery.OrderExpr("coalesce(updated_at, created_at) DESC")

	count, err := dbQuery.ScanAndCount(ctx, &results)
	if err != nil {
		return nil, 0, validation.HandlePGError(err)
	}

	return results, int64(count), nil
}

func (repository *repositoryImpl) IsCreator(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	ok, err := repository.db.NewSelect().
		Model((*Model)(nil)).
		Where("id = ?", postID).
		Where("user_id = ?", userID).
		Exists(ctx)
	if err != nil {
		return false, validation.HandlePGError(err)
	}

	return ok, nil
}

func (repository *repositoryImpl) GetPreviews(ctx context.Context, ids []uuid.UUID) ([]*Model, error) {
	var results []*Model

	dbQuery := repository.selectPreview().Where("id IN (?)", bun.In(ids))

	if err := dbQuery.Scan(ctx, &results); err != nil {
		return nil, validation.HandlePGError(err)
	}

	return results, nil
}
