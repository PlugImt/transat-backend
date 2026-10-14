// Package devices keeps one row per app installation and tells whether a user is connected.
//
// Three different things are kept apart on purpose:
//
//   - A push token only addresses an installation. Having one says nothing about who uses the phone.
//   - A session is what the authentication gives: the JWT is stateless and lives for years, so the server cannot
//     list or revoke sessions. The server only knows that an installation proved a valid session at last_seen_at.
//   - A connection is a recent proof. Each authenticated registration call is a heartbeat that refreshes last_seen_at.
//
// The state of a device is computed from last_seen_at at read time, never stored, so nothing has to run when a
// session ends. In this order, the first matching rule wins:
//
//	disabled   the app signed out, or another installation took over the push token
//	revoked    the last proof is older than the account's last password change (the JWT middleware refuses the
//	           tokens issued before it, so the installation cannot heartbeat any more)
//	expired    no proof for SessionTimeout (DEVICE_SESSION_TIMEOUT, default 60 days)
//	connected  a proof within PresenceWindow (DEVICE_PRESENCE_WINDOW, default 10 minutes): the app is open now
//	active     a valid session that is not in use right now (the app is in the background or closed)
//
// A user is connected when at least one device is connected. A user has an active session when at least one device
// is connected or active; only those devices may be pushed to, and an expired or revoked one never is.
package devices

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/lib/pq"
)

// State is the connection state of one device, see the package documentation.
type State string

const (
	StateConnected State = "connected"
	StateActive    State = "active"
	StateExpired   State = "expired"
	StateRevoked   State = "revoked"
	StateDisabled  State = "disabled"
)

const (
	defaultSessionTimeout    = 60 * 24 * time.Hour
	defaultPresenceWindow    = 10 * time.Minute
	defaultHeartbeatInterval = 5 * time.Minute
	defaultWriteThrottle     = time.Minute

	maxRegisterAttempts = 5
)

var (
	ErrInvalidInstallationID = errors.New("invalid installation id")
	ErrInvalidPlatform       = errors.New("invalid platform")
	ErrInvalidPushToken      = errors.New("invalid push token")
	ErrInvalidPermission     = errors.New("invalid permission status")
	ErrInvalidAppVersion     = errors.New("invalid app version")

	installationIDPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	expoTokenPattern      = regexp.MustCompile(`^(?:Expo|Exponent)PushToken\[[A-Za-z0-9_-]{1,200}\]$`)
)

// Config holds the timings of the connection state.
type Config struct {
	// SessionTimeout is how long an installation stays a valid session without any proof.
	SessionTimeout time.Duration
	// PresenceWindow is how recent a proof must be for the app to count as open right now.
	PresenceWindow time.Duration
	// HeartbeatInterval is how often an open app is asked to send a proof; it is returned to the app so it can
	// be tuned without a release. It must leave room for a missed heartbeat inside PresenceWindow.
	HeartbeatInterval time.Duration
	// WriteThrottle is the shortest time between two writes for an unchanged installation.
	WriteThrottle time.Duration
}

// DefaultConfig returns the timings used when nothing is configured.
func DefaultConfig() Config {
	return Config{
		SessionTimeout:    defaultSessionTimeout,
		PresenceWindow:    defaultPresenceWindow,
		HeartbeatInterval: defaultHeartbeatInterval,
		WriteThrottle:     defaultWriteThrottle,
	}
}

// ConfigFromEnv reads DEVICE_SESSION_TIMEOUT, DEVICE_PRESENCE_WINDOW, DEVICE_HEARTBEAT_INTERVAL and
// DEVICE_WRITE_THROTTLE (Go durations such as "90m" or "1440h"). An inconsistent set falls back to the defaults.
func ConfigFromEnv() Config {
	cfg := Config{
		SessionTimeout:    envDuration("DEVICE_SESSION_TIMEOUT", defaultSessionTimeout),
		PresenceWindow:    envDuration("DEVICE_PRESENCE_WINDOW", defaultPresenceWindow),
		HeartbeatInterval: envDuration("DEVICE_HEARTBEAT_INTERVAL", defaultHeartbeatInterval),
		WriteThrottle:     envDuration("DEVICE_WRITE_THROTTLE", defaultWriteThrottle),
	}
	if err := cfg.Validate(); err != nil {
		log.Printf("Warning: invalid device timings (%v), using the defaults", err)
		return DefaultConfig()
	}
	return cfg
}

