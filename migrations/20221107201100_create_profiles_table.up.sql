CREATE TABLE IF NOT EXISTS profiles (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,

    username VARCHAR(64),
    slug VARCHAR(64),

    UNIQUE(slug)
);

--bun:split

CREATE UNIQUE INDEX profiles_slug ON profiles (slug);
