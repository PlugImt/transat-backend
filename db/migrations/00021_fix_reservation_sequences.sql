-- +goose Up
-- +goose StatementBegin
-- Fix sequences for reservation tables to prevent duplicate key errors
-- This ensures the sequences are set to the maximum existing ID + 1

SELECT setval(
    pg_get_serial_sequence('reservation_category', 'id_reservation_category'),
    COALESCE((SELECT MAX(id_reservation_category) FROM reservation_category), 0) + 1,
    false
);

SELECT setval(
    pg_get_serial_sequence('reservation_element', 'id_reservation_element'),
    COALESCE((SELECT MAX(id_reservation_element) FROM reservation_element), 0) + 1,
    false
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- No need to revert sequence changes
-- +goose StatementEnd
