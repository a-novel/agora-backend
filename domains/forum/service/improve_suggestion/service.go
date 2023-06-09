package improve_suggestion_service

import (
	"context"
	"fmt"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"regexp"
	"time"
)

var (
	// Just prevent line breaks in title.
	titleRegexp = regexp.MustCompile(`^[^\n\r]+$`)
)

const (
	MinTitleLength   = 4
	MaxTitleLength   = 256
	MinContentLength = 4
	MaxContentLength = 4096
)

// Service of the current layer. You can instantiate a new one with NewService.
type Service interface {
	// Read returns the improvement suggestion with the given ID.
	Read(ctx context.Context, id uuid.UUID) (*models.ImproveSuggestion, error)
	// Create creates a new improvement suggestion for a given improvement request revision.
	Create(ctx context.Context, data *models.ImproveSuggestionUpsert, userID, sourceID, id uuid.UUID, now time.Time) (*models.ImproveSuggestion, error)
	// Update updates an existing improvement suggestion.
	Update(ctx context.Context, data *models.ImproveSuggestionUpsert, id uuid.UUID, now time.Time) (*models.ImproveSuggestion, error)
	// Delete deletes an existing improvement suggestion.
	Delete(ctx context.Context, id uuid.UUID) error

	// Validate validates an existing improvement suggestion.
	Validate(ctx context.Context, validated bool, id uuid.UUID) (*models.ImproveSuggestion, error)

	// List returns a list of improvement suggestions, matching the provided query. Results must be paginated using
	// the limit and offset parameters.
	// It also returns the total number of available results, to help with pagination.
	List(ctx context.Context, query models.ImproveSuggestionsList, limit, offset int) ([]*models.ImproveSuggestion, int64, error)

	// IsCreator returns whether the user is the creator of the improvement suggestion.
	IsCreator(ctx context.Context, userID, postID uuid.UUID) (bool, error)

	GetPreviews(ctx context.Context, ids []uuid.UUID) ([]*models.ImproveSuggestion, error)

	// StorageToModel converts a storage model to a service model.
	StorageToModel(source *improve_suggestion_storage.Model) *models.ImproveSuggestion
}

type serviceImpl struct {
	repository improve_suggestion_storage.Repository
}

// NewService returns a new Service instance.
// To use a mocked one, call NewMockService.
func NewService(repository improve_suggestion_storage.Repository) Service {
	return &serviceImpl{repository: repository}
}

func (service *serviceImpl) Read(ctx context.Context, id uuid.UUID) (*models.ImproveSuggestion, error) {
	storageModel, err := service.repository.Read(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get improve suggestion: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) Create(ctx context.Context, data *models.ImproveSuggestionUpsert, userID, sourceID, id uuid.UUID, now time.Time) (*models.ImproveSuggestion, error) {
	if err := validation.CheckRequire("data", data); err != nil {
		return nil, err
	}
	if err := validation.CheckRequire("title", data.Title); err != nil {
		return nil, err
	}
	if err := validation.CheckRequire("content", data.Content); err != nil {
		return nil, err
	}
	if err := validation.CheckMinMax("title", data.Title, MinTitleLength, MaxTitleLength); err != nil {
		return nil, err
	}
	if err := validation.CheckMinMax("content", data.Content, MinContentLength, MaxContentLength); err != nil {
		return nil, err
	}
	if err := validation.CheckRegexp("title", data.Title, titleRegexp); err != nil {
		return nil, err
	}

	storageModel, err := service.repository.Create(ctx, &improve_suggestion_storage.Core{
		RequestID: data.RequestID,
		Title:     data.Title,
		Content:   data.Content,
	}, userID, sourceID, id, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create improve suggestion: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) Update(ctx context.Context, data *models.ImproveSuggestionUpsert, id uuid.UUID, now time.Time) (*models.ImproveSuggestion, error) {
	if err := validation.CheckRequire("data", data); err != nil {
		return nil, err
	}
	if err := validation.CheckRequire("title", data.Title); err != nil {
		return nil, err
	}
	if err := validation.CheckRequire("content", data.Content); err != nil {
		return nil, err
	}
	if err := validation.CheckMinMax("title", data.Title, MinTitleLength, MaxTitleLength); err != nil {
		return nil, err
	}
	if err := validation.CheckMinMax("content", data.Content, MinContentLength, MaxContentLength); err != nil {
		return nil, err
	}
	if err := validation.CheckRegexp("title", data.Title, titleRegexp); err != nil {
		return nil, err
	}

	storageModel, err := service.repository.Update(ctx, &improve_suggestion_storage.Core{
		RequestID: data.RequestID,
		Title:     data.Title,
		Content:   data.Content,
	}, id, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create improve suggestion: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := service.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete improve suggestion: %w", err)
	}

	return nil
}

func (service *serviceImpl) Validate(ctx context.Context, validated bool, id uuid.UUID) (*models.ImproveSuggestion, error) {
	storageModel, err := service.repository.Validate(ctx, validated, id)
	if err != nil {
		return nil, fmt.Errorf("failed to validate improve suggestion: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) List(ctx context.Context, query models.ImproveSuggestionsList, limit, offset int) ([]*models.ImproveSuggestion, int64, error) {
	storageQuery := improve_suggestion_storage.ListQuery{
		UserID:    query.UserID,
		SourceID:  query.SourceID,
		RequestID: query.RequestID,
		Validated: query.Validated,
	}

	if query.Order != nil {
		storageQuery.Order = &improve_suggestion_storage.SearchQueryOrder{
			Created: query.Order.Created,
			Score:   query.Order.Score,
		}
	}

	storageModels, total, err := service.repository.List(ctx, storageQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list improve suggestions: %w", err)
	}

	serviceModels := make([]*models.ImproveSuggestion, len(storageModels))
	for i, storageModel := range storageModels {
		serviceModels[i] = service.StorageToModel(storageModel)
	}

	return serviceModels, total, nil
}

func (service *serviceImpl) IsCreator(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	ok, err := service.repository.IsCreator(ctx, userID, postID)
	if err != nil {
		return false, fmt.Errorf("failed to check improve suggestions: %w", err)
	}

	return ok, nil
}

func (service *serviceImpl) GetPreviews(ctx context.Context, ids []uuid.UUID) ([]*models.ImproveSuggestion, error) {
	storageModels, err := service.repository.GetPreviews(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get improve suggestion previews: %w", err)
	}

	serviceModels := make([]*models.ImproveSuggestion, len(storageModels))
	for i, storageModel := range storageModels {
		serviceModels[i] = service.StorageToModel(storageModel)
	}

	return serviceModels, nil
}

func (service *serviceImpl) StorageToModel(source *improve_suggestion_storage.Model) *models.ImproveSuggestion {
	if source == nil {
		return nil
	}

	return &models.ImproveSuggestion{
		ID:        source.ID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		SourceID:  source.SourceID,
		UserID:    source.UserID,
		Validated: source.Validated,
		UpVotes:   source.UpVotes,
		DownVotes: source.DownVotes,
		RequestID: source.RequestID,
		Title:     source.Title,
		Content:   source.Content,
	}
}
