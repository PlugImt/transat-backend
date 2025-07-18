-- +goose Up
-- +goose StatementBegin

ALTER TABLE newf
    ADD COLUMN formation_name VARCHAR(50) CHECK ( formation_name IN ('FISE', 'FIL', 'FIP', 'FIT', 'FID'));

-- +goose StatementEnd
