package bookmark_storage

// Level specifies the importance level of the bookmark.
// See Model for reference.
type Level string

const (
	// LevelBookmark is the reference bookmark level.
	LevelBookmark Level = "bookmark"
	// LevelFavorite specifies that the bookmark is a favorite bookmark. Favorite bookmarks are of higher
	// importance than regular bookmarks.
	LevelFavorite Level = "favorite"
)
