package votes_storage

import (
	"context"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"testing"
	"time"
)

var (
	baseTime   = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	updateTime = time.Date(2020, time.May, 4, 8, 10, 0, 0, time.UTC)
)

var Fixtures = []interface{}{
	// Requests.
	&improve_request_storage.Model{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		Source:    test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(200),
		UpVotes:   10,
		DownVotes: 3, // 4
		Title:     "Test",
		Content:   "Dummy content.",
	},
	// Suggestions.
	&improve_suggestion_storage.Model{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		UpdatedAt: &updateTime,
		SourceID:  test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(201),
		UpVotes:   4, // 5
		DownVotes: 1,
		Core: improve_suggestion_storage.Core{
			RequestID: test_utils.NumberUUID(1000),
			Title:     "Ipsum Lorem",
			Content:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a.",
		},
	},
	// Votes.
	&Model{
		UpdatedAt: baseTime,
		PostID:    test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(210),
		Target:    TargetImproveSuggestion,
		Vote:      VoteUp,
	},
	&Model{
		UpdatedAt: baseTime,
		PostID:    test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(211),
		Target:    TargetImproveRequest,
		Vote:      VoteDown,
	},
}

// Generates votes for a given post.
func generateVotesFor(src interface{}, upVotes map[int]time.Time, downVotes map[int]time.Time) []interface{} {
	var (
		target Target
		id     uuid.UUID
		votes  []interface{}
	)

	if req, ok := src.(*improve_request_storage.Model); ok {
		target = TargetImproveRequest
		id = req.ID
	} else if sug, ok := src.(*improve_suggestion_storage.Model); ok {
		target = TargetImproveSuggestion
		id = sug.ID
	} else {
		panic("invalid source")
	}

	for userID, updatedAt := range upVotes {
		votes = append(votes, &Model{
			UpdatedAt: updatedAt,
			PostID:    id,
			UserID:    test_utils.NumberUUID(userID),
			Target:    target,
			Vote:      VoteUp,
		})
	}

	for userID, updatedAt := range downVotes {
		votes = append(votes, &Model{
			UpdatedAt: updatedAt,
			PostID:    id,
			UserID:    test_utils.NumberUUID(userID),
			Target:    target,
			Vote:      VoteDown,
		})
	}

	return votes
}

