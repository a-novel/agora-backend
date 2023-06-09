package votes_service

import (
	"context"
	"errors"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	baseTime = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	fooErr   = errors.New("it broken")
)

func TestVotesService_Vote(t *testing.T) {
	data := []struct {
		name string

		postID uuid.UUID
		userID uuid.UUID
		target models.VoteTarget
		vote   models.VoteValue
		now    time.Time

		shouldCallRepository bool
		voteData             votes_storage.Vote
		voteErr              error

		expect    models.VoteValue
		expectErr error
	}{
		{
			name:                 "Success",
			postID:               test_utils.NumberUUID(1),
			userID:               test_utils.NumberUUID(2),
			target:               models.VoteTargetImproveRequest,
			vote:                 models.VoteUp,
			now:                  baseTime,
			shouldCallRepository: true,
			voteData:             votes_storage.VoteUp,
			expect:               models.VoteUp,
		},
		{
			name:      "Error/InvalidVote",
			postID:    test_utils.NumberUUID(1),
			userID:    test_utils.NumberUUID(2),
			target:    models.VoteTargetImproveRequest,
			vote:      models.VoteValue("foo"),
			now:       baseTime,
			expectErr: validation.ErrNotAllowed,
		},
		{
			name:      "Error/InvalidTarget",
			postID:    test_utils.NumberUUID(1),
			userID:    test_utils.NumberUUID(2),
			target:    models.VoteTarget("foo"),
			vote:      models.VoteUp,
			now:       baseTime,
			expectErr: validation.ErrNotAllowed,
		},
		{
			name:                 "Error/RepositoryFailure",
			postID:               test_utils.NumberUUID(1),
			userID:               test_utils.NumberUUID(2),
			target:               models.VoteTargetImproveRequest,
			vote:                 models.VoteUp,
			now:                  baseTime,
			shouldCallRepository: true,
			voteErr:              fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := votes_storage.NewMockRepository(st)
			if d.shouldCallRepository {
				repository.
					On("Vote", context.TODO(), d.postID, d.userID, votes_storage.Target(d.target), votes_storage.Vote(d.vote), d.now).
					Return(d.voteData, d.voteErr)
			}

			service := NewService(repository)

			res, err := service.Vote(context.TODO(), d.postID, d.userID, d.target, d.vote, d.now)
			test_utils.RequireError(st, d.expectErr, err)
			require.Equal(st, d.expect, res)

			require.True(st, repository.AssertExpectations(st))
		})
	}
}

func TestVotesService_HasVoted(t *testing.T) {
	data := []struct {
		name string

		postID uuid.UUID
		userID uuid.UUID
		target models.VoteTarget

		shouldCallRepository bool
		voteData             votes_storage.Vote
		voteErr              error

		expect    models.VoteValue
		expectErr error
	}{
		{
			name:                 "Success",
			postID:               test_utils.NumberUUID(1),
			userID:               test_utils.NumberUUID(2),
			target:               models.VoteTargetImproveRequest,
			shouldCallRepository: true,
			voteData:             votes_storage.VoteUp,
			expect:               models.VoteUp,
		},
		{
			name:      "Error/InvalidTarget",
			postID:    test_utils.NumberUUID(1),
			userID:    test_utils.NumberUUID(2),
			target:    models.VoteTarget("foo"),
			expectErr: validation.ErrNotAllowed,
		},
		{
			name:                 "Error/RepositoryFailure",
			postID:               test_utils.NumberUUID(1),
			userID:               test_utils.NumberUUID(2),
			target:               models.VoteTargetImproveRequest,
			shouldCallRepository: true,
			voteErr:              fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := votes_storage.NewMockRepository(st)
			if d.shouldCallRepository {
				repository.
					On("HasVoted", context.TODO(), d.postID, d.userID, votes_storage.Target(d.target)).
					Return(d.voteData, d.voteErr)
			}

			service := NewService(repository)

			res, err := service.HasVoted(context.TODO(), d.postID, d.userID, d.target)
			test_utils.RequireError(st, d.expectErr, err)
			require.Equal(st, d.expect, res)

			require.True(st, repository.AssertExpectations(st))
		})
	}
}

func TestVotesService_GetVotedPosts(t *testing.T) {
	data := []struct {
		name string

		userID uuid.UUID
		target models.VoteTarget
		limit  int
		offset int

		shouldCallRepository bool
		votedPostsData       []*votes_storage.VotedPost
		votedPostsCount      int64
		votedPostsErr        error

		expect      []*models.VotedPost
		expectCount int64
		expectErr   error
	}{
		{
			name:                 "Success",
			userID:               test_utils.NumberUUID(1),
			target:               models.VoteTargetImproveRequest,
			limit:                10,
			offset:               20,
			shouldCallRepository: true,
			votedPostsData: []*votes_storage.VotedPost{
				{
					PostID:    test_utils.NumberUUID(1),
					UpdatedAt: baseTime,
					Vote:      votes_storage.VoteUp,
				},
				{
					PostID:    test_utils.NumberUUID(2),
					UpdatedAt: baseTime,
					Vote:      votes_storage.VoteDown,
				},
			},
			votedPostsCount: 123,
			expect: []*models.VotedPost{
				{
					PostID:    test_utils.NumberUUID(1),
					UpdatedAt: baseTime,
					Vote:      models.VoteUp,
				},
				{
					PostID:    test_utils.NumberUUID(2),
					UpdatedAt: baseTime,
					Vote:      models.VoteDown,
				},
			},
			expectCount: 123,
		},
		{
			name:                 "Success/NoResults",
			userID:               test_utils.NumberUUID(1),
			target:               models.VoteTargetImproveRequest,
			limit:                10,
			offset:               20,
			shouldCallRepository: true,
			votedPostsCount:      123,
			expect:               []*models.VotedPost{},
			expectCount:          123,
		},
		{
			name:      "Error/InvalidTarget",
			userID:    test_utils.NumberUUID(1),
			target:    models.VoteTarget("foo"),
			limit:     10,
			offset:    20,
			expectErr: validation.ErrNotAllowed,
		},
		{
			name:                 "Error/RepositoryFailure",
			userID:               test_utils.NumberUUID(1),
			target:               models.VoteTargetImproveRequest,
			limit:                10,
			offset:               20,
			shouldCallRepository: true,
			votedPostsErr:        fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := votes_storage.NewMockRepository(st)
			if d.shouldCallRepository {
				repository.
					On("GetVotedPosts", context.TODO(), d.userID, votes_storage.Target(d.target), d.limit, d.offset).
					Return(d.votedPostsData, d.votedPostsCount, d.votedPostsErr)
			}

			service := NewService(repository)

			res, count, err := service.GetVotedPosts(context.TODO(), d.userID, d.target, d.limit, d.offset)
			test_utils.RequireError(st, d.expectErr, err)
			require.Equal(st, d.expect, res)
			require.Equal(st, d.expectCount, count)

			require.True(st, repository.AssertExpectations(st))
		})
	}
}
