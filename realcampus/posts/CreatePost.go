package posts

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// CreatePost creates a new post with two images.
func CreatePost(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Locals("email").(string)

		// Parse request body.
		var request struct {
			FileID1  int    `json:"file_id_1" validate:"required"`
			FileID2  int    `json:"file_id_2" validate:"required"`
			Location string `json:"location"`
			Privacy  string `json:"privacy"` // PUBLIC or PRIVATE.
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		// Validate privacy level according to enum.
		if request.Privacy != "PUBLIC" && request.Privacy != "PRIVATE" {
			request.Privacy = "PRIVATE" // Default to PRIVATE.
		}

		// Verify the files belong to the user.
		var count int
		err := db.QueryRow(`
			SELECT COUNT(*) 
			FROM files 
			WHERE id_files IN ($1, $2) AND email = $3`,
			request.FileID1, request.FileID2, email).Scan(&count)

		if err != nil || count != 2 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You can only create posts with your own files",
			})
		}

		// Create the post.
		var postID int
		err = db.QueryRow(`
			INSERT INTO realcampus_posts (id_file_1, id_file_2, author_email, location, privacy)
			VALUES ($1, $2, $3, $4, $5::realcampus_privacy_level)
			RETURNING id_post`,
			request.FileID1, request.FileID2, email, request.Location, request.Privacy).Scan(&postID)

		if err != nil {
			fmt.Println("Error creating post:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create post",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"post_id": postID,
			"message": "Post created successfully",
		})
	}
}
