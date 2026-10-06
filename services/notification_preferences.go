package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/plugimt/transat-backend/models"
)

// ErrUnknownNotificationCategory is returned when a category has no row in the `services` table.
var ErrUnknownNotificationCategory = errors.New("unknown notification category")

func categoryNames() []string {
	names := make([]string, len(models.NotificationCategories))
	for i, c := range models.NotificationCategories {
		names[i] = string(c)
	}
	return names
}

// GetPreferences returns every notification category with whether the user receives it.
func (ns *NotificationService) GetPreferences(email string) ([]models.NotificationPreference, error) {
	rows, err := ns.db.Query(`
		SELECT s.name, n.email IS NOT NULL
		FROM services s
		LEFT JOIN notifications n ON n.id_services = s.id_services AND n.email = $1
		WHERE s.name = ANY($2)
	`, email, pq.Array(categoryNames()))
	if err != nil {
		return nil, fmt.Errorf("failed to query notification preferences: %w", err)
	}
	defer rows.Close()

	enabled := make(map[models.NotificationCategory]bool, len(models.NotificationCategories))
	for rows.Next() {
		var name string
		var on bool
		if err := rows.Scan(&name, &on); err != nil {
			return nil, fmt.Errorf("failed to scan notification preference: %w", err)
		}
		enabled[models.NotificationCategory(name)] = on
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read notification preferences: %w", err)
	}

	prefs := make([]models.NotificationPreference, 0, len(enabled))
	for _, category := range models.NotificationCategories {
		if on, seeded := enabled[category]; seeded {
			prefs = append(prefs, models.NotificationPreference{Service: category, Enabled: on})
		}
	}
	return prefs, nil
}

// SetPreference switches a category on or off for a user. It is idempotent.
func (ns *NotificationService) SetPreference(email string, category models.NotificationCategory, enabled bool) error {
	var serviceID int
	err := ns.db.QueryRow(`SELECT id_services FROM services WHERE name = $1`, string(category)).Scan(&serviceID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUnknownNotificationCategory
	}
	if err != nil {
		return fmt.Errorf("failed to look up notification category: %w", err)
	}

	if enabled {
		_, err = ns.db.Exec(
			`INSERT INTO notifications (email, id_services) VALUES ($1, $2) ON CONFLICT (email, id_services) DO NOTHING`,
			email, serviceID,
		)
	} else {
		_, err = ns.db.Exec(`DELETE FROM notifications WHERE email = $1 AND id_services = $2`, email, serviceID)
	}
	if err != nil {
		return fmt.Errorf("failed to save notification preference: %w", err)
	}
	return nil
}

// IsSubscribed reports whether the user receives the given category.
func (ns *NotificationService) IsSubscribed(email string, category models.NotificationCategory) (bool, error) {
	var subscribed bool
	err := ns.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM notifications n
			JOIN services s ON s.id_services = n.id_services
			WHERE n.email = $1 AND s.name = $2
		)
	`, email, string(category)).Scan(&subscribed)
	if err != nil {
		return false, fmt.Errorf("failed to check notification subscription: %w", err)
	}
	return subscribed, nil
}
