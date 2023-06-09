package models

import (
	"github.com/google/uuid"
	"time"
)

// Vote represents a vote on a forum post.
//
// A vote is used to measure the relevance of an improvement suggestion, or the quality of an improvement request.
// VoteValue can either be VoteUp or VoteDown. Combined, the amount of VoteDown subtracted to the amount of VoteUp gives
// the net score of the target.
//
// Only one vote per target can be cast by a user.
type Vote struct {
	// UpdatedAt stores the time at which the vote was last updated.
	UpdatedAt time.Time `json:"updatedAt"`
	// PostID is the ID of the target of the vote.
	PostID uuid.UUID `json:"postID"`
	// UserID is the ID of the user who cast the vote.
	UserID uuid.UUID `json:"userID"`
	// Target is the target of the vote.
	Target VoteTarget `json:"target"`
	// Vote is the value of the vote.
	Vote VoteValue `json:"vote"`
}

// VotedPost represents a post voted by a user for a specific target.
type VotedPost struct {
	// PostID is the ID of the target of the vote.
	PostID uuid.UUID `json:"postID"`
	// UpdatedAt stores the time at which the vote was last updated.
	UpdatedAt time.Time `json:"updatedAt"`
	// Vote is the value of the vote.
	Vote VoteValue `json:"vote"`
}

// VoteTarget specifies the target table of the vote.
type VoteTarget string

const (
	// VoteTargetImproveRequest is the target for votes on improvement requests.
	VoteTargetImproveRequest VoteTarget = "improve_request"
	// VoteTargetImproveSuggestion is the target for votes on improvement suggestions.
	VoteTargetImproveSuggestion VoteTarget = "improve_suggestion"
)

// VoteValue specifies how a resource was voted by the user.
type VoteValue string

const (
	// VoteUp is a vote in favor of the target.
	VoteUp VoteValue = "up"
	// VoteDown is a vote against the target.
	VoteDown VoteValue = "down"
	// NoVote is a special value, that does not exist in the database, but indicates the intent to
	// undo an existing vote.
	NoVote VoteValue = ""
)
