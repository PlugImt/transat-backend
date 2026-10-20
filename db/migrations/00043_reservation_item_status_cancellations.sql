-- +goose Up
-- +goose StatementBegin
-- active: bookable and listed in the app; maintenance and inactive: hidden from the app and not bookable.
ALTER TABLE reservation_element ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active';
ALTER TABLE reservation_element ADD COLUMN IF NOT EXISTS status_note VARCHAR(200);
ALTER TABLE reservation_element
    ADD CONSTRAINT reservation_element_status_check CHECK (status IN ('active', 'maintenance', 'inactive'));

-- Slot reservations are deleted when cancelled, so the history of what was cancelled (and of loans an
-- administrator marked as returned) lives here.
CREATE TABLE reservation_cancellations
(
    id                     BIGSERIAL PRIMARY KEY,
    email                  VARCHAR(100) NOT NULL,
    id_reservation_element INTEGER      NOT NULL REFERENCES reservation_element (id_reservation_element) ON DELETE CASCADE,
    start_date             TIMESTAMP    NOT NULL,
    end_date               TIMESTAMP,
    kind                   VARCHAR(20)  NOT NULL CHECK (kind IN ('cancelled', 'cancelled_by_admin', 'returned_by_admin')),
    cancelled_by           VARCHAR(100),
    cancelled_at           TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX reservation_cancellations_item ON reservation_cancellations (id_reservation_element, cancelled_at DESC);
CREATE INDEX reservation_cancellations_recent ON reservation_cancellations (cancelled_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reservation_cancellations;
ALTER TABLE reservation_element DROP CONSTRAINT IF EXISTS reservation_element_status_check;
ALTER TABLE reservation_element DROP COLUMN IF EXISTS status_note;
ALTER TABLE reservation_element DROP COLUMN IF EXISTS status;
-- +goose StatementEnd
