// Package push is the delivery pipeline for push notifications.
//
// A notification is a push_messages row, written by Enqueue and never sent from the request that created it.
// A worker then does the rest, in four steps that can each be repeated safely:
//
//  1. Fan-out: the audience is resolved into devices, applying the notification preferences and the session
//     policy of package devices. One push_deliveries row per eligible device is created.
//  2. Dispatch: due deliveries are claimed under a lease, checked again (the device may have signed out, changed
//     owner, or lost its token since the fan-out), and handed to the provider in batches. Transient failures are
//     retried with a bounded, jittered exponential backoff; a token is deactivated only when the provider
//     confirms that it is dead.
//  3. Receipts: the provider's receipts are read later and recorded.
//  4. Completion: a message whose deliveries are all settled is marked completed.
//
// What the statuses mean, from weakest to strongest evidence: a delivery that is "sending" was attempted;
// "accepted" means the provider took it; "handed_off" means its receipt says APNs or FCM took it. None of them
// means that the phone received or showed the notification: no provider gives reliable evidence of that, so the
// pipeline has no "delivered" status.
//
// Limits of the session policy: eligibility is evaluated when the message is fanned out and again right before
// each attempt, which is the strongest server-side filter there is. A push already accepted by APNs or FCM is out
// of our hands and can still appear after the user has left; the TTL given to the provider (the message's
// deadline, 2 hours by default) bounds how late that can be.
package push

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"regexp"
	"strings"
	"time"

	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services/devices"
)

// DefaultLanguage is the language whose text is used for a user whose own language has none.
const DefaultLanguage = "fr"

// The providers accept 4096 bytes of payload; this leaves room for what they add.
const maxPayloadBytes = 3800

var validActionCategory = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

var (
	ErrNoContent       = errors.New("the message has no text in the default language")
	ErrPayloadTooLarge = errors.New("the notification payload is too large")
)

// Text is the title and body of a notification in one language.
type Text struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
}

// Content maps a language code to its text. DefaultLanguage must be present.
type Content map[string]Text

// Message is a notification to queue.
type Message struct {
	// Kind labels the message in the delivery log ("restaurant_menu", "admin"...).
	Kind string
	// Category is the notification category whose per-user preference applies; empty for none.
	Category models.NotificationCategory
	// DedupeKey makes Enqueue idempotent: a second message with the same key is not queued. Empty means never.
	DedupeKey  string
	Audience   Audience
	Content    Content
	Navigation *models.NavigationTarget
	// ImageURL is an optional https image shown with the notification. Phones that cannot show it (an app build
	// without the iOS service extension, a failed download) show the text alone.
	ImageURL string
	// ActionCategory names the interactive notification category (the action buttons) the app registered; empty
	// for a notification without buttons. Apps that did not register it just show the notification.
	ActionCategory string
	// NotBefore delays the message; zero means now.
	NotBefore time.Time
	// ExpiresAt is the deadline after which the message is no longer worth sending; zero means DefaultTTL from now.
	ExpiresAt time.Time
	CreatedBy string
}

// EnqueueResult says what Enqueue did.
type EnqueueResult struct {
	ID int64
	// Duplicate: a message with the same DedupeKey was already queued, nothing was added.
	Duplicate bool
}

// Pipeline queues and delivers push notifications.
type Pipeline struct {
	db        *sql.DB
	provider  Provider
	cfg       Config
	deviceCfg devices.Config
	deviceSvc *devices.Service
	wake      chan struct{}
	// jitter returns a number in [0, 1); replaced in tests.
	jitter func() float64
}

// New builds a pipeline. The device config gives the session timings of the eligibility policy.
func New(db *sql.DB, provider Provider, cfg Config, deviceCfg devices.Config) *Pipeline {
	if cfg.BatchSize <= 0 || cfg.BatchSize > provider.MaxBatch() {
		cfg.BatchSize = provider.MaxBatch()
	}
	if cfg.ClaimSize <= 0 {
		cfg.ClaimSize = DefaultConfig().ClaimSize
	}
	return &Pipeline{
		db: db, provider: provider, cfg: cfg, deviceCfg: deviceCfg, deviceSvc: devices.NewService(db, deviceCfg),
		wake:   make(chan struct{}, 1),
		jitter: rand.Float64,
	}
}

