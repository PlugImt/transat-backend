-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (name, description)
VALUES ('STAFF', 'IMT Atlantique staff member (@imt-atlantique.fr)');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM roles WHERE name = 'STAFF';
-- +goose StatementEnd