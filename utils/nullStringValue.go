package utils

import "database/sql"

// Helper function to handle null strings.
func NullStringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
