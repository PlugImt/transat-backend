-- +goose Up
-- +goose StatementBegin
-- Drive folders already resolved for an event ("<year>/<date> - <name>/Photobooth"), so each shoot costs no lookup.
CREATE TABLE IF NOT EXISTS photobooth_event_folders (
    event_date        DATE        NOT NULL,
    event_key         TEXT        NOT NULL, -- normalized event name
    event_folder_id   TEXT        NOT NULL,
    photobooth_folder TEXT        NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (event_date, event_key)
);

-- One guest session on a booth. status: open (uploading), finalized (share link issued), cancelled (files trashed).
CREATE TABLE IF NOT EXISTS photobooth_shoots (
    id                BIGSERIAL PRIMARY KEY,
    device_key_id     TEXT        NOT NULL,
    client_session_id TEXT        NOT NULL,
    event_name        TEXT        NOT NULL,
    event_date        DATE        NOT NULL,
    layout_type       TEXT        NOT NULL DEFAULT '',
    drive_folder_id   TEXT        NOT NULL,
    status            TEXT        NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'finalized', 'cancelled')),
    share_token       TEXT UNIQUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    finalized_at      TIMESTAMPTZ,
    cancelled_at      TIMESTAMPTZ,
    UNIQUE (device_key_id, client_session_id)
);

-- One image of a shoot. kind photo (number 1..n) or collage (number 0). status: pending (upload session open),
-- uploaded (verified on Drive), trash (must be removed from Drive), deleted (gone from Drive).
CREATE TABLE IF NOT EXISTS photobooth_files (
    id             BIGSERIAL PRIMARY KEY,
    shoot_id       BIGINT      NOT NULL REFERENCES photobooth_shoots (id) ON DELETE CASCADE,
    kind           TEXT        NOT NULL CHECK (kind IN ('photo', 'collage')),
    number         INT         NOT NULL CHECK (number >= 0 AND number <= 50),
    file_name      TEXT        NOT NULL,
    mime_type      TEXT        NOT NULL,
    size_bytes     BIGINT      NOT NULL,
    md5            TEXT        NOT NULL,
    drive_file_id  TEXT,
    width          INT,
    height         INT,
    status         TEXT        NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'uploaded', 'trash', 'deleted')),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    uploaded_at    TIMESTAMPTZ
);
-- At most one live image per slot: a retake replaces the previous one, which moves to trash first.
CREATE UNIQUE INDEX IF NOT EXISTS photobooth_files_live_slot
    ON photobooth_files (shoot_id, kind, number) WHERE status IN ('pending', 'uploaded');
CREATE INDEX IF NOT EXISTS photobooth_files_to_trash
    ON photobooth_files (id) WHERE status = 'trash';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS photobooth_files;
DROP TABLE IF EXISTS photobooth_shoots;
DROP TABLE IF EXISTS photobooth_event_folders;
-- +goose StatementEnd
