-- +goose Up
-- +goose StatementBegin
ALTER TABLE newf ADD COLUMN last_activity TIMESTAMP NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE newf DROP COLUMN last_activity;
-- +goose StatementEnd
