-- +goose Up
-- +goose StatementBegin
ALTER TABLE newf
    ADD COLUMN phone_number_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN phone_number_verification_code VARCHAR(6),
    ADD COLUMN phone_number_verification_code_expiration TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE newf
    DROP COLUMN phone_number_verified,
    DROP COLUMN phone_number_verification_code,
    DROP COLUMN phone_number_verification_code_expiration;
-- +goose StatementEnd
