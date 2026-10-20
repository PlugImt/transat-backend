package push

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/lib/pq"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services/devices"
)

const (
	maxFanOutPerRun  = 20
	maxRoundsPerRun  = 20
	maxReceiptRounds = 5
	purgeEvery       = time.Hour
)

// Stats counts what one run of the worker did.
type Stats struct {
	FannedOut   int
	Claimed     int
	Skipped     int
	Accepted    int
	Retried     int
	Failed      int
	Invalidated int
	Receipts    int
	Completed   int
}

// Worker runs the pipeline in the background.
type Worker struct {
	cancel context.CancelFunc
	done   chan struct{}
}

// Stop ends the worker and waits for the round in progress. Deliveries it had claimed are taken back once their
// lease runs out.
func (w *Worker) Stop() {
	w.cancel()
	<-w.done
}

// Start runs the pipeline every PollInterval, and at once after Wake. Several instances can run it at the same
// time: every step claims its rows with FOR UPDATE SKIP LOCKED.
func (p *Pipeline) Start() *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	w := &Worker{cancel: cancel, done: make(chan struct{})}
	go func() {
		defer close(w.done)
		ticker := time.NewTicker(p.cfg.PollInterval)
		defer ticker.Stop()
		var lastPurge time.Time
		for {
			p.runSafely(ctx)
			if time.Since(lastPurge) >= purgeEvery {
				lastPurge = time.Now()
				if err := p.purge(ctx); err != nil && ctx.Err() == nil {
					slog.Error("push.purge_failed", "error", err)
				}
			}
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			case <-p.wake:
			}
		}
	}()
	slog.Info("push.worker_started", "provider", p.provider.Name(), "poll", p.cfg.PollInterval.String())
	return w
}

func (p *Pipeline) runSafely(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("push.worker_panic", "panic", fmt.Sprint(r))
		}
	}()
	if _, err := p.RunOnce(ctx); err != nil && ctx.Err() == nil {
		slog.Error("push.run_failed", "error", err)
	}
}

// RunOnce does every step once: fan-out, dispatch, receipts, completion. It is what the worker repeats, and what
// tests call. The first error is returned after the other steps have run.
func (p *Pipeline) RunOnce(ctx context.Context) (Stats, error) {
	var st Stats
	var first error
	keep := func(err error) {
		if err != nil && first == nil {
			first = err
		}
	}
	keep(p.fanOut(ctx, &st))
	keep(p.dispatch(ctx, &st))
	keep(p.collectReceipts(ctx, &st))
	keep(p.complete(ctx, &st))
	return st, first
}

// ---- Fan-out ----

func (p *Pipeline) fanOut(ctx context.Context, st *Stats) error {
	for i := 0; i < maxFanOutPerRun; i++ {
		more, err := p.fanOutOne(ctx)
		if err != nil {
			return err
		}
		if !more {
			return nil
		}
		st.FannedOut++
	}
	return nil
}

