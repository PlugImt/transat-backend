-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS bassine_scores CASCADE;

CREATE TABLE bassine_history
(
    id_bassine_history SERIAL,
    email              VARCHAR(100) NOT NULL,
    date               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id_bassine_history),
    FOREIGN KEY (email) REFERENCES newf (email)
);

CREATE VIEW bassine_scores AS
SELECT
    email,
    COUNT(*) AS score
FROM
    bassine_history
GROUP BY
    email;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS bassine_scores CASCADE;
DROP TABLE IF EXISTS bassine_history CASCADE;
-- +goose StatementEnd