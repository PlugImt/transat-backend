package friendships

import (
	"Transat_2.0_Backend/models"
	"Transat_2.0_Backend/utils"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// GetUserFriendships retrieves all friendships for a user, separated by type.
func GetUserFriendships(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Locals("email").(string)

		// Prepare our response structure
		response := fiber.Map{
			"success":           true,
			"received_requests": []fiber.Map{},
			"sent_requests":     []fiber.Map{},
			"friends":           []fiber.Map{},
		}

		// Query for incoming friend requests (where user is the friend_id and status is PENDING)
		incomingRows, err := db.Query(`
			SELECT 
				f.id_friendship,
				f.user_id,
				f.friend_id,
				f.status,
				f.request_date,
				f.acceptance_date,
				n.first_name,
				n.last_name,
				n.profile_picture
			FROM realcampus_friendships f
			JOIN newf n ON f.user_id = n.email
			WHERE f.friend_id = $1 AND f.status = 'PENDING'
			ORDER BY f.request_date DESC`,
			email)

		if err != nil {
			fmt.Println("Database error (incoming requests):", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve incoming friend requests",
			})
		}
		defer incomingRows.Close()

		// Process incoming friend requests
		var receivedRequests []fiber.Map
		for incomingRows.Next() {
			var friendship models.RealCampusFriendship
			var firstName, lastName, profilePicture sql.NullString

			if err := incomingRows.Scan(
				&friendship.ID,
				&friendship.UserID,
				&friendship.FriendID,
				&friendship.Status,
				&friendship.RequestDate,
				&friendship.AcceptanceDate,
				&firstName,
				&lastName,
				&profilePicture,
			); err != nil {
				fmt.Println("Error scanning incoming row:", err)
				continue
			}

			requestData := fiber.Map{
				"id":              friendship.ID,
				"user_id":         friendship.UserID,
				"friend_id":       friendship.FriendID,
				"status":          friendship.Status,
				"request_date":    friendship.RequestDate,
				"first_name":      utils.NullStringValue(firstName),
				"last_name":       utils.NullStringValue(lastName),
				"profile_picture": utils.NullStringValue(profilePicture),
			}

			receivedRequests = append(receivedRequests, requestData)
		}
		response["received_requests"] = receivedRequests
		response["received_count"] = len(receivedRequests)

		// Query for outgoing friend requests (where user is the user_id and status is PENDING)
		outgoingRows, err := db.Query(`
			SELECT 
				f.id_friendship,
				f.user_id,
				f.friend_id,
				f.status,
				f.request_date,
				f.acceptance_date,
				n.first_name,
				n.last_name,
				n.profile_picture
			FROM realcampus_friendships f
			JOIN newf n ON f.friend_id = n.email
			WHERE f.user_id = $1 AND f.status = 'PENDING'
			ORDER BY f.request_date DESC`,
			email)

		if err != nil {
			fmt.Println("Database error (outgoing requests):", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve outgoing friend requests",
			})
		}
		defer outgoingRows.Close()

		// Process outgoing friend requests
		var sentRequests []fiber.Map
		for outgoingRows.Next() {
			var friendship models.RealCampusFriendship
			var firstName, lastName, profilePicture sql.NullString

			if err := outgoingRows.Scan(
				&friendship.ID,
				&friendship.UserID,
				&friendship.FriendID,
				&friendship.Status,
				&friendship.RequestDate,
				&friendship.AcceptanceDate,
				&firstName,
				&lastName,
				&profilePicture,
			); err != nil {
				fmt.Println("Error scanning outgoing row:", err)
				continue
			}

			requestData := fiber.Map{
				"id":              friendship.ID,
				"user_id":         friendship.UserID,
				"friend_id":       friendship.FriendID,
				"status":          friendship.Status,
				"request_date":    friendship.RequestDate,
				"first_name":      utils.NullStringValue(firstName),
				"last_name":       utils.NullStringValue(lastName),
				"profile_picture": utils.NullStringValue(profilePicture),
			}

			sentRequests = append(sentRequests, requestData)
		}
		response["sent_requests"] = sentRequests
		response["sent_count"] = len(sentRequests)

		// Query for accepted friendships (where user is either user_id or friend_id and status is ACCEPTED)
		friendsRows, err := db.Query(`
			SELECT 
				f.id_friendship,
				f.user_id,
				f.friend_id,
				f.status,
				f.request_date,
				f.acceptance_date,
				n.first_name,
				n.last_name,
				n.profile_picture,
				CASE WHEN f.user_id = $1 THEN 'outgoing' ELSE 'incoming' END as direction
			FROM realcampus_friendships f
			JOIN newf n ON (
				CASE 
					WHEN f.user_id = $1 THEN f.friend_id = n.email
					ELSE f.user_id = n.email
				END
			)
			WHERE (f.user_id = $1 OR f.friend_id = $1) AND f.status = 'ACCEPTED'
			ORDER BY f.acceptance_date DESC`,
			email)

		if err != nil {
			fmt.Println("Database error (friends):", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve friends",
			})
		}
		defer friendsRows.Close()

		// Process accepted friendships.
		var friends []fiber.Map
		for friendsRows.Next() {
			var friendship models.RealCampusFriendship
			var firstName, lastName, profilePicture sql.NullString
			var direction string

			if err := friendsRows.Scan(
				&friendship.ID,
				&friendship.UserID,
				&friendship.FriendID,
				&friendship.Status,
				&friendship.RequestDate,
				&friendship.AcceptanceDate,
				&firstName,
				&lastName,
				&profilePicture,
				&direction,
			); err != nil {
				fmt.Println("Error scanning friends row:", err)
				continue
			}

			friendData := fiber.Map{
				"id":              friendship.ID,
				"user_id":         friendship.UserID,
				"friend_id":       friendship.FriendID,
				"status":          friendship.Status,
				"request_date":    friendship.RequestDate,
				"acceptance_date": friendship.AcceptanceDate.Time,
				"direction":       direction,
				"first_name":      utils.NullStringValue(firstName),
				"last_name":       utils.NullStringValue(lastName),
				"profile_picture": utils.NullStringValue(profilePicture),
			}

			friends = append(friends, friendData)
		}
		response["friends"] = friends
		response["friends_count"] = len(friends)

		// Set total count for all relationships.
		response["total_count"] = len(receivedRequests) + len(sentRequests) + len(friends)

		return c.JSON(response)
	}
}
