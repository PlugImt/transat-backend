-- +goose Up
-- +goose StatementBegin
UPDATE newf SET profile_picture = REPLACE(profile_picture, 'https://transat.destimt.fr/api/data/', 'https://assets.prod.transat.dev/') WHERE profile_picture LIKE 'https://transat.destimt.fr/api/data/%';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE newf SET profile_picture = REPLACE(profile_picture, 'https://assets.prod.transat.dev/', 'https://transat.destimt.fr/api/data/') WHERE profile_picture LIKE 'https://assets.prod.transat.dev/%';
-- +goose StatementEnd
