-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
    ADD CONSTRAINT events_end_date_gte_start_date
    CHECK (end_date >= start_date);
ALTER TABLE events
    ALTER COLUMN end_date SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
    DROP CONSTRAINT events_end_date_gte_start_date;
ALTER TABLE events
    ALTER COLUMN end_date DROP NOT NULL;
-- +goose StatementEnd