// Validate rejects timings under which an open app would flap between connected and not connected.
func (c Config) Validate() error {
	switch {
	case c.SessionTimeout <= 0 || c.PresenceWindow <= 0 || c.HeartbeatInterval <= 0 || c.WriteThrottle <= 0:
		return errors.New("every duration must be positive")
	case c.PresenceWindow >= c.SessionTimeout:
		return errors.New("the presence window must be shorter than the session timeout")
	case c.HeartbeatInterval*2 > c.PresenceWindow:
		return errors.New("the heartbeat interval must be at most half of the presence window")
	case c.WriteThrottle >= c.HeartbeatInterval:
		return errors.New("the write throttle must be shorter than the heartbeat interval")
	}
	return nil
}

func envDuration(name string, def time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	d, err := time.ParseDuration(raw)
	if err != nil || d <= 0 {
		log.Printf("Warning: ignoring %s=%q", name, raw)
		return def
	}
	return d
}

// Service reads and writes the devices table.
type Service struct {
	db  *sql.DB
	cfg Config
}

func NewService(db *sql.DB, cfg Config) *Service {
	return &Service{db: db, cfg: cfg}
}

// Config returns the timings in use.
func (s *Service) Config() Config { return s.cfg }

// Registration is what an installation reports about itself on every heartbeat.
type Registration struct {
	InstallationID string
	// Email is the authenticated user, never taken from the request body.
	Email string
	// Platform is "ios" or "android".
	Platform string
	// PushToken is empty when the installation has none (permission refused, token not obtained).
	PushToken        string
	PermissionStatus string
	AppVersion       string
}

func (r Registration) normalized() (Registration, error) {
	r.InstallationID = strings.ToLower(strings.TrimSpace(r.InstallationID))
	if !installationIDPattern.MatchString(r.InstallationID) {
		return r, ErrInvalidInstallationID
	}
	r.Platform = strings.ToLower(strings.TrimSpace(r.Platform))
	if r.Platform != "ios" && r.Platform != "android" {
		return r, ErrInvalidPlatform
	}
	r.PushToken = strings.TrimSpace(r.PushToken)
	if r.PushToken != "" && !expoTokenPattern.MatchString(r.PushToken) {
		return r, ErrInvalidPushToken
	}
	r.PermissionStatus = strings.ToLower(strings.TrimSpace(r.PermissionStatus))
	switch r.PermissionStatus {
	case "":
		r.PermissionStatus = "unknown"
	case "granted", "denied", "undetermined", "unknown":
	default:
		return r, ErrInvalidPermission
	}
	r.AppVersion = strings.TrimSpace(r.AppVersion)
	if len(r.AppVersion) > 32 {
		return r, ErrInvalidAppVersion
	}
	return r, nil
}

// Register records that the authenticated user is using this installation right now. It creates the row, takes
// it over from a previous account, follows a refreshed token, and refreshes last_seen_at. It reports whether the
// database was written: an unchanged installation that checked in less than WriteThrottle ago is left alone.
//
// The token and the installation are both unique: a token that another installation still holds is taken away
// from it (that installation is stale, a phone cannot have two), so a token never reaches two installations.
//
// user_notification_tokens is the table the current dispatch reads. It mirrors the tokens of active devices here
// so that pushes keep working, and so that a token never stays attached to the previous account of a phone.
func (s *Service) Register(ctx context.Context, in Registration) (bool, error) {
	r, err := in.normalized()
	if err != nil {
		return false, err
	}
	for attempt := 1; ; attempt++ {
		written, err := s.register(ctx, r)
		if err != nil && attempt < maxRegisterAttempts && isRaceError(err) {
			continue
		}
		return written, err
	}
}

type deviceRow struct {
	owner      string
	token      string
	platform   string
	permission string
	version    string
	status     string
	recent     bool
}

