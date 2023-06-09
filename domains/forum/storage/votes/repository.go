package votes_storage

import (
	"context"
	"errors"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// Repository of the current layer. You can instantiate a new one with NewRepository.
type Repository interface {
	// Vote creates, update or cancels an existing vote. It returns the new vote value.
	Vote(ctx context.Context, postID, userID uuid.UUID, target Target, vote Vote, now time.Time) (Vote, error)
	// HasVoted returns whether the user has voted for the targeted post. It returns VoteUp or VoteDown if the user
	// has voted, NoVote otherwise.
	HasVoted(ctx context.Context, postID, userID uuid.UUID, target Target) (Vote, error)
	// GetVotedPosts returns the IDs of the posts that the user has voted for, for a specific target. Results must be
	// paginated using the limit and offset parameters.
	// It also returns the total number of available results, to help with pagination.
	GetVotedPosts(ctx context.Context, userID uuid.UUID, target Target, limit, offset int) ([]*VotedPost, int64, error)
}

// NewRepository returns a new Repository instance.
// To use a mocked one, call NewMockRepository.
func NewRepository(db bun.IDB) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db bun.IDB
}

func getSourceTable(target Target) (string, error) {
	switch target {
	case TargetImproveRequest:
		return "improve_requests", nil
	case TargetImproveSuggestion:
		return "improve_suggestions", nil
	default:
		return "", validation.ErrInvalidEntity
	}
}

func (repository *repositoryImpl) Vote(ctx context.Context, postID, userID uuid.UUID, target Target, vote Vote, now time.Time) (Vote, error) {
	if vote == NoVote {
		if _, err := repository.db.NewDelete().
			Model(new(Model)).
			Where("post_id = ? AND user_id = ? AND target = ?", postID, userID, target).
			Exec(ctx); err != nil {
			return NoVote, validation.HandlePGError(err)
		}

		return NoVote, nil
	}

	model := &Model{
		UpdatedAt: now,
		PostID:    postID,
		UserID:    userID,
		Target:    target,
		Vote:      vote,
	}

	sourceTable, err := getSourceTable(target)
	if err != nil {
		return NoVote, err
	}

	count, err := repository.db.NewSelect().Table(sourceTable).Where("id = ?", postID).Count(ctx)
	if err != nil {
		return NoVote, validation.HandlePGError(err)
	}
	if count == 0 {
		return NoVote, validation.ErrMissingRelation
	}

	if _, err := repository.db.NewInsert().
		Model(model).
		On("conflict (post_id, user_id, target) do update").
		Set("updated_at = ?updated_at").
		Set("vote = ?vote").
		Exec(ctx); err != nil {
		return NoVote, validation.HandlePGError(err)
	}

	return vote, nil
}

func (repository *repositoryImpl) HasVoted(ctx context.Context, postID, userID uuid.UUID, target Target) (Vote, error) {
	model := new(Model)
	if err := repository.db.NewSelect().
		Model(model).
		Column("vote").
		Where("post_id = ? AND user_id = ? AND target = ?", postID, userID, target).
		Scan(ctx); err != nil {
		err = validation.HandlePGError(err)
		if errors.Is(err, validation.ErrNotFound) {
			err = nil
		}

		return NoVote, err
	}

	return model.Vote, nil
}

func (repository *repositoryImpl) GetVotedPosts(ctx context.Context, userID uuid.UUID, target Target, limit, offset int) ([]*VotedPost, int64, error) {
	var models []*VotedPost
	count, err := repository.db.NewSelect().
		Model(&models).
		Column("post_id", "updated_at", "vote").
		Where("user_id = ? AND target = ?", userID, target).
		OrderExpr("updated_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, validation.HandlePGError(err)
	}

	return models, int64(count), nil
}
