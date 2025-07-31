-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS reservations CASCADE;

CREATE TABLE reservation_category
(
    id_reservation_category SERIAL,
    name                    VARCHAR(100),
    id_parent_category      INTEGER,
    id_clubs                INTEGER NOT NULL,
    PRIMARY KEY (id_reservation_category),
    FOREIGN KEY (id_clubs) REFERENCES clubs (id_clubs)
);


CREATE TABLE reservation_element
(
    id_reservation_element  SERIAL,
    name                    VARCHAR(100),
    description             VARCHAR(100),
    location                VARCHAR(50),
    slot                    BOOLEAN,
    id_clubs                INTEGER NOT NULL,
    id_reservation_category INTEGER,
    PRIMARY KEY (id_reservation_element),
    FOREIGN KEY (id_clubs) REFERENCES clubs (id_clubs),
    FOREIGN KEY (id_reservation_category) REFERENCES reservation_category (id_reservation_category)
);
CREATE TABLE reservation
(
    email                  VARCHAR(100),
    id_reservation_element INTEGER,
    start_date             TIMESTAMP,
    end_date               TIMESTAMP,
    CONSTRAINT valid_date_range CHECK (start_date < end_date),
    CONSTRAINT non_null_dates CHECK (start_date IS NOT NULL AND end_date IS NOT NULL),
    PRIMARY KEY (email, id_reservation_element, start_date),
    FOREIGN KEY (email) REFERENCES newf (email),
    FOREIGN KEY (id_reservation_element) REFERENCES reservation_element (id_reservation_element) ON DELETE CASCADE
);

CREATE INDEX idx_reservation_element_dates ON reservation (id_reservation_element, start_date, end_date);
CREATE INDEX idx_reservation_user_dates ON reservation (email, start_date, end_date);

-- +goose StatementEnd
