# Bookmarks

A bookmark is an object that links a user to a specific post. It is a way to save a post for later reading.

A bookmark is always linking a single user to a single object. The object must be uniquely retrievable, through a link
that reference the said object alone. For example, if a user bookmarks a post on a thread, then there must exist a link
to that specific post; it cannot just be a link to the parent thread.

Bookmarks are managed by the user and only the user. They cannot be automatically created.

# Bookmark levels

Bookmarks may come in different levels:

 - **Bookmark**: the standard bookmarking level.
 - **Favorite**: a bookmark of higher importance.

Those levels are referenced through a generic type `Level`.

```go
var storageLevel bookmark_storage.Level
var serviceLevel bookmark_service.Level
```

One post can only have a single bookmark level. It is not possible to have a post bookmarked and favorited at the 
same time.

# Bookmarking

Each bookmarking service/repository presents similar accessors:
 - `Bookmark`: bookmarks a post.
 - `Unbookmark`: unbookmarks a post.
 - `IsBookmarked`: check if a post is bookmarked.
 - `List`: get all the bookmarks for a user, paginated.

## Improve post bookmarking

Bookmarking service for forum improvement posts (request / suggestion).

To differentiate between a request and a suggestion, the `BookmarkTarget` type is used.

```go
package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	bookmark_service "github.com/a-novel/agora-backend/v2/domains/bookmark/service"
	improve_post_service "github.com/a-novel/agora-backend/v2/domains/bookmark/service/improve_post"
	improve_post_storage "github.com/a-novel/agora-backend/v2/domains/bookmark/storage/improve_post"
	"time"
)

func main() {
	// Requires a bun.DB for the repository.
	db := CreateDB()

	// Prepare variables. 
	ctx := context.Background()
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	requestID := uuid.MustParse("10001000-1000-1000-1000-100010001000")

	// Create the service.
	repository := improve_post_storage.NewRepository(db)
	service := improve_post_service.NewService(repository)

	// Bookmark an improvement request.
	_, err := service.Bookmark(
		ctx, userID, requestID,
		improve_post_service.BookmarkTargetImproveRequest,
		bookmark_service.LevelBookmark,
		time.Now(),
	)

	// Update the bookmark level.
	_, err = service.Bookmark(
		ctx, userID, requestID,
		improve_post_service.BookmarkTargetImproveRequest,
		bookmark_service.LevelFavorite,
		time.Now(),
	)

	// Check if the request is bookmarked.
	isBookmarked, err := service.IsBookmarked(ctx, userID, requestID, improve_post_service.BookmarkTargetImproveRequest)
	fmt.Println(isBookmarked) // True

	// UnBookmark the request.
	err = service.UnBookmark(ctx, userID, requestID, improve_post_service.BookmarkTargetImproveRequest)
	
	isBookmarked, err = service.IsBookmarked(ctx, userID, requestID, improve_post_service.BookmarkTargetImproveRequest)
	fmt.Println(isBookmarked) // False
	
	// List all favorite bookmarks for a user, page 3 (10 rows/page).
	bookmarks, total, err := service.List(ctx, userID, bookmark_service.LevelFavorite, 10, 20)
}
```
