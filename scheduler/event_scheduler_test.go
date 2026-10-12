package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/devices"
	"github.com/plugimt/transat-backend/services/push"
)

type reminderProvider struct {
	mu     sync.Mutex
	tokens []string
}

func (p *reminderProvider) Name() string  { return "test" }
func (p *reminderProvider) MaxBatch() int { return 100 }
func (p *reminderProvider) Send(_ context.Context, msgs []push.Outbound) ([]push.SendResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]push.SendResult, len(msgs))
	for i, m := range msgs {
		p.tokens = append(p.tokens, m.Token)
		out[i] = push.SendResult{Kind: push.ResultAccepted, TicketID: "t"}
	}
	return out, nil
}
func (p *reminderProvider) Receipts(context.Context, []string) (map[string]push.ReceiptResult, error) {
	return nil, nil
}

type reminderFx struct {
	t        *testing.T
	db       *sql.DB
	es       *EventScheduler
	pipeline *push.Pipeline
	provider *reminderProvider
}

func newReminderFx(t *testing.T) *reminderFx {
	db := dbtest.Open(t)
	// The scheduler compares UTC wall-clock times with NOW(), which only holds on a UTC session (as in
	// production): pin the test's single connection to UTC.
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`SET TIME ZONE 'UTC'`); err != nil {
		t.Fatal(err)
	}
	ns := services.NewNotificationService(db)
	provider := &reminderProvider{}
	cfg := push.DefaultConfig()
	cfg.MaxPerSecond = 0
	pipeline := push.New(db, provider, cfg, devices.DefaultConfig())
	ns.SetPipeline(pipeline)
	return &reminderFx{t: t, db: db, es: NewEventScheduler(db, ns), pipeline: pipeline, provider: provider}
}

func (f *reminderFx) exec(q string, args ...any) {
	f.t.Helper()
	if _, err := f.db.Exec(q, args...); err != nil {
		f.t.Fatalf("%s: %v", q, err)
	}
}

func (f *reminderFx) user(n int, reminders bool) string {
	f.t.Helper()
	email := fmt.Sprintf("u%d@x.test", n)
	f.exec(`INSERT INTO newf (email, password, first_name, last_name, language, password_updated_date)
		VALUES ($1, 'x', 'A', 'B', 1, now() - interval '3 years')`, email)
	if reminders {
		f.exec(`INSERT INTO notifications (email, id_services) SELECT $1, id_services FROM services WHERE name = 'EVENT_REMINDERS'`, email)
	}
	_, err := devices.NewService(f.db, devices.DefaultConfig()).Register(context.Background(), devices.Registration{
		InstallationID: fmt.Sprintf("00000000-0000-4000-8000-%012d", n), Email: email, Platform: "ios",
		PushToken: fmt.Sprintf("ExponentPushToken[device-%d]", n), PermissionStatus: "granted",
	})
	if err != nil {
		f.t.Fatal(err)
	}
	return email
}

func (f *reminderFx) markers() int {
	var n int
	_ = f.db.QueryRow(`SELECT count(*) FROM event_notification_sent WHERE notification_type = '1h_before'`).Scan(&n)
	return n
}

func TestEventRemindersAreQueuedOnceForTheInterestedUsersAndExpireAtTheStart(t *testing.T) {
	f := newReminderFx(t)
	attendee, member, optedOut := f.user(1, true), f.user(2, true), f.user(4, false)
	f.user(3, true) // subscribed, but neither attending nor in the club
	f.exec(`INSERT INTO clubs (name, picture) VALUES ('Club', 'p')`)
	f.exec(`INSERT INTO events (name, start_date, end_date, location, creator, id_club)
		VALUES ('Gala', (now() AT TIME ZONE 'UTC') + interval '1 hour', (now() AT TIME ZONE 'UTC') + interval '3 hours', 'Hall', $1, (SELECT id_clubs FROM clubs))`, attendee)
	f.exec(`INSERT INTO events_attendents (email, id_events) SELECT $1, id_events FROM events`, attendee)
	f.exec(`INSERT INTO events_attendents (email, id_events) SELECT $1, id_events FROM events`, optedOut)
	f.exec(`INSERT INTO clubs_members (email, id_clubs) SELECT $1, id_clubs FROM clubs`, member)

	f.es.checkAndSendEventNotifications()
	if f.markers() != 1 {
		t.Fatal("a queued reminder must be marked as handled")
	}
	var kind string
	var untilStart float64
	if err := f.db.QueryRow(`
		SELECT m.kind, EXTRACT(EPOCH FROM (e.start_date AT TIME ZONE 'UTC') - m.expires_at)
		FROM push_messages m, events e`).Scan(&kind, &untilStart); err != nil {
		t.Fatal(err)
	}
	if kind != "event_reminder" || untilStart < -2 || untilStart > 2 {
		t.Fatalf("kind %s, expires %.1fs from the start", kind, untilStart)
	}

	// The scheduler runs every minute: neither the marker nor the dedupe key lets a second reminder through.
	f.es.checkAndSendEventNotifications()
	f.exec(`DELETE FROM event_notification_sent`)
	f.es.checkAndSendEventNotifications()
	var queued int
	_ = f.db.QueryRow(`SELECT count(*) FROM push_messages`).Scan(&queued)
	if queued != 1 {
		t.Fatalf("%d reminders queued for one event", queued)
	}

	if _, err := f.pipeline.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(f.provider.tokens) != 2 {
		t.Fatalf("reminder sent to %v: only the attendee and the club member qualify", f.provider.tokens)
	}
}

func TestAReminderThatCouldNotBeQueuedIsRetriedByTheNextTick(t *testing.T) {
	f := newReminderFx(t)
	creator := f.user(1, true)
	f.exec(`INSERT INTO clubs (name, picture) VALUES ('Club', 'p')`)
	f.exec(`INSERT INTO events (name, start_date, end_date, location, creator, id_club)
		VALUES ('Gala', (now() AT TIME ZONE 'UTC') + interval '1 hour', (now() AT TIME ZONE 'UTC') + interval '3 hours', 'Hall', $1, (SELECT id_clubs FROM clubs))`, creator)
	f.exec(`DROP TABLE push_deliveries`)
	f.exec(`DROP TABLE push_messages`)

	f.es.checkAndSendEventNotifications()
	if f.markers() != 0 {
		t.Fatal("a reminder that was not queued must not be marked as handled")
	}
}
