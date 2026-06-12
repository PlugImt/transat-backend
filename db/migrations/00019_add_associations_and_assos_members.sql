-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS associations (
    id_associations SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    picture VARCHAR(500) NOT NULL,
    description VARCHAR(500),
    location VARCHAR(100),
    link VARCHAR(500)
);

CREATE TABLE IF NOT EXISTS associations_members (
    email VARCHAR(255) NOT NULL,
    id_associations INTEGER NOT NULL,
    PRIMARY KEY (email, id_associations),
    CONSTRAINT fk_association
    FOREIGN KEY (id_associations)
    REFERENCES associations(id_associations)
    ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS associations_members;
DROP TABLE IF EXISTS associations;
-- +goose StatementEnd


