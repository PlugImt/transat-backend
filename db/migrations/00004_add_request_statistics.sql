-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE request_statistics (
    id SERIAL PRIMARY KEY,
    user_email VARCHAR(100) NULL, -- Can be NULL for unauthenticated requests
    endpoint VARCHAR(200) NOT NULL, -- The endpoint that was requested
    method VARCHAR(10) NOT NULL, -- HTTP method (GET, POST, etc.)
    request_received TIMESTAMP NOT NULL, -- When the request was received
    response_sent TIMESTAMP NOT NULL, -- When the response was sent
    status_code INTEGER NOT NULL, -- HTTP status code of the response
    duration_ms INTEGER NOT NULL, -- Duration in milliseconds
    FOREIGN KEY (user_email) REFERENCES newf(email) ON DELETE SET NULL
);

-- Create an index on user_email and endpoint for faster lookups
CREATE INDEX idx_request_statistics_user_email ON request_statistics(user_email);
CREATE INDEX idx_request_statistics_endpoint ON request_statistics(endpoint);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS request_statistics; 