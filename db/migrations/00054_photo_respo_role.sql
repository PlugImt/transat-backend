-- +goose Up
-- +goose StatementBegin
-- The photo managers ("responsables des photos"): photo reports, the photobooth's review and the albums in the
-- admin panel, without being administrators. Created here so it can be given from the panel's user roles.
INSERT INTO roles (name) VALUES ('photo_respo') ON CONFLICT (name) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Kept: users may hold it, and removing it would silently take their access away.
SELECT 1;
-- +goose StatementEnd
