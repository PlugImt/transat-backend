package services

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services/devices"
	"github.com/plugimt/transat-backend/services/push"
)

// recordingProvider accepts everything and remembers what it was given.
type recordingProvider struct {
	mu   sync.Mutex
	sent []push.Outbound
}

func (p *recordingProvider) Name() string  { return "recording" }
func (p *recordingProvider) MaxBatch() int { return 100 }
func (p *recordingProvider) Send(_ context.Context, msgs []push.Outbound) ([]push.SendResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sent = append(p.sent, msgs...)
	out := make([]push.SendResult, len(msgs))
	for i := range msgs {
		out[i] = push.SendResult{Kind: push.ResultAccepted, TicketID: fmt.Sprintf("t%d", len(p.sent)-len(msgs)+i)}
	}
	return out, nil
}
func (p *recordingProvider) Receipts(context.Context, []string) (map[string]push.ReceiptResult, error) {
	return nil, nil
}

func (p *recordingProvider) byToken() map[string]push.Outbound {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := map[string]push.Outbound{}
	for _, m := range p.sent {
		out[m.Token] = m
	}
	return out
}

func (p *recordingProvider) count() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.sent)
}

func (f *fixture) usePipeline() (*push.Pipeline, *recordingProvider) {
	rp := &recordingProvider{}
	cfg := push.DefaultConfig()
	cfg.MaxPerSecond = 0
	p := push.New(f.db, rp, cfg, devices.DefaultConfig())
	f.ns.SetPipeline(p)
	return p, rp
}

// device registers an installation for a user, as the app does.
func (f *fixture) device(n int, email, token string) {
	f.t.Helper()
	_, err := devices.NewService(f.db, devices.DefaultConfig()).Register(context.Background(), devices.Registration{
		InstallationID: fmt.Sprintf("00000000-0000-4000-8000-%012d", n), Email: email, Platform: "ios",
		PushToken: token, PermissionStatus: "granted",
	})
	if err != nil {
		f.t.Fatal(err)
	}
}

func TestPhotoAlbumIsQueuedNotSentFromTheCall(t *testing.T) {
	f := newFixture(t)
	pipeline, provider := f.usePipeline()
	f.user("fr@imt.fr", 1, true)
	f.user("en@imt.fr", 2, true)
	f.user("declined@imt.fr", 1, false)
	f.user("nodevice@imt.fr", 1, true)
	f.device(1, "fr@imt.fr", tok("fr"))
	f.device(2, "fr@imt.fr", tok("fr-tablet"))
	f.device(3, "en@imt.fr", tok("en"))
	f.device(4, "declined@imt.fr", tok("declined"))

	date := time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC)
	delivery, err := f.ns.NotifyNewPhotoAlbum(42, "Soirée de gala", &date)
	if err != nil {
		t.Fatal(err)
	}
	if delivery.Users != 3 || delivery.Devices != 3 {
		t.Fatalf("delivery = %+v: three opted-in users, three devices", delivery)
	}
	if provider.count() != 0 {
		t.Fatal("the call itself must not send anything")
	}

	// Asking again for the same album queues nothing more.
	if _, err := f.ns.NotifyNewPhotoAlbum(42, "Soirée de gala", &date); err != nil {
		t.Fatal(err)
	}
	if _, err := pipeline.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	sent := provider.byToken()
	if provider.count() != 3 || len(sent) != 3 {
		t.Fatalf("%d notifications to %d devices", provider.count(), len(sent))
	}
	if _, reached := sent[tok("declined")]; reached {
		t.Fatal("a user who did not opt in was notified")
	}
	if m := sent[tok("fr")]; m.Title != "Nouvel album photo" || m.Body != "Soirée de gala · 12/03/2026" {
		t.Fatalf("fr = %q / %q", m.Title, m.Body)
	}
	if m := sent[tok("fr-tablet")]; m.Title != "Nouvel album photo" {
		t.Fatalf("second device of the same user = %q", m.Title)
	}
	if m := sent[tok("en")]; m.Title != "New photo album" || m.Body != "Soirée de gala · Mar 12, 2026" {
		t.Fatalf("en = %q / %q", m.Title, m.Body)
	}
	nav := sent[tok("fr")].Data["navigation"].(*models.NavigationTarget)
	if nav.Type != models.NavigationPhotoAlbum || nav.ID != "42" || nav.Params["title"] != "Soirée de gala" {
		t.Fatalf("navigation = %+v", nav)
	}
	for _, m := range sent {
		if strings.Contains(fmt.Sprint(m.Data), "@imt.fr") {
			t.Fatalf("payload leaks an email: %v", m.Data)
		}
	}
}

func TestAnnouncementsAreWrittenInEveryLanguage(t *testing.T) {
	f := newFixture(t)
	pipeline, _ := f.usePipeline()
	var langs []string
	res, err := f.ns.Announce(context.Background(), Announcement{
		Kind: "test", Category: models.NotificationRestaurant, DedupeKey: "k",
		Compose: func(lang string, _ *i18n.Localizer) (string, string) {
			langs = append(langs, lang)
			return "title " + lang, "body"
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	var stored int
	_ = f.db.QueryRow(`SELECT count(*) FROM languages`).Scan(&stored)
	if len(langs) != stored || stored < 2 {
		t.Fatalf("composed %v for %d languages", langs, stored)
	}
	sum, _ := pipeline.Summary(context.Background(), res.ID)
	if sum.Kind != "test" || sum.Status != "pending" {
		t.Fatalf("summary = %+v", sum)
	}

	// Without a pipeline nothing is queued: the announcement is sent directly, here to nobody.
	f.ns.SetPipeline(nil)
	if _, err := f.ns.Announce(context.Background(), Announcement{Kind: "x", Category: models.NotificationRestaurant,
		Audience: usersOf("nobody@x.test"), Compose: func(string, *i18n.Localizer) (string, string) { return "t", "" }}); err != nil {
		t.Fatalf("Announce without a pipeline and without recipient: %v", err)
	}
}

func TestOutboxModeKeepsTheLegacySenderUntouchedWhenOff(t *testing.T) {
	f := newFixture(t)
	if f.ns.OutboxEnabled() || f.ns.Pipeline() != nil {
		t.Fatal("the service must send directly until a pipeline is set")
	}
	f.user("a@imt.fr", 1, true, tok("a"))
	if _, err := f.ns.NotifyNewPhotoAlbum(1, "Gala", nil); err != nil {
		t.Fatal(err)
	}
	if len(f.ex.tokens()) != 1 {
		t.Fatalf("the direct sender addressed %v", f.ex.tokens())
	}
}
