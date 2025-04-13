package files

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log" // Using standard log
	"os"
	"path/filepath"
	"strings"
	"time"

	"Transat_2.0_Backend/utils" // For logger if needed, EnsureDataFolder moved here potentially?
	"github.com/gofiber/fiber/v2"
)

// FileHandler handles file upload, serving, listing, and deletion.
type FileHandler struct {
	DB         *sql.DB
	DataFolder string // Make data folder path configurable
}

// NewFileHandler creates a new FileHandler.
// It ensures the data folder exists upon creation.
func NewFileHandler(db *sql.DB) (*FileHandler, error) {
	dataFolder := os.Getenv("DATA_FOLDER")
	if dataFolder == "" {
		dataFolder = "/data" // Default path
		log.Println("DATA_FOLDER environment variable not set, using default '/data'")
	}

	// Ensure data folder exists during initialization
	// Using utils.EnsureDataFolder assumes it's defined there
	if err := utils.EnsureDataFolder(); err != nil {
		return nil, fmt.Errorf("failed to ensure data folder '%s': %w", dataFolder, err)
	}

	return &FileHandler{
		DB:         db,
		DataFolder: dataFolder,
	}, nil
}

// UploadFile handles file uploads, saves them, and records them in the database.
func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	utils.LogHeader("📄 Upload File") // Using custom logger

	// Get user email from JWT (set by middleware)
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		utils.LogMessage(utils.LevelWarn, "User email not found in token during upload")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User email not found in token",
		})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "User", email)
	utils.LogLineKeyValue(utils.LevelInfo, "Content-Type", c.Get("Content-Type"))

	// --- File Handling ---
	// Prioritize FormFile for standard uploads
	fileHeader, err := c.FormFile("file") // Use "file" as the default key, adjust if needed
	if err != nil {
		// Fallback for different key or manual parsing if needed
		if fileHeader, err = c.FormFile("image"); err != nil { // Check common alternative "image"
			utils.LogMessage(utils.LevelError, "No file uploaded with key 'file' or 'image'")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			// Try manual multipart form parsing as a last resort? (Can be complex)
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No file uploaded (use key 'file' or 'image')",
			})
		}
		utils.LogMessage(utils.LevelInfo, "File uploaded using key 'image'")
	} else {
		utils.LogMessage(utils.LevelInfo, "File uploaded using key 'file'")
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Original Filename", fileHeader.Filename)
	utils.LogLineKeyValue(utils.LevelInfo, "File Size", fileHeader.Size)
	utils.LogLineKeyValue(utils.LevelInfo, "MIME Header", fileHeader.Header.Get("Content-Type"))

	// --- File Validation (Example: Image Check) ---
	// Keep the image check, but make it optional or configurable?
	// For a generic file handler, you might remove this or check based on allowed types.
	// if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/") {
	// 	utils.LogMessage(utils.LevelWarn, "Uploaded file is not an image")
	// 	utils.LogFooter()
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Uploaded file type is not allowed (expected image)",
	// 	})
	// }

	// --- Filename Generation ---
	originalFilename := filepath.Base(fileHeader.Filename)
	fileExt := filepath.Ext(originalFilename)
	// Sanitize filename (remove potentially harmful characters, replace spaces, etc.) - IMPORTANT
	safeBaseName := sanitizeFilename(strings.TrimSuffix(originalFilename, fileExt))

	// Generate a unique filename to prevent collisions and overwrites
	// Using user email + timestamp + hash of original name seems robust
	timestamp := time.Now().UnixNano()
	hashInput := fmt.Sprintf("%s_%s_%d", email, originalFilename, timestamp)
	hash := sha256.Sum256([]byte(hashInput))
	uniqueID := hex.EncodeToString(hash[:8]) // Use first 8 bytes (16 hex chars)
	// Construct final filename: safeBaseName_uniqueID.ext
	finalFilename := fmt.Sprintf("%s_%s%s", safeBaseName, uniqueID, fileExt)

	// --- Saving File ---
	destinationPath := filepath.Join(h.DataFolder, finalFilename)
	utils.LogLineKeyValue(utils.LevelInfo, "Destination Path", destinationPath)

	// Open the uploaded file stream
	srcFile, err := fileHeader.Open()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to open uploaded file stream")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read uploaded file",
		})
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(destinationPath)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to create destination file")
		utils.LogLineKeyValue(utils.LevelError, "Path", destinationPath)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create file on server",
		})
	}
	defer dstFile.Close()

	// Copy the file contents
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to save file contents")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		// Clean up partially written file
		dstFile.Close()            // Close it first
		os.Remove(destinationPath) // Attempt removal
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file contents",
		})
	}

	utils.LogMessage(utils.LevelInfo, "File saved successfully")

	// --- Database Record ---
	var fileID int
	insertQuery := `
		INSERT INTO files (name, path, email, creation_date, last_access_date)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id_files`

	// Store the *original* filename in the DB for user display?
	// Store the *final* (unique) filename as the path identifier.
	// Note: 'path' column now stores the unique path in the data folder.
	err = h.DB.QueryRow(insertQuery, originalFilename, destinationPath, email).Scan(&fileID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to record file in database")
		utils.LogLineKeyValue(utils.LevelError, "Original Name", originalFilename)
		utils.LogLineKeyValue(utils.LevelError, "Path", destinationPath)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		// Critical: If DB fails, remove the orphaned file from disk
		os.Remove(destinationPath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to record file in database",
		})
	}

	utils.LogMessage(utils.LevelInfo, "File recorded in database")
	utils.LogLineKeyValue(utils.LevelInfo, "File ID", fileID)
	utils.LogFooter()

	// --- Response ---
	// Construct the URL based on the *final* filename used for serving
	// Assuming /api/data/ serves files based on the filename part of the path
	serveURL := fmt.Sprintf("/api/data/%s", finalFilename)

	return c.JSON(fiber.Map{
		"success":  true,
		"url":      serveURL,         // URL to access the file
		"file_id":  fileID,           // Database ID
		"filename": originalFilename, // Original filename for display
	})
}

