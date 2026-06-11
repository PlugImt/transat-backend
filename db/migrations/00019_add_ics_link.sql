-- +goose Up
-- +goose StatementBegin
ALTER TABLE newf ADD COLUMN ics_link text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE newf DROP COLUMN IF EXISTS ics_link;
-- +goose StatementEnd
