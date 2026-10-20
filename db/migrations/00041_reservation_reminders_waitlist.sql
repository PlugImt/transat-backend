-- +goose Up
-- +goose StatementBegin
-- When a reservation was made. Rows that exist today stay NULL ("before the reminder time"), new rows get now().
ALTER TABLE reservation ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;
ALTER TABLE reservation ALTER COLUMN created_at SET DEFAULT now();

-- One row per reminder already queued for a reservation: the primary key is what keeps the scheduler from ever
-- sending the same milestone twice. Rows go away with the reservation (cancelled, or the item deleted).
CREATE TABLE reservation_reminders
(
    email                  VARCHAR(100) NOT NULL,
    id_reservation_element INTEGER      NOT NULL,
    start_date             TIMESTAMP    NOT NULL,
    milestone              VARCHAR(20)  NOT NULL,
    sent_at                TIMESTAMPTZ  NOT NULL DEFAULT now(),
    PRIMARY KEY (email, id_reservation_element, start_date, milestone),
    FOREIGN KEY (email, id_reservation_element, start_date)
        REFERENCES reservation (email, id_reservation_element, start_date) ON DELETE CASCADE ON UPDATE CASCADE
);

-- People interested in an item (slot_start NULL) or in one slot of an item (slot_start set) that is already
-- reserved. The queue is first come first served, ordered by id. notified_at is set while the user holds the
-- baton: they were told the resource is free and have not answered yet.
CREATE TABLE reservation_waitlist
(
    id                     BIGSERIAL PRIMARY KEY,
    email                  VARCHAR(100) NOT NULL REFERENCES newf (email) ON DELETE CASCADE ON UPDATE CASCADE,
    id_reservation_element INTEGER      NOT NULL REFERENCES reservation_element (id_reservation_element) ON DELETE CASCADE,
    slot_start             TIMESTAMP,
    created_at             TIMESTAMPTZ  NOT NULL DEFAULT now(),
    notified_at            TIMESTAMPTZ
);

CREATE UNIQUE INDEX reservation_waitlist_item_unique
    ON reservation_waitlist (email, id_reservation_element) WHERE slot_start IS NULL;
CREATE UNIQUE INDEX reservation_waitlist_slot_unique
    ON reservation_waitlist (email, id_reservation_element, slot_start) WHERE slot_start IS NOT NULL;
CREATE INDEX reservation_waitlist_queue
    ON reservation_waitlist (id_reservation_element, slot_start, id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reservation_waitlist;
DROP TABLE IF EXISTS reservation_reminders;
ALTER TABLE reservation DROP COLUMN IF EXISTS created_at;
-- +goose StatementEnd
