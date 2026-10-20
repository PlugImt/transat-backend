-- +goose Up
-- +goose StatementBegin
-- Identifier of the interactive notification category (the action buttons the app registered under that name).
-- Null for a notification without buttons.
ALTER TABLE push_messages ADD COLUMN IF NOT EXISTS action_category TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE push_messages DROP COLUMN IF EXISTS action_category;
-- +goose StatementEnd
