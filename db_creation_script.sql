CREATE TABLE newf
(
    id_newf                      SERIAL,
    email                        VARCHAR(100),
    password                     VARCHAR(200) NOT NULL,
    password_updated_date        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verification_code            VARCHAR(6)            DEFAULT LPAD(FLOOR(RANDOM() * 1000000)::TEXT, 6, '0'),
    verification_code_expiration TIMESTAMP             DEFAULT CURRENT_TIMESTAMP + INTERVAL '10 minutes',
    creation_date                TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    first_name                   VARCHAR(50)  NOT NULL,
    last_name                    VARCHAR(50)  NOT NULL,
    phone_number                 VARCHAR(20),
    profile_picture              VARCHAR(500),
    notification_token           VARCHAR(50),
    graduation_year              SMALLINT,
    campus                       VARCHAR(50) CHECK (campus IN ('NANTES', 'BREST', 'RENNES') ),
    PRIMARY KEY (email)
);

CREATE TABLE roles
(
    id_roles    SERIAL,
    name        VARCHAR(50) NOT NULL,
    description VARCHAR(100),
    PRIMARY KEY (id_roles),
    UNIQUE (name)
);

CREATE TABLE games
(
    id_games    SERIAL,
    name        VARCHAR(50) NOT NULL,
    description VARCHAR(200),
    rules       VARCHAR(1000),
    picture     VARCHAR(500),
    PRIMARY KEY (id_games),
    UNIQUE (name)
);

CREATE TABLE caps
(
    id_caps           SERIAL,
    game_code         VARCHAR(6)   NOT NULL,
    number_of_players SMALLINT     NOT NULL,
    finished          BOOLEAN,
    creation_date     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finish_date       TIMESTAMP,
    creator           VARCHAR(100) NOT NULL,
    PRIMARY KEY (id_caps),
    FOREIGN KEY (creator) REFERENCES newf (email) ON DELETE CASCADE
);

CREATE TABLE bassine
(
    id_bassine    SERIAL,
    creation_date TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    creator       VARCHAR(100) NOT NULL,
    PRIMARY KEY (id_bassine),
    FOREIGN KEY (creator) REFERENCES newf (email) ON DELETE CASCADE
);

CREATE TABLE events
(
    id_events     SERIAL,
    name          VARCHAR(100) NOT NULL,
    description   VARCHAR(200),
    link          VARCHAR(100),
    start_date    TIMESTAMP    NOT NULL,
    end_date      TIMESTAMP,
    location      VARCHAR(100) NOT NULL,
    creation_date TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    picture       VARCHAR(500),
    slots         INTEGER,
    price         NUMERIC(15, 2),
    PRIMARY KEY (id_events)
);

CREATE TABLE rooms
(
    id_rooms    SERIAL,
    name        VARCHAR(100) NOT NULL,
    description VARCHAR(200),
    picture     VARCHAR(500),
    location    VARCHAR(100) NOT NULL,
    PRIMARY KEY (id_rooms),
    UNIQUE (name),
    UNIQUE (location)
);

CREATE TABLE articles_types
(
    id_articles_types SERIAL,
    name              VARCHAR(50) NOT NULL,
    PRIMARY KEY (id_articles_types),
    UNIQUE (name)
);

CREATE TABLE clubs_types
(
    id_clubs_types SERIAL,
    name           VARCHAR(50) NOT NULL,
    PRIMARY KEY (id_clubs_types),
    UNIQUE (name)
);

CREATE TABLE traq_types
(
    id_traq_types SERIAL,
    name          VARCHAR(100) NOT NULL,
    PRIMARY KEY (id_traq_types),
    UNIQUE (name)
);

CREATE TABLE washing_machines
(
    id_washing_machines SERIAL,
    type                VARCHAR(50) NOT NULL,
    broken              BOOLEAN     NOT NULL,
    PRIMARY KEY (id_washing_machines)
);

CREATE TABLE restaurant
(
    id_restaurant SERIAL,
    articles       VARCHAR(5000) NOT NULL,
    updated_date  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id_restaurant)
);

CREATE TABLE services
(
    id_services SERIAL,
    name        VARCHAR(50) NOT NULL,
    type        VARCHAR(50),
    PRIMARY KEY (id_services),
    UNIQUE (name)
);

CREATE TABLE traq
(
    id_traq       SERIAL,
    name          VARCHAR(50)    NOT NULL,
    disabled      BOOLEAN,
    limited       BOOLEAN,
    alcohol        NUMERIC(4, 2),
    out_of_stock  BOOLEAN,
    creation_date TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    picture       VARCHAR(500)   NOT NULL,
    description   VARCHAR(200),
    price         NUMERIC(10, 2) NOT NULL,
    price_half    NUMERIC(10, 2),
    id_traq_types INTEGER        NOT NULL,
    PRIMARY KEY (id_traq),
    FOREIGN KEY (id_traq_types) REFERENCES traq_types (id_traq_types) ON DELETE CASCADE
);

CREATE TABLE clubs
(
    id_clubs       SERIAL,
    name           VARCHAR(50)  NOT NULL,
    picture        VARCHAR(500),
    link           VARCHAR(100),
    subtitle       VARCHAR(150) NOT NULL,
    description    VARCHAR(500),
    location       VARCHAR(100) NOT NULL,
    practice_time  TIME,
    practice_day   VARCHAR(20),
    id_clubs_types INTEGER      NOT NULL,
    PRIMARY KEY (id_clubs),
    UNIQUE (name),
    FOREIGN KEY (id_clubs_types) REFERENCES clubs_types (id_clubs_types) ON DELETE CASCADE
);

