package improve_request_service

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
	MaxTitleLength   = 128
	MinContentLength = 4
	MaxContentLength = 4096
)

// Service of the current layer. You can instantiate a new one with NewService.
type Service interface {
	// Read reads a single post revision, based on its ID.
	Read(ctx context.Context, id uuid.UUID) (*models.ImproveRequest, error)
	// ReadRevisions reads every revision, related to a source. The ID must be the one of the source post.
	ReadRevisions(ctx context.Context, id uuid.UUID) ([]*models.ImproveRequest, error)

	// Create creates a brand-new post. The returned model will have matching ImproveRequest.Source and ImproveRequest.ID.
	Create(ctx context.Context, userID uuid.UUID, title, content string, id uuid.UUID, now time.Time) (*models.ImproveRequest, error)
	// CreateRevision creates a new revision for a given post. The ID must be the one of the source post.
	CreateRevision(ctx context.Context, userID, sourceID uuid.UUID, title, content string, id uuid.UUID, now time.Time) (*models.ImproveRequest, error)

	// Delete a single revision for a post. If the provided id is the source id, then all associated revisions will
	// also be deleted.
	Delete(ctx context.Context, requestID uuid.UUID) error

	// Search returns a list of posts, matching the provided query. Results must be paginated using the limit and
	// offset parameters.
	// It also returns the total number of available results, to help with pagination.
	Search(ctx context.Context, query models.ImproveRequestSearch, limit, offset int) ([]*models.ImproveRequestPreview, int64, error)

	// IsCreator returns whether the user is a creator of the improvement suggestion. ID can be the id of any revision.
	// To only check if the user is the creator of the specific revision, set strict flag to true.
	IsCreator(ctx context.Context, userID, postID uuid.UUID, strict bool) (bool, error)

	GetPreviews(ctx context.Context, ids []uuid.UUID) ([]*models.ImproveRequestPreview, error)

	// StorageToModel converts a storage model to a service model.
	StorageToModel(source *improve_request_storage.Model) *models.ImproveRequest
}

type serviceImpl struct {
	repository improve_request_storage.Repository
}

// NewService returns a new Service instance.
// To use a mocked one, call NewMockService.
func NewService(repository improve_request_storage.Repository) Service {
	return &serviceImpl{repository: repository}
}

