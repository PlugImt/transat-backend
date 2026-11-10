-- +goose Up
-- +goose StatementBegin
-- A file named "cover" (any extension, any case) in an album folder is that album's cover: shown on the album card,
-- never among the album's photos or in the gallery.
ALTER TABLE photos ADD COLUMN IF NOT EXISTS is_cover BOOLEAN
    GENERATED ALWAYS AS (lower(regexp_replace(name, '\.[^.]*$', '')) = 'cover') STORED;
CREATE INDEX IF NOT EXISTS idx_photos_album_cover ON photos (album_id) WHERE is_cover AND deleted_at IS NULL;
-- +goose StatementEnd
-- +goose StatementBegin
-- Albums already indexed take their cover file now instead of at their next change.
UPDATE photo_albums a SET
    cover_photo_id = c.drive_file_id,
    photo_count = GREATEST(a.photo_count - c.n, 0)
FROM (
    SELECT album_id, (ARRAY_AGG(drive_file_id ORDER BY drive_file_id))[1] AS drive_file_id, COUNT(*) AS n
    FROM photos
    WHERE is_cover AND deleted_at IS NULL AND media_type = 'image' AND album_id IS NOT NULL
    GROUP BY album_id
) c
WHERE a.id_photo_albums = c.album_id;
UPDATE photo_sync_state SET data_version = data_version + 1 WHERE id = 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_photos_album_cover;
ALTER TABLE photos DROP COLUMN IF EXISTS is_cover;
-- +goose StatementEnd
