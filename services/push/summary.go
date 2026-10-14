package push

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

// ErrMessageNotFound is returned when a message id does not exist.
var ErrMessageNotFound = errors.New("push message not found")

// Summary is the diagnostic view of one message: what was asked, who it could go to, and how each delivery
// stands. It holds counts and codes only, never tokens, emails or the text of the notification.
type Summary struct {
	ID        int64  `json:"id"`
	Kind      string `json:"kind"`
	Category  string `json:"category,omitempty"`
	Status    string `json:"status"`
	CreatedBy string `json:"createdBy,omitempty"`

	CreatedAt   time.Time  `json:"createdAt"`
	ExpiresAt   time.Time  `json:"expiresAt"`
	FannedOutAt *time.Time `json:"fannedOutAt,omitempty"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`

	// Users is the audience after the preferences; EligibleDevices the devices that passed the session policy at
	// fan-out; Filtered why the others did not (user-level reason: no_device).
	Users           *int           `json:"users,omitempty"`
	EligibleDevices *int           `json:"eligibleDevices,omitempty"`
	Filtered        map[string]int `json:"filtered,omitempty"`

	// Deliveries counts the devices by status. "accepted" and "handed_off" mean the provider took the message;
	// no status says that a phone received it.
	Deliveries map[string]int `json:"deliveries"`
	// Attempts is the number of calls made to the provider for this message, retries included.
	Attempts int `json:"attempts"`
	// Skipped counts deliveries not sent because the device stopped qualifying, by reason.
	Skipped map[string]int `json:"skipped,omitempty"`
	// Errors counts the failed and retried deliveries by error code.
	Errors map[string]int `json:"errors,omitempty"`
}

// Summary returns the diagnostics of one message.
func (p *Pipeline) Summary(ctx context.Context, id int64) (Summary, error) {
	var s Summary
	var filtered []byte
	var users, eligible sql.NullInt64
	var createdBy sql.NullString
	err := p.db.QueryRowContext(ctx, `
		SELECT id, kind, category, status, created_by, created_at, expires_at, fanned_out_at, completed_at,
		       recipients, eligible_devices, filtered
		FROM push_messages WHERE id = $1`, id).
		Scan(&s.ID, &s.Kind, &s.Category, &s.Status, &createdBy, &s.CreatedAt, &s.ExpiresAt, &s.FannedOutAt,
			&s.CompletedAt, &users, &eligible, &filtered)
	if errors.Is(err, sql.ErrNoRows) {
		return Summary{}, ErrMessageNotFound
	}
	if err != nil {
		return Summary{}, err
	}
	s.CreatedBy = createdBy.String
	if users.Valid {
		n := int(users.Int64)
		s.Users = &n
	}
	if eligible.Valid {
		n := int(eligible.Int64)
		s.EligibleDevices = &n
	}
	if len(filtered) > 0 {
		_ = json.Unmarshal(filtered, &s.Filtered)
	}

	if s.Deliveries, err = p.groupCount(ctx, `SELECT status, count(*) FROM push_deliveries WHERE message_id = $1 GROUP BY status`, id); err != nil {
		return Summary{}, err
	}
	if s.Skipped, err = p.groupCount(ctx, `SELECT skip_reason, count(*) FROM push_deliveries WHERE message_id = $1 AND skip_reason IS NOT NULL GROUP BY skip_reason`, id); err != nil {
		return Summary{}, err
	}
	if s.Errors, err = p.groupCount(ctx, `SELECT error_code, count(*) FROM push_deliveries WHERE message_id = $1 AND error_code IS NOT NULL GROUP BY error_code`, id); err != nil {
		return Summary{}, err
	}
	err = p.db.QueryRowContext(ctx, `SELECT COALESCE(sum(attempts), 0) FROM push_deliveries WHERE message_id = $1`, id).Scan(&s.Attempts)
	return s, err
}

// Recent returns the diagnostics of the latest messages, newest first.
func (p *Pipeline) Recent(ctx context.Context, limit int) ([]Summary, error) {
	limit = max(1, min(limit, 100))
	rows, err := p.db.QueryContext(ctx, `SELECT id FROM push_messages ORDER BY id DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	out := make([]Summary, 0, len(ids))
	for _, id := range ids {
		s, err := p.Summary(ctx, id)
		if errors.Is(err, ErrMessageNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

func (p *Pipeline) groupCount(ctx context.Context, query string, id int64) (map[string]int, error) {
	rows, err := p.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var key string
		var n int
		if err := rows.Scan(&key, &n); err != nil {
			return nil, err
		}
		out[key] = n
	}
	return out, rows.Err()
}