CREATE TABLE articles
(
    id_articles       SERIAL,
    name              VARCHAR(100)   NOT NULL,
    price             NUMERIC(15, 2) NOT NULL,
    picture           VARCHAR(500)   NOT NULL,
    description       VARCHAR(200),
    id_articles_types INTEGER        NOT NULL,
    seller            VARCHAR(100)   NOT NULL,
    buyer             VARCHAR(100),
    creation_date     TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    sell_date         TIMESTAMP,
    PRIMARY KEY (id_articles),
    FOREIGN KEY (id_articles_types) REFERENCES articles_types (id_articles_types) ON DELETE CASCADE,
    FOREIGN KEY (seller) REFERENCES newf (email) ON DELETE CASCADE,
    FOREIGN KEY (buyer) REFERENCES newf (email) ON DELETE CASCADE
);

CREATE TABLE newf_roles
(
    email    VARCHAR(100),
    id_roles INTEGER,
    PRIMARY KEY (email, id_roles),
    FOREIGN KEY (email) REFERENCES newf (email) ON DELETE CASCADE,
    FOREIGN KEY (id_roles) REFERENCES roles (id_roles) ON DELETE CASCADE
);

CREATE TABLE players_caps
(
    email    VARCHAR(100),
    id_games INTEGER,
    id_caps  INTEGER,
    score    INTEGER NOT NULL,
    PRIMARY KEY (email, id_games, id_caps),
    FOREIGN KEY (email) REFERENCES newf (email) ON DELETE CASCADE,
    FOREIGN KEY (id_games) REFERENCES games (id_games) ON DELETE CASCADE,
    FOREIGN KEY (id_caps) REFERENCES caps (id_caps) ON DELETE CASCADE
);

CREATE TABLE players_bassine
(
    email      VARCHAR(100),
    id_games   INTEGER,
    id_bassine INTEGER,
    score      INTEGER,
    PRIMARY KEY (email, id_games, id_bassine),
    FOREIGN KEY (email) REFERENCES newf (email) ON DELETE CASCADE,
    FOREIGN KEY (id_games) REFERENCES games (id_games) ON DELETE CASCADE,
    FOREIGN KEY (id_bassine) REFERENCES bassine (id_bassine) ON DELETE CASCADE
);

CREATE TABLE clubs_respos
(
    email    VARCHAR(100),
    id_clubs INTEGER,
    PRIMARY KEY (email, id_clubs),
    FOREIGN KEY (email) REFERENCES newf (email) ON DELETE CASCADE,
    FOREIGN KEY (id_clubs) REFERENCES clubs (id_clubs) ON DELETE CASCADE
);

CREATE TABLE clubs_members
(
    email    VARCHAR(100),
    id_clubs INTEGER,
    PRIMARY KEY (email, id_clubs),
    FOREIGN KEY (email) REFERENCES newf (email) ON DELETE CASCADE,
    FOREIGN KEY (id_clubs) REFERENCES clubs (id_clubs) ON DELETE CASCADE
);

CREATE TABLE events_hosts
(
    email     VARCHAR(100),
    id_events INTEGER,
    PRIMARY KEY (email, id_events),
    FOREIGN KEY (email) REFERENCES newf (email) ON DELETE CASCADE,
    FOREIGN KEY (id_events) REFERENCES events (id_events) ON DELETE CASCADE
);

CREATE TABLE events_attendents
(
    email     VARCHAR(100),
    id_events INTEGER,
    PRIMARY KEY (email, id_events),
    FOREIGN KEY (email) REFERENCES newf (email) ON DELETE CASCADE,
    FOREIGN KEY (id_events) REFERENCES events (id_events) ON DELETE CASCADE
);

CREATE TABLE events_clubs
(
    id_clubs  INTEGER,
    id_events INTEGER,
    PRIMARY KEY (id_clubs, id_events),
    FOREIGN KEY (id_clubs) REFERENCES clubs (id_clubs) ON DELETE CASCADE,
    FOREIGN KEY (id_events) REFERENCES events (id_events) ON DELETE CASCADE
);

CREATE TABLE reservations
(
    email      VARCHAR(100),
    id_rooms   INTEGER,
    start_date TIMESTAMP NOT NULL,
    end_date   TIMESTAMP NOT NULL,
    PRIMARY KEY (email, id_rooms),
    FOREIGN KEY (email) REFERENCES newf (email) ON DELETE CASCADE,
    FOREIGN KEY (id_rooms) REFERENCES rooms (id_rooms) ON DELETE CASCADE
);

CREATE TABLE washing_machines_reports
(
    email               VARCHAR(100),
    id_washing_machines INTEGER,
    report_date         TIMESTAMP NOT NULL,
    PRIMARY KEY (email, id_washing_machines),
    FOREIGN KEY (email) REFERENCES newf (email) ON DELETE CASCADE,
    FOREIGN KEY (id_washing_machines) REFERENCES washing_machines (id_washing_machines) ON DELETE CASCADE
);

CREATE TABLE notifications
(
    email       VARCHAR(100),
    id_services INTEGER,
    PRIMARY KEY (email, id_services),
    FOREIGN KEY (email) REFERENCES newf (email) ON DELETE CASCADE,
    FOREIGN KEY (id_services) REFERENCES services (id_services) ON DELETE CASCADE
);

INSERT INTO roles (name, description)
VALUES ('ADMIN', 'Global administrator with all privileges'),
       ('NEWF', 'Transat user'),
       ('VALIDATED', 'Transat user with validated email'),
       ('BANNED', 'Transat user banned from the platform'),
       ('VERIFYING', 'Transat user with email not yet verified');