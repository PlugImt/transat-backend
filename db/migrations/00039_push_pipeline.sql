-- +goose Up
-- +goose StatementBegin
-- One logical notification. The audience is stored, not the recipients: they are worked out when the message is
-- fanned out, and every device is checked again right before it is sent. dedupe_key makes enqueueing idempotent.
--   pending -> fanned_out -> completed, or cancelled when it expired before it was fanned out.
CREATE TABLE IF NOT EXISTS push_messages
(
    id               BIGSERIAL PRIMARY KEY,
    kind             TEXT        NOT NULL,
    -- The notification category whose per-user preference applies; empty when none does.
    category         TEXT        NOT NULL DEFAULT '',
    dedupe_key       TEXT UNIQUE,
    audience         JSONB       NOT NULL,
    -- {"fr": {"title": "...", "body": "..."}, "en": {...}}: the text of each language, "fr" being the fallback.
    content          JSONB       NOT NULL,
    navigation       JSONB,
    created_by       TEXT,
    not_before       TIMESTAMPTZ NOT NULL DEFAULT now(),
    -- A notification still queued after this is dropped, and the providers are told not to keep it longer.
    expires_at       TIMESTAMPTZ NOT NULL,
    status           TEXT        NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'fanned_out', 'completed', 'cancelled')),
    -- Outcome of the fan-out: users in the audience, devices that may be pushed to, and why the others may not.
    recipients       INT,
    eligible_devices INT,
    filtered         JSONB,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    fanned_out_at    TIMESTAMPTZ,
    completed_at     TIMESTAMPTZ
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_push_messages_pending ON push_messages (not_before) WHERE status = 'pending';
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_push_messages_open ON push_messages (id) WHERE status = 'fanned_out';
-- +goose StatementEnd
-- +goose StatementBegin
-- One row per message and device: the delivery log. The token is never stored here, only a fingerprint.
--   pending    waiting for an attempt or a retry (next_attempt_at)
--   sending    claimed by a worker (lease_until)
--   accepted   the push provider accepted it and gave a ticket; nothing says it reached the phone
--   handed_off the provider's receipt says APNs/FCM took it; still not a delivery confirmation
--   failed     permanent failure, or the attempts ran out (error_code)
--   skipped    the device was no longer eligible when the message was due (skip_reason)
--   expired    the message expired before it was sent
--   unconfirmed accepted, but no receipt could be read
-- No status says "delivered": Expo, APNs and FCM give no reliable evidence of it.
CREATE TABLE IF NOT EXISTS push_deliveries
(
    id               BIGSERIAL PRIMARY KEY,
    message_id       BIGINT       NOT NULL REFERENCES push_messages (id) ON DELETE CASCADE,
    device_id        UUID         NOT NULL REFERENCES devices (installation_id) ON DELETE CASCADE,
    -- The account the device belonged to at fan-out: a device handed to another account is skipped.
    user_email       VARCHAR(100) NOT NULL,
    language         TEXT         NOT NULL,
    status           TEXT         NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'sending', 'accepted', 'handed_off', 'failed', 'skipped', 'expired', 'unconfirmed')),
    attempts         INT          NOT NULL DEFAULT 0,
    next_attempt_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    lease_until      TIMESTAMPTZ,
    token_fingerprint TEXT,
    ticket_id        TEXT,
    error_code       TEXT,
    error_detail     TEXT,
    skip_reason      TEXT,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    first_attempt_at TIMESTAMPTZ,
    accepted_at      TIMESTAMPTZ,
    receipt_due_at   TIMESTAMPTZ,
    receipt_at       TIMESTAMPTZ,
    finished_at      TIMESTAMPTZ,
    UNIQUE (message_id, device_id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_push_deliveries_due ON push_deliveries (next_attempt_at) WHERE status = 'pending';
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_push_deliveries_leased ON push_deliveries (lease_until) WHERE status = 'sending';
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_push_deliveries_receipt ON push_deliveries (receipt_due_at) WHERE status = 'accepted';
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_push_deliveries_message ON push_deliveries (message_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS push_deliveries;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS push_messages;
-- +goose StatementEnd
