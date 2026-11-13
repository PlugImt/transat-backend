-- +goose Up
-- +goose StatementBegin
-- Set by an administrator, never by the Drive sync: a hidden album disappears with its photos from the app's gallery
-- and from its public page on transat.dev, whatever the classifier decides at the next sync.
ALTER TABLE photo_albums
    ADD COLUMN IF NOT EXISTS hidden    BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS hidden_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS hidden_by TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE photo_albums
    DROP COLUMN IF EXISTS hidden_by,
    DROP COLUMN IF EXISTS hidden_at,
    DROP COLUMN IF EXISTS hidden;
-- +goose StatementEnd
