-- +goose Up
-- +goose StatementBegin
-- A hidden category disappears from the app with everything inside it (sub-categories and items), without touching
-- the items' own status: showing the category again brings them back as they were.
ALTER TABLE reservation_category ADD COLUMN IF NOT EXISTS hidden BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE VIEW reservation_hidden_categories AS
WITH RECURSIVE h AS (
    SELECT id_reservation_category AS id FROM reservation_category WHERE hidden
    UNION
    SELECT c.id_reservation_category FROM reservation_category c JOIN h ON c.id_parent_category = h.id
)
SELECT id FROM h;
-- +goose StatementEnd
-- +goose StatementBegin
-- The status that decides whether an item can be booked: its own, or inactive when one of its categories is hidden.
CREATE OR REPLACE FUNCTION reservation_item_status(item INT) RETURNS TEXT
    LANGUAGE sql STABLE AS
$$
SELECT CASE
           WHEN re.id_reservation_category IN (SELECT id FROM reservation_hidden_categories) THEN 'inactive'
           ELSE re.status
       END
FROM reservation_element re
WHERE re.id_reservation_element = item
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS reservation_item_status(INT);
DROP VIEW IF EXISTS reservation_hidden_categories;
ALTER TABLE reservation_category DROP COLUMN IF EXISTS hidden;
-- +goose StatementEnd
