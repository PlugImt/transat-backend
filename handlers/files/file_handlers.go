package files

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

type File struct {
	File string `format:"binary"`
}

// FileHandler handles file upload, serving, listing, and deletion.
type FileHandler struct {
	DB        *sql.DB
	R2Service *services.R2Service
}

// NewFileHandler creates a new FileHandler.
func NewFileHandler(db *sql.DB, r2Service *services.R2Service) (*FileHandler, error) {
	return &FileHandler{
		DB:        db,
		R2Service: r2Service,
	}, nil
}

// UploadFile handles file uploads, saves them, and records them in the database.
// @Summary		Télécharger un fichier
// @Description	Télécharge un fichier et l'enregistre dans le stockage cloud
// @Tags			Files
// @Produce		json
// @Security		BearerAuth
// @Accept	multipart/form-data
// @Param 		file 	formData 	File true "Fichier à uploader"
// @Success		200		{object}	map[string]interface{}	"Fichier téléchargé avec succès"
// @Failure		400		{object}	models.ErrorResponse	"Données invalides"
// @Failure		401		{object}	models.ErrorResponse	"Non autorisé"
// @Failure		500		{object}	models.ErrorResponse	"Erreur serveur"
// @Router			/upload [post]
func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	utils.LogHeader("📄 Upload File")

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

	// Check if the request is multipart form
	if !strings.HasPrefix(c.Get("Content-Type"), "multipart/form-data") {
		utils.LogMessage(utils.LevelError, "Invalid Content-Type, expected multipart/form-data")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Content-Type, expected multipart/form-data",
		})
	}

	// Get the uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		if fileHeader, err = c.FormFile("image"); err != nil {
			utils.LogMessage(utils.LevelError, "No file uploaded with key 'file' or 'image'")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
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

	// Generate unique filename
	originalFilename := filepath.Base(fileHeader.Filename)
	fileExt := filepath.Ext(originalFilename)
	safeBaseName := sanitizeFilename(strings.TrimSuffix(originalFilename, fileExt))

	timestamp := time.Now().UnixNano()
	hashInput := fmt.Sprintf("%s_%s_%d", email, originalFilename, timestamp)
	hash := sha256.Sum256([]byte(hashInput))
	uniqueID := hex.EncodeToString(hash[:8])
	finalFilename := fmt.Sprintf("%s_%s%s", safeBaseName, uniqueID, fileExt)

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to open uploaded file")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process uploaded file",
		})
	}
	defer file.Close()

	// Upload to R2
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	publicURL, err := h.R2Service.UploadFile(finalFilename, file, contentType)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to upload file to R2")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store file",
		})
	}

	// Store in database - only store the filename, not the full URL
	var fileID int
	insertQuery := `
		INSERT INTO files (name, path, email, creation_date, last_access_date)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id_files`

	err = h.DB.QueryRow(insertQuery, originalFilename, finalFilename, email).Scan(&fileID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to record file in database")
		utils.LogLineKeyValue(utils.LevelError, "Original Name", originalFilename)
		utils.LogLineKeyValue(utils.LevelError, "Path", finalFilename)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()

		// Try to delete the uploaded file from R2
		if delErr := h.R2Service.DeleteFile(finalFilename); delErr != nil {
			utils.LogMessage(utils.LevelError, "Failed to delete file from R2 after DB error")
			utils.LogLineKeyValue(utils.LevelError, "Error", delErr)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to record file in database",
		})
	}

	utils.LogMessage(utils.LevelInfo, "File uploaded successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "File ID", fileID)
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"success":  true,
		"url":      publicURL,
		"file_id":  fileID,
		"filename": originalFilename,
	})
}

// ServeFile serves files stored in R2.
// @Summary		Servir un fichier
// @Description	Sert un fichier stocké dans le cloud storage
// @Tags			Files
// @Produce		application/octet-stream
// @Param			filename	path	string	true	"Nom du fichier à servir"
// @Success		200			{object}	File	"Fichier servi avec succès"
// @Failure		400			{string}	string	"Nom de fichier manquant"
// @Failure		404			{string}	string	"Fichier non trouvé"
// @Failure		500			{string}	string	"Erreur serveur"
// @Router			/data/{filename} [get]
func (h *FileHandler) ServeFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Filename parameter is missing")
	}

	utils.LogHeader("📄 Serve File")
	utils.LogLineKeyValue(utils.LevelInfo, "Requested Filename", filename)

	// Update last access date in database
	updateQuery := `UPDATE files SET last_access_date = CURRENT_TIMESTAMP WHERE path = $1`
	go func() {
		if _, err := h.DB.Exec(updateQuery, filename); err != nil {
			log.Printf("Warning: Failed to update last access date for '%s': %v", filename, err)
		}
	}()

	// Get the file from R2
	reader, err := h.R2Service.GetObject(filename)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get file from R2")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).SendString("File not found")
	}
	defer reader.Close()

	// Stream the file to the response
	return c.SendStream(reader)
}

