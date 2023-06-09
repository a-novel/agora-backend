package models

import (
	"github.com/google/uuid"
	"time"
)

// ImproveRequest represents an improvement request.
// An improvement request contains a novel scene that the user wants to improve.
//
// An improvement request is followed by a series of ImproveSuggestion, which are stored in another repository. On
// each update, a new revision is created, in order to keep existing suggestions relevant.
type ImproveRequest struct {
	// ID of the request.
	ID uuid.UUID `json:"id"`
	// Source points to the first revision. It is equal to the ID if no other revision exist.
	// An improvement request is never updated. Instead, new revisions are created every time.
	Source uuid.UUID `json:"source"`
	// CreatedAt stores the time at which the request (or the current revision) was created.
	CreatedAt time.Time `json:"createdAt"`

	// UserID is the ID of the user who created the request, or edited the revision.
	UserID uuid.UUID `json:"userID"`
	// Title is a quick summary of the Content, and the goal it tries to achieve.
	Title string `json:"title"`
	// Content is a novel scene that the user wants to improve.
	Content string `json:"content"`

	// UpVotes is the number of up votes the request has received. This value is indirectly updated from the
	// votes table.
	UpVotes int64 `json:"upVotes"`
	// DownVotes is the number of down votes the request has received. This value is indirectly updated from the
	// votes table.
	DownVotes int64 `json:"downVotes"`
}

// ImproveRequestPreview merges together different metrics about an improvement request, for display in a preview
// list (eg. search results).
type ImproveRequestPreview struct {
	// ID of the request.
	ID uuid.UUID `json:"id"`
	// Source points to the first revision. It is equal to the ID if no other revision exist.
	// An improvement request is never updated. Instead, new revisions are created every time.
	Source uuid.UUID `json:"source"`
	// CreatedAt stores the time at which the request (or the current revision) was created.
	CreatedAt time.Time `json:"createdAt"`

	// UserID is the ID of the user who created the request, or edited the revision.
	UserID uuid.UUID `json:"userID"`
	// Title is a quick summary of the Content, and the goal it tries to achieve.
	Title string `json:"title"`
	// Content is a novel scene that the user wants to improve.
	Content string `json:"content"`

	// UpVotes is the number of up votes the request has received. This value is indirectly updated from the
	// votes table.
	UpVotes int64 `json:"upVotes"`
	// DownVotes is the number of down votes the request has received. This value is indirectly updated from the
	// votes table.
	DownVotes int64 `json:"downVotes"`

	// RevisionCount is the number of revisions the request has.
	RevisionCount int64 `json:"revisionCount"`
	// MoreRecentRevisions is the number of revisions that were created after the current one.
	MoreRecentRevisions int `json:"moreRecentRevisions"`
	// SuggestionsCount is the number of suggestions the request has received, for every revision combined.
	SuggestionsCount int `json:"suggestionsCount"`
	// AcceptedSuggestionsCount is the number of suggestions that were accepted by the author, for every
	// revision combined.
	AcceptedSuggestionsCount int `json:"acceptedSuggestionsCount"`
}

// ImproveRequestSearchOrder allows to order improve requests in a search query.
type ImproveRequestSearchOrder struct {
	// Created puts more recent requests first.
	Created bool `json:"created"`
	// Score puts requests with the highest score first.
	Score bool `json:"score"`
}

// ImproveRequestSearch allows to filter improve requests.
type ImproveRequestSearch struct {
	// UserID is an optional parameter, to only target requests that were created/revised by a specific author.
	UserID *uuid.UUID `json:"userID"`
	// Query is an optional parameter, to filter requests based on their title or content.
	Query string `json:"query"`
	// Order is an optional parameter, to order requests based on a specific criteria.
	Order *ImproveRequestSearchOrder `json:"order"`
}
