-- +goose Up
-- +goose StatementBegin
-- The booth's layouts as its app draws them, one set per photo ratio. The booths publish it (their app owns the
-- geometry); the admin panel draws the templates and checks the frames against it.
CREATE TABLE IF NOT EXISTS photobooth_layout_catalog (
    photo_aspect       TEXT        NOT NULL,
    layout_type        TEXT        NOT NULL,
    position           INT         NOT NULL,
    enabled_by_default BOOLEAN     NOT NULL DEFAULT FALSE,
    layout             JSONB       NOT NULL,
    device_key_id      TEXT        NOT NULL,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (photo_aspect, layout_type)
);

-- Events prepared in the admin panel; the booths download them with their settings, layouts and frames.
CREATE TABLE IF NOT EXISTS photobooth_events (
    id                     BIGSERIAL PRIMARY KEY,
    name                   TEXT        NOT NULL,
    event_date             DATE        NOT NULL,
    photo_aspect           TEXT        NOT NULL DEFAULT 'RATIO_4_3',
    countdown_seconds      INT         NOT NULL DEFAULT 5,
    max_prints_per_session INT         NOT NULL DEFAULT 5,
    allow_retakes          BOOLEAN     NOT NULL DEFAULT TRUE,
    digital_share          BOOLEAN     NOT NULL DEFAULT TRUE,
    attract_slideshow      BOOLEAN     NOT NULL DEFAULT TRUE,
    mirror_live_preview    BOOLEAN     NOT NULL DEFAULT TRUE,
    created_by             TEXT        NOT NULL DEFAULT '',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_photobooth_events_date ON photobooth_events (event_date);

-- One layout of an event: switched on or off, and the frame drawn over it (a PNG of the layout's shape).
CREATE TABLE IF NOT EXISTS photobooth_event_layouts (
    event_id          BIGINT      NOT NULL REFERENCES photobooth_events (id) ON DELETE CASCADE,
    layout_type       TEXT        NOT NULL,
    enabled           BOOLEAN     NOT NULL DEFAULT TRUE,
    frame_png         BYTEA,
    frame_sha256      TEXT,
    frame_width       INT,
    frame_height      INT,
    frame_uploaded_at TIMESTAMPTZ,
    PRIMARY KEY (event_id, layout_type)
);

-- A QR-code link opened for one of these events also files the frames it receives in the event.
ALTER TABLE photobooth_frame_sessions
    ADD COLUMN IF NOT EXISTS event_id BIGINT REFERENCES photobooth_events (id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE photobooth_frame_sessions DROP COLUMN IF EXISTS event_id;
DROP TABLE IF EXISTS photobooth_event_layouts;
DROP TABLE IF EXISTS photobooth_events;
DROP TABLE IF EXISTS photobooth_layout_catalog;
-- +goose StatementEnd
