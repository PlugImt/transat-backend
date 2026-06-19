-- +goose Up
-- +goose StatementBegin

CREATE TABLE inte_schedule
(
    raw_ics  TEXT,
    ics_json JSONB
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS inte_schedule CASCADE;

-- +goose StatementEnd
