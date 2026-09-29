-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (name, description)
VALUES ('ACADEMICS', 'IMT Atlantique academic staff member (@imt-atlantique.fr)');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM roles WHERE name = 'ACADEMICS';
-- +goose StatementEnd
