-- +goose Up
-- +goose StatementBegin
ALTER TABLE newf ADD COLUMN last_activity TIMESTAMP NULL;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_newf_last_activity ON newf (last_activity);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_newf_last_activity;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE newf DROP COLUMN last_activity;
-- +goose StatementEnd
