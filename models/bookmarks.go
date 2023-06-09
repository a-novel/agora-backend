package models

import (
	"github.com/google/uuid"
	"time"
)

// Bookmark represents a bookmark object.
// A bookmark allows a user to save a specific resource for later reference.
//
// The below example specifies a regular bookmark for the improvement request
// "b3f5ae63-b499-4946-b130-512606927607" for the user "090c481a-b967-49a5-9e89-d3659a0a0433".
//
//	bookmark := &improve_post_service.Bookmark{
//		UserID:        uuid.MustParse("090c481a-b967-49a5-9e89-d3659a0a0433"),
//		RequestID:     uuid.MustParse("b3f5ae63-b499-4946-b130-512606927607"),
//		CreatedAt:     time.Now(),
//		VoteTarget:    improve_post_service.BookmarkTargetImproveRequest,
//		BookmarkLevel: improve_post_service.BookmarkLevelBookmark,
//	}
type Bookmark struct {
	// UserID is the ID of the user who created the bookmark.
	UserID uuid.UUID `json:"userID"`
	// RequestID is the ID of the improvement request or suggestion that is bookmarked.
	RequestID uuid.UUID `json:"requestID"`
	// CreatedAt stores the time at which the post was bookmarked.
	CreatedAt time.Time `json:"createdAt"`
	// Target specifies the target table of the bookmark.
	Target BookmarkTarget `json:"target"`
	// Level specifies the importance level of the bookmark.
	Level BookmarkLevel `json:"level"`
}

// BookmarkLevel specifies the importance level of the bookmark.
type BookmarkLevel string

const (
	// BookmarkLevelBookmark is the reference bookmark level.
	BookmarkLevelBookmark BookmarkLevel = "bookmark"
	// BookmarkLevelFavorite specifies that the bookmark is a favorite bookmark. Favorite bookmarks are of higher
	// importance than regular bookmarks.
	BookmarkLevelFavorite BookmarkLevel = "favorite"
)

// BookmarkTarget specifies the target table of the bookmark.
// See Bookmark for reference.
type BookmarkTarget string

const (
	// BookmarkTargetImproveRequest specifies that the bookmark is related to an improvement request.
	BookmarkTargetImproveRequest BookmarkTarget = "improve_request"
	// BookmarkTargetImproveSuggestion specifies that the bookmark is related to an improvement suggestion.
	BookmarkTargetImproveSuggestion BookmarkTarget = "improve_suggestion"
)