// ServeFile serves files stored in the data directory.
func (h *FileHandler) ServeFile(c *fiber.Ctx) error {
	// Get filename from URL parameter
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Filename parameter is missing")
	}
	// Basic sanitization: prevent directory traversal
	filename = filepath.Base(filename) // Ensures we only have the filename part

	filePath := filepath.Join(h.DataFolder, filename)
	utils.LogHeader("📄 Serve File")
	utils.LogLineKeyValue(utils.LevelInfo, "Requested Filename", filename)
	utils.LogLineKeyValue(utils.LevelInfo, "Serving Path", filePath)

	// Check if file exists BEFORE trying to update DB
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.LogMessage(utils.LevelWarn, "File not found on disk")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).SendString("File not found")
	} else if err != nil {
		utils.LogMessage(utils.LevelError, "Error accessing file")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).SendString("Error accessing file")
	}

	// Update the last_access_date in the database asynchronously
	// We update based on the full path stored in the DB.
	go func(path string) {
		updateQuery := `UPDATE files SET last_access_date = CURRENT_TIMESTAMP WHERE path = $1`
		_, err := h.DB.Exec(updateQuery, path)
		if err != nil {
			// Log error but don't fail the file serving
			log.Printf("Warning: Failed to update last access date for '%s': %v", path, err)
		}
	}(filePath) // Pass the full path to the goroutine

	utils.LogMessage(utils.LevelInfo, "Serving file")
	utils.LogFooter()
	// Send the file using Fiber's efficient SendFile
	return c.SendFile(filePath) // Let Fiber handle Content-Type, etc.
}

// ListUserFiles lists all files uploaded by the logged-in user.
func (h *FileHandler) ListUserFiles(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	utils.LogHeader("📄 List User Files")
	utils.LogLineKeyValue(utils.LevelInfo, "User", email)

	query := `
		SELECT id_files, name, path, creation_date, last_access_date
		FROM files
		WHERE email = $1
		ORDER BY creation_date DESC;`

	rows, err := h.DB.Query(query, email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to retrieve user files from DB")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve files",
		})
	}
	defer rows.Close()

	var userFiles []fiber.Map
	for rows.Next() {
		var id int
		var name, path string
		var creationDate, lastAccessDate time.Time // Use time.Time

		if err := rows.Scan(&id, &name, &path, &creationDate, &lastAccessDate); err != nil {
			utils.LogMessage(utils.LevelError, "Error scanning user file row")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue // Skip this file
		}

		// Extract the filename from the path to build the serving URL
		serveFilename := filepath.Base(path)
		serveURL := fmt.Sprintf("/api/data/%s", serveFilename) // Assuming this route structure

		userFiles = append(userFiles, fiber.Map{
			"id":            id,
			"name":          name, // Original name
			"url":           serveURL,
			"created":       creationDate.Format(time.RFC3339),   // Format timestamp
			"last_accessed": lastAccessDate.Format(time.RFC3339), // Format timestamp
		})
	}

	if err = rows.Err(); err != nil {
		utils.LogMessage(utils.LevelError, "Error iterating user file rows")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}

	utils.LogMessage(utils.LevelInfo, "User files listed successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(userFiles))
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"success": true,
		"files":   userFiles,
	})
}

