package improve_request_storage

import (
	"context"
	"fmt"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// Repository of the current layer. You can instantiate a new one with NewRepository.
type Repository interface {
	// Read reads a single post revision, based on its ID.
	Read(ctx context.Context, id uuid.UUID) (*Model, error)
	// ReadRevisions reads every revision, related to a source. The ID must be the one of the source post.
	ReadRevisions(ctx context.Context, id uuid.UUID) ([]*Model, error)

	// Create creates a brand-new post. The returned model will have matching Model.Source and Model.ID.
	Create(ctx context.Context, userID uuid.UUID, title, content string, id uuid.UUID, now time.Time) (*Model, error)
	// CreateRevision creates a new revision for a given post. The ID must be the one of the source post.
	CreateRevision(ctx context.Context, userID, sourceID uuid.UUID, title, content string, id uuid.UUID, now time.Time) (*Model, error)

	// Delete a single revision for a post. If the provided id is the source id, then all associated revisions will
	// also be deleted.
	Delete(ctx context.Context, requestID uuid.UUID) error

	// Search returns a list of posts, matching the provided query. Results must be paginated using the limit and
	// offset parameters.
	// It also returns the total number of available results, to help with pagination.
	Search(ctx context.Context, query SearchQuery, limit, offset int) ([]*Preview, int64, error)

	// IsCreator returns whether the user is a creator of the improvement suggestion. ID can be the id of any revision.
	// To only check if the user is the creator of the specific revision, set strict flag to true.
	IsCreator(ctx context.Context, userID, postID uuid.UUID, strict bool) (bool, error)

	GetPreviews(ctx context.Context, ids []uuid.UUID) ([]*Preview, error)
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

// Return a table for the current model, with aggregated stats. This can be used in a TableExpr clause.
func (repository *repositoryImpl) selectModelWithStats() *bun.SelectQuery {
	return repository.db.NewSelect().
		Model((*Model)(nil)).
		// Manually add columns hidden from the model.
		Column("*").
		// Stats.
		ColumnExpr("SUM(up_votes) OVER (PARTITION BY source) AS total_up_votes").
		ColumnExpr("SUM(down_votes) OVER (PARTITION BY source) AS total_down_votes").
		ColumnExpr("COUNT(*) OVER (PARTITION BY source) AS revision_count")
}

// Return a column selector, that will count the more recent revisions of a given post. Alias is the name of the
// source table.
func (repository *repositoryImpl) selectMoreRecentRevisions(alias string) *bun.SelectQuery {
	return repository.db.NewSelect().
		ColumnExpr("COUNT(*)").
		TableExpr("improve_requests AS more_recent").
		Where(fmt.Sprintf("more_recent.source = %s.source", alias)).
		Where(fmt.Sprintf("more_recent.created_at > %s.created_at", alias))
}

func (repository *repositoryImpl) selectSuggestions(alias string) *bun.SelectQuery {
	return repository.db.NewSelect().
		ColumnExpr("COUNT(*)").
		TableExpr("improve_suggestions").
		Where(fmt.Sprintf("improve_suggestions.source_id = %s.source", alias))
}

func (repository *repositoryImpl) selectValidatedSuggestions(alias string) *bun.SelectQuery {
	return repository.selectSuggestions(alias).Where("improve_suggestions.validated = TRUE")
}

// Select columns for a Preview model. Alias is the name of the source table. The source table must contain
// the stats columns (see selectModelWithStats).
func (repository *repositoryImpl) selectPreview(alias string) *bun.SelectQuery {
	return repository.db.NewSelect().
		Column(
			fmt.Sprintf("%s.id", alias),
			fmt.Sprintf("%s.source", alias),
			fmt.Sprintf("%s.created_at", alias),
			fmt.Sprintf("%s.user_id", alias),
			fmt.Sprintf("%s.title", alias),
			fmt.Sprintf("%s.revision_count", alias),
		).
		ColumnExpr("(?) AS more_recent_revisions", repository.selectMoreRecentRevisions(alias)).
		ColumnExpr("(?) AS suggestions_count", repository.selectSuggestions(alias)).
		ColumnExpr("(?) AS accepted_suggestions_count", repository.selectValidatedSuggestions(alias)).
		ColumnExpr(fmt.Sprintf("%s.content::VARCHAR(?)", alias), repository.cropPreviewContent).
		ColumnExpr(fmt.Sprintf("%s.total_up_votes as up_votes", alias)).
		ColumnExpr(fmt.Sprintf("%s.total_down_votes as down_votes", alias))
}

func (repository *repositoryImpl) Read(ctx context.Context, id uuid.UUID) (*Model, error) {
	model := &Model{ID: id}
	if err := repository.db.NewSelect().Model(model).Column(exposedColumns...).WherePK().Scan(ctx); err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) ReadRevisions(ctx context.Context, id uuid.UUID) ([]*Model, error) {
	var models []*Model
	if err := repository.db.NewSelect().
		Model(&models).
		Column(exposedColumns...).
		ColumnExpr("(?) AS suggestions_count", repository.selectSuggestions("i")).
		ColumnExpr("(?) AS accepted_suggestions_count", repository.selectValidatedSuggestions("i")).
		Where("source = ?", id).
		Order("created_at DESC").
		Scan(ctx); err != nil {
		return nil, validation.HandlePGError(err)
	}

	if len(models) == 0 {
		return nil, validation.ErrNotFound
	}

	return models, nil
}

func (repository *repositoryImpl) Create(ctx context.Context, userID uuid.UUID, title, content string, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{
		ID:        id,
		Source:    id,
		CreatedAt: now,
		UserID:    userID,
		Title:     title,
		Content:   content,
	}

	if err := repository.db.NewInsert().Model(model).Returning(exposedColumnsSTR).Scan(ctx); err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) CreateRevision(ctx context.Context, userID, sourceID uuid.UUID, title, content string, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{
		ID:        id,
		Source:    sourceID,
		CreatedAt: now,
		UserID:    userID,
		Title:     title,
		Content:   content,
	}

	// Ensure source exists
	count, err := repository.db.NewSelect().Model(new(Model)).Where("id = ?", sourceID).Where("source = ?", sourceID).Count(ctx)
	if err != nil {
		return nil, validation.HandlePGError(err)
	}
	if count == 0 {
		return nil, validation.ErrMissingRelation
	}

	if err := repository.db.NewInsert().Model(model).Returning(exposedColumnsSTR).Scan(ctx); err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) Delete(ctx context.Context, requestID uuid.UUID) error {
	model := &Model{ID: requestID}
	if _, err := repository.db.NewDelete().Model(model).Where("source = ?0::uuid OR id = ?0::uuid", requestID).Exec(ctx); err != nil {
		return validation.HandlePGError(err)
	}

	return nil
}

func (repository *repositoryImpl) Search(ctx context.Context, query SearchQuery, limit, offset int) ([]*Preview, int64, error) {
	var results []*Preview

	queryRequestWithStats := repository.selectModelWithStats()

	// Filter previous query to keep only the latest revisions for the current query.
	queryLatestRevision := repository.db.NewSelect().
		TableExpr("(?) as with_stats", queryRequestWithStats).
		Column("*").
		// Filter latest revision.
		DistinctOn("with_stats.source").
		Order("with_stats.source", "with_stats.created_at DESC")

	// Apply filters.
	if query.UserID != nil {
		queryLatestRevision = queryLatestRevision.Where("user_id = ?", query.UserID)
	}

	queryPreviews := repository.selectPreview("i").
		TableExpr("(?) as i", queryLatestRevision).
		Limit(limit).
		Offset(offset)

	// Use FullText search filter.
	if query.Query != "" {
		queryFullText := repository.db.NewSelect().
			ColumnExpr("to_tsquery('french', string_agg(lexeme || ':*', ' & ' order by positions)) AS query").
			TableExpr("unnest(to_tsvector('french', unaccent(?)))", query.Query)

		queryPreviews = queryPreviews.
			TableExpr("(?) AS search", queryFullText).
			Where("i.text_searchable_index_col @@ search.query").
			OrderExpr("ts_rank_cd(i.text_searchable_index_col, search.query) DESC")
	}

	if query.Order != nil {
		if query.Order.Score {
			queryPreviews = queryPreviews.OrderExpr("i.up_votes - i.down_votes DESC")
		}
	}

	// Order by date by default.
	queryPreviews = queryPreviews.OrderExpr("i.created_at DESC")

	count, err := queryPreviews.ScanAndCount(ctx, &results)
	if err != nil {
		return nil, 0, validation.HandlePGError(err)
	}

	return results, int64(count), nil
}

func (repository *repositoryImpl) IsCreator(ctx context.Context, userID, postID uuid.UUID, strict bool) (bool, error) {
	var (
		err error
		ok  bool
	)

	if strict {
		// The user is the author of the targeted revision.
		ok, err = repository.db.NewSelect().
			Model((*Model)(nil)).
			Where("id = ?", postID).
			Where("user_id = ?", userID).
			Exists(ctx)
		if err != nil {
			return false, validation.HandlePGError(err)
		}
	} else {
		querySameSource := repository.db.NewSelect().
			ColumnExpr("1").
			TableExpr("improve_requests AS same_source").
			Where("same_source.source = improve_requests.source").
			Where("same_source.user_id = ?", userID)

		// The user is the author of any revision on the same source as the targeted revision.
		ok, err = repository.db.NewSelect().
			Model((*Model)(nil)).
			Where("id = ?", postID).
			Where("EXISTS(?)", querySameSource).
			Exists(ctx)
		if err != nil {
			return false, validation.HandlePGError(err)
		}
	}

	return ok, nil
}

func (repository *repositoryImpl) GetPreviews(ctx context.Context, ids []uuid.UUID) ([]*Preview, error) {
	var results []*Preview

	queryRequestWithStats := repository.selectModelWithStats()

	err := repository.selectPreview("i").
		TableExpr("(?) as i", queryRequestWithStats).
		Where("i.id IN (?)", bun.In(ids)).
		Scan(ctx, &results)
	if err != nil {
		return nil, validation.HandlePGError(err)
	}

	return results, nil
}
