DROP INDEX IF EXISTS user_profile_username_search;
DROP INDEX IF EXISTS user_profile_slug_search;
DROP INDEX IF EXISTS user_identity_name_search;

--bun:split

DROP FUNCTION IF EXISTS format_user_search;