// DeleteFile handles the deletion of a specific file owned by the user.
func (h *FileHandler) DeleteFile(c *fiber.Ctx) error {
	// Filename here likely refers to the *original* filename the user knows.
	// We need to fetch the actual stored path from the DB based on the original name and user email.
	originalFilename := c.Params("filename")
	email := c.Locals("email").(string)

	utils.LogHeader("📄 Delete File")
	utils.LogLineKeyValue(utils.LevelInfo, "User", email)
	utils.LogLineKeyValue(utils.LevelInfo, "Requested Filename (Original)", originalFilename)

	if originalFilename == "" {
		utils.LogMessage(utils.LevelWarn, "Missing filename parameter for delete")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Filename parameter is required"})
	}

	// --- Transaction for safety ---
	tx, err := h.DB.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction for file deletion")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	// Defer rollback
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; rollback
		}
	}()

	// 1. Get the stored path from the database FOR UPDATE to lock the row
	var storedPath string
	query := `SELECT path FROM files WHERE name = $1 AND email = $2 FOR UPDATE;`
	err = tx.QueryRow(query, originalFilename, email).Scan(&storedPath)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogMessage(utils.LevelWarn, "File not found or permission denied")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "File not found or you don't have permission to delete it",
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to query file path for deletion")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error finding file"})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Stored Path", storedPath)

	// 2. Delete record from the database (within transaction)
	deleteQuery := `DELETE FROM files WHERE name = $1 AND email = $2;`
	var result sql.Result
	result, err = tx.Exec(deleteQuery, originalFilename, email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete file record from DB")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete file record"})
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		// Error getting rows affected (should be rare)
		utils.LogMessage(utils.LevelError, "Failed to get rows affected after delete")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to confirm file record deletion"})
	}
	if rowsAffected == 0 {
		// Should not happen if FOR UPDATE query succeeded, but check anyway
		utils.LogMessage(utils.LevelError, "Delete query affected 0 rows after finding the file")
		utils.LogFooter()
		err = fmt.Errorf("inconsistent state during file deletion") // Cause rollback
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete file record consistency error"})
	}
	utils.LogMessage(utils.LevelInfo, "File record deleted from database")

	// 3. Delete the actual file from disk *after* DB record is marked for deletion
	// This happens *before* commit to ensure consistency. If file delete fails, TX rollbacks.
	err = os.Remove(storedPath)
	if err != nil && !os.IsNotExist(err) { // Don't error if file was already gone
		utils.LogMessage(utils.LevelError, "Failed to delete file from disk")
		utils.LogLineKeyValue(utils.LevelError, "Path", storedPath)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		// Error is already set, defer will rollback the DB change
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete file from storage"})
	} else if err == nil {
		utils.LogMessage(utils.LevelInfo, "File deleted from disk")
	} else {
		utils.LogMessage(utils.LevelWarn, "File was already deleted from disk")
		err = nil // Reset error since it's not a failure condition for the operation
	}

	// 4. Commit transaction
	err = tx.Commit()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction for file deletion")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to finalize file deletion"})
	}

	utils.LogMessage(utils.LevelInfo, "File deleted successfully (DB and Disk)")
	utils.LogFooter()
	return c.JSON(fiber.Map{
		"success": true,
		"message": "File deleted successfully",
	})
}

// ListAllFiles (potentially admin-only) lists all files in the data directory.
// WARNING: This bypasses ownership checks and lists everything. Secure appropriately.
func (h *FileHandler) ListAllFiles(c *fiber.Ctx) error {
	// --- Permission Check Placeholder ---
	// role := c.Locals("role").(string) // Get role from JWT
	// if role != "ADMIN" { // Example admin check
	// 	 return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
	// }
	// --- End Permission Check ---

	utils.LogHeader("📄 List All Files (Admin)")

	// Read the configured data directory
	files, err := os.ReadDir(h.DataFolder)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to read data directory")
		utils.LogLineKeyValue(utils.LevelError, "Path", h.DataFolder)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to read data directory '%s'", h.DataFolder),
		})
	}

	var fileList []fiber.Map
	for _, file := range files {
		if !file.IsDir() { // List only files
			serveURL := fmt.Sprintf("/api/data/%s", file.Name()) // Assuming route structure
			fileList = append(fileList, fiber.Map{
				"name": file.Name(), // This is the unique stored name
				"url":  serveURL,
				// Add more info if needed (size, mod time - requires os.Stat)
			})
		}
	}

	utils.LogMessage(utils.LevelInfo, "Listed all files successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(fileList))
	utils.LogFooter()
	return c.JSON(fiber.Map{
		"success": true,
		"files":   fileList,
	})
}

// --- Helper Functions ---

// sanitizeFilename removes or replaces characters potentially harmful in filenames.
// This is a basic example; enhance as needed (e.g., handle unicode, length limits).
func sanitizeFilename(filename string) string {
	// Replace common problematic characters with underscore
	replacer := strings.NewReplacer(
		" ", "_",
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		"&", "_", // Add others as needed
		"..", "_", // Prevent path traversal component
	)
	sanitized := replacer.Replace(filename)

	// Remove leading/trailing dots or underscores if desired
	sanitized = strings.Trim(sanitized, "._")

	// Limit length?
	// maxLen := 100
	// if len(sanitized) > maxLen {
	//     sanitized = sanitized[:maxLen]
	// }

	// Ensure filename is not empty after sanitization
	if sanitized == "" {
		sanitized = "default_filename"
	}

	return sanitized
}
