-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS photo_sync_runs
(
    id          BIGSERIAL PRIMARY KEY,
    started_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMPTZ,
    complete    BOOLEAN     NOT NULL DEFAULT FALSE,
    stats       JSONB       NOT NULL DEFAULT '{}',
    error       TEXT
);
-- +goose StatementEnd
-- +goose StatementBegin
-- An album is an event folder of the shared Drive; root_folder_id is its stable identity.
CREATE TABLE IF NOT EXISTS photo_albums
(
    id_photo_albums   SERIAL PRIMARY KEY,
    root_folder_id    TEXT        NOT NULL UNIQUE,
    parent_album_id   INT REFERENCES photo_albums (id_photo_albums) ON DELETE SET NULL,
    title             TEXT        NOT NULL,
    raw_folder_name   TEXT        NOT NULL,
    event_date        DATE,
    event_date_source TEXT        NOT NULL DEFAULT '',
    academic_year     TEXT        NOT NULL DEFAULT '',
    category          TEXT        NOT NULL DEFAULT '',
    drive_url         TEXT        NOT NULL,
    visibility        TEXT        NOT NULL CHECK (visibility IN ('published', 'needs_review', 'hidden')),
    confidence        NUMERIC(3, 2) NOT NULL,
    rule              TEXT        NOT NULL,
    cover_photo_id    TEXT,
    photo_count       INT         NOT NULL DEFAULT 0,
    first_taken_at    TIMESTAMPTZ,
    last_taken_at     TIMESTAMPTZ,
    last_seen_run     BIGINT      NOT NULL,
    deleted_at        TIMESTAMPTZ
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS photo_folders
(
    drive_folder_id    TEXT PRIMARY KEY,
    parent_drive_id    TEXT,
    name               TEXT          NOT NULL,
    effective_name     TEXT          NOT NULL,
    shortcut_id        TEXT,
    web_url            TEXT          NOT NULL,
    drive_modified_at  TIMESTAMPTZ,
    kind               TEXT          NOT NULL,
    rule               TEXT          NOT NULL,
    confidence         NUMERIC(3, 2) NOT NULL,
    status             TEXT          NOT NULL,
    signals            TEXT[]        NOT NULL DEFAULT '{}',
    classifier_version INT           NOT NULL,
    academic_year      TEXT          NOT NULL DEFAULT '',
    category           TEXT          NOT NULL DEFAULT '',
    section_path       TEXT          NOT NULL DEFAULT '',
    album_id           INT REFERENCES photo_albums (id_photo_albums) ON DELETE SET NULL,
    last_seen_run      BIGINT        NOT NULL,
    deleted_at         TIMESTAMPTZ
);
-- +goose StatementEnd
-- +goose StatementBegin
-- One row per Drive file ID; the album is derived from the folder and denormalised for listing speed.
CREATE TABLE IF NOT EXISTS photos
(
    drive_file_id     TEXT PRIMARY KEY,
    folder_id         TEXT REFERENCES photo_folders (drive_folder_id) ON DELETE SET NULL,
    album_id          INT REFERENCES photo_albums (id_photo_albums) ON DELETE SET NULL,
    name              TEXT        NOT NULL,
    mime_type         TEXT        NOT NULL,
    media_type        TEXT        NOT NULL CHECK (media_type IN ('image', 'video', 'other')),
    size_bytes        BIGINT      NOT NULL DEFAULT 0,
    md5               TEXT        NOT NULL DEFAULT '',
    width             INT         NOT NULL DEFAULT 0,
    height            INT         NOT NULL DEFAULT 0,
    taken_at          TIMESTAMPTZ,
    taken_at_source   TEXT        NOT NULL DEFAULT '',
    taken_at_suspect  BOOLEAN     NOT NULL DEFAULT FALSE,
    drive_created_at  TIMESTAMPTZ,
    drive_modified_at TIMESTAMPTZ,
    sort_at           TIMESTAMPTZ NOT NULL,
    sort_at_source    TEXT        NOT NULL,
    web_url           TEXT        NOT NULL,
    -- Exact duplicate (same md5): only the canonical file is listed. album_duplicate limits the hiding to
    -- album views where the canonical file sits in the same album.
    duplicate_of      TEXT,
    album_duplicate   BOOLEAN     NOT NULL DEFAULT FALSE,
    last_seen_run     BIGINT      NOT NULL,
    deleted_at        TIMESTAMPTZ
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_photos_gallery ON photos (sort_at DESC, drive_file_id DESC)
    WHERE deleted_at IS NULL AND duplicate_of IS NULL AND media_type <> 'other';
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_photos_album ON photos (album_id, sort_at, drive_file_id)
    WHERE deleted_at IS NULL AND media_type <> 'other';
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_photos_md5 ON photos (md5) WHERE md5 <> '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS photos;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS photo_folders;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS photo_albums;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS photo_sync_runs;
-- +goose StatementEnd
