-- +goose Up
-- +goose StatementBegin
-- When an inactive account was told it will be deleted; reset by any activity after it.
ALTER TABLE newf ADD COLUMN IF NOT EXISTS deletion_warned_at TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS idx_newf_last_activity ON newf (COALESCE(last_activity, creation_date));
-- +goose StatementEnd
-- +goose StatementBegin
-- Records that must outlive a deleted account (events it created, its reservation history) are reassigned to this
-- placeholder. Its password is not a bcrypt hash, so nobody can sign in with it.
INSERT INTO newf (email, password, first_name, last_name, password_updated_date)
VALUES ('deleted-account@transat.invalid', '!', 'Compte', 'supprimé', now())
ON CONFLICT (email) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_newf_last_activity;
ALTER TABLE newf DROP COLUMN IF EXISTS deletion_warned_at;
-- +goose StatementEnd
