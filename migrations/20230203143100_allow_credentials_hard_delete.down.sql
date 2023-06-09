ALTER TABLE credentials
    ADD COLUMN deleted_at TIMESTAMP,
    DROP CONSTRAINT email_user_filled,
    DROP CONSTRAINT email_domain_filled,
    DROP CONSTRAINT require_full_new_email;
