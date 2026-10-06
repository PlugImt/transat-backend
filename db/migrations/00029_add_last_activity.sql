-- +goose Up
-- +goose StatementBegin
ALTER TABLE newf ADD COLUMN IF NOT EXISTS last_activity TIMESTAMP NULL;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_newf_last_activity ON newf (last_activity);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_newf_last_activity;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE newf DROP COLUMN IF EXISTS last_activity;
-- +goose StatementEnd
