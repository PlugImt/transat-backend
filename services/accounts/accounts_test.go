package accounts

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/models"
)

type fx struct {
	t   *testing.T
	svc *Service
}

func (f fx) exec(q string, args ...any) {
	f.t.Helper()
	if _, err := f.svc.db.Exec(q, args...); err != nil {
		f.t.Fatalf("%s: %v", q, err)
	}
}

func (f fx) count(q string, args ...any) int {
	f.t.Helper()
	var n int
	if err := f.svc.db.QueryRow(q, args...).Scan(&n); err != nil {
		f.t.Fatalf("%s: %v", q, err)
	}
	return n
}

func (f fx) user(email string, lastActivity time.Time) {
	f.exec(`INSERT INTO newf (email, password, first_name, last_name, last_activity, creation_date)
		VALUES ($1, 'x', 'Jo', 'Doe', $2, $2)`, email, lastActivity)
}

type files struct{ deleted []string }

func (r *files) DeleteFile(key string) error {
	r.deleted = append(r.deleted, key)
	return nil
}

func setup(t *testing.T) (fx, *files) {
	db := dbtest.Open(t)
	r := &files{}
	return fx{t: t, svc: NewService(db, r)}, r
}

func TestDeletingAnAccountRemovesPersonalDataAndAnonymisesWhatMustStay(t *testing.T) {
	f, stored := setup(t)
	now := time.Now()
	f.user("gone@x.test", now)
	f.user("other@x.test", now)
	f.exec(`UPDATE newf SET profile_picture = 'https://cdn.example/profiles/gone.jpg' WHERE email = 'gone@x.test'`)
	f.exec(`INSERT INTO clubs (name, picture) VALUES ('Club', 'p')`)
	f.exec(`INSERT INTO clubs_members (email, id_clubs) SELECT 'gone@x.test', id_clubs FROM clubs`)
	f.exec(`INSERT INTO events (name, start_date, end_date, location, creation_date, creator, id_club)
		SELECT 'Gala', now(), now(), 'Here', now(), 'gone@x.test', id_clubs FROM clubs`)
	f.exec(`INSERT INTO bassine_history (email, date) VALUES ('gone@x.test', now())`)
	f.exec(`INSERT INTO reservation_element (name, slot, id_clubs) SELECT 'Salle', TRUE, id_clubs FROM clubs`)
	f.exec(`INSERT INTO reservation (email, id_reservation_element, start_date, end_date)
		SELECT 'gone@x.test', id_reservation_element, now() - interval '3 days', now() - interval '2 days' FROM reservation_element`)
	f.exec(`INSERT INTO reservation (email, id_reservation_element, start_date, end_date)
		SELECT 'gone@x.test', id_reservation_element, now() + interval '3 days', now() + interval '4 days' FROM reservation_element`)
	f.exec(`INSERT INTO files (name, path, email, creation_date, last_access_date) VALUES ('a', 'uploads/a.jpg', 'gone@x.test', now(), now())`)

	if err := f.svc.Delete(context.Background(), "GONE@x.test"); err != nil {
		t.Fatal(err)
	}

	if n := f.count(`SELECT count(*) FROM newf WHERE email = 'gone@x.test'`); n != 0 {
		t.Fatal("account still there")
	}
	for _, q := range []string{
		`SELECT count(*) FROM clubs_members WHERE email = 'gone@x.test'`,
		`SELECT count(*) FROM bassine_history WHERE email = 'gone@x.test'`,
		`SELECT count(*) FROM files WHERE email = 'gone@x.test'`,
		`SELECT count(*) FROM reservation WHERE email = 'gone@x.test'`,
		`SELECT count(*) FROM reservation WHERE start_date > now()`,
	} {
		if n := f.count(q); n != 0 {
			t.Fatalf("%s = %d, want 0", q, n)
		}
	}
	if f.count(`SELECT count(*) FROM events WHERE creator = $1`, DeletedAccountEmail) != 1 {
		t.Fatal("the event the account created must stay, anonymised")
	}
	if f.count(`SELECT count(*) FROM reservation WHERE email = $1`, DeletedAccountEmail) != 1 {
		t.Fatal("past reservations must stay, anonymised")
	}
	if f.count(`SELECT count(*) FROM newf WHERE email = 'other@x.test'`) != 1 {
		t.Fatal("another account was touched")
	}
	if len(stored.deleted) != 2 || stored.deleted[0] != "uploads/a.jpg" || stored.deleted[1] != "profiles/gone.jpg" {
		t.Fatalf("stored files removed: %v", stored.deleted)
	}

	if err := f.svc.Delete(context.Background(), "gone@x.test"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("second delete: %v", err)
	}
	if err := f.svc.Delete(context.Background(), DeletedAccountEmail); !errors.Is(err, ErrProtected) {
		t.Fatalf("placeholder delete: %v", err)
	}
}

type mails struct{ to []string }

func (m *mails) SendEmail(mail models.Email, _ interface{}) error {
	m.to = append(m.to, mail.Recipient)
	return nil
}

func TestInactiveAccountsAreWarnedThenDeletedUnlessTheyComeBack(t *testing.T) {
	f, _ := setup(t)
	ctx := context.Background()
	now := time.Now()
	years := func(n float64) time.Time { return now.Add(-time.Duration(n * float64(InactivityLimit/5))) }
	f.user("active@x.test", now.Add(-24*time.Hour))
	f.user("almost@x.test", years(4.95)) // within the 30-day notice
	f.user("old@x.test", years(6))
	f.user("returns@x.test", years(6))

	m := &mails{}
	rep, err := f.svc.PurgeInactive(ctx, m)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Warned != 3 || rep.Deleted != 0 || len(m.to) != 3 {
		t.Fatalf("first run: %+v, mailed %v", rep, m.to)
	}
	// Nobody is deleted on the day of the warning, and nobody is warned twice.
	if rep, _ := f.svc.PurgeInactive(ctx, m); rep.Warned != 0 || rep.Deleted != 0 {
		t.Fatalf("second run: %+v", rep)
	}

	// 31 days later; one account signed in after its warning.
	f.svc.now = func() time.Time { return now.Add(31 * 24 * time.Hour) }
	f.exec(`UPDATE newf SET last_activity = $1 WHERE email = 'returns@x.test'`, now.Add(24*time.Hour))
	rep, err = f.svc.PurgeInactive(ctx, m)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Deleted != 2 {
		t.Fatalf("deleted %d, want the two accounts past 5 years that stayed away", rep.Deleted)
	}
	for email, want := range map[string]int{"active@x.test": 1, "almost@x.test": 0, "old@x.test": 0, "returns@x.test": 1} {
		if got := f.count(`SELECT count(*) FROM newf WHERE email = $1`, email); got != want {
			t.Errorf("%s: %d, want %d", email, got, want)
		}
	}
}
