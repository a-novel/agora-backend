DROP TRIGGER IF EXISTS improve_request_source_must_exist ON improve_requests;
DROP TRIGGER IF EXISTS format_searchable_content ON improve_requests;
DROP TRIGGER IF EXISTS improve_suggestion_source_validation ON improve_suggestions;
DROP TRIGGER IF EXISTS vote_source_must_exist ON votes;
DROP TRIGGER IF EXISTS update_score ON votes;

--bun:split

DROP FUNCTION IF EXISTS improve_request_source_must_exist;
DROP FUNCTION IF EXISTS format_searchable_content;
DROP FUNCTION IF EXISTS improve_suggestion_source_validation;
DROP FUNCTION IF EXISTS vote_source_must_exist;

--bun:split

DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS improve_request_comments;
DROP TABLE IF EXISTS improve_requests;

--bun:split

DROP TYPE IF EXISTS vote;
DROP TYPE IF EXISTS vote_target;
