-- +goose Up
-- +goose StatementBegin
-- Who opened a public page: a photobooth share link (the QR code of a print, target = the shoot) or a public
-- album (target = the album). visitor is the SHA-256 of a random id the visitor's browser keeps, never an address;
-- one row per visitor and page, views counts the times they opened it.
CREATE TABLE IF NOT EXISTS public_visits (
    kind       TEXT        NOT NULL CHECK (kind IN ('photobooth_share', 'photo_album')),
    target_id  BIGINT      NOT NULL,
    visitor    TEXT        NOT NULL,
    views      INT         NOT NULL DEFAULT 1,
    first_seen TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (kind, target_id, visitor)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public_visits;
-- +goose StatementEnd
