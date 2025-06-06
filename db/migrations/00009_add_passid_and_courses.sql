-- +goose Up
-- +goose StatementBegin
ALTER TABLE newf ADD COLUMN pass_id INTEGER;

CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    title VARCHAR(255) NOT NULL,
    start_time DATE NOT NULL,
    end_time DATE NOT NULL,
    teacher VARCHAR(255),
    room VARCHAR(255),
    "group" VARCHAR(255),
    user_email VARCHAR(255) REFERENCES newf(email) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS courses;
ALTER TABLE newf DROP COLUMN IF EXISTS pass_id;
-- +goose StatementEnd
