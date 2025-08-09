-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bassine_scores (
    email VARCHAR(100) PRIMARY KEY REFERENCES newf(email) ON DELETE CASCADE,
    score INTEGER NOT NULL DEFAULT 0 CHECK (score >= 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bassine_scores;
-- +goose StatementEnd


