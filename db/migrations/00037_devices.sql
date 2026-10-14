-- +goose Up
-- +goose StatementBegin
-- One row per app installation, owned by whoever is signed in on it right now. The installation id is generated
-- by the app and survives logout, so a second account on the same phone takes the row over instead of adding one.
-- A NULL push_token means the installation cannot receive pushes (permission refused, token not obtained, signed out).
CREATE TABLE IF NOT EXISTS devices
(
    installation_id   UUID PRIMARY KEY,
    user_email        VARCHAR(100) NOT NULL REFERENCES newf (email) ON DELETE CASCADE,
    provider          TEXT         NOT NULL DEFAULT 'expo' CHECK (provider IN ('expo')),
    platform          TEXT         NOT NULL CHECK (platform IN ('ios', 'android')),
    push_token        TEXT,
    permission_status TEXT         NOT NULL DEFAULT 'unknown'
        CHECK (permission_status IN ('granted', 'denied', 'undetermined', 'unknown')),
    app_version       TEXT,
    status            TEXT         NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'disabled')),
    disabled_reason   TEXT CHECK (disabled_reason IN ('signed_out', 'token_replaced')),
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT now(),
    -- Last authenticated heartbeat: the only input of the connection state, see services/devices.
    last_seen_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    CHECK ((status = 'disabled') = (disabled_reason IS NOT NULL))
);
-- +goose StatementEnd
-- +goose StatementBegin
-- A push token addresses exactly one installation.
CREATE UNIQUE INDEX IF NOT EXISTS uq_devices_provider_token
    ON devices (provider, push_token)
    WHERE push_token IS NOT NULL;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_devices_user_seen
    ON devices (user_email, last_seen_at DESC)
    WHERE status = 'active';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS devices;
-- +goose StatementEnd
