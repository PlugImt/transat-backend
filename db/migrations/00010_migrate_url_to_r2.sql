-- +goose Up
-- +goose StatementBegin
UPDATE newf SET profile_picture = REPLACE(profile_picture, 'https://transat.destimt.fr/api/data/', 'https://assets.dev.transat.dev/') WHERE profile_picture LIKE 'https%';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE newf SET profile_picture = REPLACE(profile_picture, 'https://assets.dev.transat.dev/', 'https://transat.destimt.fr/api/data/') WHERE profile_picture LIKE 'https%';
-- +goose StatementEnd
