-- +goose Up
-- +goose StatementBegin

ALTER TABLE newf
    ADD CONSTRAINT newf_id_newf_unique UNIQUE (id_newf);

CREATE TABLE user_schedule
(
    user_id       INTEGER   NOT NULL,
    ics_url       TEXT,
    last_sync_at  TIMESTAMP,
    calendar_data JSONB,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES newf (id_newf) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS user_schedule CASCADE;

ALTER TABLE newf
    DROP CONSTRAINT IF EXISTS newf_id_newf_unique;

-- +goose StatementEnd
