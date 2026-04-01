-- +goose Up
-- +goose StatementBegin
-- Add custom warning and confirmation messages to reservation_element table
-- These allow club owners to customize the dialog messages shown when users reserve items

ALTER TABLE reservation_element
ADD COLUMN IF NOT EXISTS warning_message TEXT,
ADD COLUMN IF NOT EXISTS confirmation_message TEXT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove custom message columns
ALTER TABLE reservation_element
DROP COLUMN IF EXISTS warning_message,
DROP COLUMN IF EXISTS confirmation_message;

-- +goose StatementEnd
