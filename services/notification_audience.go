package services

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/plugimt/transat-backend/models"
)

// Recipients are the devices an audience resolves to.
type Recipients struct {
	Users  int
	Tokens []string
}

// ResolveAudience returns the devices of the audience's users. With a category, users who opted
// out of it are excluded.
func (ns *NotificationService) ResolveAudience(audience models.NotificationAudience, category models.NotificationCategory) (Recipients, error) {
	query := `
		SELECT n.email, unt.token
		FROM user_notification_tokens unt
		JOIN newf n ON n.email = unt.email
		WHERE unt.token != ''`
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

	users := make(map[string]struct{})
	var tokens []string
	for rows.Next() {
		var email, token string
		if err := rows.Scan(&email, &token); err != nil {
			return Recipients{}, fmt.Errorf("failed to scan notification recipient: %w", err)
		}
		users[email] = struct{}{}
		tokens = append(tokens, token)
	}
	if err := rows.Err(); err != nil {
		return Recipients{}, fmt.Errorf("failed to read notification recipients: %w", err)
	}

	return Recipients{Users: len(users), Tokens: dedupe(tokens)}, nil
}
