CREATE TYPE improve_posts_bookmark_target AS ENUM ('improve_request', 'improve_suggestion');
CREATE TYPE bookmark_level AS ENUM ('bookmark', 'favorite');

--bun:split

CREATE TABLE IF NOT EXISTS improve_posts_bookmarks (
    user_id uuid NOT NULL,
    request_id uuid NOT NULL,
    created_at TIMESTAMP NOT NULL,
    target improve_posts_bookmark_target NOT NULL,
    level bookmark_level NOT NULL,

    UNIQUE(user_id, request_id, target)
);

--bun:split

CREATE INDEX IF NOT EXISTS improve_posts_bookmarked_by_user ON improve_posts_bookmarks (user_id, level, target);
