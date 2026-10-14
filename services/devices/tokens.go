package devices

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
)

// TokenFingerprint identifies a push token without storing it: the delivery log keeps this, never the token.
func TokenFingerprint(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:8])
}

// InvalidateToken disables a device whose push token the provider confirmed as permanently dead. It does nothing
// when the device holds another token by now (the app refreshed it since the failed send), so a good token is
// never disabled by news about an old one. The account can register the device again at any time.
func (s *Service) InvalidateToken(ctx context.Context, installationID, tokenFingerprint string) (bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback() }()

	var owner string
	var token sql.NullString
	err = tx.QueryRowContext(ctx,
		`SELECT user_email, push_token FROM devices WHERE installation_id = $1 AND status = 'active' FOR UPDATE`,
		installationID).Scan(&owner, &token)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !token.Valid || TokenFingerprint(token.String) != tokenFingerprint {
		return false, nil
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE devices
		SET status = 'disabled', disabled_reason = 'invalid_token', push_token = NULL, updated_at = now()
		WHERE installation_id = $1`, installationID); err != nil {
		return false, err
	}
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM user_notification_tokens WHERE email = $1 AND token = $2`, owner, token.String); err != nil {
		return false, err
	}
	return true, tx.Commit()
}

// RegisterLegacyToken keeps the token sent by an app version without device registration (PATCH /newf/me) in the
// devices table. The call is the only proof of session such a version gives, so it refreshes last_seen_at. A token
// a registered installation already holds is left to that installation.
func RegisterLegacyToken(ctx context.Context, db *sql.DB, email, token string) error {
	token = strings.TrimSpace(token)
	if !expoTokenPattern.MatchString(token) {
		return ErrInvalidPushToken
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var source string
	err = tx.QueryRowContext(ctx,
		`SELECT source FROM devices WHERE provider = 'expo' AND push_token = $1 FOR UPDATE`, token).Scan(&source)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		_, err = tx.ExecContext(ctx, `
			INSERT INTO devices (installation_id, user_email, platform, push_token, permission_status, source)
			VALUES (gen_random_uuid(), $1, 'unknown', $2, 'unknown', 'legacy')`, email, token)
	case err != nil:
	case source == "legacy":
		_, err = tx.ExecContext(ctx, `
			UPDATE devices
			SET user_email = $2, status = 'active', disabled_reason = NULL, last_seen_at = now(), updated_at = now()
			WHERE provider = 'expo' AND push_token = $1`, token, email)
	default:
		return nil
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}
