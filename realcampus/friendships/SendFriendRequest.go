package friendships

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

// SendFriendRequest handles creating a new friend request.
func SendFriendRequest(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		utils.LogHeader("Sending Friend Request")

		// Get the authenticated user's email.
		userEmail := c.Locals("email").(string)
		utils.LogLineKeyValue(utils.LevelInfo, "Authenticated User", userEmail)

		// Get the target friend's email from request body.
		var request struct {
			FriendEmail string `json:"friend_email"`
		}

		if err := c.BodyParser(&request); err != nil {
			utils.LogLineKeyValue(utils.LevelError, "Body Parsing Error", err.Error())
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid request format",
			})
		}

		// Ensure the friend email is provided.
		if request.FriendEmail == "" {
			utils.LogMessage(utils.LevelError, "Friend email is missing in request")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Friend email is required",
			})
		}

		utils.LogLineKeyValue(utils.LevelInfo, "Target Friend", request.FriendEmail)

		// Ensure user is not trying to send a request to themselves.
		if userEmail == request.FriendEmail {
			utils.LogMessage(utils.LevelError, "User attempted to send friend request to self.")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "You cannot send a friend request to yourself",
			})
		}

		// Check if the friend exists.
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM newf WHERE email = $1)", request.FriendEmail).Scan(&exists)
		if err != nil {
			utils.LogLineKeyValue(utils.LevelError, "Database Error", err.Error())
			utils.LogMessage(utils.LevelError, "Error checking if friend exists")
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Database error occurred : Error checking if friend exists.",
			})
		}

		if !exists {
			utils.LogMessage(utils.LevelError, "Target friend not found in database.")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Target friend not found.",
			})
		}

		utils.LogMessage(utils.LevelInfo, "Target friend exists, checking existing friendship")

		// Check if there's already a friendship between these users.
		var existingFriendship models.RealCampusFriendship
		err = db.QueryRow(`
			SELECT id_friendship, status
			FROM realcampus_friendships
			WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
		`, userEmail, request.FriendEmail).Scan(&existingFriendship.ID, &existingFriendship.Status)

		if err != nil && err != sql.ErrNoRows {
			utils.LogLineKeyValue(utils.LevelError, "Database Error", err.Error())
			utils.LogMessage(utils.LevelError, "Error checking existing friendship.")
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Database error occurred : Error checking existing friendship.",
			})
		}

		// If a friendship exists, handle accordingly.
		if err != sql.ErrNoRows {
			utils.LogLineKeyValue(utils.LevelInfo, "Existing Friendship Status", existingFriendship.Status)

			if existingFriendship.Status == "ACCEPTED" {
				utils.LogMessage(utils.LevelError, "Users are already friends.")
				utils.LogFooter()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"error":   "You are already friends with this user.",
				})
			} else if existingFriendship.Status == "PENDING" {
				utils.LogMessage(utils.LevelError, "Friend request already exists.")
				utils.LogFooter()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"error":   "A friend request already exists between you and this user.",
				})
			} else if existingFriendship.Status == "REJECTED" {
				utils.LogMessage(utils.LevelInfo, "Previously rejected request, updating to PENDING.")
				// If previously rejected, we'll allow sending a new request by updating the existing one.
				_, err = db.Exec(`
					UPDATE realcampus_friendships
					SET status = 'PENDING', request_date = CURRENT_TIMESTAMP, acceptance_date = NULL
					WHERE id_friendship = $1
				`, existingFriendship.ID)

				if err != nil {
					utils.LogLineKeyValue(utils.LevelError, "Database Error", err.Error())
					utils.LogMessage(utils.LevelError, "Failed to update friend request.")
					utils.LogFooter()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"success": false,
						"error":   "Database error occurred : Failed to update friend request.",
					})
				}

				utils.LogMessage(utils.LevelInfo, "Friend request successfully sent.")
				utils.LogFooter()

				return c.JSON(fiber.Map{
					"success": true,
					"message": "Friend request sent.",
				})
			}
		}

		utils.LogMessage(utils.LevelInfo, "Creating new friend request.")

		// Create a new friend request.
		_, err = db.Exec(`
			INSERT INTO realcampus_friendships (user_id, friend_id, status)
			VALUES ($1, $2, 'PENDING')
		`, userEmail, request.FriendEmail)

		if err != nil {
			utils.LogLineKeyValue(utils.LevelError, "Database Error", err.Error())
			utils.LogMessage(utils.LevelError, "Failed to send friend request.")
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Database error occurred : Failed to send friend request.",
			})
		}

		utils.LogMessage(utils.LevelInfo, "Friend request successfully sent!")
		utils.LogFooter()
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Friend request sent!",
		})
	}
}
