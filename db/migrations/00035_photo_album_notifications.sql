-- +goose Up
-- +goose StatementBegin
-- Opt-in category: nobody is subscribed here, neither existing users nor new ones (see models.OptInNotificationCategories).
INSERT INTO services (name)
VALUES ('PHOTOS')
ON CONFLICT (name) DO NOTHING;
-- +goose StatementEnd
-- +goose StatementBegin
-- One row per album folder that was, or is about to be, announced. It is keyed by the Drive folder id, which never
-- changes when an album is renamed or re-indexed, so an album can only ever be announced once.
--   baseline: existed when notifications were switched on; never announced
--   pending:  new and available; announced once its photo count has been stable for the settle delay
--   sending:  claimed by one process (leased until claimed_at + lease)
--   sent / failed: final
CREATE TABLE IF NOT EXISTS photo_album_notifications
(
    root_folder_id   TEXT PRIMARY KEY,
    album_id         INT REFERENCES photo_albums (id_photo_albums) ON DELETE SET NULL,
    status           TEXT        NOT NULL CHECK (status IN ('baseline', 'pending', 'sending', 'sent', 'failed')),
    first_seen_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_count       INT         NOT NULL DEFAULT 0,
    count_changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    attempts         INT         NOT NULL DEFAULT 0,
    claimed_at       TIMESTAMPTZ,
    retry_at         TIMESTAMPTZ,
    notified_at      TIMESTAMPTZ,
    devices          INT,
    last_error       TEXT
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_photo_album_notifications_open
    ON photo_album_notifications (status)
    WHERE status IN ('pending', 'sending');
-- +goose StatementEnd
-- +goose StatementBegin
-- Single row: set once the existing albums have been recorded as baseline, after the first complete sync.
CREATE TABLE IF NOT EXISTS photo_notify_state
(
    id           SMALLINT PRIMARY KEY CHECK (id = 1),
    baselined_at TIMESTAMPTZ
);
-- +goose StatementEnd
-- +goose StatementBegin
INSERT INTO photo_notify_state (id) VALUES (1) ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS photo_notify_state;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS photo_album_notifications;
-- +goose StatementEnd
-- +goose StatementBegin
DELETE FROM services WHERE name = 'PHOTOS';
-- +goose StatementEnd
