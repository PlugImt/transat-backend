package posts

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

// GetUserTodayPosts retrieves all posts created by the user today.
func GetUserTodayPosts(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		utils.LogHeader("GetUserTodayPosts")

		emailIfc := c.Locals("email")
		if emailIfc == nil {
			utils.LogMessage(utils.LevelError, "Missing email in context.")
			utils.LogFooter()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: email not found in context",
			})
		}

		email, ok := emailIfc.(string)
		if !ok || email == "" {
			utils.LogMessage(utils.LevelError, "Invalid email format in context.")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid email in context",
			})
		}

		utils.LogLineKeyValue(utils.LevelInfo, "User", email)

		now := time.Now().UTC()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		endOfDay := startOfDay.Add(24 * time.Hour)

		// Query for posts created today by the user.
		rows, err := db.Query(`
			SELECT 
				p.id_post, 
				p.id_file_1, 
				p.id_file_2, 
				p.location, 
				p.privacy, 
				p.creation_date,
				f1.path AS path1,
				f1.name AS name1,
				f2.path AS path2,
				f2.name AS name2
			FROM realcampus_posts p
			JOIN files f1 ON p.id_file_1 = f1.id_files
			JOIN files f2 ON p.id_file_2 = f2.id_files
			WHERE p.author_email = $1 
			AND p.creation_date >= $2 
			AND p.creation_date < $3
			ORDER BY p.creation_date DESC`,
			email, startOfDay, endOfDay)

		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to query today's posts.")
			utils.LogLineKeyValue(utils.LevelError, "Error", err.Error())
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve today's posts.",
			})
		}
		defer rows.Close()

		var posts []fiber.Map
		for rows.Next() {
			var post models.RealCampusPost
			var path1, path2, name1, name2 string

			if err := rows.Scan(
				&post.ID,
				&post.FileID1,
				&post.FileID2,
				&post.Location,
				&post.Privacy,
				&post.CreationDate,
				&path1,
				&name1,
				&path2,
				&name2,
			); err != nil {
				utils.LogMessage(utils.LevelWarn, "Error scanning row, skipping.")
				utils.LogLineKeyValue(utils.LevelError, "Error", err.Error())
				continue
			}

			posts = append(posts, fiber.Map{
				"id":       post.ID,
				"location": post.Location,
				"privacy":  post.Privacy,
				"created":  post.CreationDate,
				"file1": fiber.Map{
					"id":   post.FileID1,
					"name": name1,
					"url":  strings.Replace(path1, "/data/", "/api/data/", 1),
				},
				"file2": fiber.Map{
					"id":   post.FileID2,
					"name": name2,
					"url":  strings.Replace(path2, "/data/", "/api/data/", 1),
				},
			})
		}

		utils.LogLineKeyValue(utils.LevelInfo, "Posts Found, amount", len(posts))
		utils.LogFooter()

		return c.JSON(fiber.Map{
			"success": true,
			"count":   len(posts),
			"posts":   posts,
		})
	}
}
