package models

import (
	"github.com/google/uuid"
	"time"
)

// ImproveSuggestion represents an improvement suggestion.
// An improvement suggestion is a response to an ImproveRequest. It proposes modifications to improve the source
// request.
//
// To remain relevant, an improvement suggestion is tied to an improvement request revision. When updated, the revision
// can also be changed, to point to a more recent revision.
//
// When the improvement request creator up-votes a suggestion, the suggestion becomes validated. It then has a special
// display in the thread.
type ImproveSuggestion struct {
	// ID of the suggestion.
	ID uuid.UUID `json:"id"`
	// CreatedAt stores the time at which the suggestion was created.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt stores the time at which the suggestion was last updated.
	UpdatedAt *time.Time `json:"updatedAt"`

	// SourceID is the ID of the first revision of the related improvement request. It cannot be changed.
	SourceID uuid.UUID `json:"sourceID"`
	// UserID is the ID of the user who created the suggestion.
	UserID uuid.UUID `json:"userID"`
	// Validated is true if the suggestion has been validated by the improvement request creator.
	Validated bool `json:"validated"`

	// UpVotes is the number of up votes the suggestion has received. This value is indirectly updated from the
	// votes table.
	UpVotes int64 `json:"upVotes"`
	// DownVotes is the number of down votes the suggestion has received. This value is indirectly updated from the
	// votes table.
	DownVotes int64 `json:"downVotes"`

	// RequestID is the ID of the improvement request revision the suggestion is tied to. It must point to a revision
	// of the improvement request with the SourceID.
	RequestID uuid.UUID `json:"requestID"`
	// Title an improved version of the source Title. It should match it if no modifications are intended.
	Title string `json:"title"`
	// Content contains the updated content of the source request.
	Content string `json:"content"`
}

// ImproveSuggestionUpsert is the data required to create or update an improvement suggestion.
type ImproveSuggestionUpsert struct {
	// RequestID is the ID of the improvement request revision the suggestion is tied to. It must point to a revision
	// of the improvement request with the source ID.
	RequestID uuid.UUID `json:"requestID"`
	// Title an improved version of the source Title. It should match it if no modifications are intended.
	Title string `json:"title"`
	// Content contains the updated content of the source request.
	Content string `json:"content"`
}

// ImproveSuggestionSearchOrder allows to order suggestions requests in a search query.
type ImproveSuggestionSearchOrder struct {
	// Created puts more recent suggestions first.
	Created bool `json:"created"`
	// Score puts suggestions with the highest score first.
	Score bool `json:"score"`
}

// ImproveSuggestionsList allows to filter improvement suggestions.
type ImproveSuggestionsList struct {
	// UserID is an optional parameter, to only target suggestions that were created by a specific author.
	UserID *uuid.UUID `json:"userID"`
	// SourceID is an optional parameter, to only target suggestions that were created for a specific improvement
	// request.
	SourceID *uuid.UUID `json:"sourceID"`
	// RequestID is an optional parameter, to only target suggestions that were created for a specific improvement
	// request revision.
	RequestID *uuid.UUID `json:"requestID"`
	// Validated is an optional parameter, to only target suggestions that have been validated by the improvement
	// request creator.
	Validated *bool `json:"validated"`
	// Order is an optional parameter, to order suggestions based on a specific criteria.
	Order *ImproveSuggestionSearchOrder `json:"order"`
}