// ListUserFiles lists all files uploaded by the logged-in user.
// @Summary		Lister les fichiers de l'utilisateur
// @Description	Récupère la liste de tous les fichiers téléchargés par l'utilisateur connecté
// @Tags			Files
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	map[string]interface{}	"Liste des fichiers récupérée avec succès"
// @Failure		401	{object}	models.ErrorResponse	"Non autorisé"
// @Failure		500	{object}	models.ErrorResponse	"Erreur serveur"
// @Router			/files [get]
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
		var creationDate, lastAccessDate time.Time

		if err := rows.Scan(&id, &name, &path, &creationDate, &lastAccessDate); err != nil {
			utils.LogMessage(utils.LevelError, "Error scanning user file row")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}

		// Si le chemin contient déjà l'URL complète, l'utiliser directement
		// Sinon, construire l'URL avec le domaine public
		var publicURL string
		if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
			publicURL = path
		} else {
			publicURL = h.R2Service.GetPublicURL(path)
		}

		userFiles = append(userFiles, fiber.Map{
			"id":            id,
			"name":          name,
			"url":           publicURL,
			"created":       creationDate.Format(time.RFC3339),
			"last_accessed": lastAccessDate.Format(time.RFC3339),
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
// @Summary		Supprimer un fichier
// @Description	Supprime un fichier spécifique appartenant à l'utilisateur connecté
// @Tags			Files
// @Produce		json
// @Security		BearerAuth
// @Param			filename	path	string	true	"Nom du fichier à supprimer"
// @Success		200			{object}	models.Response			"Fichier supprimé avec succès"
// @Failure		400			{object}	models.ErrorResponse	"Nom de fichier manquant"
// @Failure		401			{object}	models.ErrorResponse	"Non autorisé"
// @Failure		404			{object}	models.ErrorResponse	"Fichier non trouvé"
// @Failure		500			{object}	models.ErrorResponse	"Erreur serveur"
// @Router			/files/{filename} [delete]
func (h *FileHandler) DeleteFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	email := c.Locals("email").(string)

	utils.LogHeader("📄 Delete File")
	utils.LogLineKeyValue(utils.LevelInfo, "User", email)
	utils.LogLineKeyValue(utils.LevelInfo, "Requested Filename", filename)

	if filename == "" {
		utils.LogMessage(utils.LevelWarn, "Missing filename parameter for delete")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Filename parameter is required"})
	}

	tx, err := h.DB.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction for file deletion")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	defer tx.Rollback()

	var storedPath string
	query := `SELECT path FROM files WHERE name = $1 AND email = $2 FOR UPDATE;`
	err = tx.QueryRow(query, filename, email).Scan(&storedPath)
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

	// Delete from R2 first
	if err := h.R2Service.DeleteFile(storedPath); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete file from R2")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete file from storage"})
	}

	// Then delete from database
	deleteQuery := `DELETE FROM files WHERE name = $1 AND email = $2;`
	result, err := tx.Exec(deleteQuery, filename, email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete file record from DB")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete file record"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get rows affected after delete")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to confirm file record deletion"})
	}

	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelError, "Delete query affected 0 rows after finding the file")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete file record consistency error"})
	}

	if err = tx.Commit(); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction for file deletion")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to finalize file deletion"})
	}

	utils.LogMessage(utils.LevelInfo, "File deleted successfully")
	utils.LogFooter()
	return c.JSON(fiber.Map{
		"success": true,
		"message": "File deleted successfully",
	})
}

// ListAllFiles (potentially admin-only) lists all files.
func (h *FileHandler) ListAllFiles(c *fiber.Ctx) error {
	utils.LogHeader("📄 List All Files (Admin)")

	query := `
		SELECT id_files, name, path, email, creation_date, last_access_date
		FROM files
		ORDER BY creation_date DESC;`

	rows, err := h.DB.Query(query)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to retrieve files from DB")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve files",
		})
	}
	defer rows.Close()

	var files []fiber.Map
	for rows.Next() {
		var id int
		var name, path, email string
		var creationDate, lastAccessDate time.Time

		if err := rows.Scan(&id, &name, &path, &email, &creationDate, &lastAccessDate); err != nil {
			utils.LogMessage(utils.LevelError, "Error scanning file row")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}

		publicURL := h.R2Service.GetPublicURL(path)

		files = append(files, fiber.Map{
			"id":            id,
			"name":          name,
			"url":           publicURL,
			"email":         email,
			"created":       creationDate.Format(time.RFC3339),
			"last_accessed": lastAccessDate.Format(time.RFC3339),
		})
	}

	if err = rows.Err(); err != nil {
		utils.LogMessage(utils.LevelError, "Error iterating file rows")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}

	utils.LogMessage(utils.LevelInfo, "Listed all files successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(files))
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"success": true,
		"files":   files,
	})
}

// sanitizeFilename removes or replaces characters potentially harmful in filenames.
func sanitizeFilename(filename string) string {
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
		"&", "_",
		"..", "_",
	)
	sanitized := replacer.Replace(filename)
	sanitized = strings.Trim(sanitized, "._")
	if sanitized == "" {
		sanitized = "default_filename"
	}
	return sanitized
}
