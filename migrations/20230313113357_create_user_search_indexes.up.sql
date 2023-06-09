CREATE FUNCTION format_user_search(a text) RETURNS text
    LANGUAGE sql IMMUTABLE
    RETURNS NULL ON NULL INPUT
    RETURN unaccent(lower(a));

--bun:split

CREATE INDEX IF NOT EXISTS user_profile_username_search ON profiles USING gin (format_user_search(username) gin_trgm_ops);
CREATE INDEX IF NOT EXISTS user_profile_slug_search ON profiles USING gin (slug gin_trgm_ops);
CREATE INDEX IF NOT EXISTS user_identity_name_search ON identities USING gin (format_user_search(first_name || ' ' || last_name) gin_trgm_ops);
