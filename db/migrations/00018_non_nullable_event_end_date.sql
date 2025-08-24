-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
    ALTER COLUMN end_date SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
    ALTER COLUMN end_date DROP NOT NULL;
-- +goose StatementEnd
