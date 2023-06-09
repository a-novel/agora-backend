package votes_storage

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// Target specifies the target table of the vote.
type Target string

const (
	// TargetImproveRequest is the target for votes on improvement requests.
	TargetImproveRequest Target = "improve_request"
	// TargetImproveSuggestion is the target for votes on improvement suggestions.
	TargetImproveSuggestion Target = "improve_suggestion"
)

// Vote is the value of the vote.
type Vote string

const (
	// VoteUp is a vote in favor of the target.
	VoteUp Vote = "up"
	// VoteDown is a vote against the target.
	VoteDown Vote = "down"
	// NoVote is a special value, that does not exist in the database, but indicates the intent to
	// undo an existing vote.
	NoVote Vote = ""
)

// Model is the database model for the votes table.
// A vote is used to measure the relevance of an improvement suggestion, or the quality of an improvement request.
// Vote can either be VoteUp or VoteDown. Combined, the amount of VoteDown subtracted to the amount of VoteUp gives
// the net score of the target.
// Only one vote per target can be cast by a user.
type Model struct {
	bun.BaseModel `bun:"table:votes"`

	// UpdatedAt stores the time at which the vote was last updated.
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,notnull"`
	// PostID is the ID of the target of the vote.
	PostID uuid.UUID `json:"post_id" bun:"post_id,type:uuid"`
	// UserID is the ID of the user who cast the vote.
	UserID uuid.UUID `json:"user_id" bun:"user_id,type:uuid"`
	// Target is the target of the vote.
	Target Target `json:"target" bun:"target"`
	// Vote is the value of the vote.
	Vote Vote `json:"vote" bun:"vote"`
}

// VotedPost represents a post voted by a user for a specific target.
type VotedPost struct {
	bun.BaseModel `bun:"table:votes"`

	// PostID is the ID of the target of the vote.
	PostID uuid.UUID `json:"post_id" bun:"post_id,type:uuid"`
	// UpdatedAt stores the time at which the vote was last updated.
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,notnull"`
	// Vote is the value of the vote.
	Vote Vote `json:"vote" bun:"vote"`
}