func TestVotesRepository_Vote(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	type checkTarget struct {
		postID    uuid.UUID
		target    Target
		upVotes   int64
		downVotes int64
	}

	data := []struct {
		name string

		postID uuid.UUID
		userID uuid.UUID
		target Target
		vote   Vote
		now    time.Time

		expect    Vote
		expectErr error

		expectScores []checkTarget
	}{
		{
			name:   "Success/UpVote",
			postID: test_utils.NumberUUID(1000),
			userID: test_utils.NumberUUID(212),
			target: TargetImproveRequest,
			vote:   VoteUp,
			now:    baseTime.Add(1 * time.Minute),
			expect: VoteUp,
			expectScores: []checkTarget{
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveRequest,
					upVotes:   11,
					downVotes: 4,
				},
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveSuggestion,
					upVotes:   5,
					downVotes: 1,
				},
			},
		},
		{
			name:   "Success/DownVote",
			postID: test_utils.NumberUUID(1000),
			userID: test_utils.NumberUUID(212),
			target: TargetImproveRequest,
			vote:   VoteDown,
			now:    baseTime.Add(1 * time.Minute),
			expect: VoteDown,
			expectScores: []checkTarget{
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveRequest,
					upVotes:   10,
					downVotes: 5,
				},
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveSuggestion,
					upVotes:   5,
					downVotes: 1,
				},
			},
		},
		{
			name:   "Success/TargetImproveSuggestion",
			postID: test_utils.NumberUUID(1000),
			userID: test_utils.NumberUUID(212),
			target: TargetImproveSuggestion,
			vote:   VoteUp,
			now:    baseTime.Add(1 * time.Minute),
			expect: VoteUp,
			expectScores: []checkTarget{
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveRequest,
					upVotes:   10,
					downVotes: 4,
				},
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveSuggestion,
					upVotes:   6,
					downVotes: 1,
				},
			},
		},
		{
			name:   "Success/DownToUpVote",
			postID: test_utils.NumberUUID(1000),
			userID: test_utils.NumberUUID(211),
			target: TargetImproveRequest,
			vote:   VoteUp,
			now:    baseTime.Add(1 * time.Minute),
			expect: VoteUp,
			expectScores: []checkTarget{
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveRequest,
					upVotes:   11,
					downVotes: 3,
				},
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveSuggestion,
					upVotes:   5,
					downVotes: 1,
				},
			},
		},
		{
			name:   "Success/UpToDownVote",
			postID: test_utils.NumberUUID(1000),
			userID: test_utils.NumberUUID(210),
			target: TargetImproveSuggestion,
			vote:   VoteDown,
			now:    baseTime.Add(1 * time.Minute),
			expect: VoteDown,
			expectScores: []checkTarget{
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveRequest,
					upVotes:   10,
					downVotes: 4,
				},
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveSuggestion,
					upVotes:   4,
					downVotes: 2,
				},
			},
		},
		{
			name:   "Success/RemoveDownVote",
			postID: test_utils.NumberUUID(1000),
			userID: test_utils.NumberUUID(211),
			target: TargetImproveRequest,
			vote:   NoVote,
			now:    baseTime.Add(1 * time.Minute),
			expect: NoVote,
			expectScores: []checkTarget{
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveRequest,
					upVotes:   10,
					downVotes: 3,
				},
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveSuggestion,
					upVotes:   5,
					downVotes: 1,
				},
			},
		},
		{
			name:   "Success/RemoveUpVote",
			postID: test_utils.NumberUUID(1000),
			userID: test_utils.NumberUUID(210),
			target: TargetImproveSuggestion,
			vote:   NoVote,
			now:    baseTime.Add(1 * time.Minute),
			expect: NoVote,
			expectScores: []checkTarget{
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveRequest,
					upVotes:   10,
					downVotes: 4,
				},
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveSuggestion,
					upVotes:   4,
					downVotes: 1,
				},
			},
		},
		{
			name:   "Success/UpVoteNoChange",
			postID: test_utils.NumberUUID(1000),
			userID: test_utils.NumberUUID(210),
			target: TargetImproveSuggestion,
			vote:   VoteUp,
			now:    baseTime.Add(1 * time.Minute),
			expect: VoteUp,
			expectScores: []checkTarget{
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveRequest,
					upVotes:   10,
					downVotes: 4,
				},
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveSuggestion,
					upVotes:   5,
					downVotes: 1,
				},
			},
		},
		{
			name:   "Success/DownVoteNoChange",
			postID: test_utils.NumberUUID(1000),
			userID: test_utils.NumberUUID(211),
			target: TargetImproveRequest,
			vote:   VoteDown,
			now:    baseTime.Add(1 * time.Minute),
			expect: VoteDown,
			expectScores: []checkTarget{
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveRequest,
					upVotes:   10,
					downVotes: 4,
				},
				{
					postID:    test_utils.NumberUUID(1000),
					target:    TargetImproveSuggestion,
					upVotes:   5,
					downVotes: 1,
				},
			},
		},
		{
			name:      "Error/MissingPost",
			postID:    test_utils.NumberUUID(1010),
			userID:    test_utils.NumberUUID(212),
			target:    TargetImproveRequest,
			vote:      VoteUp,
			now:       baseTime.Add(1 * time.Minute),
			expectErr: validation.ErrMissingRelation,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.Begin()
				require.NoError(st, err)
				defer stx.Rollback()
				repository := NewRepository(stx)

				res, err := repository.Vote(ctx, d.postID, d.userID, d.target, d.vote, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)

				for _, check := range d.expectScores {
					var post interface{}
					switch check.target {
					case TargetImproveRequest:
						post = new(improve_request_storage.Model)
					case TargetImproveSuggestion:
						post = new(improve_suggestion_storage.Model)
					}

					require.NoError(st, stx.NewSelect().Model(post).Where("id = ?", check.postID).Scan(ctx))

					switch check.target {
					case TargetImproveRequest:
						require.Equal(st, check.upVotes, post.(*improve_request_storage.Model).UpVotes)
						require.Equal(st, check.downVotes, post.(*improve_request_storage.Model).DownVotes)
					case TargetImproveSuggestion:
						require.Equal(st, check.upVotes, post.(*improve_suggestion_storage.Model).UpVotes)
						require.Equal(st, check.downVotes, post.(*improve_suggestion_storage.Model).DownVotes)
					}
				}
			})
		}
	})
	require.NoError(t, err)
}

