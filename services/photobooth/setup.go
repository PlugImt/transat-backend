package photobooth

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
)

const defaultShareBaseURL = "https://transat.dev/photobooth"

// FromEnv builds the photobooth service. It returns nil, nil, nil when PHOTOBOOTH_DEVICE_KEYS is unset, so the
// rest of the API keeps working without it.
//
//	PHOTOBOOTH_DEVICE_KEYS             booth keys, "id:secret,id:secret" (required)
//	PHOTOBOOTH_DRIVE_CREDENTIALS_JSON  base64 service-account key with write access (default GOOGLE_DRIVE_CREDENTIALS_JSON)
//	PHOTOBOOTH_DRIVE_ROOT_ID           folder holding the school-year folders (default PHOTOS_DRIVE_ID)
//	PHOTOBOOTH_SHARE_BASE_URL          share page prefix (default https://transat.dev/photobooth)
func FromEnv(ctx context.Context, db *sql.DB) (*Service, *Verifier, error) {
	rawKeys := os.Getenv("PHOTOBOOTH_DEVICE_KEYS")
	if strings.TrimSpace(rawKeys) == "" {
		return nil, nil, nil
	}
	keys, err := ParseDeviceKeys(rawKeys)
	if err != nil {
		return nil, nil, err
	}
	if len(keys) == 0 {
		return nil, nil, errors.New("PHOTOBOOTH_DEVICE_KEYS holds no key")
	}
	creds := firstSet(os.Getenv("PHOTOBOOTH_DRIVE_CREDENTIALS_JSON"), os.Getenv("GOOGLE_DRIVE_CREDENTIALS_JSON"))
	root := firstSet(os.Getenv("PHOTOBOOTH_DRIVE_ROOT_ID"), os.Getenv("PHOTOS_DRIVE_ID"))
	if creds == "" || root == "" {
		return nil, nil, errors.New("photobooth needs Drive credentials and a root folder (PHOTOBOOTH_DRIVE_* or GOOGLE_DRIVE_CREDENTIALS_JSON/PHOTOS_DRIVE_ID)")
	}
	drive, err := NewDriveClient(ctx, creds)
	if err != nil {
		return nil, nil, err
	}
	share := firstSet(os.Getenv("PHOTOBOOTH_SHARE_BASE_URL"), defaultShareBaseURL)
	return NewService(db, drive, root, share), NewVerifier(keys), nil
}

func firstSet(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
