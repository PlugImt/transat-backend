-- +goose Up
-- +goose StatementBegin
-- Truncate the table to remove all data
DROP TABLE IF EXISTS restaurant CASCADE;

CREATE TABLE restaurant
(
    id_restaurant SERIAL UNIQUE,
    type          VARCHAR(50),
    PRIMARY KEY (id_restaurant, type)
);

CREATE TABLE restaurant_articles
(
    id_restaurant_articles SERIAL,
    first_time_served      DATE         NOT NULL DEFAULT CURRENT_DATE,
    last_time_served       DATE,
    name                   VARCHAR(500) NOT NULL,
    PRIMARY KEY (id_restaurant_articles),
    UNIQUE (name)
);

CREATE TABLE restaurant_articles_notes
(
    email                  VARCHAR(100),
    id_restaurant_articles INTEGER,
    note                   SMALLINT NOT NULL,
    comment                VARCHAR(500),
    PRIMARY KEY (email, id_restaurant_articles),
    FOREIGN KEY (email) REFERENCES newf (email),
    FOREIGN KEY (id_restaurant_articles) REFERENCES restaurant_articles (id_restaurant_articles),
    FOREIGN KEY (email) REFERENCES newf (email)
);

CREATE TABLE restaurant_meals
(
    id_restaurant          INTEGER,
    id_restaurant_articles INTEGER,
    date_served            DATE NOT NULL DEFAULT CURRENT_DATE,
    PRIMARY KEY (id_restaurant, id_restaurant_articles, date_served),
    FOREIGN KEY (id_restaurant) REFERENCES restaurant (id_restaurant),
    FOREIGN KEY (id_restaurant_articles) REFERENCES restaurant_articles (id_restaurant_articles)
);

INSERT INTO restaurant (type)
VALUES ('grilladesMidi'),
       ('migrateurs'),
       ('cibo'),
       ('accompMidi'),
       ('grilladesSoir'),
       ('accompSoir');

-- +goose StatementEnd
