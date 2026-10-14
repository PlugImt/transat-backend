-- +goose Up
-- +goose StatementBegin
-- Tokens registered by app versions that predate device registration (PATCH /newf/me) become devices too, so
-- the delivery pipeline reads one table. They have no heartbeat: their session is read from newf.last_activity.
ALTER TABLE devices
    ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT 'client' CHECK (source IN ('client', 'legacy'));
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE devices
    DROP CONSTRAINT IF EXISTS devices_platform_check;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE devices
    ADD CONSTRAINT devices_platform_check CHECK (platform IN ('ios', 'android', 'unknown'));
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE devices
    DROP CONSTRAINT IF EXISTS devices_disabled_reason_check;
-- +goose StatementEnd
-- +goose StatementBegin
-- invalid_token: the push provider confirmed that the token is dead.
ALTER TABLE devices
    ADD CONSTRAINT devices_disabled_reason_check CHECK (disabled_reason IN ('signed_out', 'token_replaced', 'invalid_token'));
-- +goose StatementEnd
-- +goose StatementBegin
-- A token shared by several accounts (a phone used after a logout) goes to the one that registered it last.
INSERT INTO devices (installation_id, user_email, platform, push_token, permission_status, source, created_at, updated_at, last_seen_at)
SELECT gen_random_uuid(), x.email, 'unknown', x.token, 'unknown', 'legacy', x.created_at, now(), x.created_at
FROM (SELECT DISTINCT ON (token) email, token, created_at
      FROM user_notification_tokens
      WHERE token ~ '^(Expo|Exponent)PushToken\[[A-Za-z0-9_-]{1,200}\]$'
      ORDER BY token, created_at DESC) x
WHERE NOT EXISTS (SELECT 1 FROM devices d WHERE d.provider = 'expo' AND d.push_token = x.token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM devices WHERE source = 'legacy' OR platform = 'unknown';
-- +goose StatementEnd
-- +goose StatementBegin
UPDATE devices SET disabled_reason = 'token_replaced' WHERE disabled_reason = 'invalid_token';
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE devices DROP CONSTRAINT IF EXISTS devices_disabled_reason_check;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE devices ADD CONSTRAINT devices_disabled_reason_check CHECK (disabled_reason IN ('signed_out', 'token_replaced'));
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE devices DROP CONSTRAINT IF EXISTS devices_platform_check;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE devices ADD CONSTRAINT devices_platform_check CHECK (platform IN ('ios', 'android'));
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE devices DROP COLUMN IF EXISTS source;
-- +goose StatementEnd
