-- +goose Up
-- +goose StatementBegin
-- Create table to track sent event notifications to avoid duplicates
CREATE TABLE IF NOT EXISTS event_notification_sent (
    id SERIAL PRIMARY KEY,
    id_events INTEGER NOT NULL,
    notification_type VARCHAR(50) NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(id_events, notification_type),
    FOREIGN KEY (id_events) REFERENCES events(id_events) ON DELETE CASCADE
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_event_notification_sent_event ON event_notification_sent(id_events);
CREATE INDEX IF NOT EXISTS idx_event_notification_sent_type ON event_notification_sent(notification_type);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop table and indexes
DROP INDEX IF EXISTS idx_event_notification_sent_type;
DROP INDEX IF EXISTS idx_event_notification_sent_event;
DROP TABLE IF EXISTS event_notification_sent;

-- +goose StatementEnd
