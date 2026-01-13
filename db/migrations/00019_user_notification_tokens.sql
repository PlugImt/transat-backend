-- +goose Up
-- Migration: Create user_notification_tokens table for multiple tokens per user
-- This allows users to have multiple notification tokens (e.g., multiple devices)

CREATE TABLE IF NOT EXISTS user_notification_tokens (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (email) REFERENCES newf(email) ON DELETE CASCADE,
    UNIQUE(email, token) -- Prevent duplicate tokens for the same user
);

-- Create index for faster lookups by email
CREATE INDEX idx_user_notification_tokens_email ON user_notification_tokens(email);

-- Create index for faster lookups by token
CREATE INDEX idx_user_notification_tokens_token ON user_notification_tokens(token);

-- Migrate existing notification tokens from newf table to the new table
-- Only migrate non-null tokens
INSERT INTO user_notification_tokens (email, token, created_at)
SELECT email, notification_token, CURRENT_TIMESTAMP
FROM newf
WHERE notification_token IS NOT NULL AND notification_token != ''
ON CONFLICT (email, token) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS user_notification_tokens;