// Enqueue queues a message and returns at once: nothing is sent here. It fails only when the message is invalid
// or the database is down, in which case the caller should report the failure and may try again.
func (p *Pipeline) Enqueue(ctx context.Context, m Message) (EnqueueResult, error) {
	if m.Kind == "" {
		return EnqueueResult{}, errors.New("a message needs a kind")
	}
	if err := m.Audience.Validate(); err != nil {
		return EnqueueResult{}, err
	}
	if m.Category != "" && !m.Category.IsValid() {
		return EnqueueResult{}, fmt.Errorf("unknown notification category %q", m.Category)
	}
	if m.Navigation != nil {
		if err := m.Navigation.Validate(); err != nil {
			return EnqueueResult{}, err
		}
	}
	if m.ImageURL != "" {
		if err := models.ValidateImageURL(m.ImageURL); err != nil {
			return EnqueueResult{}, err
		}
	}
	if m.ActionCategory != "" && !validActionCategory.MatchString(m.ActionCategory) {
		return EnqueueResult{}, errors.New("action category must be 1 to 64 letters, digits, dots, dashes or underscores")
	}
	content, err := fitContent(m.Content, m.Kind, m.Navigation, m.ImageURL)
	if err != nil {
		return EnqueueResult{}, err
	}

	audience, _ := json.Marshal(m.Audience)
	contentJSON, _ := json.Marshal(content)
	var navigation any
	if m.Navigation != nil {
		raw, _ := json.Marshal(m.Navigation)
		navigation = string(raw)
	}
	var image any
	if m.ImageURL != "" {
		image = m.ImageURL
	}
	var action any
	if m.ActionCategory != "" {
		action = m.ActionCategory
	}
	var dedupe, createdBy any
	if m.DedupeKey != "" {
		dedupe = m.DedupeKey
	}
	if m.CreatedBy != "" {
		createdBy = m.CreatedBy
	}
	var notBefore, expiresAt any
	if !m.NotBefore.IsZero() {
		notBefore = m.NotBefore.UTC()
	}
	if !m.ExpiresAt.IsZero() {
		expiresAt = m.ExpiresAt.UTC()
	}

	var id int64
	err = p.db.QueryRowContext(ctx, `
		INSERT INTO push_messages (kind, category, dedupe_key, audience, content, navigation, image_url, action_category,
		                           created_by, not_before, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $12, $8, COALESCE($9::timestamptz, now()),
		        COALESCE($10::timestamptz, COALESCE($9::timestamptz, now()) + make_interval(secs => $11)))
		ON CONFLICT (dedupe_key) DO NOTHING
		RETURNING id`,
		m.Kind, string(m.Category), dedupe, string(audience), string(contentJSON), navigation, image, createdBy, notBefore, expiresAt,
		p.cfg.DefaultTTL.Seconds(), action).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		if err := p.db.QueryRowContext(ctx, `SELECT id FROM push_messages WHERE dedupe_key = $1`, m.DedupeKey).Scan(&id); err != nil {
			return EnqueueResult{}, err
		}
		slog.Info("push.duplicate", "message", id, "kind", m.Kind)
		return EnqueueResult{ID: id, Duplicate: true}, nil
	}
	if err != nil {
		return EnqueueResult{}, err
	}
	slog.Info("push.created", "message", id, "kind", m.Kind, "category", string(m.Category), "languages", len(content))
	p.Wake()
	return EnqueueResult{ID: id}, nil
}

// Wake asks the worker to run now instead of at its next tick.
func (p *Pipeline) Wake() {
	select {
	case p.wake <- struct{}{}:
	default:
	}
}

// fitContent checks the default language is there and shrinks every text until the payload fits the providers'
// limit. A payload that cannot fit even without text (the data and the image link alone) is refused.
func fitContent(in Content, kind string, nav *models.NavigationTarget, image string) (Content, error) {
	if strings.TrimSpace(in[DefaultLanguage].Title) == "" {
		return nil, ErrNoContent
	}
	data := models.NotificationPayload{Kind: kind, Navigation: nav}.ExpoData()
	limit := maxPayloadBytes - imageOverhead(image)
	out := make(Content, len(in))
	for lang, t := range in {
		t.Title, t.Body = strings.TrimSpace(t.Title), strings.TrimSpace(t.Body)
		if t.Title == "" {
			continue
		}
		fitted, err := fitText(t, data, limit)
		if err != nil {
			return nil, err
		}
		out[lang] = fitted
	}
	return out, nil
}

// imageOverhead is what the image adds to the payload: its link plus the keys that carry it.
func imageOverhead(image string) int {
	if image == "" {
		return 0
	}
	return len(image) + 64
}

func payloadSize(t Text, data map[string]any) int {
	raw, _ := json.Marshal(map[string]any{"title": t.Title, "body": t.Body, "data": data, "sound": "default", "channelId": "default"})
	return len(raw)
}

func fitText(t Text, data map[string]any, limit int) (Text, error) {
	for payloadSize(t, data) > limit {
		switch {
		case t.Body != "":
			t.Body = truncateRunes(t.Body, len([]rune(t.Body))/2)
		case len([]rune(t.Title)) > 1:
			t.Title = truncateRunes(t.Title, len([]rune(t.Title))/2)
		default:
			return Text{}, ErrPayloadTooLarge
		}
	}
	return t, nil
}

func truncateRunes(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	if n <= 0 {
		return ""
	}
	return strings.TrimSpace(string(runes[:n-1])) + "…"
}

// Backoff is the wait before the attempt that follows the given one (1 is the first): the base, doubled at each
// attempt, capped by max, with a random half so that many retries do not return together. The result lies in
// [d/2, d] where d is the capped exponential delay.
func Backoff(attempt int, base, max time.Duration, jitter float64) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	d := base
	for i := 1; i < attempt && d < max; i++ {
		d *= 2
	}
	if d > max {
		d = max
	}
	return d/2 + time.Duration(jitter*float64(d/2))
}
