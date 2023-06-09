package improve_post

import (
	"context"
	"fmt"
	"github.com/a-novel/agora-backend/domains/forum/service/improve_request"
	"github.com/a-novel/agora-backend/domains/forum/service/improve_suggestion"
	"github.com/a-novel/agora-backend/domains/forum/service/votes"
	"github.com/a-novel/agora-backend/domains/keys/service/jwk"
	"github.com/a-novel/agora-backend/domains/user/service/token"
	user_service "github.com/a-novel/agora-backend/domains/user/service/user"
	"github.com/a-novel/agora-backend/environment/user/authentication"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"time"
)

type Provider interface {
	ReadImproveRequest(ctx context.Context, id uuid.UUID) ([]*models.ImproveRequest, error)
	ReadImproveSuggestion(ctx context.Context, id uuid.UUID) (*models.ImproveSuggestion, error)

	CreateImproveRequest(ctx context.Context, token, title, content string) (*models.ImproveRequest, error)
	CreateImproveRequestRevision(ctx context.Context, token string, sourceID uuid.UUID, title, content string) (*models.ImproveRequest, error)
	CreateImproveSuggestion(ctx context.Context, token string, requestID, sourceID uuid.UUID, title, content string) (*models.ImproveSuggestion, error)
	UpdateImproveSuggestion(ctx context.Context, token string, postID, requestID uuid.UUID, title, content string) (*models.ImproveSuggestion, error)

	DeleteImproveRequest(ctx context.Context, token string, requestID uuid.UUID) error
	DeleteImproveSuggestion(ctx context.Context, token string, id uuid.UUID) error

	ListImproveSuggestions(ctx context.Context, query models.ImproveSuggestionsList, limit, offset int) ([]*models.ImproveSuggestion, int64, error)
	SearchImproveRequests(ctx context.Context, query models.ImproveRequestSearch, limit, offset int) ([]*models.ImproveRequestPreview, int64, error)

	Vote(ctx context.Context, token string, postID uuid.UUID, target models.VoteTarget, vote models.VoteValue) (models.VoteValue, error)
	HasVoted(ctx context.Context, token string, postID uuid.UUID, target models.VoteTarget) (models.VoteValue, error)
	GetVotedPosts(ctx context.Context, userID uuid.UUID, target models.VoteTarget, limit, offset int) ([]*models.VotedPost, int64, error)

	GetImproveRequestPreviews(ctx context.Context, ids []uuid.UUID) ([]*models.ImproveRequestPreview, error)
	GetImproveSuggestionPreviews(ctx context.Context, ids []uuid.UUID) ([]*models.ImproveSuggestion, error)
}

type Config struct {
	ImproveRequestService    improve_request_service.Service
	ImproveSuggestionService improve_suggestion_service.Service
	VotesService             votes_service.Service
	TokenService             token_service.Service
	KeysService              jwk_service.ServiceCached
	UserService              user_service.Service

	Time func() time.Time
	ID   func() uuid.UUID
}

type providerImpl struct {
	improveRequestService    improve_request_service.Service
	improveSuggestionService improve_suggestion_service.Service
	votesService             votes_service.Service
	tokenService             token_service.Service
	keysService              jwk_service.ServiceCached
	userService              user_service.Service

	time func() time.Time
	id   func() uuid.UUID
}

func NewProvider(config Config) Provider {
	return &providerImpl{
		improveRequestService:    config.ImproveRequestService,
		improveSuggestionService: config.ImproveSuggestionService,
		votesService:             config.VotesService,
		tokenService:             config.TokenService,
		keysService:              config.KeysService,
		userService:              config.UserService,

		time: config.Time,
		id:   config.ID,
	}
}

func (provider *providerImpl) ReadImproveRequest(ctx context.Context, id uuid.UUID) ([]*models.ImproveRequest, error) {
	revisions, err := provider.improveRequestService.ReadRevisions(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch revisions for improve request %q: %w", id, err)
	}

	return revisions, nil
}

