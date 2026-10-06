-- +goose Up
-- +goose StatementBegin
INSERT INTO services (name)
VALUES ('RESTAURANT'),
       ('TRAQ'),
       ('EVENTS'),
       ('EVENT_REMINDERS'),
       ('RESERVATIONS')
ON CONFLICT (name) DO NOTHING;

-- Events and reservations used to notify everyone; keep that until each user opts out.
INSERT INTO notifications (email, id_services)
SELECT n.email, s.id_services
FROM newf n
         CROSS JOIN services s
WHERE s.name IN ('EVENTS', 'EVENT_REMINDERS', 'RESERVATIONS')
ON CONFLICT (email, id_services) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM services WHERE name IN ('EVENTS', 'EVENT_REMINDERS', 'RESERVATIONS');
-- +goose StatementEnd
