-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_profiles
    ADD COLUMN decoration_image_ids INTEGER[] NOT NULL DEFAULT '{}',
    ADD COLUMN decoration_emojis TEXT[] NOT NULL DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_profiles
    DROP COLUMN decoration_emojis,
    DROP COLUMN decoration_image_ids;
-- +goose StatementEnd
