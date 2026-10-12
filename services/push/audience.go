package push

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/lib/pq"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services/devices"
)

// Audience selects the users a message is for. The recipients are worked out when the message is fanned out.
type Audience struct {
	models.NotificationAudience
	// InterestedEventID and InterestedClubID narrow the audience to the attendees of the event or the members of
	// the club (event reminders). They go with the audience type "all".
	InterestedEventID int `json:"interestedEventId,omitempty"`
	InterestedClubID  int `json:"interestedClubId,omitempty"`
}

func (a Audience) interest() bool { return a.InterestedEventID > 0 || a.InterestedClubID > 0 }

// Validate checks the audience the way the admin API does, plus the interest filter.
func (a Audience) Validate() error {
	if a.interest() {
		if a.Type != models.AudienceAll {
			return errors.New("an interest audience must have the type all")
		}
		return nil
	}
	return a.NotificationAudience.Validate()
}

// Target is one device a message may be sent to.
type Target struct {
	DeviceID string
	Email    string
	Language string
}

// Why a device is not a target, as recorded in the message's `filtered` counts. The user-level reason is
// reasonNoDevice: the account never registered any device.
const (
	reasonNoDevice          = "no_device"
	reasonNoToken           = "no_token"
	reasonPermissionDenied  = "permission_denied"
	reasonMessageExpired    = "message_expired"
	reasonPreferenceOff     = "preference_off"
	reasonOwnerChanged      = "owner_changed"
	reasonInvalidAudience   = "invalid_audience"
	reasonRetriesExhausted  = "retries_exhausted"
	reasonInterruptedWorker = "interrupted"
)

// Resolution is who a message goes to right now.
type Resolution struct {
	// Users is the size of the audience after the notification preferences were applied.
	Users int
	// Targets are the devices that may be pushed to: a valid session (connected or active), a token, and a
	// permission that is not refused.
	Targets []Target
	// WithoutTarget lists the users in the audience who have no such device.
	WithoutTarget []string
	// Filtered counts the devices left out, by reason (disabled, revoked, expired, no_token, permission_denied),
	// and the users with no device at all (no_device).
	Filtered map[string]int
}

type queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// eligibility applies the session policy to one device row: the state comes from devices.StateSQL.
func eligibility(state string, hasToken bool, permission string) (ok bool, reason string) {
	switch devices.State(state) {
	case devices.StateConnected, devices.StateActive:
	default:
		return false, state
	}
	switch {
	case !hasToken:
		return false, reasonNoToken
	case permission == "denied":
		return false, reasonPermissionDenied
	}
	return true, ""
}

// Resolve works out who a message with this audience and category goes to now: the users of the audience who
// did not switch the category off, and, among their devices, those that pass the session policy. A token alone
// is never enough.
func (p *Pipeline) Resolve(ctx context.Context, aud Audience, category models.NotificationCategory) (Resolution, error) {
	return p.resolve(ctx, p.db, aud, category)
}

func (p *Pipeline) resolve(ctx context.Context, q queryer, aud Audience, category models.NotificationCategory) (Resolution, error) {
	args := []any{p.deviceCfg.SessionTimeout.Seconds(), p.deviceCfg.PresenceWindow.Seconds()}
	arg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}

	var where strings.Builder
	if aud.interest() {
		var either []string
		if aud.InterestedEventID > 0 {
			either = append(either, "n.email IN (SELECT email FROM events_attendents WHERE id_events = "+arg(aud.InterestedEventID)+")")
		}
		if aud.InterestedClubID > 0 {
			either = append(either, "n.email IN (SELECT email FROM clubs_members WHERE id_clubs = "+arg(aud.InterestedClubID)+")")
		}
		where.WriteString(" AND (" + strings.Join(either, " OR ") + ")")
	} else {
		switch aud.Type {
		case models.AudienceClub:
			where.WriteString(" AND n.email IN (SELECT email FROM clubs_members WHERE id_clubs = " + arg(aud.ClubID) + ")")
		case models.AudienceCampus:
			where.WriteString(" AND n.campus = " + arg(strings.ToUpper(aud.Campus)))
		case models.AudienceUsers:
			where.WriteString(" AND n.email = ANY(" + arg(pq.Array(aud.Emails)) + ")")
		case models.AudienceCohort:
			if aud.Formation != "" {
				where.WriteString(" AND n.formation_name = " + arg(strings.ToUpper(aud.Formation)))
			}
			if aud.GraduationYear != 0 {
				where.WriteString(" AND n.graduation_year = " + arg(aud.GraduationYear))
			}
		case models.AudienceAll:
		default:
			return Resolution{}, fmt.Errorf("unknown audience type %q", aud.Type)
		}
	}
	if category != "" {
		where.WriteString(` AND n.email IN (
			SELECT nt.email FROM notifications nt
			JOIN services s ON s.id_services = nt.id_services
			WHERE s.name = ` + arg(string(category)) + ")")
	}

	rows, err := q.QueryContext(ctx, `
		SELECT n.email, COALESCE(l.code, 'fr'), d.installation_id, COALESCE(d.push_token IS NOT NULL, FALSE),
		       COALESCE(d.permission_status, ''), `+devices.StateSQL("$1", "$2")+`
		FROM newf n
		LEFT JOIN languages l ON l.id_languages = n.language
		LEFT JOIN devices d ON d.user_email = n.email
		WHERE TRUE`+where.String()+`
		ORDER BY n.email, d.installation_id`, args...)
	if err != nil {
		return Resolution{}, fmt.Errorf("resolve audience: %w", err)
	}
	defer rows.Close()

	res := Resolution{Filtered: map[string]int{}}
	hasTarget := map[string]bool{}
	for rows.Next() {
		var email, lang, permission, state string
		var deviceID sql.NullString
		var hasToken bool
		if err := rows.Scan(&email, &lang, &deviceID, &hasToken, &permission, &state); err != nil {
			return Resolution{}, fmt.Errorf("read audience: %w", err)
		}
		if _, seen := hasTarget[email]; !seen {
			hasTarget[email] = false
			res.Users++
		}
		if !deviceID.Valid {
			res.Filtered[reasonNoDevice]++
			continue
		}
		ok, reason := eligibility(state, hasToken, permission)
		if !ok {
			res.Filtered[reason]++
			continue
		}
		hasTarget[email] = true
		res.Targets = append(res.Targets, Target{DeviceID: deviceID.String, Email: email, Language: lang})
	}
	if err := rows.Err(); err != nil {
		return Resolution{}, fmt.Errorf("read audience: %w", err)
	}
	for email, ok := range hasTarget {
		if !ok {
			res.WithoutTarget = append(res.WithoutTarget, email)
		}
	}
	sort.Strings(res.WithoutTarget)
	return res, nil
}
