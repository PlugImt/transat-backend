-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE carpools
(
    id_carpools     SERIAL,
    creator_email   VARCHAR(100) NOT NULL,
    trip_type       VARCHAR(20)  NOT NULL CHECK (trip_type IN ('SHOPPING', 'WEEKEND', 'OTHER')),
    departure_place VARCHAR(100) NOT NULL,
    destination     VARCHAR(100) NOT NULL,
    departure_time  TIMESTAMP    NOT NULL,
    contact_details VARCHAR(100) NOT NULL,
    description     VARCHAR(500),
    status          VARCHAR(20)  NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'FULL', 'ARCHIVED')),
    creation_date   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_date    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id_carpools),
    FOREIGN KEY (creator_email) REFERENCES newf (email) ON DELETE CASCADE
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE carpools;