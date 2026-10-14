package database_test

import (
	"database/sql"
	"testing"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/pressly/goose/v3"
)

// A production database is upgraded from where it is, with data in it, and may have to go back. This runs the
// notification migrations (37 to 40) on a database that already holds legacy push tokens.
func TestNotificationMigrationsUpgradeKeepTheLegacyTokensAndCanBeRolledBack(t *testing.T) {
	db, dir := dbtest.OpenEmpty(t)
	if err := goose.UpTo(db, dir, 36); err != nil {
		t.Fatal(err)
	}

	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	for _, email := range []string{"alice@x.test", "bob@x.test"} {
		exec(`INSERT INTO newf (email, password, first_name, last_name, password_updated_date)
			VALUES ($1, 'x', 'A', 'B', now() - interval '3 years')`, email)
	}
	exec(`INSERT INTO user_notification_tokens (email, token, created_at) VALUES
		('alice@x.test', 'ExponentPushToken[alice-only]', now() - interval '2 days'),
		('alice@x.test', 'ExponentPushToken[shared-phone]', now() - interval '3 days'),
		('bob@x.test', 'ExponentPushToken[shared-phone]', now() - interval '1 day'),
		('bob@x.test', 'not-a-token', now())`)

	if err := goose.UpTo(db, dir, 40); err != nil {
		t.Fatalf("upgrade: %v", err)
	}

	owners := func() map[string]string {
		t.Helper()
		rows, err := db.Query(`SELECT push_token, user_email FROM devices WHERE source = 'legacy'`)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		out := map[string]string{}
		for rows.Next() {
			var token, email string
			if err := rows.Scan(&token, &email); err != nil {
				t.Fatal(err)
			}
			out[token] = email
		}
		return out
	}
	got := owners()
	if len(got) != 2 || got["ExponentPushToken[alice-only]"] != "alice@x.test" || got["ExponentPushToken[shared-phone]"] != "bob@x.test" {
		t.Fatalf("legacy devices = %v: a token shared by two accounts goes to the one that registered it last, and a malformed one is dropped", got)
	}
	if n := count(t, db, `SELECT count(*) FROM user_notification_tokens`); n != 4 {
		t.Fatalf("the legacy table lost rows: %d", n)
	}
	// Every table of the pipeline exists and takes a message with an image.
	exec(`INSERT INTO push_messages (kind, audience, content, expires_at, image_url)
		VALUES ('test', '{"type":"all"}', '{"fr":{"title":"t"}}', now() + interval '1 hour', 'https://cdn.x.test/a.jpg')`)

	// Rolling back leaves the legacy table, which the previous version of the server reads, untouched.
	if err := goose.DownTo(db, dir, 36); err != nil {
		t.Fatalf("rollback: %v", err)
	}
	for _, table := range []string{"devices", "push_messages", "push_deliveries"} {
		if n := count(t, db, `SELECT count(*) FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = $1`, table); n != 0 {
			t.Fatalf("%s survived the rollback", table)
		}
	}
	if n := count(t, db, `SELECT count(*) FROM user_notification_tokens`); n != 4 {
		t.Fatalf("the rollback touched the legacy table: %d", n)
	}

	// And the upgrade can be run again.
	if err := goose.Up(db, dir); err != nil {
		t.Fatalf("second upgrade: %v", err)
	}
	if got := owners(); len(got) != 2 {
		t.Fatalf("legacy devices after the second upgrade = %v", got)
	}
}

func count(t *testing.T, db *sql.DB, q string, args ...any) int {
	t.Helper()
	var n int
	if err := db.QueryRow(q, args...).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}