func (s *Service) register(ctx context.Context, r Registration) (written bool, err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback() }()

	var cur deviceRow
	var token sql.NullString
	err = tx.QueryRowContext(ctx, `
		SELECT d.user_email, d.push_token, d.platform, d.permission_status, COALESCE(d.app_version, ''), d.status,
		       -- A revoked session is never "recent": the first registration after a password change must write.
		       d.last_seen_at > now() - make_interval(secs => $2)
		       AND d.last_seen_at > COALESCE(n.password_updated_date::timestamptz, '-infinity'::timestamptz)
		FROM devices d LEFT JOIN newf n ON n.email = d.user_email
		WHERE d.installation_id = $1 FOR UPDATE OF d`,
		r.InstallationID, s.cfg.WriteThrottle.Seconds()).
		Scan(&cur.owner, &token, &cur.platform, &cur.permission, &cur.version, &cur.status, &cur.recent)
	exists := err == nil
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	cur.token = token.String

	unchanged := exists && cur.status == "active" && cur.owner == r.Email && cur.token == r.PushToken &&
		cur.platform == r.Platform && cur.permission == r.PermissionStatus && cur.version == r.AppVersion
	if unchanged && cur.recent {
		return false, tx.Commit()
	}

	if r.PushToken != "" {
		if _, err := tx.ExecContext(ctx, `
			UPDATE devices
			SET push_token = NULL, status = 'disabled', disabled_reason = 'token_replaced', updated_at = now()
			WHERE provider = 'expo' AND push_token = $1 AND installation_id <> $2`,
			r.PushToken, r.InstallationID); err != nil {
			return false, err
		}
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM user_notification_tokens WHERE token = $1 AND email <> $2`, r.PushToken, r.Email); err != nil {
			return false, err
		}
	}
	if exists && cur.token != "" && (cur.token != r.PushToken || cur.owner != r.Email) {
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM user_notification_tokens WHERE email = $1 AND token = $2`, cur.owner, cur.token); err != nil {
			return false, err
		}
	}

	var pushToken any
	if r.PushToken != "" {
		pushToken = r.PushToken
	}
	if exists {
		_, err = tx.ExecContext(ctx, `
			UPDATE devices
			SET user_email = $2, platform = $3, push_token = $4, permission_status = $5, app_version = NULLIF($6, ''),
			    status = 'active', disabled_reason = NULL, last_seen_at = now(),
			    updated_at = CASE WHEN $7 THEN updated_at ELSE now() END
			WHERE installation_id = $1`,
			r.InstallationID, r.Email, r.Platform, pushToken, r.PermissionStatus, r.AppVersion, unchanged)
	} else {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO devices (installation_id, user_email, platform, push_token, permission_status, app_version)
			VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''))`,
			r.InstallationID, r.Email, r.Platform, pushToken, r.PermissionStatus, r.AppVersion)
	}
	if err != nil {
		return false, err
	}
	if r.PushToken != "" {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO user_notification_tokens (email, token) VALUES ($1, $2) ON CONFLICT (email, token) DO NOTHING`,
			r.Email, r.PushToken); err != nil {
			return false, err
		}
	}
	return true, tx.Commit()
}

// SignOut ends the session of an installation: it is disabled and loses its token, so nothing more can be pushed
// to it for the account that was using it. It is keyed by the installation id alone, which only the installation
// and its owner's authenticated calls know, because it must work after the session has already expired or been
// revoked. It is idempotent and reports whether a device was disabled.
func (s *Service) SignOut(ctx context.Context, installationID string) (bool, error) {
	id := strings.ToLower(strings.TrimSpace(installationID))
	if !installationIDPattern.MatchString(id) {
		return false, ErrInvalidInstallationID
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback() }()

	var owner string
	var token sql.NullString
	err = tx.QueryRowContext(ctx,
		`SELECT user_email, push_token FROM devices WHERE installation_id = $1 AND status = 'active' FOR UPDATE`, id).
		Scan(&owner, &token)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE devices
		SET status = 'disabled', disabled_reason = 'signed_out', push_token = NULL, updated_at = now()
		WHERE installation_id = $1`, id); err != nil {
		return false, err
	}
	if token.Valid {
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM user_notification_tokens WHERE email = $1 AND token = $2`, owner, token.String); err != nil {
			return false, err
		}
	}
	return true, tx.Commit()
}

