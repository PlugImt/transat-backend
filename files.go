package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ensureDataFolder checks if the /data folder exists, creates it if not
func ensureDataFolder() error {
	path := "/data"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("/data folder does not exist. Creating...")
		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Println("Failed to create /data folder:", err)
			return err
		}
		fmt.Println("/data folder created successfully.")
	} else {
		fmt.Println("/data folder exists.")
	}
	return nil
}

// uploadImage handles file uploads and saves them to the /data directory
func uploadImage(c *fiber.Ctx) error {
	fmt.Println("Received file upload request")

	fmt.Println("Received file upload request")
	fmt.Println("Content-Type:", c.Get("Content-Type"))

	// Get user email from JWT
	email := c.Locals("email").(string)
	if email == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User email not found in token",
		})
	}

	// Try direct file parsing first
	file, err := c.FormFile("image")
	if err != nil {
		// If standard parsing fails, try to parse multipart form manually
		if form, err := c.MultipartForm(); err == nil {
			if files := form.File["image"]; len(files) > 0 {
				file = files[0]
			} else {
				fmt.Println("No files found in multipart form")
			}
		} else {
			fmt.Println("Error parsing multipart form:", err)
		}
	}

	// If we still don't have a file, return error
	if file == nil {
		fmt.Println("No image uploaded or invalid file")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No image uploaded or invalid file",
		})
	}
	fmt.Println("Uploaded file name:", file.Filename)

	// Validate file is an image
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		fmt.Println("Uploaded file is not an image")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Uploaded file is not an image",
		})
	}

	// Ensure the /data directory exists
	if err := ensureDataFolder(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to ensure storage directory",
		})
	}

	// Get original filename and extension
	originalFilename := filepath.Base(file.Filename)
	fileExt := filepath.Ext(originalFilename)
	fileNameWithoutExt := strings.TrimSuffix(originalFilename, fileExt)

	// Check if a file with this name already exists for this user
	var existingId int
	var existingPath string
	err = db.QueryRow("SELECT id_files, path FROM files WHERE name = $1 AND email = $2", originalFilename, email).Scan(&existingId, &existingPath)

	// Generate a unique filename if this filename already exists
	filename := originalFilename
	if err == nil {
		// File exists, create a unique name by adding timestamp hash
		timestamp := time.Now().UnixNano()
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s%d", originalFilename, timestamp)))
		hashStr := hex.EncodeToString(hash[:8]) // Use first 8 bytes of hash
		filename = fmt.Sprintf("%s_%s%s", fileNameWithoutExt, hashStr, fileExt)
		fmt.Printf("File '%s' already exists, using unique name: '%s'\n", originalFilename, filename)
	}

	// Set destination path
	dst := fmt.Sprintf("/data/%s", filename)
	fmt.Println("Destination path:", dst)

	// Save file to /data directory
	fileReader, err := file.Open()
	if err != nil {
		fmt.Println("Failed to read uploaded file:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read uploaded file",
		})
	}
	defer fileReader.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Println("Failed to create destination file:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create destination file",
		})
	}
	defer dstFile.Close()

	// Copy file contents
	if _, err := io.Copy(dstFile, fileReader); err != nil {
		fmt.Println("Failed to save file:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}
	fmt.Println("File uploaded successfully:", dst)

	// Insert file record into the database
	var fileID int
	err = db.QueryRow(`
        INSERT INTO files (name, path, email, creation_date, last_access_date) 
        VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
        RETURNING id_files`,
		originalFilename, dst, email,
	).Scan(&fileID)

	if err != nil {
		fmt.Println("Failed to record file in database:", err)
		// Delete the file if we couldn't record it in the database
		os.Remove(dst)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to record file in database",
		})
	}

	// Return success response with image URL and database ID
	return c.JSON(fiber.Map{
		"success": true,
		"url":     "/data/" + filename,
		"file_id": fileID,
	})
}

// serveImage serves images from the /data directory
func serveImage(c *fiber.Ctx) error {
	filename := c.Params("filename")
	filepath := fmt.Sprintf("/data/%s", filename)

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Println("Image not found:", filepath)
		return c.Status(fiber.StatusNotFound).SendString("Image not found")
	}

	// Update the last_access_date in the database
	// We're doing this concurrently so it doesn't slow down the file serving
	go func() {
		_, err := db.Exec(`
            UPDATE files 
            SET last_access_date = CURRENT_TIMESTAMP 
            WHERE path = $1`,
			filepath)
		if err != nil {
			fmt.Printf("Failed to update last access date for %s: %v\n", filepath, err)
		}
	}()

	fmt.Println("Serving image:", filepath)
	// Send the file
	return c.SendFile(filepath)
}

func listUserFiles(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	rows, err := db.Query(`
        SELECT id_files, name, path, creation_date, last_access_date 
        FROM files 
        WHERE email = $1 
        ORDER BY creation_date DESC`,
		email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve files",
		})
	}
	defer rows.Close()

	var files []fiber.Map
	for rows.Next() {
		var id int
		var name, path string
		var creationDate, lastAccessDate time.Time

		if err := rows.Scan(&id, &name, &path, &creationDate, &lastAccessDate); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		files = append(files, fiber.Map{
			"id":            id,
			"name":          name,
			"url":           strings.Replace(path, "/data/", "/api/data/", 1),
			"created":       creationDate,
			"last_accessed": lastAccessDate,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"files":   files,
	})
}

func deleteFile(c *fiber.Ctx) error {
	fileID := c.Params("id")
	email := c.Locals("email").(string)

	// First get the filepath to delete the actual file
	var filePath string
	err := db.QueryRow("SELECT path FROM files WHERE id_files = $1 AND email = $2", fileID, email).Scan(&filePath)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "File not found or you don't have permission to delete it",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Delete from database
	_, err = db.Exec("DELETE FROM files WHERE id_files = $1 AND email = $2", fileID, email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete file record",
		})
	}

	// Delete the actual file
	if err := os.Remove(filePath); err != nil {
		fmt.Printf("Warning: Could not delete file at %s: %v\n", filePath, err)
		// We don't return an error since the DB record is already deleted
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "File deleted successfully",
	})
}
