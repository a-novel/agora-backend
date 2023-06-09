package votes_service

import (
	"context"
	"fmt"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"time"
)

var (
	voteValues   = []models.VoteValue{models.VoteUp, models.VoteDown, models.NoVote}
	targetValues = []models.VoteTarget{models.VoteTargetImproveRequest, models.VoteTargetImproveSuggestion}
)

// Service of the current layer. You can instantiate a new one with NewService.
type Service interface {
	// Vote creates, update or cancels an existing vote. It returns the new vote value.
	Vote(ctx context.Context, postID, userID uuid.UUID, target models.VoteTarget, vote models.VoteValue, now time.Time) (models.VoteValue, error)
	// HasVoted returns whether the user has voted for the targeted post. It returns VoteUp or VoteDown if the user
	// has voted, NoVote otherwise.
	HasVoted(ctx context.Context, postID, userID uuid.UUID, target models.VoteTarget) (models.VoteValue, error)
	// GetVotedPosts returns the IDs of the posts that the user has voted for, for a specific target. Results must be
	// paginated using the limit and offset parameters.
	// It also returns the total number of available results, to help with pagination.
	GetVotedPosts(ctx context.Context, userID uuid.UUID, target models.VoteTarget, limit, offset int) ([]*models.VotedPost, int64, error)

	// StorageToModel converts a storage model to a service model.
	StorageToModel(source *votes_storage.Model) *models.Vote
}

type serviceImpl struct {
	repository votes_storage.Repository
}

// NewService returns a new Service instance.
// To use a mocked one, call NewMockService.
func NewService(repository votes_storage.Repository) Service {
	return &serviceImpl{repository: repository}
}

func (serviceImpl *serviceImpl) Vote(ctx context.Context, postID, userID uuid.UUID, target models.VoteTarget, vote models.VoteValue, now time.Time) (models.VoteValue, error) {
	if err := validation.CheckRestricted("vote", vote, voteValues...); err != nil {
		return models.NoVote, err
	}
	if err := validation.CheckRestricted("target", target, targetValues...); err != nil {
		return models.NoVote, err
	}

	storageModel, err := serviceImpl.repository.Vote(
		ctx, postID, userID, votes_storage.Target(target), votes_storage.Vote(vote), now,
	)
	if err != nil {
		return models.NoVote, fmt.Errorf("failed to cast vote: %w", err)
	}

	return models.VoteValue(storageModel), nil
}

func (serviceImpl *serviceImpl) HasVoted(ctx context.Context, postID, userID uuid.UUID, target models.VoteTarget) (models.VoteValue, error) {
	if err := validation.CheckRestricted("target", target, targetValues...); err != nil {
		return models.NoVote, err
	}

	storageModel, err := serviceImpl.repository.HasVoted(ctx, postID, userID, votes_storage.Target(target))
	if err != nil {
		return models.NoVote, fmt.Errorf("failed to check vote status: %w", err)
	}

	return models.VoteValue(storageModel), nil
}

func (serviceImpl *serviceImpl) GetVotedPosts(ctx context.Context, userID uuid.UUID, target models.VoteTarget, limit, offset int) ([]*models.VotedPost, int64, error) {
	if err := validation.CheckRestricted("target", target, targetValues...); err != nil {
		return nil, 0, err
	}

	storageModels, total, err := serviceImpl.repository.GetVotedPosts(ctx, userID, votes_storage.Target(target), limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get voted posts: %w", err)
	}

	posts := make([]*models.VotedPost, len(storageModels))
	for i, storageModel := range storageModels {
		posts[i] = &models.VotedPost{
			PostID:    storageModel.PostID,
			UpdatedAt: storageModel.UpdatedAt,
			Vote:      models.VoteValue(storageModel.Vote),
		}
	}

	return posts, total, nil
}

func (serviceImpl *serviceImpl) StorageToModel(source *votes_storage.Model) *models.Vote {
	if source == nil {
		return nil
	}

	return &models.Vote{
		UpdatedAt: source.UpdatedAt,
		PostID:    source.PostID,
		UserID:    source.UserID,
		Target:    models.VoteTarget(source.Target),
		Vote:      models.VoteValue(source.Vote),
	}
}