func (provider *providerImpl) CreateImproveRequest(ctx context.Context, token, title, content string) (*models.ImproveRequest, error) {
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

	request, err := provider.improveRequestService.Create(ctx, claims.Payload.ID, title, content, provider.id(), now)
	if err != nil {
		return nil, fmt.Errorf("failed to create improve request %q, for user %q: %w", title, claims.Payload.ID, err)
	}

	return request, nil
}

func (provider *providerImpl) CreateImproveRequestRevision(ctx context.Context, token string, sourceID uuid.UUID, title, content string) (*models.ImproveRequest, error) {
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

	// Only creator is allowed to edit its own post.
	source, err := provider.improveRequestService.Read(ctx, sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source improve request %q: %w", sourceID, err)
	}

	if source.UserID != claims.Payload.ID {
		return nil, fmt.Errorf(
			"%w: user %q is not allowed to create a revision for the post %q (created by %q)",
			validation.ErrInvalidCredentials, claims.Payload.ID, sourceID, source.UserID,
		)
	}

	request, err := provider.improveRequestService.CreateRevision(ctx, claims.Payload.ID, sourceID, title, content, provider.id(), now)
	if err != nil {
		return nil, fmt.Errorf("failed to create revision on improve request %q for user %q: %w", source.Title, claims.Payload.ID, err)
	}

	return request, nil
}

func (provider *providerImpl) DeleteImproveRequest(ctx context.Context, token string, requestID uuid.UUID) error {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return err
	}

	// Force revision to be from the same user.
	source, err := provider.improveRequestService.Read(ctx, requestID)
	if err != nil {
		return fmt.Errorf("failed to fetch improve request %q: %w", requestID, err)
	}

	if source.UserID != claims.Payload.ID {
		return fmt.Errorf(
			"%w: user %q is not allowed to delete the post %q (created by %q)", validation.ErrInvalidCredentials,
			claims.Payload.ID, requestID, source.UserID,
		)
	}

	if err := provider.improveRequestService.Delete(ctx, requestID); err != nil {
		return fmt.Errorf("failed to delete improve request %q: %w", requestID, err)
	}

	return nil
}

func (provider *providerImpl) SearchImproveRequests(ctx context.Context, query models.ImproveRequestSearch, limit, offset int) ([]*models.ImproveRequestPreview, int64, error) {
	requests, total, err := provider.improveRequestService.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search improve requests: %w", err)
	}

	return requests, total, nil
}

func (provider *providerImpl) GetImproveRequestPreviews(ctx context.Context, ids []uuid.UUID) ([]*models.ImproveRequestPreview, error) {
	requests, err := provider.improveRequestService.GetPreviews(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get improve request previews: %w", err)
	}

	return requests, nil
}

func (provider *providerImpl) ReadImproveSuggestion(ctx context.Context, id uuid.UUID) (*models.ImproveSuggestion, error) {
	suggestion, err := provider.improveSuggestionService.Read(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read improve suggestion %q: %w", id, err)
	}

	return suggestion, nil
}

func (provider *providerImpl) CreateImproveSuggestion(ctx context.Context, token string, requestID, sourceID uuid.UUID, title, content string) (*models.ImproveSuggestion, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return nil, err
	}

	suggestion, err := provider.improveSuggestionService.Create(
		ctx, &models.ImproveSuggestionUpsert{
			RequestID: requestID,
			Title:     title,
			Content:   content,
		}, claims.Payload.ID, sourceID, provider.id(), now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create improve suggestion %q, on improve request %q: %w", title, requestID, err)
	}

	return suggestion, nil
}

func (provider *providerImpl) UpdateImproveSuggestion(ctx context.Context, token string, postID, requestID uuid.UUID, title, content string) (*models.ImproveSuggestion, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return nil, err
	}

	// Force suggestion to be from the same user.
	source, err := provider.improveSuggestionService.Read(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source improve suggestion %q: %w", postID, err)
	}

	if source.UserID != claims.Payload.ID {
		return nil, fmt.Errorf(
			"%w: user %q is not allowed to update improve suggestion %q (created by %q)",
			validation.ErrInvalidCredentials, claims.Payload.ID, postID, source.UserID,
		)
	}

	suggestion, err := provider.improveSuggestionService.Update(
		ctx, &models.ImproveSuggestionUpsert{
			RequestID: requestID,
			Title:     title,
			Content:   content,
		}, postID, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update improve suggestion %q for user %q: %w", source.ID, claims.Payload.ID, err)
	}

	return suggestion, nil
}

