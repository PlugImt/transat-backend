-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- First, drop the existing foreign key constraint
ALTER TABLE request_statistics DROP CONSTRAINT IF EXISTS request_statistics_user_email_fkey;

-- Then add a new constraint that properly handles NULL values
ALTER TABLE request_statistics ADD CONSTRAINT request_statistics_user_email_fkey
    FOREIGN KEY (user_email) REFERENCES newf(email) ON DELETE SET NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
-- Restore original constraint
ALTER TABLE request_statistics DROP CONSTRAINT IF EXISTS request_statistics_user_email_fkey;
ALTER TABLE request_statistics ADD CONSTRAINT request_statistics_user_email_fkey
    FOREIGN KEY (user_email) REFERENCES newf(email) ON DELETE SET NULL; 