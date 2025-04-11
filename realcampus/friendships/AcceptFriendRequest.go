package friendships

import (
	"Transat_2.0_Backend/utils"
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

// AcceptFriendRequest handles accepting a friend request.
func AcceptFriendRequest(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		utils.LogHeader("Friend Request Acceptance")

		// Get the authenticated user's email.
		userEmail := c.Locals("email").(string)
		utils.LogLineKeyValue(utils.LevelInfo, "Authenticated User", userEmail)

		// Get the friendship ID from request body.
		var request struct {
			FriendshipID int `json:"friendship_id"`
		}

		if err := c.BodyParser(&request); err != nil {
			utils.LogLineKeyValue(utils.LevelError, "Body Parsing Error", err.Error())
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid request format",
			})
		}

		utils.LogLineKeyValue(utils.LevelInfo, "Friendship ID", request.FriendshipID)

		// Ensure the friendship ID is valid
		if request.FriendshipID <= 0 {
			utils.LogMessage(utils.LevelError, "Invalid friendship ID provided.")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid friendship ID.",
			})
		}

		// Check if the friendship exists and user is authorized to accept it.
		var exists bool
		var status string
		err := db.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM realcampus_friendships WHERE id_friendship = $1 AND friend_id = $2),
			       (SELECT status FROM realcampus_friendships WHERE id_friendship = $1)
		`, request.FriendshipID, userEmail).Scan(&exists, &status)

		if err != nil {
			utils.LogLineKeyValue(utils.LevelError, "Database Error", err.Error())
			utils.LogMessage(utils.LevelError, "Error checking friendship.")
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Database error occurred.",
			})
		}

		if !exists {
			utils.LogMessage(utils.LevelError, "Friend request not found or user not authorized.")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Friend request not found or you're not authorized to accept it",
			})
		}

		utils.LogLineKeyValue(utils.LevelInfo, "Friendship Status", status)

		if status != "PENDING" {
			utils.LogMessage(utils.LevelError, "Friendship is not in pending state.")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "This friendship is not in pending state.",
			})
		}

		// Accept the friend request.
		_, err = db.Exec(`
			UPDATE realcampus_friendships
			SET status = 'ACCEPTED', acceptance_date = CURRENT_TIMESTAMP
			WHERE id_friendship = $1
		`, request.FriendshipID)

		if err != nil {
			utils.LogLineKeyValue(utils.LevelError, "Database Error", err.Error())
			utils.LogMessage(utils.LevelError, "Failed to accept friend request.")
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Database error occurred.",
			})
		}

		utils.LogMessage(utils.LevelInfo, "Friend request successfully accepted.")
		utils.LogFooter()
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Friend request accepted.",
		})
	}
}