// fanOutOne turns the oldest due message into deliveries. The message row stays locked until the commit, so two
// workers never fan out the same message, and a crash leaves it pending for the next run.
func (p *Pipeline) fanOutOne(ctx context.Context) (bool, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback() }()

	var id int64
	var category string
	var audienceJSON []byte
	var expired bool
	err = tx.QueryRowContext(ctx, `
		SELECT id, category, audience, expires_at <= now()
		FROM push_messages
		WHERE status = 'pending' AND not_before <= now()
		ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED`).Scan(&id, &category, &audienceJSON, &expired)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	var aud Audience
	if expired || json.Unmarshal(audienceJSON, &aud) != nil || aud.Validate() != nil {
		reason := reasonMessageExpired
		if !expired {
			reason = reasonInvalidAudience
		}
		filtered, _ := json.Marshal(map[string]int{reason: 1})
		if _, err := tx.ExecContext(ctx, `
			UPDATE push_messages SET status = 'cancelled', filtered = $2, completed_at = now() WHERE id = $1`,
			id, string(filtered)); err != nil {
			return false, err
		}
		slog.Warn("push.cancelled", "message", id, "reason", reason)
		return true, tx.Commit()
	}

	res, err := p.resolve(ctx, tx, aud, models.NotificationCategory(category))
	if err != nil {
		return false, err
	}
	if len(res.Targets) > 0 {
		ids, emails, langs := make([]string, len(res.Targets)), make([]string, len(res.Targets)), make([]string, len(res.Targets))
		for i, t := range res.Targets {
			ids[i], emails[i], langs[i] = t.DeviceID, t.Email, t.Language
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO push_deliveries (message_id, device_id, user_email, language)
			SELECT $1, d, e, l FROM unnest($2::uuid[], $3::text[], $4::text[]) AS t(d, e, l)
			ON CONFLICT (message_id, device_id) DO NOTHING`,
			id, pq.Array(ids), pq.Array(emails), pq.Array(langs)); err != nil {
			return false, err
		}
	}
	filtered, _ := json.Marshal(res.Filtered)
	if _, err := tx.ExecContext(ctx, `
		UPDATE push_messages
		SET status = 'fanned_out', recipients = $2, eligible_devices = $3, filtered = $4, fanned_out_at = now()
		WHERE id = $1`, id, res.Users, len(res.Targets), string(filtered)); err != nil {
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	slog.Info("push.fanout", "message", id, "users", res.Users, "eligible_devices", len(res.Targets),
		"users_without_target", len(res.WithoutTarget), "filtered", string(filtered))
	return true, nil
}

// ---- Dispatch ----

// item is a delivery that passed the eligibility check and is ready to go to the provider.
type item struct {
	id          int64
	messageID   int64
	deviceID    string
	attempts    int
	fingerprint string
	out         Outbound
	expired     bool
}

type claim struct {
	id, messageID int64
	deviceID      string
	email         string
	language      string
	attempts      int
}

func (p *Pipeline) dispatch(ctx context.Context, st *Stats) error {
	if _, err := p.db.ExecContext(ctx, `
		UPDATE push_deliveries
		SET status = 'failed', error_code = $2, lease_until = NULL, finished_at = now()
		WHERE status = 'sending' AND lease_until < now() AND attempts >= $1`,
		p.cfg.MaxAttempts, reasonInterruptedWorker); err != nil {
		return err
	}
	for round := 0; round < maxRoundsPerRun; round++ {
		claims, err := p.claim(ctx)
		if err != nil {
			return err
		}
		if len(claims) == 0 {
			return nil
		}
		if err := p.process(ctx, claims, st); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return nil
}

// claim takes the deliveries that are due, and the ones whose worker died, and counts an attempt for each so a
// delivery that keeps killing its worker cannot loop forever.
func (p *Pipeline) claim(ctx context.Context) ([]claim, error) {
	rows, err := p.db.QueryContext(ctx, `
		UPDATE push_deliveries d
		SET status = 'sending', attempts = d.attempts + 1, first_attempt_at = COALESCE(d.first_attempt_at, now()),
		    lease_until = now() + make_interval(secs => $1)
		WHERE d.id IN (
			SELECT id FROM push_deliveries
			WHERE attempts < $2
			  AND ((status = 'pending' AND next_attempt_at <= now()) OR (status = 'sending' AND lease_until < now()))
			ORDER BY next_attempt_at, id
			LIMIT $3
			FOR UPDATE SKIP LOCKED)
		RETURNING d.id, d.message_id, d.device_id, d.user_email, d.language, d.attempts`,
		p.cfg.Lease.Seconds(), p.cfg.MaxAttempts, p.cfg.ClaimSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []claim
	for rows.Next() {
		var c claim
		if err := rows.Scan(&c.id, &c.messageID, &c.deviceID, &c.email, &c.language, &c.attempts); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

type recheck struct {
	owner      string
	token      string
	hasToken   bool
	platform   string
	permission string
	state      string
	content    []byte
	navigation []byte
	kind       string
	imageURL   string
	action     string
	expired    bool
	remaining  float64
	prefOn     bool
}

// process checks every claimed delivery again against the current state of its device, its owner's preferences
// and the message's deadline, then sends the ones that still qualify.
func (p *Pipeline) process(ctx context.Context, claims []claim, st *Stats) error {
	st.Claimed += len(claims)
	ids := make([]int64, len(claims))
	for i, c := range claims {
		ids[i] = c.id
	}
	rows, err := p.db.QueryContext(ctx, `
		SELECT del.id, d.user_email, COALESCE(d.push_token, ''), d.push_token IS NOT NULL, d.platform, d.permission_status,
		       `+devices.StateSQL("$2", "$3")+`,
		       m.content, m.navigation, m.kind, COALESCE(m.image_url, ''), COALESCE(m.action_category, ''), m.expires_at <= now(),
		       GREATEST(EXTRACT(EPOCH FROM (m.expires_at - now())), 0)::float8,
		       (m.category = '' OR EXISTS (
		           SELECT 1 FROM notifications nt JOIN services s ON s.id_services = nt.id_services
		           WHERE nt.email = d.user_email AND s.name = m.category))
		FROM push_deliveries del
		JOIN push_messages m ON m.id = del.message_id
		JOIN devices d ON d.installation_id = del.device_id
		JOIN newf n ON n.email = d.user_email
		WHERE del.id = ANY($1)`,
		pq.Array(ids), p.deviceCfg.SessionTimeout.Seconds(), p.deviceCfg.PresenceWindow.Seconds())
	if err != nil {
		return err
	}
	current := make(map[int64]recheck, len(claims))
	for rows.Next() {
		var id int64
		var r recheck
		if err := rows.Scan(&id, &r.owner, &r.token, &r.hasToken, &r.platform, &r.permission, &r.state,
			&r.content, &r.navigation, &r.kind, &r.imageURL, &r.action, &r.expired, &r.remaining, &r.prefOn); err != nil {
			rows.Close()
			return err
		}
		current[id] = r
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	skipped := map[string]int{}
	var items []*item
	for _, c := range claims {
		r, ok := current[c.id]
		if !ok {
			// The device or the account is gone.
			if err := p.settle(ctx, c.id, "skipped", reasonOwnerChanged); err != nil {
				return err
			}
			skipped[reasonOwnerChanged]++
			continue
		}
		status, reason := "skipped", ""
		switch ok, why := eligibility(r.state, r.hasToken, r.permission); {
		case r.expired:
			status, reason = "expired", reasonMessageExpired
		case r.owner != c.email:
			reason = reasonOwnerChanged
		case !ok:
			reason = why
		case !r.prefOn:
			reason = reasonPreferenceOff
		}
		if reason != "" {
			if err := p.settle(ctx, c.id, status, reason); err != nil {
				return err
			}
			skipped[reason]++
			continue
		}

		text, nav, err := decodeContent(r, c.language)
		if err != nil {
			if err := p.fail(ctx, c.id, "no_content", err.Error(), ""); err != nil {
				return err
			}
			st.Failed++
			continue
		}
		items = append(items, &item{
			id: c.id, messageID: c.messageID, deviceID: c.deviceID, attempts: c.attempts,
			fingerprint: devices.TokenFingerprint(r.token),
			out: Outbound{
				Token: r.token, Platform: r.platform, Title: text.Title, Body: text.Body,
				Data:           models.NotificationPayload{Kind: r.kind, Navigation: nav}.ExpoData(),
				ImageURL:       r.imageURL,
				ActionCategory: r.action,
				TTL:            time.Duration(r.remaining * float64(time.Second)),
				DedupeTag:      fmt.Sprintf("m%d", c.messageID),
			},
		})
	}
	for _, n := range skipped {
		st.Skipped += n
	}
	if len(skipped) > 0 {
		slog.Info("push.skipped", "reasons", fmt.Sprint(skipped))
	}

	for start := 0; start < len(items); start += p.cfg.BatchSize {
		end := min(start+p.cfg.BatchSize, len(items))
		if err := p.sendChunk(ctx, items[start:end], st); err != nil {
			return err
		}
	}
	return nil
}

func decodeContent(r recheck, language string) (Text, *models.NavigationTarget, error) {
	var content Content
	if err := json.Unmarshal(r.content, &content); err != nil {
		return Text{}, nil, err
	}
	text, ok := content[language]
	if !ok || text.Title == "" {
		text, ok = content[DefaultLanguage]
	}
	if !ok || text.Title == "" {
		return Text{}, nil, ErrNoContent
	}
	var nav *models.NavigationTarget
	if len(r.navigation) > 0 {
		nav = &models.NavigationTarget{}
		if err := json.Unmarshal(r.navigation, nav); err != nil {
			return Text{}, nil, err
		}
	}
	return text, nav, nil
}

// sendChunk gives one batch to the provider and records what came back. A request the provider refuses for good
// is split in two until the message that causes it is alone, so one bad message cannot hold the others back.
func (p *Pipeline) sendChunk(ctx context.Context, items []*item, st *Stats) error {
	messages := make([]Outbound, len(items))
	for i, it := range items {
		messages[i] = it.out
	}
	start := time.Now()
	results, err := p.provider.Send(ctx, messages)
	p.pace(ctx, len(items), time.Since(start))

	if err != nil {
		var re *RequestError
		switch {
		case errors.As(err, &re) && re.Isolate && len(items) > 1:
			mid := len(items) / 2
			if err := p.sendChunk(ctx, items[:mid], st); err != nil {
				return err
			}
			return p.sendChunk(ctx, items[mid:], st)
		case errors.As(err, &re) && !re.Retryable:
			slog.Warn("push.request_rejected", "provider", p.provider.Name(), "code", re.Code, "messages", len(items))
			for _, it := range items {
				if err := p.fail(ctx, it.id, re.Code, re.Detail, it.fingerprint); err != nil {
					return err
				}
				st.Failed++
			}
		default:
			code, detail := "provider_error", redact(err.Error())
			if re != nil {
				code, detail = re.Code, re.Detail
			}
			slog.Warn("push.request_failed", "provider", p.provider.Name(), "code", code, "messages", len(items))
			for _, it := range items {
				if err := p.retryOrFail(ctx, it, code, detail, st); err != nil {
					return err
				}
			}
		}
		return ctx.Err()
	}
	if len(results) != len(items) {
		return fmt.Errorf("provider %s returned %d results for %d messages", p.provider.Name(), len(results), len(items))
	}

	var accepted, retry, invalid, rejected int
	for i, r := range results {
		it := items[i]
		switch r.Kind {
		case ResultAccepted:
			accepted++
			st.Accepted++
			if err := p.accept(ctx, it, r.TicketID); err != nil {
				return err
			}
		case ResultInvalidDevice:
			invalid++
			if err := p.invalidate(ctx, it, r.Code, r.Detail, st); err != nil {
				return err
			}
		case ResultRetry:
			retry++
			if err := p.retryOrFail(ctx, it, r.Code, r.Detail, st); err != nil {
				return err
			}
		default:
			rejected++
			st.Failed++
			if err := p.fail(ctx, it.id, r.Code, r.Detail, it.fingerprint); err != nil {
				return err
			}
		}
	}
	slog.Info("push.attempt", "provider", p.provider.Name(), "messages", len(items),
		"accepted", accepted, "retry", retry, "invalid_device", invalid, "rejected", rejected)
	return nil
}

// pace keeps the rate handed to the provider under MaxPerSecond.
func (p *Pipeline) pace(ctx context.Context, sent int, took time.Duration) {
	if p.cfg.MaxPerSecond <= 0 {
		return
	}
	wait := time.Duration(float64(sent)/float64(p.cfg.MaxPerSecond)*float64(time.Second)) - took
	if wait <= 0 {
		return
	}
	select {
	case <-ctx.Done():
	case <-time.After(wait):
	}
}

// The updates below all end with "status = 'sending'": a worker whose lease was taken over must not overwrite
// what the new owner recorded.

func (p *Pipeline) accept(ctx context.Context, it *item, ticket string) error {
	_, err := p.db.ExecContext(ctx, `
		UPDATE push_deliveries
		SET status = CASE WHEN $2::text = '' THEN 'unconfirmed' ELSE 'accepted' END,
		    ticket_id = NULLIF($2::text, ''), token_fingerprint = $3, accepted_at = now(),
		    receipt_due_at = CASE WHEN $2::text = '' THEN NULL ELSE now() + make_interval(secs => $4) END,
		    finished_at = CASE WHEN $2::text = '' THEN now() END,
		    lease_until = NULL, error_code = NULL, error_detail = NULL
		WHERE id = $1 AND status = 'sending'`,
		it.id, ticket, it.fingerprint, p.cfg.ReceiptDelay.Seconds())
	return err
}

func (p *Pipeline) fail(ctx context.Context, id int64, code, detail, fingerprint string) error {
	_, err := p.db.ExecContext(ctx, `
		UPDATE push_deliveries
		SET status = 'failed', error_code = $2, error_detail = $3, token_fingerprint = NULLIF($4, ''),
		    lease_until = NULL, finished_at = now()
		WHERE id = $1 AND status = 'sending'`, id, code, redact(detail), fingerprint)
	return err
}

func (p *Pipeline) settle(ctx context.Context, id int64, status, reason string) error {
	_, err := p.db.ExecContext(ctx, `
		UPDATE push_deliveries
		SET status = $2, skip_reason = $3, lease_until = NULL, finished_at = now()
		WHERE id = $1 AND status = 'sending'`, id, status, reason)
	return err
}

// retryOrFail schedules the next attempt after a transient failure, or gives up when the attempts are used up
// or the message has expired.
func (p *Pipeline) retryOrFail(ctx context.Context, it *item, code, detail string, st *Stats) error {
	if it.attempts >= p.cfg.MaxAttempts {
		st.Failed++
		return p.fail(ctx, it.id, reasonRetriesExhausted+":"+code, detail, it.fingerprint)
	}
	delay := Backoff(it.attempts, p.cfg.RetryBase, p.cfg.RetryMax, p.jitter())
	res, err := p.db.ExecContext(ctx, `
		UPDATE push_deliveries d
		SET status = 'pending', next_attempt_at = now() + make_interval(secs => $2), lease_until = NULL,
		    error_code = $3, error_detail = $4, token_fingerprint = NULLIF($5, '')
		WHERE d.id = $1 AND d.status = 'sending'`,
		it.id, delay.Seconds(), code, redact(detail), it.fingerprint)
	if err == nil {
		if n, _ := res.RowsAffected(); n > 0 {
			st.Retried++
		}
	}
	return err
}

// invalidate records that the provider confirmed the token as dead and deactivates the device, unless the app
// has refreshed its token in the meantime.
func (p *Pipeline) invalidate(ctx context.Context, it *item, code, detail string, st *Stats) error {
	st.Failed++
	if err := p.fail(ctx, it.id, code, detail, it.fingerprint); err != nil {
		return err
	}
	disabled, err := p.deviceSvc.InvalidateToken(ctx, it.deviceID, it.fingerprint)
	if err != nil {
		return err
	}
	if disabled {
		st.Invalidated++
		slog.Info("push.device_invalidated", "device", it.deviceID, "code", code)
	}
	return nil
}

// ---- Receipts ----

type pendingReceipt struct {
	id          int64
	ticket      string
	deviceID    string
	fingerprint string
	attempts    int
	old         bool
	expired     bool
}

func (p *Pipeline) collectReceipts(ctx context.Context, st *Stats) error {
	for round := 0; round < maxReceiptRounds; round++ {
		rows, err := p.db.QueryContext(ctx, `
			UPDATE push_deliveries d
			SET lease_until = now() + make_interval(secs => $1)
			WHERE d.id IN (
				SELECT id FROM push_deliveries
				WHERE status = 'accepted' AND ticket_id IS NOT NULL AND receipt_due_at <= now()
				  AND (lease_until IS NULL OR lease_until < now())
				ORDER BY receipt_due_at
				LIMIT $2
				FOR UPDATE SKIP LOCKED)
			RETURNING d.id, d.ticket_id, d.device_id, COALESCE(d.token_fingerprint, ''), d.attempts,
			          d.accepted_at < now() - make_interval(secs => $3),
			          (SELECT m.expires_at <= now() FROM push_messages m WHERE m.id = d.message_id)`,
			p.cfg.Lease.Seconds(), MaxReceiptBatch, p.cfg.ReceiptTTL.Seconds())
		if err != nil {
			return err
		}
		var batch []pendingReceipt
		for rows.Next() {
			var r pendingReceipt
			if err := rows.Scan(&r.id, &r.ticket, &r.deviceID, &r.fingerprint, &r.attempts, &r.old, &r.expired); err != nil {
				rows.Close()
				return err
			}
			batch = append(batch, r)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return err
		}
		if len(batch) == 0 {
			return nil
		}
		if err := p.applyReceipts(ctx, batch, st); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pipeline) applyReceipts(ctx context.Context, batch []pendingReceipt, st *Stats) error {
	tickets := make([]string, len(batch))
	ids := make([]int64, len(batch))
	for i, r := range batch {
		tickets[i], ids[i] = r.ticket, r.id
	}
	got, err := p.provider.Receipts(ctx, tickets)
	if err != nil {
		slog.Warn("push.receipts_failed", "provider", p.provider.Name(), "tickets", len(batch), "error", redact(err.Error()))
		// Ask again later; the wait is not a retry of the message.
		_, uerr := p.db.ExecContext(ctx, `
			UPDATE push_deliveries SET lease_until = NULL, receipt_due_at = now() + make_interval(secs => $2)
			WHERE id = ANY($1) AND status = 'accepted'`, pq.Array(ids), p.cfg.ReceiptRetry.Seconds())
		if uerr != nil {
			return uerr
		}
		return err
	}

	var handedOff, failed, missing, retry int
	for _, r := range batch {
		res, ok := got[r.ticket]
		var uerr error
		switch {
		case !ok && r.old:
			missing++
			uerr = p.finishReceipt(ctx, r.id, "unconfirmed", "receipt_unavailable", "")
		case !ok:
			missing++
			_, uerr = p.db.ExecContext(ctx, `
				UPDATE push_deliveries SET lease_until = NULL, receipt_due_at = now() + make_interval(secs => $2)
				WHERE id = $1 AND status = 'accepted'`, r.id, p.cfg.ReceiptRetry.Seconds())
		case res.Kind == ReceiptHandedOff:
			handedOff++
			uerr = p.finishReceipt(ctx, r.id, "handed_off", "", "")
		case res.Kind == ReceiptInvalidDevice:
			failed++
			uerr = p.finishReceipt(ctx, r.id, "failed", res.Code, res.Detail)
			if uerr == nil && r.fingerprint != "" {
				disabled, derr := p.deviceSvc.InvalidateToken(ctx, r.deviceID, r.fingerprint)
				if derr != nil {
					uerr = derr
				} else if disabled {
					st.Invalidated++
					slog.Info("push.device_invalidated", "device", r.deviceID, "code", res.Code)
				}
			}
		case res.Kind == ReceiptRetry && r.attempts < p.cfg.MaxAttempts && !r.expired:
			retry++
			st.Retried++
			delay := Backoff(r.attempts, p.cfg.RetryBase, p.cfg.RetryMax, p.jitter())
			_, uerr = p.db.ExecContext(ctx, `
				UPDATE push_deliveries
				SET status = 'pending', next_attempt_at = now() + make_interval(secs => $2), lease_until = NULL,
				    ticket_id = NULL, receipt_due_at = NULL, error_code = $3, error_detail = $4
				WHERE id = $1 AND status = 'accepted'`, r.id, delay.Seconds(), res.Code, redact(res.Detail))
		default:
			failed++
			code := res.Code
			if res.Kind == ReceiptRetry {
				code = reasonRetriesExhausted + ":" + code
			}
			if res.Code == "InvalidCredentials" || res.Code == "MismatchSenderId" {
				slog.Error("push.provider_credentials", "provider", p.provider.Name(), "code", res.Code)
			}
			uerr = p.finishReceipt(ctx, r.id, "failed", code, res.Detail)
		}
		if uerr != nil {
			return uerr
		}
	}
	st.Receipts += len(batch)
	slog.Info("push.receipts", "provider", p.provider.Name(), "tickets", len(batch), "handed_off", handedOff,
		"failed", failed, "retry", retry, "not_available_yet", missing)
	return nil
}

func (p *Pipeline) finishReceipt(ctx context.Context, id int64, status, code, detail string) error {
	_, err := p.db.ExecContext(ctx, `
		UPDATE push_deliveries
		SET status = $2, error_code = NULLIF($3, ''), error_detail = NULLIF($4, ''), receipt_at = now(),
		    lease_until = NULL, finished_at = now()
		WHERE id = $1 AND status = 'accepted'`, id, status, code, redact(detail))
	return err
}

// ---- Completion and retention ----

// complete marks as completed the messages whose deliveries are all settled and logs how each one ended.
func (p *Pipeline) complete(ctx context.Context, st *Stats) error {
	rows, err := p.db.QueryContext(ctx, `
		UPDATE push_messages m SET status = 'completed', completed_at = now()
		WHERE m.status = 'fanned_out'
		  AND NOT EXISTS (SELECT 1 FROM push_deliveries d WHERE d.message_id = m.id AND d.status IN ('pending', 'sending', 'accepted'))
		RETURNING m.id`)
	if err != nil {
		return err
	}
	var done []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		done = append(done, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	for _, id := range done {
		st.Completed++
		counts, err := p.countDeliveries(ctx, id)
		if err != nil {
			return err
		}
		slog.Info("push.completed", "message", id, "deliveries", fmt.Sprint(counts))
	}
	return nil
}

func (p *Pipeline) countDeliveries(ctx context.Context, messageID int64) (map[string]int, error) {
	rows, err := p.db.QueryContext(ctx, `SELECT status, count(*) FROM push_deliveries WHERE message_id = $1 GROUP BY status`, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := map[string]int{}
	for rows.Next() {
		var status string
		var n int
		if err := rows.Scan(&status, &n); err != nil {
			return nil, err
		}
		counts[status] = n
	}
	return counts, rows.Err()
}

// purge deletes finished messages, and with them their deliveries, once they are older than the retention.
func (p *Pipeline) purge(ctx context.Context) error {
	res, err := p.db.ExecContext(ctx, `
		DELETE FROM push_messages
		WHERE status IN ('completed', 'cancelled') AND COALESCE(completed_at, created_at) < now() - make_interval(secs => $1)`,
		p.cfg.Retention.Seconds())
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		slog.Info("push.purged", "messages", n)
	}
	return nil
}
