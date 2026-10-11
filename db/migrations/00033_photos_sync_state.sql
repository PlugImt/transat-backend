-- +goose Up
-- +goose StatementBegin
-- Single-row state: Drive changes checkpoint and the resume point of an interrupted full sync.
CREATE TABLE IF NOT EXISTS photo_sync_state
(
    id                    SMALLINT PRIMARY KEY CHECK (id = 1),
    changes_page_token    TEXT,
    last_full_at          TIMESTAMPTZ,
    last_incremental_at   TIMESTAMPTZ,
    full_run_id           BIGINT,
    full_phase            TEXT,
    full_page_token       TEXT,
    pending_changes_token TEXT,
    -- Bumped whenever indexed data changes: invalidates the API response caches of every instance.
    data_version          BIGINT NOT NULL DEFAULT 0
);
-- +goose StatementEnd
-- +goose StatementBegin
INSERT INTO photo_sync_state (id) VALUES (1) ON CONFLICT DO NOTHING;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE photo_sync_runs ADD COLUMN IF NOT EXISTS mode TEXT NOT NULL DEFAULT 'full';
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE photo_folders ADD COLUMN IF NOT EXISTS logical_path TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd
-- +goose StatementBegin
-- Drive's hasThumbnail: false for RAW files stored as octet-stream, which have no preview at all.
ALTER TABLE photos ADD COLUMN IF NOT EXISTS has_thumbnail BOOLEAN NOT NULL DEFAULT TRUE;
-- +goose StatementEnd
-- +goose StatementBegin
-- Album ordering key: event date, else first capture date, else epoch (undated and empty albums sink).
ALTER TABLE photo_albums ADD COLUMN IF NOT EXISTS sort_date DATE NOT NULL DEFAULT '1970-01-01';
-- +goose StatementEnd
-- +goose StatementBegin
UPDATE photo_albums SET sort_date = COALESCE(event_date, first_taken_at::date, '1970-01-01');
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_photo_albums_list ON photo_albums (sort_date DESC, id_photo_albums DESC)
    WHERE visibility = 'published' AND deleted_at IS NULL AND photo_count > 0;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_photos_folder ON photos (folder_id) WHERE deleted_at IS NULL;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_photo_folders_shortcut ON photo_folders (shortcut_id) WHERE shortcut_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_photo_folders_shortcut;
-- +goose StatementEnd
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_photos_folder;
-- +goose StatementEnd
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_photo_albums_list;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE photo_albums DROP COLUMN IF EXISTS sort_date;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE photos DROP COLUMN IF EXISTS has_thumbnail;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE photo_folders DROP COLUMN IF EXISTS logical_path;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE photo_sync_runs DROP COLUMN IF EXISTS mode;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS photo_sync_state;
-- +goose StatementEnd
