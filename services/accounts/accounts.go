// Package accounts deletes user accounts, on request or after five years without activity, removing what can be
// removed and anonymising the rest.
package accounts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// DeletedAccountEmail is the placeholder that keeps records which must outlive an account: the events it created
// and its reservation history (a club must still see that an item is out). It cannot sign in.
const DeletedAccountEmail = "deleted-account@transat.invalid"

var (
	ErrNotFound  = errors.New("account not found")
	ErrProtected = errors.New("this account cannot be deleted")
)

// FileRemover deletes stored files (profile pictures, uploads) once the account is gone.
type FileRemover interface {
	DeleteFile(key string) error
}

type Service struct {
	db    *sql.DB
	files FileRemover // optional
	now   func() time.Time
}

func NewService(db *sql.DB, files FileRemover) *Service {
	return &Service{db: db, files: files, now: time.Now}
}

// Delete removes an account and its personal data in one transaction:
//   - deleted: profile, roles, devices, notification tokens and history, club memberships, event participations,
//     waiting lists, upcoming reservations, reviews, game history (scores are a view of it), uploaded files;
//   - anonymised (reassigned to DeletedAccountEmail): events it created and its past or running reservations, with
//     their cancellation and reminder records;
//   - unlinked: photo reports it made or decided, request statistics.
//
// Stored files are removed after the commit, best effort.
func (s *Service) Delete(ctx context.Context, email string) error {
	email = strings.TrimSpace(email)
	if email == "" || strings.EqualFold(email, DeletedAccountEmail) {
		return ErrProtected
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var picture sql.NullString
	err = tx.QueryRowContext(ctx, `SELECT email, profile_picture FROM newf WHERE lower(email) = lower($1) FOR UPDATE`, email).
		Scan(&email, &picture)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	paths, err := strings1(ctx, tx, `SELECT path FROM files WHERE email = $1`, email)
	if err != nil {
		return err
	}

	now := s.now().UTC()
	steps := []struct {
		what  string
		query string
		args  []any
	}{
		// Upcoming bookings free their slot; history and running loans are kept, anonymised.
		{"upcoming reservations", `DELETE FROM reservation WHERE email = $1 AND start_date > $2`, []any{email, now}},
		{"reservation reminders", `DELETE FROM reservation_reminders WHERE email = $1 AND start_date > $2`, []any{email, now}},
		{"reservations", `UPDATE reservation SET email = $2 WHERE email = $1`, []any{email, DeletedAccountEmail}},
		{"reservation reminders", `UPDATE reservation_reminders SET email = $2 WHERE email = $1`, []any{email, DeletedAccountEmail}},
		{"reservation cancellations", `UPDATE reservation_cancellations SET email = $2 WHERE email = $1`, []any{email, DeletedAccountEmail}},
		{"cancellations made", `UPDATE reservation_cancellations SET cancelled_by = $2 WHERE cancelled_by = $1`, []any{email, DeletedAccountEmail}},
		{"events created", `UPDATE events SET creator = $2 WHERE creator = $1`, []any{email, DeletedAccountEmail}},

		{"club memberships", `DELETE FROM clubs_members WHERE email = $1`, []any{email}},
		{"reviews", `DELETE FROM restaurant_articles_notes WHERE email = $1`, []any{email}},
		{"game history", `DELETE FROM bassine_history WHERE email = $1`, []any{email}},
		{"push deliveries", `DELETE FROM push_deliveries WHERE user_email = $1`, []any{email}},
		{"report emails", `DELETE FROM photo_report_deliveries WHERE recipient_email = $1`, []any{email}},
		{"report decisions", `UPDATE photo_reports SET decided_by = NULL WHERE decided_by = $1`, []any{email}},
		{"push authors", `UPDATE push_messages SET created_by = NULL WHERE created_by = $1`, []any{email}},
		// Everything else references newf with ON DELETE CASCADE or SET NULL.
		{"account", `DELETE FROM newf WHERE email = $1`, []any{email}},
	}
	for _, st := range steps {
		if _, err := tx.ExecContext(ctx, st.query, st.args...); err != nil {
			return fmt.Errorf("%s: %w", st.what, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	if s.files != nil {
		keys := append([]string{}, paths...)
		if picture.Valid && picture.String != "" {
			keys = append(keys, picture.String)
		}
		for _, key := range keys {
			if k := storageKey(key); k != "" {
				if err := s.files.DeleteFile(k); err != nil {
					log.Printf("Accounts: stored file of a deleted account not removed: %v", err)
				}
			}
		}
	}
	return nil
}

// storageKey turns a stored path or public URL into the object key, ignoring anything that is not ours.
func storageKey(v string) string {
	v = strings.TrimSpace(v)
	if i := strings.Index(v, "://"); i >= 0 {
		rest := v[i+3:]
		j := strings.IndexByte(rest, '/')
		if j < 0 {
			return ""
		}
		v = rest[j+1:]
	}
	return strings.TrimLeft(v, "/")
}

func strings1(ctx context.Context, tx *sql.Tx, q string, args ...any) ([]string, error) {
	rows, err := tx.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
