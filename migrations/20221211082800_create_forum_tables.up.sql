CREATE TYPE vote AS ENUM ('up', 'down');
CREATE TYPE vote_target AS ENUM ('improve_request', 'improve_suggestion');

--bun:split

CREATE TABLE IF NOT EXISTS improve_requests (
    id uuid PRIMARY KEY NOT NULL,
    source uuid NOT NULL,
    created_at TIMESTAMP NOT NULL,

    user_id uuid NOT NULL,
    title VARCHAR(256) NOT NULL,
    content TEXT NOT NULL,

    up_votes BIGINT,
    down_votes BIGINT,
    text_searchable_index_col tsvector,

    CONSTRAINT title_filled CHECK ( title <> '' ),
    CONSTRAINT content_filled CHECK ( content <> '' ),
    CONSTRAINT content_length CHECK ( char_length(content) <= 4096 )
);

CREATE TABLE IF NOT EXISTS improve_suggestions (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,

    validated boolean,
    user_id uuid NOT NULL,
    source_id uuid NOT NULL,
    request_id uuid NOT NULL,

    title VARCHAR(256) NOT NULL,
    content TEXT NOT NULL,

    up_votes BIGINT,
    down_votes BIGINT,

    CONSTRAINT title_filled CHECK ( title <> '' ),
    CONSTRAINT content_filled CHECK ( content <> '' ),
    CONSTRAINT content_length CHECK ( char_length(content) <= 4096 )
);

CREATE TABLE IF NOT EXISTS votes (
    updated_at TIMESTAMP NOT NULL,
    user_id uuid NOT NULL,
    post_id uuid NOT NULL,
    vote vote NOT NULL,
    target vote_target NOT NULL,

    UNIQUE(user_id, post_id, target)
);

--bun:split

/* unaccent cannot be used in a standard constraint, as it is not immutable */
CREATE FUNCTION format_searchable_content()
    RETURNS trigger AS $format_searchable_content$
BEGIN
    NEW.text_searchable_index_col :=
                setweight(to_tsvector('french',  unaccent(NEW.title)), 'A') ||
                setweight(to_tsvector('french', unaccent(NEW.content)), 'B');
    RETURN NEW;
END;
$format_searchable_content$ LANGUAGE plpgsql;

CREATE FUNCTION update_score()
RETURNS trigger AS $update_score$
DECLARE target vote_target; DECLARE target_id uuid; DECLARE downdiff BIGINT; DECLARE updiff BIGINT;
BEGIN
    target := CASE WHEN NEW IS NULL THEN OLD.target ELSE NEW.target END;
    target_id := CASE WHEN NEW IS NULL THEN OLD.post_id ELSE NEW.post_id END;
    updiff := 0;
    downdiff := 0;

    IF OLD IS NOT NULL THEN
        IF OLD.vote = 'up' THEN
            updiff := updiff - 1;
        ELSIF OLD.vote = 'down' THEN
            downdiff := downdiff - 1;
        END IF;
    END IF;

    IF NEW IS NOT NULL THEN
        IF NEW.vote = 'up' THEN
            updiff := updiff + 1;
        ELSIF NEW.vote = 'down' THEN
            downdiff := downdiff + 1;
        END IF;
    END IF;

    IF target = 'improve_request' THEN
        UPDATE improve_requests SET up_votes = up_votes + updiff, down_votes = down_votes + downdiff WHERE id = target_id;
    ELSIF target = 'improve_suggestion' THEN
        UPDATE improve_suggestions SET up_votes = up_votes + updiff, down_votes = down_votes + downdiff WHERE id = target_id;
    ELSE
        RAISE EXCEPTION 'Invalid vote target';
    END IF;

    RETURN NEW;
END;
$update_score$ LANGUAGE plpgsql;

--bun:split
/*
Use a trigger so that we can use non-immutable functions for proper formatting.
*/
CREATE TRIGGER format_searchable_content
    BEFORE INSERT ON improve_requests
    FOR EACH ROW
    EXECUTE FUNCTION format_searchable_content();

CREATE TRIGGER update_score
    AFTER INSERT OR UPDATE OR DELETE ON votes
    FOR EACH ROW
    EXECUTE FUNCTION update_score();
