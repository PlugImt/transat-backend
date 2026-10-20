-- +goose Up
-- +goose StatementBegin
-- Until when the person holding the baton (notified_at set) may answer; after that the scheduler passes the turn on.
ALTER TABLE reservation_waitlist ADD COLUMN IF NOT EXISTS turn_expires_at TIMESTAMPTZ;
-- Entries already holding the baton when this runs get a fresh half hour.
UPDATE reservation_waitlist SET turn_expires_at = now() + interval '30 minutes' WHERE notified_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS reservation_waitlist_turn_expiry
    ON reservation_waitlist (turn_expires_at) WHERE turn_expires_at IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS reservation_waitlist_turn_expiry;
ALTER TABLE reservation_waitlist DROP COLUMN IF EXISTS turn_expires_at;
-- +goose StatementEnd