func TestVotesRepository_HasVoted(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		postID uuid.UUID
		userID uuid.UUID
		target Target

		expect    Vote
		expectErr error
	}{
		{
			name:   "Success/UpVote",
			postID: test_utils.NumberUUID(1000),
			userID: test_utils.NumberUUID(210),
			target: TargetImproveSuggestion,
			expect: VoteUp,
		},
		{
			name:   "Success/DownVote",
			postID: test_utils.NumberUUID(1000),
			userID: test_utils.NumberUUID(211),
			target: TargetImproveRequest,
			expect: VoteDown,
		},
		{
			name:   "Success/NoVote",
			postID: test_utils.NumberUUID(1000),
			userID: test_utils.NumberUUID(212),
			target: TargetImproveRequest,
			expect: NoVote,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.HasVoted(ctx, d.postID, d.userID, d.target)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestVotesRepository_GetVotedPosts(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		userID uuid.UUID
		target Target
		limit  int
		offset int

		expect      []*VotedPost
		expectCount int64
		expectErr   error
	}{
		{
			name:        "Success/ImproveRequest",
			userID:      test_utils.NumberUUID(201),
			target:      TargetImproveRequest,
			limit:       10,
			offset:      0,
			expectCount: 3,
			expect: []*VotedPost{
				{
					PostID:    test_utils.NumberUUID(6000),
					UpdatedAt: baseTime.Add(13 * time.Minute),
					Vote:      VoteUp,
				},
				{
					PostID:    test_utils.NumberUUID(1000),
					UpdatedAt: baseTime.Add(2 * time.Minute),
					Vote:      VoteUp,
				},
				{
					PostID:    test_utils.NumberUUID(1002),
					UpdatedAt: baseTime,
					Vote:      VoteDown,
				},
			},
		},
		{
			name:        "Success/ImproveSuggestion",
			userID:      test_utils.NumberUUID(201),
			target:      TargetImproveSuggestion,
			limit:       10,
			offset:      0,
			expectCount: 3,
			expect: []*VotedPost{
				{
					PostID:    test_utils.NumberUUID(1004),
					UpdatedAt: baseTime.Add(13 * time.Minute),
					Vote:      VoteUp,
				},
				{
					PostID:    test_utils.NumberUUID(1000),
					UpdatedAt: baseTime.Add(2 * time.Minute),
					Vote:      VoteUp,
				},
				{
					PostID:    test_utils.NumberUUID(1001),
					UpdatedAt: baseTime,
					Vote:      VoteDown,
				},
			},
		},
		{
			name:        "Success/Pagination",
			userID:      test_utils.NumberUUID(201),
			target:      TargetImproveSuggestion,
			limit:       1,
			offset:      1,
			expectCount: 3,
			expect: []*VotedPost{
				{
					PostID:    test_utils.NumberUUID(1000),
					UpdatedAt: baseTime.Add(2 * time.Minute),
					Vote:      VoteUp,
				},
			},
		},
	}

	// Target user: 201 - Requests
	fixtures := test_utils.Concat(
		improve_suggestion_storage.Fixtures,
		// T+2m
		generateVotesFor(
			improve_suggestion_storage.Fixtures[0],
			map[int]time.Time{200: baseTime.Add(10 * time.Minute), 201: baseTime.Add(2 * time.Minute), 210: baseTime},
			map[int]time.Time{203: baseTime.Add(8 * time.Minute)},
		),
		// T
		generateVotesFor(
			improve_suggestion_storage.Fixtures[2],
			map[int]time.Time{200: baseTime.Add(7 * time.Minute), 207: baseTime.Add(30 * time.Minute)},
			map[int]time.Time{203: baseTime.Add(time.Minute), 201: baseTime},
		),
		// -
		generateVotesFor(
			improve_suggestion_storage.Fixtures[3],
			map[int]time.Time{204: baseTime.Add(10 * time.Minute), 206: baseTime.Add(16 * time.Minute), 208: baseTime.Add(3 * time.Minute)},
			map[int]time.Time{205: baseTime.Add(4 * time.Minute), 207: baseTime.Add(5 * time.Minute)},
		),
		// T+13m
		generateVotesFor(
			improve_suggestion_storage.Fixtures[4],
			map[int]time.Time{201: baseTime.Add(13 * time.Minute), 210: baseTime.Add(10 * time.Minute)},
			map[int]time.Time{},
		),
	)

	// Target user: 201 - Suggestions
	fixtures = test_utils.Concat(
		fixtures,
		// T+2m
		generateVotesFor(
			improve_suggestion_storage.Fixtures[5],
			map[int]time.Time{200: baseTime.Add(10 * time.Minute), 201: baseTime.Add(2 * time.Minute), 210: baseTime},
			map[int]time.Time{203: baseTime.Add(8 * time.Minute)},
		),
		// T
		generateVotesFor(
			improve_suggestion_storage.Fixtures[6],
			map[int]time.Time{200: baseTime.Add(7 * time.Minute), 207: baseTime.Add(30 * time.Minute)},
			map[int]time.Time{203: baseTime.Add(time.Minute), 201: baseTime},
		),
		// -
		generateVotesFor(
			improve_suggestion_storage.Fixtures[8],
			map[int]time.Time{204: baseTime.Add(10 * time.Minute), 206: baseTime.Add(16 * time.Minute), 208: baseTime.Add(3 * time.Minute)},
			map[int]time.Time{205: baseTime.Add(4 * time.Minute), 207: baseTime.Add(5 * time.Minute)},
		),
		// T+13m
		generateVotesFor(
			improve_suggestion_storage.Fixtures[9],
			map[int]time.Time{201: baseTime.Add(13 * time.Minute), 210: baseTime.Add(10 * time.Minute)},
			map[int]time.Time{},
		),
	)

	err := test_utils.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, count, err := repository.GetVotedPosts(ctx, d.userID, d.target, d.limit, d.offset)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
				require.Equal(t, d.expectCount, count)
			})
		}
	})
	require.NoError(t, err)
}
