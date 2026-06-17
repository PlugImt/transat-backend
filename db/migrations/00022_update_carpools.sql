-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE carpools DROP CONSTRAINT IF EXISTS carpools_trip_type_check;
TRUNCATE TABLE carpools CASCADE;

ALTER TABLE carpools ADD CONSTRAINT carpools_trip_type_check 
    CHECK (trip_type IN ('SHOPPING', 'LONG_TRIP', 'OTHER'));


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE carpools DROP CONSTRAINT IF EXISTS carpools_trip_type_check;

ALTER TABLE carpools ADD CONSTRAINT carpools_trip_type_check 
    CHECK (trip_type IN ('SHOPPING', 'WEEKEND', 'OTHER'));