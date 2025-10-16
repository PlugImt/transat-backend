-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_notifications_log
(
    id_event_notifications_log SERIAL,
    id_events                  INTEGER      NOT NULL,
    email                      VARCHAR(100) NOT NULL,
    notification_type          VARCHAR(50)  NOT NULL, -- 'CREATED' | 'REMINDER_1H'
    sent_at                    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id_event_notifications_log),
    CONSTRAINT fk_event_notifications_event FOREIGN KEY (id_events) REFERENCES events (id_events) ON DELETE CASCADE,
    CONSTRAINT fk_event_notifications_user FOREIGN KEY (email) REFERENCES newf (email) ON DELETE CASCADE,
    CONSTRAINT uq_event_notification_once UNIQUE (id_events, email, notification_type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event_notifications_log;
-- +goose StatementEnd


