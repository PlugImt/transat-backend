package services

import (
	"os"
	"testing"

	"github.com/plugimt/transat-backend/internal/database"
	"github.com/plugimt/transat-backend/models"
)

// Needs a disposable Postgres: TEST_DATABASE_URL=postgres://user@host/db?sslmode=disable go test ./services
func TestNotificationPreferencesIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	db, err := database.Open(dsn, "../db/migrations")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ns := NewNotificationService(db)

	mustExec := func(query string, args ...interface{}) {
		t.Helper()
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("%s: %v", query, err)
		}
	}

	const (
		alice = "alice.prefs@test.fr"
		bob   = "bob.prefs@test.fr"
	)
	mustExec(`DELETE FROM newf WHERE email IN ($1, $2)`, alice, bob)
	for _, email := range []string{alice, bob} {
		mustExec(`INSERT INTO newf (email, password, first_name, last_name) VALUES ($1, 'x', 'T', 'T')`, email)
		mustExec(`INSERT INTO user_notification_tokens (email, token) VALUES ($1, $2)`, email, "ExponentPushToken["+email+"]")
		// Same default as registration: subscribed to every service.
		mustExec(`INSERT INTO notifications (email, id_services) SELECT $1, id_services FROM services`, email)
	}
	defer mustExec(`DELETE FROM newf WHERE email IN ($1, $2)`, alice, bob)

	prefs, err := ns.GetPreferences(alice)
	if err != nil {
		t.Fatal(err)
	}
	if len(prefs) != len(models.NotificationCategories) {
		t.Fatalf("got %d preferences, want %d: %v", len(prefs), len(models.NotificationCategories), prefs)
	}
	for i, p := range prefs {
		if p.Service != models.NotificationCategories[i] || !p.Enabled {
			t.Fatalf("preference %d = %+v, want %s enabled", i, p, models.NotificationCategories[i])
		}
	}

	// Disabling is idempotent and scoped to one user and one category.
	for range 2 {
		if err := ns.SetPreference(alice, models.NotificationEvents, false); err != nil {
			t.Fatal(err)
		}
	}
	assertEnabled := func(email string, category models.NotificationCategory, want bool) {
		t.Helper()
		got, err := ns.IsSubscribed(email, category)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("IsSubscribed(%s, %s) = %v, want %v", email, category, got, want)
		}
	}
	assertEnabled(alice, models.NotificationEvents, false)
	assertEnabled(alice, models.NotificationRestaurant, true)
	assertEnabled(bob, models.NotificationEvents, true)

	// Senders only reach subscribers.
	subscribers, err := ns.GetSubscribedUsersWithLanguage(models.NotificationEvents)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range subscribers {
		if s.Email == alice {
			t.Fatal("opted-out user returned as events subscriber")
		}
	}

	var clubID, eventID int
	if err := db.QueryRow(`INSERT INTO clubs (name, picture) VALUES ('prefs', 'x') RETURNING id_clubs`).Scan(&clubID); err != nil {
		t.Fatal(err)
	}
	defer func() {
		mustExec(`DELETE FROM events WHERE id_events = $1`, eventID)
		mustExec(`DELETE FROM clubs_members WHERE id_clubs = $1`, clubID)
		mustExec(`DELETE FROM clubs WHERE id_clubs = $1`, clubID)
	}()
	mustExec(`INSERT INTO clubs_members (email, id_clubs) VALUES ($1, $3), ($2, $3)`, alice, bob, clubID)
	if err := db.QueryRow(
		`INSERT INTO events (name, start_date, end_date, location, creator, id_club) VALUES ('e', now(), now(), 'l', $1, $2) RETURNING id_events`,
		bob, clubID,
	).Scan(&eventID); err != nil {
		t.Fatal(err)
	}

	recipients := func(category models.NotificationCategory) map[string]bool {
		t.Helper()
		users, err := ns.GetUsersInterestedInEventOrClubWithLanguage(eventID, clubID, category)
		if err != nil {
			t.Fatal(err)
		}
		set := map[string]bool{}
		for _, u := range users {
			set[u.Email] = true
		}
		return set
	}
	if got := recipients(models.NotificationEventReminders); !got[alice] || !got[bob] {
		t.Fatalf("reminders should reach both users, got %v", got)
	}
	if err := ns.SetPreference(alice, models.NotificationEventReminders, false); err != nil {
		t.Fatal(err)
	}
	if got := recipients(models.NotificationEventReminders); got[alice] || !got[bob] {
		t.Fatalf("reminders should skip the opted-out user, got %v", got)
	}

	// Admin audiences resolve to the right devices, and a category excludes opted-out users.
	mustExec(`UPDATE newf SET campus = 'NANTES' WHERE email = $1`, bob)
	devices := func(a models.NotificationAudience, category models.NotificationCategory) map[string]bool {
		t.Helper()
		r, err := ns.ResolveAudience(a, category)
		if err != nil {
			t.Fatal(err)
		}
		set := map[string]bool{}
		for _, token := range r.Tokens {
			set[token] = true
		}
		return set
	}
	aliceToken, bobToken := "ExponentPushToken["+alice+"]", "ExponentPushToken["+bob+"]"
	if got := devices(models.NotificationAudience{Type: models.AudienceClub, ClubID: clubID}, ""); len(got) != 2 || !got[aliceToken] || !got[bobToken] {
		t.Fatalf("club audience = %v", got)
	}
	if got := devices(models.NotificationAudience{Type: models.AudienceCampus, Campus: "nantes"}, ""); !got[bobToken] || got[aliceToken] {
		t.Fatalf("campus audience = %v", got)
	}
	if got := devices(models.NotificationAudience{Type: models.AudienceUsers, Emails: []string{alice}}, ""); len(got) != 1 || !got[aliceToken] {
		t.Fatalf("users audience = %v", got)
	}
	if got := devices(models.NotificationAudience{Type: models.AudienceAll}, ""); !got[aliceToken] || !got[bobToken] {
		t.Fatalf("all audience = %v", got)
	}
	// alice opted out of reminders above.
	if got := devices(models.NotificationAudience{Type: models.AudienceClub, ClubID: clubID}, models.NotificationEventReminders); got[aliceToken] || !got[bobToken] {
		t.Fatalf("club audience with category = %v", got)
	}

	// Re-enabling is idempotent too.
	for range 2 {
		if err := ns.SetPreference(alice, models.NotificationEvents, true); err != nil {
			t.Fatal(err)
		}
	}
	assertEnabled(alice, models.NotificationEvents, true)
}

func TestSetPreferenceUnknownCategory(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	db, err := database.Open(dsn, "../db/migrations")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := NewNotificationService(db).SetPreference("x@test.fr", "NOPE", true); err != ErrUnknownNotificationCategory {
		t.Fatalf("err = %v, want ErrUnknownNotificationCategory", err)
	}
}