func (service *serviceImpl) Read(ctx context.Context, id uuid.UUID) (*models.ImproveRequest, error) {
	storageModel, err := service.repository.Read(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get improve request: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) ReadRevisions(ctx context.Context, id uuid.UUID) ([]*models.ImproveRequest, error) {
	storageModels, err := service.repository.ReadRevisions(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get improve request revisions: %w", err)
	}

	serviceModels := make([]*models.ImproveRequest, len(storageModels))
	for i, storageModel := range storageModels {
		serviceModels[i] = service.StorageToModel(storageModel)
	}

	return serviceModels, nil
}

func (service *serviceImpl) Create(ctx context.Context, userID uuid.UUID, title, content string, id uuid.UUID, now time.Time) (*models.ImproveRequest, error) {
	if err := validation.CheckRequire("title", title); err != nil {
		return nil, err
	}
	if err := validation.CheckRequire("content", content); err != nil {
		return nil, err
	}
	if err := validation.CheckMinMax("title", title, MinTitleLength, MaxTitleLength); err != nil {
		return nil, err
	}
	if err := validation.CheckMinMax("content", content, MinContentLength, MaxContentLength); err != nil {
		return nil, err
	}
	if err := validation.CheckRegexp("title", title, titleRegexp); err != nil {
		return nil, err
	}

	storageModel, err := service.repository.Create(ctx, userID, title, content, id, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create improve request: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) CreateRevision(ctx context.Context, userID, sourceID uuid.UUID, title, content string, id uuid.UUID, now time.Time) (*models.ImproveRequest, error) {
	if err := validation.CheckRequire("title", title); err != nil {
		return nil, err
	}
	if err := validation.CheckRequire("content", content); err != nil {
		return nil, err
	}
	if err := validation.CheckMinMax("title", title, MinTitleLength, MaxTitleLength); err != nil {
		return nil, err
	}
	if err := validation.CheckMinMax("content", content, MinContentLength, MaxContentLength); err != nil {
		return nil, err
	}
	if err := validation.CheckRegexp("title", title, titleRegexp); err != nil {
		return nil, err
	}

	storageModel, err := service.repository.CreateRevision(ctx, userID, sourceID, title, content, id, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create improve request: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) Delete(ctx context.Context, requestID uuid.UUID) error {
	if err := service.repository.Delete(ctx, requestID); err != nil {
		return fmt.Errorf("failed to delete improve request: %w", err)
	}

	return nil
}

func (service *serviceImpl) Search(ctx context.Context, query models.ImproveRequestSearch, limit, offset int) ([]*models.ImproveRequestPreview, int64, error) {
	storageQuery := improve_request_storage.SearchQuery{
		UserID: query.UserID,
		Query:  query.Query,
	}

	if query.Order != nil {
		storageQuery.Order = &improve_request_storage.SearchQueryOrder{
			Created: query.Order.Created,
			Score:   query.Order.Score,
		}
	}

	storageModels, total, err := service.repository.Search(ctx, storageQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search improve requests: %w", err)
	}

	serviceModels := make([]*models.ImproveRequestPreview, len(storageModels))
	for i, storageModel := range storageModels {
		serviceModels[i] = &models.ImproveRequestPreview{
			ID:                       storageModel.ID,
			Source:                   storageModel.Source,
			CreatedAt:                storageModel.CreatedAt,
			UserID:                   storageModel.UserID,
			Title:                    storageModel.Title,
			Content:                  storageModel.Content,
			UpVotes:                  storageModel.UpVotes,
			DownVotes:                storageModel.DownVotes,
			RevisionCount:            storageModel.RevisionCount,
			MoreRecentRevisions:      storageModel.MoreRecentRevisions,
			SuggestionsCount:         storageModel.SuggestionsCount,
			AcceptedSuggestionsCount: storageModel.AcceptedSuggestionsCount,
		}
	}

	return serviceModels, total, nil
}

func (service *serviceImpl) IsCreator(ctx context.Context, userID, postID uuid.UUID, strict bool) (bool, error) {
	ok, err := service.repository.IsCreator(ctx, userID, postID, strict)
	if err != nil {
		return false, fmt.Errorf("failed to check improve requests: %w", err)
	}

	return ok, nil
}

func (service *serviceImpl) GetPreviews(ctx context.Context, ids []uuid.UUID) ([]*models.ImproveRequestPreview, error) {
	storageModels, err := service.repository.GetPreviews(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get improve request previews: %w", err)
	}

	serviceModels := make([]*models.ImproveRequestPreview, len(storageModels))
	for i, storageModel := range storageModels {
		serviceModels[i] = &models.ImproveRequestPreview{
			ID:                       storageModel.ID,
			Source:                   storageModel.Source,
			CreatedAt:                storageModel.CreatedAt,
			UserID:                   storageModel.UserID,
			Title:                    storageModel.Title,
			Content:                  storageModel.Content,
			UpVotes:                  storageModel.UpVotes,
			DownVotes:                storageModel.DownVotes,
			RevisionCount:            storageModel.RevisionCount,
			MoreRecentRevisions:      storageModel.MoreRecentRevisions,
			SuggestionsCount:         storageModel.SuggestionsCount,
			AcceptedSuggestionsCount: storageModel.AcceptedSuggestionsCount,
		}
	}

	return serviceModels, nil
}

func (service *serviceImpl) StorageToModel(source *improve_request_storage.Model) *models.ImproveRequest {
	if source == nil {
		return nil
	}

	return &models.ImproveRequest{
		ID:        source.ID,
		Source:    source.Source,
		CreatedAt: source.CreatedAt,
		UserID:    source.UserID,
		Title:     source.Title,
		Content:   source.Content,
		UpVotes:   source.UpVotes,
		DownVotes: source.DownVotes,
	}
}
