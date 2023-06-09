package improve_post_storage

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// BookmarkTarget specifies the target table of the bookmark.
// See Model for reference.
type BookmarkTarget string

const (
	// BookmarkTargetImproveRequest specifies that the bookmark is related to an improvement request.
	BookmarkTargetImproveRequest BookmarkTarget = "improve_request"
	// BookmarkTargetImproveSuggestion specifies that the bookmark is related to an improvement suggestion.
	BookmarkTargetImproveSuggestion BookmarkTarget = "improve_suggestion"
)

// Model is the database model for the bookmarks table.
// A bookmark allows a user to save a specific post for later reference.
//
// The below example specifies a regular bookmark for the improvement request
// "b3f5ae63-b499-4946-b130-512606927607" for the user "090c481a-b967-49a5-9e89-d3659a0a0433".
//
//	bookmark := &improve_post_storage.Model{
//		UserID:     uuid.MustParse("090c481a-b967-49a5-9e89-d3659a0a0433"),
//		RequestID:  uuid.MustParse("b3f5ae63-b499-4946-b130-512606927607"),
//		CreatedAt:  time.Now(),
//		Target:     improve_post_storage.BookmarkTargetImproveRequest,
//		Level:      improve_post_storage.LevelBookmark,
//	}
type Model struct {
	bun.BaseModel `bun:"table:improve_posts_bookmarks"`

	// UserID is the ID of the user who created the bookmark.
	UserID uuid.UUID `json:"user_id" bun:"user_id,pk,type:uuid"`
	// RequestID is the ID of the improvement request or suggestion that is bookmarked.
	RequestID uuid.UUID `json:"request_id" bun:"request_id,pk,type:uuid"`
	// CreatedAt stores the time at which the post was bookmarked.
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull"`
	// Target specifies the target table of the bookmark.
	Target BookmarkTarget `json:"target" bun:"target,pk"`
	// Level specifies the importance level of the bookmark.
	Level bookmark_storage.Level `json:"level" bun:"level"`
}
