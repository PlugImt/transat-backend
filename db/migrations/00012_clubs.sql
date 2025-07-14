-- +goose Up
-- +goose StatementBegin
-- Truncate the table to remove all data
DROP TABLE IF EXISTS clubs_types CASCADE;
DROP TABLE IF EXISTS clubs CASCADE;
DROP TABLE IF EXISTS clubs_respos CASCADE;
DROP TABLE IF EXISTS clubs_members CASCADE;

CREATE TABLE clubs
(
    id_clubs       SERIAL,
    name           VARCHAR(50)  NOT NULL,
    picture        VARCHAR(500) NOT NULL,
    description    VARCHAR(500),
    location       VARCHAR(100),
    link           VARCHAR(500),
    PRIMARY KEY (id_clubs),
    UNIQUE (name)
);

CREATE TABLE clubs_members
(
    email    VARCHAR(100),
    id_clubs INTEGER,
    PRIMARY KEY (email, id_clubs),
    FOREIGN KEY (email) REFERENCES newf (email),
    FOREIGN KEY (id_clubs) REFERENCES clubs (id_clubs)
);


-- +goose StatementEnd
