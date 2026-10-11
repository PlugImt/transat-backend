-- +goose Up
-- +goose StatementBegin
-- Gallery order: newest first, then events kept together when several share a timestamp.
CREATE INDEX IF NOT EXISTS idx_photos_gallery_event
    ON photos (sort_at DESC, (COALESCE(album_id, 0)) DESC, drive_file_id DESC)
    WHERE deleted_at IS NULL AND duplicate_of IS NULL AND media_type <> 'other';
-- +goose StatementEnd
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_photos_gallery;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_photos_gallery
    ON photos (sort_at DESC, drive_file_id DESC)
    WHERE deleted_at IS NULL AND duplicate_of IS NULL AND media_type <> 'other';
-- +goose StatementEnd
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_photos_gallery_event;
-- +goose StatementEnd
