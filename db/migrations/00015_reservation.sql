-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS reservations CASCADE;

CREATE TABLE reservation_category
(
    id_reservation_category SERIAL,
    name                    VARCHAR(100),
    id_parent_category      INTEGER NOT NULL,
    id_clubs                INTEGER NOT NULL,
    PRIMARY KEY (id_reservation_category),
    FOREIGN KEY (id_clubs) REFERENCES clubs (id_clubs)
);


CREATE TABLE reservation_element
(
    id_reservation_element  SERIAL,
    name                     VARCHAR(100),
    description             VARCHAR(100),
    location                VARCHAR(50),
    slot                 BOOLEAN,
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
    PRIMARY KEY (email, id_reservation_element),
    FOREIGN KEY (email) REFERENCES newf (email),
    FOREIGN KEY (id_reservation_element) REFERENCES reservation_element (id_reservation_element)
);

-- +goose StatementEnd
