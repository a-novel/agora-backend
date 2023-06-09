CREATE TYPE sex AS ENUM ('male', 'female');

--bun:split

CREATE TABLE IF NOT EXISTS identities (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    first_name VARCHAR(32),
    last_name VARCHAR(32),
    sex sex,
    birthday TIMESTAMP
);
