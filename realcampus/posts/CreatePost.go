package posts

import (
	"Transat_2.0_Backend/utils"
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

// CreatePost creates a new post with two images.
func CreatePost(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		utils.LogHeader("CreatePost")

		email := c.Locals("email").(string)
		utils.LogLineKeyValue(utils.LevelInfo, "User", email)

		// Parse request body.
		var request struct {
			FileID1  int    `json:"file_id_1" validate:"required"`
			FileID2  int    `json:"file_id_2" validate:"required"`
			Location string `json:"location"`
			Privacy  string `json:"privacy"` // PUBLIC or PRIVATE.
		}

		if err := c.BodyParser(&request); err != nil {
			utils.LogMessage(utils.LevelError, "Failed to parse request body.")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		utils.LogLineKeyValue(utils.LevelDebug, "Parsed Request", request)

		// Validate privacy level.
		if request.Privacy != "PUBLIC" && request.Privacy != "PRIVATE" {
			utils.LogMessage(utils.LevelWarn, "Invalid or missing privacy value; defaulting to PRIVATE.")
			request.Privacy = "PRIVATE"
		}

		// Check file ownership.
		var ownershipCount int
		err := db.QueryRow(`
			SELECT COUNT(*) 
			FROM files 
			WHERE id_files IN ($1, $2) AND email = $3`,
			request.FileID1, request.FileID2, email).Scan(&ownershipCount)

		if err != nil || ownershipCount != 2 {
			utils.LogMessage(utils.LevelError, "Files ownership verification failed or not all files are owned by user.")
			utils.LogFooter()
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You can only create posts with your own files",
			})
		}

		// Check if any file is already used in a post.
		var usageCount int
		err = db.QueryRow(`
			SELECT COUNT(*) 
			FROM realcampus_posts 
			WHERE id_file_1 IN ($1, $2) OR id_file_2 IN ($1, $2)`,
			request.FileID1, request.FileID2).Scan(&usageCount)

		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to check file usage in posts.")
			utils.LogLineKeyValue(utils.LevelError, "Error", err.Error())
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to verify file usage.",
			})
		}

		if usageCount > 0 {
			utils.LogMessage(utils.LevelError, "One or both files are already used in another post.")
			utils.LogFooter()
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "One or both files are already used in another post.",
			})
		}

		// Insert new post.
		var postID int
		err = db.QueryRow(`
			INSERT INTO realcampus_posts (id_file_1, id_file_2, author_email, location, privacy)
			VALUES ($1, $2, $3, $4, $5::realcampus_privacy_level)
			RETURNING id_post`,
			request.FileID1, request.FileID2, email, request.Location, request.Privacy).Scan(&postID)

		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to insert post into database.")
			utils.LogLineKeyValue(utils.LevelError, "Error", err.Error())
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create post.",
			})
		}

		utils.LogLineKeyValue(utils.LevelInfo, "Post Created", postID)
		utils.LogFooter()

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"post_id": postID,
			"message": "Post created successfully",
		})
	}
}
