-- +goose Up
-- +goose StatementBegin
-- Face recognition ("my photos"). Biometric data (GDPR art. 9): everything tied to a user hangs off
-- face_consents, so withdrawing consent, deleting the account or the inactivity purge removes it in one cascade.

-- One row while a user has opted in; withdrawing deletes it.
CREATE TABLE IF NOT EXISTS face_consents (
    email           VARCHAR(100) PRIMARY KEY REFERENCES newf (email) ON DELETE CASCADE ON UPDATE CASCADE,
    consent_version TEXT        NOT NULL,
    consented_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- A user's reference face, computed from a selfie that is never stored. matched_through is the last
-- photo_faces.id already compared with it: a new reference starts at 0 and sweeps the whole library once.
CREATE TABLE IF NOT EXISTS face_references (
    id              BIGSERIAL PRIMARY KEY,
    email           VARCHAR(100) NOT NULL REFERENCES face_consents (email) ON DELETE CASCADE ON UPDATE CASCADE,
    model           TEXT        NOT NULL,
    embedding       BYTEA       NOT NULL, -- L2-normalised float32 little-endian
    matched_through BIGINT      NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_face_references_email ON face_references (email);

-- Which indexed photos went through the face detector, with which model and on which file content.
CREATE TABLE IF NOT EXISTS face_photo_scans (
    drive_file_id TEXT PRIMARY KEY REFERENCES photos (drive_file_id) ON DELETE CASCADE,
    model         TEXT        NOT NULL,
    md5           TEXT        NOT NULL DEFAULT '',
    status        TEXT        NOT NULL CHECK (status IN ('done', 'failed', 'skipped')),
    attempts      INT         NOT NULL DEFAULT 1,
    face_count    INT         NOT NULL DEFAULT 0,
    last_error    TEXT,
    scanned_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Every face found on a photo. Removed with the photo (or when it leaves the gallery).
CREATE TABLE IF NOT EXISTS photo_faces (
    id            BIGSERIAL PRIMARY KEY,
    drive_file_id TEXT        NOT NULL REFERENCES photos (drive_file_id) ON DELETE CASCADE,
    model         TEXT        NOT NULL,
    -- Bounding box relative to the image (0..1).
    box_x1        REAL        NOT NULL,
    box_y1        REAL        NOT NULL,
    box_x2        REAL        NOT NULL,
    box_y2        REAL        NOT NULL,
    score         REAL        NOT NULL,
    embedding     BYTEA       NOT NULL, -- L2-normalised float32 little-endian
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_photo_faces_file ON photo_faces (drive_file_id);

-- A face that looks like a consenting user. rejected_at is set when the user says "not me": the photo leaves
-- their album and the same face is never matched to them again.
CREATE TABLE IF NOT EXISTS face_matches (
    email         VARCHAR(100) NOT NULL REFERENCES face_consents (email) ON DELETE CASCADE ON UPDATE CASCADE,
    face_id       BIGINT      NOT NULL REFERENCES photo_faces (id) ON DELETE CASCADE,
    drive_file_id TEXT        NOT NULL REFERENCES photos (drive_file_id) ON DELETE CASCADE,
    distance      REAL        NOT NULL,
    matched_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    rejected_at   TIMESTAMPTZ,
    PRIMARY KEY (email, face_id)
);
CREATE INDEX IF NOT EXISTS idx_face_matches_photo ON face_matches (email, drive_file_id) WHERE rejected_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS face_matches;
DROP TABLE IF EXISTS photo_faces;
DROP TABLE IF EXISTS face_photo_scans;
DROP TABLE IF EXISTS face_references;
DROP TABLE IF EXISTS face_consents;
-- +goose StatementEnd
