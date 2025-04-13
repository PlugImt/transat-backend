package traq

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"Transat_2.0_Backend/models"
	"Transat_2.0_Backend/utils" // Assuming logger is here
	"github.com/gofiber/fiber/v2"
)

// TraqHandler handles requests related to Traq articles and types.
type TraqHandler struct {
	DB *sql.DB
}

// NewTraqHandler creates a new TraqHandler.
func NewTraqHandler(db *sql.DB) *TraqHandler {
	return &TraqHandler{DB: db}
}

// CreateTraqType handles the creation of a new Traq type.
func (h *TraqHandler) CreateTraqType(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Create Traq Type")
	var traqType models.TraqType
	if err := c.BodyParser(&traqType); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if traqType.Name == "" {
		utils.LogMessage(utils.LevelWarn, "Missing name field for Traq type")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Traq type name is required",
		})
	}

	query := `INSERT INTO traq_types (name) VALUES ($1);`
	_, err := h.DB.Exec(query, traqType.Name)
	if err != nil {
		// Handle potential duplicate error more gracefully if needed
		utils.LogMessage(utils.LevelError, "Failed to create Traq type")
		utils.LogLineKeyValue(utils.LevelError, "Name", traqType.Name)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create Traq type",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Traq type created successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Name", traqType.Name)
	utils.LogFooter()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Traq type created successfully",
	})
}

// GetAllTraqTypes retrieves all available Traq types.
func (h *TraqHandler) GetAllTraqTypes(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Get All Traq Types")

	query := `SELECT id_traq_types, name FROM traq_types ORDER BY name;` // Added ID and ordering
	rows, err := h.DB.Query(query)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to query Traq types")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve Traq types",
		})
	}
	defer rows.Close()

	var traqTypes []models.TraqType
	for rows.Next() {
		var traqType models.TraqType
		// Scan both ID and Name
		if err := rows.Scan(&traqType.IDType, &traqType.Name); err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan Traq type row")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			// Continue to next row? Or return error?
			// Let's log and continue for now.
			continue
		}
		traqTypes = append(traqTypes, traqType)
	}

	if err := rows.Err(); err != nil {
		utils.LogMessage(utils.LevelError, "Error iterating Traq type rows")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		// Return potentially partial data or error?
	}

	utils.LogMessage(utils.LevelInfo, "Traq types retrieved successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(traqTypes))
	utils.LogFooter()

	return c.JSON(traqTypes)
}

// CreateTraqArticle handles the creation of a new Traq article.
func (h *TraqHandler) CreateTraqArticle(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Create Traq Article")
	var article models.TraqArticle
	if err := c.BodyParser(&article); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if article.Name == "" || article.Description == "" || article.Picture == "" || article.TraqType == "" {
		utils.LogMessage(utils.LevelWarn, "Missing required fields for Traq article")
		utils.LogLineKeyValue(utils.LevelWarn, "Received", article) // Log received data for debugging
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, description, picture, and type are required fields",
		})
	}

	// Query to insert article, getting the type ID via subquery
	query := `
		INSERT INTO traq (name, description, picture, id_traq_types, limited, price, disabled, alcohol, out_of_stock, price_half)
		VALUES ($1, $2, $3, (SELECT id_traq_types FROM traq_types WHERE name = $4), $5, $6, $7, $8, $9, $10)
		RETURNING id_traq; -- Return the newly created article ID
	`

	var newArticleID int
	err := h.DB.QueryRow(
		query,
		article.Name, article.Description, article.Picture, article.TraqType,
		article.Limited, article.Price, article.Disabled, article.Alcohol, article.OutOfStock, article.PriceHalf,
	).Scan(&newArticleID)

	if err != nil {
		// Check if error is due to non-existent traq_type
		// This might involve checking the error string or doing a pre-check
		if strings.Contains(err.Error(), "null value in column \"id_traq_types\"") {
			utils.LogMessage(utils.LevelWarn, "Failed to create article: TraqType not found")
			utils.LogLineKeyValue(utils.LevelWarn, "Type Name", article.TraqType)
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Traq type '%s' not found", article.TraqType),
			})
		}

		utils.LogMessage(utils.LevelError, "Failed to create Traq article")
		utils.LogLineKeyValue(utils.LevelError, "Article Name", article.Name)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create Traq article",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Traq article created successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Article ID", newArticleID)
	utils.LogLineKeyValue(utils.LevelInfo, "Name", article.Name)
	utils.LogFooter()

	// Return created status and maybe the new ID or the full article?
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Article created successfully",
		"id_traq": newArticleID,
	})
}

// GetAllTraqArticles retrieves all Traq articles.
func (h *TraqHandler) GetAllTraqArticles(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Get All Traq Articles")

	query := `
		SELECT
			t.id_traq, t.name, t.description, t.picture, t.price, t.price_half,
			t.alcohol, t.creation_date, t.limited, t.out_of_stock, t.disabled,
			COALESCE(tt.name, 'Unknown') as traq_type -- Get type name, default if missing
		FROM traq t
		LEFT JOIN traq_types tt ON t.id_traq_types = tt.id_traq_types
		ORDER BY tt.name, t.name; -- Order by type then name
	`

	rows, err := h.DB.Query(query)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to query Traq articles")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve Traq articles",
		})
	}
	defer rows.Close()

	var articles []models.TraqArticle
	for rows.Next() {
		var article models.TraqArticle
		var creationDate sql.NullTime // Use NullTime for dates

		err := rows.Scan(
			&article.ID, &article.Name, &article.Description, &article.Picture,
			&article.Price, &article.PriceHalf, &article.Alcohol,
			&creationDate, // Scan into NullTime
			&article.Limited, &article.OutOfStock, &article.Disabled,
			&article.TraqType,
		)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan Traq article row")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue // Skip this article
		}
		// Format date if valid
		if creationDate.Valid {
			article.CreationDate = creationDate.Time.Format(time.RFC3339) // Or desired format
		} else {
			article.CreationDate = "" // Or "N/A"
		}

		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		utils.LogMessage(utils.LevelError, "Error iterating Traq article rows")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		// Return potentially partial data or error?
	}

	utils.LogMessage(utils.LevelInfo, "Traq articles retrieved successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(articles))
	utils.LogFooter()

	return c.JSON(articles)
}

// --- Placeholder Handlers for potential future routes ---

// GetTraqArticle (Placeholder)
func (h *TraqHandler) GetTraqArticle(c *fiber.Ctx) error {
	// id := c.Params("id")
	// Fetch single article by ID
	return c.Status(fiber.StatusNotImplemented).SendString("GetTraqArticle not implemented")
}

// UpdateTraqArticle (Placeholder)
func (h *TraqHandler) UpdateTraqArticle(c *fiber.Ctx) error {
	// id := c.Params("id")
	// Update article by ID
	return c.Status(fiber.StatusNotImplemented).SendString("UpdateTraqArticle not implemented")
}

// DeleteTraqArticle (Placeholder)
func (h *TraqHandler) DeleteTraqArticle(c *fiber.Ctx) error {
	// id := c.Params("id")
	// Delete article by ID
	return c.Status(fiber.StatusNotImplemented).SendString("DeleteTraqArticle not implemented")
}
