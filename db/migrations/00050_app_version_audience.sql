-- +goose Up
-- +goose StatementBegin
-- Comparable form of an app version ("2.10.1" -> {2,10,1}). A device that never reported its version (it predates
-- the version logging) reads as {0}, the oldest of all.
CREATE OR REPLACE FUNCTION app_version_key(v TEXT) RETURNS INT[]
    LANGUAGE sql IMMUTABLE AS
$$
SELECT COALESCE(string_to_array(substring(COALESCE(v, '') FROM '^[0-9]+(?:\.[0-9]+)*'), '.')::INT[], '{0}')
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS app_version_key(TEXT);
-- +goose StatementEnd
