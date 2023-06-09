CREATE INDEX IF NOT EXISTS improve_requests_source ON improve_requests (source);
CREATE INDEX IF NOT EXISTS improve_requests_user ON improve_requests (user_id);
CREATE INDEX IF NOT EXISTS improve_requests_fts ON improve_requests USING GIN (text_searchable_index_col);
CREATE INDEX IF NOT EXISTS improve_requests_last_rev ON improve_requests (source, created_at DESC NULLS LAST);

CREATE INDEX IF NOT EXISTS improve_suggestions_source ON improve_suggestions (source_id);
CREATE INDEX IF NOT EXISTS improve_suggestions_request ON improve_suggestions (request_id);
CREATE INDEX IF NOT EXISTS improve_suggestions_user ON improve_suggestions (user_id);

CREATE INDEX IF NOT EXISTS votes_on_post ON votes (post_id, vote, target);
CREATE INDEX IF NOT EXISTS votes_by_user ON votes (user_id);
