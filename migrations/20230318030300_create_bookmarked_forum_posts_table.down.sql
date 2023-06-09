DROP INDEX IF EXISTS improve_posts_bookmarked_by_user;

--bun:split

DROP TABLE IF EXISTS improve_posts_bookmarks;

--bun:split

DROP TYPE IF EXISTS bookmark_level;
DROP TYPE IF EXISTS improve_posts_bookmark_target;
