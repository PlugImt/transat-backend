package services

import (
	"fmt"
	"sort"
	"strings"

	"github.com/lib/pq"
	"github.com/plugimt/transat-backend/models"
)

// Recipients are the users an audience resolves to and their registered devices.
type Recipients struct {
	Users  int
	Tokens []string
	// WithoutDevice lists users who never registered a push token and cannot be reached.
	WithoutDevice []string
}

// ResolveAudience returns the audience's users and their devices. With a category, users who
// opted out of it are excluded.
func (ns *NotificationService) ResolveAudience(audience models.NotificationAudience, category models.NotificationCategory) (Recipients, error) {
	query := `
		SELECT n.email, COALESCE(unt.token, '')
		FROM newf n
		LEFT JOIN user_notification_tokens unt ON unt.email = n.email
		WHERE TRUE`
	var args []interface{}
	arg := func(v interface{}) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}

	switch audience.Type {
	case models.AudienceClub:
		query += " AND n.email IN (SELECT email FROM clubs_members WHERE id_clubs = " + arg(audience.ClubID) + ")"
	case models.AudienceCampus:
		query += " AND n.campus = " + arg(strings.ToUpper(audience.Campus))
	case models.AudienceUsers:
		query += " AND n.email = ANY(" + arg(pq.Array(audience.Emails)) + ")"
	case models.AudienceCohort:
		if audience.Formation != "" {
			query += " AND n.formation_name = " + arg(strings.ToUpper(audience.Formation))
		}
		if audience.GraduationYear != 0 {
			query += " AND n.graduation_year = " + arg(audience.GraduationYear)
		}
	}
	if category != "" {
		query += ` AND n.email IN (
			SELECT nt.email FROM notifications nt
			JOIN services s ON s.id_services = nt.id_services
			WHERE s.name = ` + arg(string(category)) + ")"
	}

	rows, err := ns.db.Query(query, args...)
	if err != nil {
		return Recipients{}, fmt.Errorf("failed to resolve notification audience: %w", err)
	}
	defer rows.Close()

	users := make(map[string]bool) // email -> has a device
	var tokens []string
	for rows.Next() {
		var email, token string
		if err := rows.Scan(&email, &token); err != nil {
			return Recipients{}, fmt.Errorf("failed to scan notification recipient: %w", err)
		}
		users[email] = users[email] || token != ""
		if token != "" {
			tokens = append(tokens, token)
		}
	}
	if err := rows.Err(); err != nil {
		return Recipients{}, fmt.Errorf("failed to read notification recipients: %w", err)
	}

	withoutDevice := []string{}
	for email, hasDevice := range users {
		if !hasDevice {
			withoutDevice = append(withoutDevice, email)
		}
	}
	sort.Strings(withoutDevice)

	return Recipients{Users: len(users), Tokens: dedupe(tokens), WithoutDevice: withoutDevice}, nil
}
