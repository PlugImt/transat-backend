-- +goose Up
-- +goose StatementBegin
-- The club that owns an album's photos. Drive folders say nothing reliable about clubs, so the link is set
-- by an administrator; sync never touches this column. Photos of an album without club go to administrators only.
ALTER TABLE photo_albums
    ADD COLUMN IF NOT EXISTS club_id INT REFERENCES clubs (id_clubs) ON DELETE SET NULL;
-- +goose StatementEnd
-- +goose StatementBegin
-- A report is a review request, never a moderation decision. Photo, event and club are copied at report time:
-- the report stays readable after a re-index, a rename or a deletion in Drive.
CREATE TABLE IF NOT EXISTS photo_reports
(
    id_photo_reports          BIGSERIAL PRIMARY KEY,
    photo_id                  TEXT        NOT NULL,
    photo_name                TEXT        NOT NULL,
    photo_url                 TEXT        NOT NULL,
    album_id                  INT REFERENCES photo_albums (id_photo_albums) ON DELETE SET NULL,
    album_title               TEXT        NOT NULL DEFAULT '',
    album_url                 TEXT        NOT NULL DEFAULT '',
    club_id                   INT REFERENCES clubs (id_clubs) ON DELETE SET NULL,
    club_name                 TEXT        NOT NULL DEFAULT '',
    reason                    TEXT        NOT NULL CHECK (reason IN ('appears_in_photo', 'inappropriate', 'poor_quality', 'other')),
    comment                   TEXT        NOT NULL DEFAULT '' CHECK (char_length(comment) <= 1000),
    -- Email goes with the account; the name and ids stay so the report can still be investigated.
    reporter_email            VARCHAR(100) REFERENCES newf (email) ON DELETE SET NULL,
    reporter_id               INT         NOT NULL,
    reporter_public_id        UUID,
    reporter_name             TEXT        NOT NULL,
    reporter_graduation_year  SMALLINT,
    reporter_formation        TEXT        NOT NULL DEFAULT '',
    -- sha256(reporter, photo, reason, normalised comment): the same report within 24h is not stored twice.
    dedupe_key                TEXT        NOT NULL,
    created_at                TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_photo_reports_dedupe ON photo_reports (dedupe_key, created_at DESC);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_photo_reports_reporter ON photo_reports (reporter_email, created_at DESC);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_photo_reports_created ON photo_reports (created_at DESC, id_photo_reports DESC);
-- +goose StatementEnd
-- +goose StatementBegin
-- One row per recipient, resolved on the server when the report is received. A row becomes 'sent' only after
-- the email service accepted the message, so a retry never repeats a mail that already left.
--   pending -> sending (leased) -> sent | failed; a failed row below the attempt limit is retried at next_attempt_at.
CREATE TABLE IF NOT EXISTS photo_report_deliveries
(
    report_id       BIGINT       NOT NULL REFERENCES photo_reports (id_photo_reports) ON DELETE CASCADE,
    recipient_email VARCHAR(100) NOT NULL,
    recipient_role  TEXT         NOT NULL CHECK (recipient_role IN ('club_owner', 'admin')),
    language        TEXT         NOT NULL DEFAULT 'fr',
    status          TEXT         NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sending', 'sent', 'failed')),
    attempts        INT          NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    lease_until     TIMESTAMPTZ,
    sent_at         TIMESTAMPTZ,
    last_error      TEXT,
    PRIMARY KEY (report_id, recipient_email)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_photo_report_deliveries_open
    ON photo_report_deliveries (next_attempt_at)
    WHERE status <> 'sent';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS photo_report_deliveries;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS photo_reports;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE photo_albums
    DROP COLUMN IF EXISTS club_id;
-- +goose StatementEnd
