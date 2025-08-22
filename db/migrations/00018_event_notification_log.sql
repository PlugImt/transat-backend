-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_notification_sends (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL,
    email VARCHAR(100) NOT NULL,
    notification_type VARCHAR(50) NOT NULL,
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_event_notification_event FOREIGN KEY (event_id) REFERENCES events(id_events) ON DELETE CASCADE,
    CONSTRAINT fk_event_notification_user FOREIGN KEY (email) REFERENCES newf(email) ON DELETE CASCADE,
    CONSTRAINT uq_event_notification_send UNIQUE (event_id, email, notification_type)
);
-- +goose StatementEnd
