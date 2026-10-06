-- +goose Up
-- +goose StatementBegin
-- Opaque identifier shown in public payloads, so emails never have to be used to point at a user.
ALTER TABLE newf ADD COLUMN IF NOT EXISTS public_id UUID NOT NULL DEFAULT gen_random_uuid();
-- +goose StatementEnd
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS idx_newf_public_id ON newf (public_id);
-- +goose StatementEnd
-- +goose StatementBegin
-- Fields a user chooses to show on their public profile (kept apart from account data).
CREATE TABLE IF NOT EXISTS user_profiles
(
    email      VARCHAR(100) PRIMARY KEY,
    bio        VARCHAR(160) NOT NULL DEFAULT '',
    emoji      VARCHAR(32)  NOT NULL DEFAULT '',
    interests  TEXT[]       NOT NULL DEFAULT '{}',
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (email) REFERENCES newf (email) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_profiles;
-- +goose StatementEnd
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_newf_public_id;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE newf DROP COLUMN IF EXISTS public_id;
-- +goose StatementEnd
