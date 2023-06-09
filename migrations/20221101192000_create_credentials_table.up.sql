CREATE TABLE IF NOT EXISTS credentials (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    email_user VARCHAR(128),
    email_domain VARCHAR(128),
    email_validation_code VARCHAR(256),
    new_email_user VARCHAR(128),
    new_email_domain VARCHAR(128),
    new_email_validation_code VARCHAR(256),
    password_hashed VARCHAR(2048),
    password_validation_code VARCHAR(256),

    UNIQUE(email_user, email_domain)
);

--bun:split

CREATE UNIQUE INDEX credentials_email ON credentials (email_user, email_domain);