func (provider *providerImpl) DeleteImproveSuggestion(ctx context.Context, token string, id uuid.UUID) error {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return err
	}

	// Force suggestion to be from the same user.
	source, err := provider.improveSuggestionService.Read(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to fetch improve suggestion %q: %w", id, err)
	}

	if source.UserID != claims.Payload.ID {
		return fmt.Errorf(
			"%w: user %q is not allowed to delete the post %q (created by %q)", validation.ErrInvalidCredentials,
			claims.Payload.ID, id, source.UserID,
		)
	}

	if err := provider.improveSuggestionService.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete improve suggestion %q: %w", id, err)
	}

	return nil
}

func (provider *providerImpl) ListImproveSuggestions(ctx context.Context, query models.ImproveSuggestionsList, limit, offset int) ([]*models.ImproveSuggestion, int64, error) {
	suggestions, total, err := provider.improveSuggestionService.List(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list improve suggestions: %w", err)
	}

	return suggestions, total, nil
}

func (provider *providerImpl) GetImproveSuggestionPreviews(ctx context.Context, ids []uuid.UUID) ([]*models.ImproveSuggestion, error) {
	suggestions, err := provider.improveSuggestionService.GetPreviews(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get improve suggestions: %w", err)
	}

	return suggestions, nil
}

func (provider *providerImpl) Vote(ctx context.Context, token string, postID uuid.UUID, target models.VoteTarget, vote models.VoteValue) (models.VoteValue, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return models.NoVote, err
	}

	// Cannot vote own post.
	switch target {
	case models.VoteTargetImproveSuggestion:
		isCreator, err := provider.improveSuggestionService.IsCreator(ctx, claims.Payload.ID, postID)
		if err != nil {
			return models.NoVote, fmt.Errorf(
				"failed to check if user %q is the creator of improve suggestion %q: %w",
				claims.Payload.ID, postID, err,
			)
		}

		if isCreator {
			return models.NoVote, fmt.Errorf(
				"%w: user %q cannot vote on its own improve suggestion %q",
				validation.ErrInvalidEntity, claims.Payload.ID, postID,
			)
		}
	case models.VoteTargetImproveRequest:
		// Cannot vote an improvement request, if you are the creator or one of the editors.
		isCreator, err := provider.improveRequestService.IsCreator(ctx, claims.Payload.ID, postID, false)
		if err != nil {
			return models.NoVote, fmt.Errorf(
				"failed to check if user %q is a creator of improve request %q: %w",
				claims.Payload.ID, postID, err,
			)
		}

		if isCreator {
			return models.NoVote, fmt.Errorf(
				"%w: user %q cannot vote on its own improve request %q",
				validation.ErrInvalidEntity, claims.Payload.ID, postID,
			)
		}
	}

	res, err := provider.votesService.Vote(
		ctx, postID, claims.Payload.ID, target, vote, now,
	)
	if err != nil {
		return models.NoVote, fmt.Errorf("failed to vote for %q %q, for user %q: %w", target, postID, claims.Payload.ID, err)
	}

	return res, nil
}

func (provider *providerImpl) HasVoted(ctx context.Context, token string, postID uuid.UUID, target models.VoteTarget) (models.VoteValue, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return models.NoVote, err
	}

	res, err := provider.votesService.HasVoted(ctx, postID, claims.Payload.ID, target)
	if err != nil {
		return models.NoVote, fmt.Errorf("failed to check vote for %q %q, for user %q: %w", target, postID, claims.Payload.ID, err)
	}

	return res, nil
}

func (provider *providerImpl) GetVotedPosts(ctx context.Context, userID uuid.UUID, target models.VoteTarget, limit, offset int) ([]*models.VotedPost, int64, error) {
	posts, total, err := provider.votesService.GetVotedPosts(ctx, userID, target, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get voted posts for user %q: %w", userID, err)
	}

	return posts, total, nil
}
