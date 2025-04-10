package users

import (
	"Transat_2.0_Backend/utils"
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

// SearchUsers handles searching for users to add as friends.
func SearchUsers(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the authenticated user's email.
		userEmail := c.Locals("email").(string)

		utils.LogHeader("Search Users")
		utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)

		// Get the search query from request.
		query := c.Query("q")
		utils.LogLineKeyValue(utils.LevelInfo, "Search Query", query)

		var querySQL string
		var args []interface{}

		// Base arguments always include the current user's email.
		args = append(args, userEmail)

		// If no query is provided, return all users.
		if query == "" {
			querySQL = `
                SELECT 
                    email, 
                    first_name, 
                    last_name, 
                    profile_picture
                FROM newf
                WHERE email != $1
                ORDER BY first_name, last_name
                LIMIT 50
            `
		} else {
			// Search for users matching the query that are not the current user.
			querySQL = `
                SELECT 
                    email, 
                    first_name, 
                    last_name, 
                    profile_picture
                FROM newf
                WHERE 
                    email != $1 AND
                    (
                        email ILIKE $2 OR
                        first_name ILIKE $2 OR
                        last_name ILIKE $2
                    )
                ORDER BY 
                    CASE 
                        WHEN first_name ILIKE $2 OR last_name ILIKE $2 THEN 1
                        ELSE 2
                    END,
                    first_name, 
                    last_name
                LIMIT 50
            `
			args = append(args, "%"+query+"%")
		}

		// Execute the search query.
		rows, err := db.Query(querySQL, args...)

		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to search users")
			utils.LogLineKeyValue(utils.LevelError, "SQL Error", err.Error())
			utils.LogFooter()

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to search users",
			})
		}
		defer rows.Close()

		// Process search results.
		var users []fiber.Map
		for rows.Next() {
			var email, firstName, lastName string
			var profilePicture sql.NullString

			if err := rows.Scan(&email, &firstName, &lastName, &profilePicture); err != nil {
				utils.LogLineKeyValue(utils.LevelError, "Scan Error", err.Error())
				continue
			}

			// Check friendship status with this user.
			var friendshipID sql.NullInt64
			var status sql.NullString
			var direction sql.NullString

			err := db.QueryRow(`
				SELECT 
					id_friendship,
					status,
					CASE 
						WHEN user_id = $1 THEN 'outgoing'
						WHEN friend_id = $1 THEN 'incoming'
					END as direction
				FROM realcampus_friendships
				WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
			`, userEmail, email).Scan(&friendshipID, &status, &direction)

			userData := fiber.Map{
				"email":           email,
				"first_name":      firstName,
				"last_name":       lastName,
				"profile_picture": utils.NullStringValue(profilePicture),
				"friendship":      nil,
			}

			// If there's a friendship, include its details.
			if err == nil && friendshipID.Valid {
				userData["friendship"] = fiber.Map{
					"id":        friendshipID.Int64,
					"status":    status.String,
					"direction": direction.String,
				}
			}

			users = append(users, userData)
		}

		utils.LogLineKeyValue(utils.LevelInfo, "Users Found", len(users))
		utils.LogFooter()

		return c.JSON(fiber.Map{
			"success": true,
			"count":   len(users),
			"users":   users,
		})
	}
}
