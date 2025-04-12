package friendships

import (
	"Transat_2.0_Backend/utils"
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

// RemoveFriend handles removing an existing friend.
func RemoveFriend(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		utils.LogHeader("Friend Removal")

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
				"error":   "Invalid request format.",
			})
		}

		utils.LogLineKeyValue(utils.LevelInfo, "Friendship ID", request.FriendshipID)

		// Ensure the friendship ID is valid.
		if request.FriendshipID <= 0 {
			utils.LogMessage(utils.LevelWarn, "Invalid friendship ID provided.")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid friendship ID.",
			})
		}

		// Check if the friendship exists and user is authorized to remove it.
		var exists bool
		var status string
		err := db.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM realcampus_friendships WHERE id_friendship = $1 AND (user_id = $2 OR friend_id = $2)),
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
			utils.LogMessage(utils.LevelWarn, "Friendship not found or user not authorized.")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Friendship not found or you're not authorized to remove it.",
			})
		}

		utils.LogLineKeyValue(utils.LevelInfo, "Friendship Status.", status)

		if status != "ACCEPTED" {
			utils.LogMessage(utils.LevelWarn, "Friendship is not in accepted state.")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "This is not an accepted friendship.",
			})
		}

		// Remove the friendship (delete it).
		utils.LogMessage(utils.LevelInfo, "Removing friendship")
		_, err = db.Exec(`
			DELETE FROM realcampus_friendships
			WHERE id_friendship = $1
		`, request.FriendshipID)

		if err != nil {
			utils.LogLineKeyValue(utils.LevelError, "Database Error", err.Error())
			utils.LogMessage(utils.LevelError, "Failed to remove friend.")
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Database error occurred.",
			})
		}

		utils.LogMessage(utils.LevelInfo, "Friend removed successfully.")
		utils.LogFooter()
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Friend removed successfully.",
		})
	}
}
