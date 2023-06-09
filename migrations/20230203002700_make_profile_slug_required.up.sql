ALTER TABLE profiles
    ALTER COLUMN slug SET NOT NULL,
    ADD CONSTRAINT slug_filled CHECK (slug <> '');
