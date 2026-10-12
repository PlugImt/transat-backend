-- +goose Up
-- +goose StatementBegin
-- Optional image shown with the notification (https link, downloaded by the phone). Null for a text notification.
ALTER TABLE push_messages ADD COLUMN IF NOT EXISTS image_url TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE push_messages DROP COLUMN IF EXISTS image_url;
-- +goose StatementEnd
