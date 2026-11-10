-- +goose Up
-- +goose StatementBegin
-- A short-lived link a booth shows as a QR code so the organisers can send it one frame per layout from their
-- phone or computer. The layouts' geometry comes from the booth, so the page can draw exact templates.
CREATE TABLE IF NOT EXISTS photobooth_frame_sessions (
    id            BIGSERIAL PRIMARY KEY,
    device_key_id TEXT        NOT NULL,
    token         TEXT        NOT NULL UNIQUE,
    event_name    TEXT        NOT NULL,
    layouts       JSONB       NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at    TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_photobooth_frame_sessions_expiry ON photobooth_frame_sessions (expires_at);

-- The frame sent for one layout of a session; a new upload replaces it. Removed with the session.
CREATE TABLE IF NOT EXISTS photobooth_frame_uploads (
    session_id  BIGINT      NOT NULL REFERENCES photobooth_frame_sessions (id) ON DELETE CASCADE,
    layout_type TEXT        NOT NULL,
    png         BYTEA       NOT NULL,
    sha256      TEXT        NOT NULL,
    width       INT         NOT NULL,
    height      INT         NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (session_id, layout_type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS photobooth_frame_uploads;
DROP TABLE IF EXISTS photobooth_frame_sessions;
-- +goose StatementEnd
