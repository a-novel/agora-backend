package improve_suggestion_storage

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// Model is the database model for the improve_suggestions table.
// An improvement suggestion is a response to an improvement request (improve_request_storage.Model). It proposes
// alterations to improve the source request, and achieve its goal.
//
// To remain relevant, an improvement suggestion is tied to an improvement request revision. When updated, the revision
// can also be changed, to point to another more recent revision.
//
// When the improvement request creator upvotes a suggestion, the suggestion becomes validated. It then has a special
// display in the thread.
type Model struct {
	bun.BaseModel `bun:"table:improve_suggestions"`

	// ID of the suggestion.
	ID uuid.UUID `json:"id" bun:"id,pk,type:uuid"`
	// CreatedAt stores the time at which the suggestion was created.
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull"`
	// UpdatedAt stores the time at which the suggestion was last updated.
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`

	// SourceID is the ID of the first revision of the related improvement request. It cannot be changed.
	SourceID uuid.UUID `json:"source_id" bun:"source_id,type:uuid"`
	// UserID is the ID of the user who created the suggestion.
	UserID uuid.UUID `json:"user_id" bun:"user_id,type:uuid"`
	// Validated is true if the suggestion has been validated by the improvement request creator.
	Validated bool `json:"validated" bun:"validated"`

	// UpVotes is the number of up votes the suggestion has received. This value is indirectly updated from the
	// votes table.
	UpVotes int64 `json:"up_votes" bun:"up_votes"`
	// DownVotes is the number of down votes the suggestion has received. This value is indirectly updated from the
	// votes table.
	DownVotes int64 `json:"down_votes" bun:"down_votes"`

	Core
}

// Core contains the explicitly mutable data of the suggestion.
type Core struct {
	// RequestID is the ID of the improvement request revision the suggestion is tied to. It must point to a revision
	// of the improvement request with the Model.SourceID.
	RequestID uuid.UUID `json:"request_id" bun:"request_id,type:uuid"`
	// Title an improved version of the source Title. It should match it if no modifications are intended.
	Title string `json:"title" bun:"title"`
	// Content contains the updated content of the source request.
	Content string `json:"content" bun:"content"`
}

type SearchQueryOrder struct {
	Created bool `json:"created"`
	Score   bool `json:"score"`
}

// ListQuery allows to filter improve suggestions.
type ListQuery struct {
	// UserID is an optional parameter, to only target suggestions that were created by a specific author.
	UserID *uuid.UUID `json:"user_id"`
	// SourceID is an optional parameter, to only target suggestions that were created for a specific improvement
	// request.
	SourceID *uuid.UUID `json:"source_id"`
	// RequestID is an optional parameter, to only target suggestions that were created for a specific improvement
	// request revision.
	RequestID *uuid.UUID `json:"request_id"`
	// Validated is an optional parameter, to only target suggestions that have been validated by the improvement
	// request creator.
	Validated *bool             `json:"validated"`
	Order     *SearchQueryOrder `json:"order"`
}
