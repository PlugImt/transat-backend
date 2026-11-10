-- +goose Up
-- +goose StatementBegin
-- A report now ends in a decision: approved (the photo is removed from the Drive and the gallery) or rejected.
ALTER TABLE photo_reports
    ADD COLUMN IF NOT EXISTS status        TEXT NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS decided_at    TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS decided_by    VARCHAR(100),
    ADD COLUMN IF NOT EXISTS decision_note TEXT NOT NULL DEFAULT '';
ALTER TABLE photo_reports DROP CONSTRAINT IF EXISTS photo_reports_status_check;
ALTER TABLE photo_reports ADD CONSTRAINT photo_reports_status_check CHECK (status IN ('pending', 'approved', 'rejected'));
ALTER TABLE photo_reports DROP CONSTRAINT IF EXISTS photo_reports_decision_note_check;
ALTER TABLE photo_reports ADD CONSTRAINT photo_reports_decision_note_check CHECK (char_length(decision_note) <= 1000);
CREATE INDEX IF NOT EXISTS idx_photo_reports_status ON photo_reports (status, created_at DESC);
-- +goose StatementEnd
-- +goose StatementBegin
-- The Club Photo's owners receive every report, whatever club the photo belongs to.
ALTER TABLE photo_report_deliveries DROP CONSTRAINT IF EXISTS photo_report_deliveries_recipient_role_check;
ALTER TABLE photo_report_deliveries ADD CONSTRAINT photo_report_deliveries_recipient_role_check
    CHECK (recipient_role IN ('club_owner', 'photo_club', 'admin'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE photo_report_deliveries SET recipient_role = 'club_owner' WHERE recipient_role = 'photo_club';
ALTER TABLE photo_report_deliveries DROP CONSTRAINT IF EXISTS photo_report_deliveries_recipient_role_check;
ALTER TABLE photo_report_deliveries ADD CONSTRAINT photo_report_deliveries_recipient_role_check
    CHECK (recipient_role IN ('club_owner', 'admin'));
DROP INDEX IF EXISTS idx_photo_reports_status;
ALTER TABLE photo_reports DROP CONSTRAINT IF EXISTS photo_reports_decision_note_check;
ALTER TABLE photo_reports DROP CONSTRAINT IF EXISTS photo_reports_status_check;
ALTER TABLE photo_reports
    DROP COLUMN IF EXISTS decision_note,
    DROP COLUMN IF EXISTS decided_by,
    DROP COLUMN IF EXISTS decided_at,
    DROP COLUMN IF EXISTS status;
-- +goose StatementEnd