// Device is one installation of a user with its computed state. The push token is never exposed.
type Device struct {
	InstallationID   string    `json:"installation_id"`
	Platform         string    `json:"platform"`
	PermissionStatus string    `json:"permission_status"`
	AppVersion       string    `json:"app_version,omitempty"`
	HasPushToken     bool      `json:"has_push_token"`
	State            State     `json:"state"`
	DisabledReason   string    `json:"disabled_reason,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	LastSeenAt       time.Time `json:"last_seen_at"`
}

// StateSQL is the one place where the rules of the package documentation are written, as a SQL expression over
// the aliases d (devices) and n (newf). timeoutArg and presenceArg are the placeholders that hold the session
// timeout and the presence window in seconds. password_updated_date has no time zone and is compared the way the
// JWT middleware compares it, in the session time zone. A legacy device (a token registered by an app version
// without heartbeat) is as recent as its account's last activity. d.installation_id may be NULL in an outer join:
// the caller must check that before reading the state.
func StateSQL(timeoutArg, presenceArg string) string {
	const seen = `GREATEST(d.last_seen_at, CASE WHEN d.source = 'legacy' THEN n.last_activity::timestamptz END)`
	return `CASE
	WHEN d.status <> 'active' THEN 'disabled'
	WHEN ` + seen + ` <= COALESCE(n.password_updated_date::timestamptz, '-infinity'::timestamptz) THEN 'revoked'
	WHEN ` + seen + ` < now() - make_interval(secs => ` + timeoutArg + `) THEN 'expired'
	WHEN ` + seen + ` >= now() - make_interval(secs => ` + presenceArg + `) THEN 'connected'
	ELSE 'active' END`
}

// List returns every device of the user, most recently seen first.
func (s *Service) List(ctx context.Context, email string) ([]Device, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT d.installation_id, d.platform, d.permission_status, COALESCE(d.app_version, ''), d.push_token IS NOT NULL,
		       `+StateSQL("$2", "$3")+`, COALESCE(d.disabled_reason, ''), d.created_at, d.updated_at, d.last_seen_at
		FROM devices d JOIN newf n ON n.email = d.user_email
		WHERE d.user_email = $1
		ORDER BY d.last_seen_at DESC, d.installation_id`,
		email, s.cfg.SessionTimeout.Seconds(), s.cfg.PresenceWindow.Seconds())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Device{}
	for rows.Next() {
		var d Device
		if err := rows.Scan(&d.InstallationID, &d.Platform, &d.PermissionStatus, &d.AppVersion, &d.HasPushToken,
			&d.State, &d.DisabledReason, &d.CreatedAt, &d.UpdatedAt, &d.LastSeenAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// IsConnected reports whether the user has at least one connected device. Zero devices is simply not connected.
func (s *Service) IsConnected(ctx context.Context, email string) (bool, error) {
	devices, err := s.List(ctx, email)
	return AnyConnected(devices), err
}

// HasActiveSession reports whether the user has at least one connected or active device.
func (s *Service) HasActiveSession(ctx context.Context, email string) (bool, error) {
	devices, err := s.List(ctx, email)
	return AnyActiveSession(devices), err
}

// AnyConnected reports whether one of the devices is connected.
func AnyConnected(devices []Device) bool {
	for _, d := range devices {
		if d.State == StateConnected {
			return true
		}
	}
	return false
}

// AnyActiveSession reports whether one of the devices is connected or active.
func AnyActiveSession(devices []Device) bool {
	for _, d := range devices {
		if d.State == StateConnected || d.State == StateActive {
			return true
		}
	}
	return false
}

// isRaceError tells whether two concurrent registrations collided: a unique violation on the installation or the
// token, or a deadlock between two transactions that took each other's token.
func isRaceError(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && (pqErr.Code == "23505" || pqErr.Code == "40P01")
}
