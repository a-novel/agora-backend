package improve_request_storage

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"strings"
	"time"
)

// Improve request table has some full text search columns we don't want to fetch.
var (
	exposedColumns = []string{
		"id",
		"source",
		"created_at",
		"user_id",
		"title",
		"content",
		"up_votes",
		"down_votes",
	}
	exposedColumnsSTR = strings.Join(exposedColumns, ",")
)

// Model is the database model for the improve_requests table.
// An improvement request contains a novel scene that the user wants to improve.
//
// An improvement request is followed by a series of improvement suggestions (improve_suggestion_storage.Model), which
// are stored in another repository. On each update, a new revision is created, in order to keep existing suggestions
// relevant.
type Model struct {
	bun.BaseModel `bun:"table:improve_requests,alias:improve_requests"`

	// ID of the request.
	ID uuid.UUID `json:"id" bun:"id,pk,type:uuid"`
	// Source points to the first revision. It is equal to the ID if no other revision exist.
	// An improvement request is never updated. Instead, new revisions are created every time.
	Source uuid.UUID `json:"source" bun:"source,type:uuid"`
	// CreatedAt stores the time at which the request (or the current revision) was created.
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull"`

	// UserID is the ID of the user who created the request, or edited the revision.
	UserID uuid.UUID `json:"user_id" bun:"user_id,type:uuid"`
	// Title is a quick summary of the Content, and the goal it tries to achieve.
	Title string `json:"title" bun:"title"`
	// Content is a novel scene that the user wants to improve.
	Content string `json:"content" bun:"content"`

	// UpVotes is the number of up votes the request has received. This value is indirectly updated from the
	// votes table.
	UpVotes int64 `json:"up_votes" bun:"up_votes"`
	// DownVotes is the number of down votes the request has received. This value is indirectly updated from the
	// votes table.
	DownVotes int64 `json:"down_votes" bun:"down_votes"`

	SuggestionsCount         int `json:"suggestions_count" bun:"suggestions_count,scanonly"`
	AcceptedSuggestionsCount int `json:"accepted_suggestions_count" bun:"accepted_suggestions_count,scanonly"`
}

type Preview struct {
	// ID of the request.
	ID uuid.UUID `json:"id" bun:"id,pk,type:uuid"`
	// Source points to the first revision. It is equal to the ID if no other revision exist.
	// An improvement request is never updated. Instead, new revisions are created every time.
	Source uuid.UUID `json:"source" bun:"source,type:uuid"`
	// CreatedAt stores the time at which the request (or the current revision) was created.
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull"`

	// UserID is the ID of the user who created the request, or edited the revision.
	UserID uuid.UUID `json:"user_id" bun:"user_id,type:uuid"`
	// Title is a quick summary of the Content, and the goal it tries to achieve.
	Title string `json:"title" bun:"title"`
	// Content is a novel scene that the user wants to improve.
	Content string `json:"content" bun:"content"`

	// UpVotes is the number of up votes the request and all its revisions has received. This value is indirectly
	// updated from the votes table.
	UpVotes int64 `json:"up_votes" bun:"up_votes"`
	// DownVotes is the number of down votes the request and all its revisions has received. This value is indirectly
	// updated from the votes table.
	DownVotes int64 `json:"down_votes" bun:"down_votes"`

	// RevisionCount is the number of revisions the request has.
	RevisionCount int64 `json:"revision_count" bun:"revision_count"`
	// MoreRecentRevisions is the number of revisions that were created after the current one.
	MoreRecentRevisions      int `json:"more_recent_revisions" bun:"more_recent_revisions"`
	SuggestionsCount         int `json:"suggestions_count" bun:"suggestions_count"`
	AcceptedSuggestionsCount int `json:"accepted_suggestions_count" bun:"accepted_suggestions_count"`
}

type SearchQueryOrder struct {
	Created bool `json:"created"`
	Score   bool `json:"score"`
}

// SearchQuery allows to filter improve requests.
type SearchQuery struct {
	// UserID is an optional parameter, to only target requests that were created/revised by a specific author.
	UserID *uuid.UUID `json:"user_id"`
	// Query is an optional parameter, to filter requests based on their title or content.
	Query string            `json:"query"`
	Order *SearchQueryOrder `json:"order"`
}
