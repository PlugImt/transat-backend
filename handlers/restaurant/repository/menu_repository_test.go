package repository

import (
	"context"
	"testing"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/devices"
	"github.com/plugimt/transat-backend/services/push"
)

type silentProvider struct{}

func (silentProvider) Name() string  { return "test" }
func (silentProvider) MaxBatch() int { return 100 }
func (silentProvider) Send(_ context.Context, msgs []push.Outbound) ([]push.SendResult, error) {
	return make([]push.SendResult, len(msgs)), nil
}
func (silentProvider) Receipts(context.Context, []string) (map[string]push.ReceiptResult, error) {
	return nil, nil
}

func newMenuRepository(t *testing.T) *MenuRepository {
	db := dbtest.Open(t)
	ns := services.NewNotificationService(db)
	ns.SetPipeline(push.New(db, silentProvider{}, push.DefaultConfig(), devices.DefaultConfig()))
	return &MenuRepository{DB: db, NotifService: ns}
}

func count(t *testing.T, r *MenuRepository, query string) int {
	t.Helper()
	var n int
	if err := r.DB.QueryRow(query).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func TestTheMenuNotificationIsQueuedOncePerDayAndRecordedAfterwards(t *testing.T) {
	r := newMenuRepository(t)
	if err := r.sendMenuUpdateNotification("2026-10-12"); err != nil {
		t.Fatal(err)
	}
	if count(t, r, `SELECT count(*) FROM restaurant_notifications WHERE notification_sent`) != 1 {
		t.Fatal("the day was not recorded")
	}
	// Even if the record had been lost, the dedupe key keeps a second notification out.
	if err := r.sendMenuUpdateNotification("2026-10-12"); err != nil {
		t.Fatal(err)
	}
	if n := count(t, r, `SELECT count(*) FROM push_messages WHERE kind = 'restaurant_menu'`); n != 1 {
		t.Fatalf("%d menu notifications queued for one day", n)
	}
	if err := r.sendMenuUpdateNotification("2026-10-13"); err != nil {
		t.Fatal(err)
	}
	if n := count(t, r, `SELECT count(*) FROM push_messages WHERE kind = 'restaurant_menu'`); n != 2 {
		t.Fatalf("%d messages after two days", n)
	}
	var category string
	_ = r.DB.QueryRow(`SELECT category FROM push_messages LIMIT 1`).Scan(&category)
	if category != "RESTAURANT" {
		t.Fatalf("category = %q: the restaurant preference must apply", category)
	}
}

func TestAMenuNotificationThatCouldNotBeQueuedIsNotRecordedAsSent(t *testing.T) {
	r := newMenuRepository(t)
	if _, err := r.DB.Exec(`DROP TABLE push_deliveries`); err != nil {
		t.Fatal(err)
	}
	if _, err := r.DB.Exec(`DROP TABLE push_messages`); err != nil {
		t.Fatal(err)
	}
	if err := r.sendMenuUpdateNotification("2026-10-12"); err == nil {
		t.Fatal("the failure was swallowed")
	}
	if count(t, r, `SELECT count(*) FROM restaurant_notifications`) != 0 {
		t.Fatal("a notification that was not queued was recorded as sent: it would never be retried")
	}
}
