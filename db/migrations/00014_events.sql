-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS events_hosts CASCADE;
DROP TABLE IF EXISTS events_clubs CASCADE;

ALTER TABLE events
    DROP COLUMN IF EXISTS slots,
    DROP COLUMN IF EXISTS price,
    ADD COLUMN creator VARCHAR(255) NOT NULL,
    ADD CONSTRAINT fk_creator FOREIGN KEY (creator) REFERENCES newf(email) ON DELETE CASCADE,
    ADD COLUMN id_club INT NOT NULL,
    ADD CONSTRAINT fk_club_id FOREIGN KEY (id_club) REFERENCES clubs(id_clubs) ON DELETE CASCADE;
-- +goose StatementEnd
