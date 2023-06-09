ALTER TABLE credentials
    DROP COLUMN deleted_at,
    ADD CONSTRAINT email_user_filled CHECK ( email_user <> '' AND email_user IS NOT NULL ),
    ADD CONSTRAINT email_domain_filled CHECK ( email_domain <> '' AND email_domain IS NOT NULL ),
    ADD CONSTRAINT require_full_new_email
        CHECK (
            /* Must either be all empty or all filled */
            (
                ( new_email_user = '' OR new_email_user IS NULL ) AND
                ( new_email_domain = '' OR new_email_domain IS NULL ) AND
                ( new_email_validation_code = '' OR new_email_validation_code IS NULL )
            ) OR (
                new_email_user IS NOT NULL AND new_email_user <> '' AND
                new_email_domain IS NOT NULL AND new_email_domain <> '' AND
                new_email_validation_code IS NOT NULL AND new_email_validation_code <> ''
            )
        );
