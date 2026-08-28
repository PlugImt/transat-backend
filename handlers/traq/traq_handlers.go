package traq

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

const traqArticleSelect = `
	SELECT
		t.id_traq, t.name, t.description, t.picture, t.price, t.price_half,
		t.alcohol, t.creation_date, t.limited, t.out_of_stock, t.disabled,
		COALESCE(tt.name, 'Unknown') as traq_type
	FROM traq t
	LEFT JOIN traq_types tt ON t.id_traq_types = tt.id_traq_types
`

type TraqHandler struct {
	DB *sql.DB
}

func NewTraqHandler(db *sql.DB) *TraqHandler {
	return &TraqHandler{DB: db}
}

func parseIDParam(c *fiber.Ctx, name string) (int, error) {
	id, err := strconv.Atoi(c.Params(name))
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid %s", name)
	}
	return id, nil
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}

func isMissingTraqType(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, `null value in column "id_traq_types"`) ||
		strings.Contains(msg, "violates not-null constraint")
}

func scanTraqArticle(scanner interface{ Scan(dest ...any) error }) (models.TraqArticle, error) {
	var article models.TraqArticle
	var creationDate sql.NullTime

	err := scanner.Scan(
		&article.ID, &article.Name, &article.Description, &article.Picture,
		&article.Price, &article.PriceHalf, &article.Alcohol,
		&creationDate,
		&article.Limited, &article.OutOfStock, &article.Disabled,
		&article.TraqType,
	)
	if err != nil {
		return article, err
	}
	if creationDate.Valid {
		article.CreationDate = creationDate.Time.Format(time.RFC3339)
	}
	return article, nil
}

func (h *TraqHandler) listTraqArticles(query string, args ...any) ([]models.TraqArticle, error) {
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]models.TraqArticle, 0)
	for rows.Next() {
		article, err := scanTraqArticle(rows)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan Traq article row")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}
		articles = append(articles, article)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return articles, nil
}

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

	if strings.TrimSpace(traqType.Name) == "" {
		utils.LogMessage(utils.LevelWarn, "Missing name field for Traq type")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Traq type name is required",
		})
	}

	query := `INSERT INTO traq_types (name) VALUES ($1) RETURNING id_traq_types;`
	err := h.DB.QueryRow(query, traqType.Name).Scan(&traqType.IDType)
	if err != nil {
		if isUniqueViolation(err) {
			utils.LogMessage(utils.LevelWarn, "Traq type already exists")
			utils.LogLineKeyValue(utils.LevelWarn, "Name", traqType.Name)
			utils.LogFooter()
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Traq type already exists",
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to create Traq type")
		utils.LogLineKeyValue(utils.LevelError, "Name", traqType.Name)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create Traq type",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Traq type created successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "ID", traqType.IDType)
	utils.LogLineKeyValue(utils.LevelInfo, "Name", traqType.Name)
	utils.LogFooter()
	return c.Status(fiber.StatusCreated).JSON(traqType)
}

func (h *TraqHandler) GetAllTraqTypes(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Get All Traq Types")

	query := `SELECT id_traq_types, name FROM traq_types ORDER BY name;`
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

	traqTypes := make([]models.TraqType, 0)
	for rows.Next() {
		var traqType models.TraqType
		if err := rows.Scan(&traqType.IDType, &traqType.Name); err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan Traq type row")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}
		traqTypes = append(traqTypes, traqType)
	}

	if err := rows.Err(); err != nil {
		utils.LogMessage(utils.LevelError, "Error iterating Traq type rows")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}

	utils.LogMessage(utils.LevelInfo, "Traq types retrieved successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(traqTypes))
	utils.LogFooter()

	return c.JSON(traqTypes)
}

func (h *TraqHandler) GetTraqType(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Get Traq Type")

	id, err := parseIDParam(c, "id")
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid Traq type ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Traq type ID"})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Type ID", id)

	var traqType models.TraqType
	err = h.DB.QueryRow(`SELECT id_traq_types, name FROM traq_types WHERE id_traq_types = $1`, id).
		Scan(&traqType.IDType, &traqType.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Traq type not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Traq type not found"})
		}
		utils.LogMessage(utils.LevelError, "Failed to retrieve Traq type")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve Traq type"})
	}

	utils.LogMessage(utils.LevelInfo, "Traq type retrieved successfully")
	utils.LogFooter()
	return c.JSON(traqType)
}

func (h *TraqHandler) UpdateTraqType(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Update Traq Type")

	id, err := parseIDParam(c, "id")
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid Traq type ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Traq type ID"})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Type ID", id)

	var req struct {
		Name *string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == nil || strings.TrimSpace(*req.Name) == "" {
		utils.LogMessage(utils.LevelWarn, "Missing name field for Traq type")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Traq type name is required"})
	}

	result, err := h.DB.Exec(`UPDATE traq_types SET name = $1 WHERE id_traq_types = $2`, *req.Name, id)
	if err != nil {
		if isUniqueViolation(err) {
			utils.LogMessage(utils.LevelWarn, "Traq type name already exists")
			utils.LogLineKeyValue(utils.LevelWarn, "Name", *req.Name)
			utils.LogFooter()
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Traq type already exists"})
		}
		utils.LogMessage(utils.LevelError, "Failed to update Traq type")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update Traq type"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "Traq type not found for update")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Traq type not found"})
	}

	utils.LogMessage(utils.LevelInfo, "Traq type updated successfully")
	utils.LogFooter()
	return c.JSON(fiber.Map{"message": "Traq type updated successfully"})
}

func (h *TraqHandler) DeleteTraqType(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Delete Traq Type")

	id, err := parseIDParam(c, "id")
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid Traq type ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Traq type ID"})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Type ID", id)

	var articleCount int
	if err := h.DB.QueryRow(`SELECT COUNT(*) FROM traq WHERE id_traq_types = $1`, id).Scan(&articleCount); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to check Traq type usage")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete Traq type"})
	}
	if articleCount > 0 {
		utils.LogMessage(utils.LevelWarn, "Cannot delete Traq type still in use")
		utils.LogLineKeyValue(utils.LevelWarn, "Article Count", articleCount)
		utils.LogFooter()
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Cannot delete Traq type that still has articles",
		})
	}

	result, err := h.DB.Exec(`DELETE FROM traq_types WHERE id_traq_types = $1`, id)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete Traq type")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete Traq type"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "Traq type not found for deletion")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Traq type not found"})
	}

	utils.LogMessage(utils.LevelInfo, "Traq type deleted successfully")
	utils.LogFooter()
	return c.SendStatus(fiber.StatusOK)
}

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

	if article.Name == "" || article.Description == "" || article.Picture == "" || article.TraqType == "" {
		utils.LogMessage(utils.LevelWarn, "Missing required fields for Traq article")
		utils.LogLineKeyValue(utils.LevelWarn, "Received", article)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, description, picture, and type are required fields",
		})
	}

	query := `
		INSERT INTO traq (name, description, picture, id_traq_types, limited, price, disabled, alcohol, out_of_stock, price_half)
		VALUES ($1, $2, $3, (SELECT id_traq_types FROM traq_types WHERE name = $4), $5, $6, $7, $8, $9, $10)
		RETURNING id_traq;
	`

	err := h.DB.QueryRow(
		query,
		article.Name, article.Description, article.Picture, article.TraqType,
		article.Limited, article.Price, article.Disabled, article.Alcohol, article.OutOfStock, article.PriceHalf,
	).Scan(&article.ID)

	if err != nil {
		if isMissingTraqType(err) {
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
	utils.LogLineKeyValue(utils.LevelInfo, "Article ID", article.ID)
	utils.LogLineKeyValue(utils.LevelInfo, "Name", article.Name)
	utils.LogFooter()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Article created successfully",
		"id_traq": article.ID,
	})
}

func (h *TraqHandler) GetAllTraqArticles(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Get All Traq Articles")

	articles, err := h.listTraqArticles(traqArticleSelect + ` ORDER BY tt.name, t.name;`)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to query Traq articles")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve Traq articles",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Traq articles retrieved successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(articles))
	utils.LogFooter()

	return c.JSON(articles)
}

func (h *TraqHandler) GetAvailableTraqArticles(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Get Available Traq Articles")

	articles, err := h.listTraqArticles(traqArticleSelect + ` WHERE t.disabled = FALSE ORDER BY tt.name, t.name;`)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to query available Traq articles")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve Traq articles",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Available Traq articles retrieved successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(articles))
	utils.LogFooter()

	return c.JSON(articles)
}

func (h *TraqHandler) GetTraqArticle(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Get Traq Article")

	id, err := parseIDParam(c, "id")
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid Traq article ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Traq article ID"})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Article ID", id)

	article, err := scanTraqArticle(h.DB.QueryRow(traqArticleSelect+` WHERE t.id_traq = $1`, id))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Traq article not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Traq article not found"})
		}
		utils.LogMessage(utils.LevelError, "Failed to retrieve Traq article")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve Traq article"})
	}

	utils.LogMessage(utils.LevelInfo, "Traq article retrieved successfully")
	utils.LogFooter()
	return c.JSON(article)
}

func (h *TraqHandler) UpdateTraqArticle(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Update Traq Article")

	id, err := parseIDParam(c, "id")
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid Traq article ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Traq article ID"})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Article ID", id)

	var req struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Picture     *string  `json:"picture"`
		TraqType    *string  `json:"traq_type"`
		Limited     *bool    `json:"limited"`
		Price       *float32 `json:"price"`
		Disabled    *bool    `json:"disabled"`
		Alcohol     *float32 `json:"alcohol"`
		OutOfStock  *bool    `json:"out_of_stock"`
		PriceHalf   *float32 `json:"price_half"`
	}
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var updateFields []string
	var updateValues []any
	paramCount := 1

	if req.Name != nil {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", paramCount))
		updateValues = append(updateValues, *req.Name)
		paramCount++
	}
	if req.Description != nil {
		updateFields = append(updateFields, fmt.Sprintf("description = $%d", paramCount))
		updateValues = append(updateValues, *req.Description)
		paramCount++
	}
	if req.Picture != nil {
		updateFields = append(updateFields, fmt.Sprintf("picture = $%d", paramCount))
		updateValues = append(updateValues, *req.Picture)
		paramCount++
	}
	if req.TraqType != nil {
		updateFields = append(updateFields, fmt.Sprintf("id_traq_types = (SELECT id_traq_types FROM traq_types WHERE name = $%d)", paramCount))
		updateValues = append(updateValues, *req.TraqType)
		paramCount++
	}
	if req.Limited != nil {
		updateFields = append(updateFields, fmt.Sprintf("limited = $%d", paramCount))
		updateValues = append(updateValues, *req.Limited)
		paramCount++
	}
	if req.Price != nil {
		updateFields = append(updateFields, fmt.Sprintf("price = $%d", paramCount))
		updateValues = append(updateValues, *req.Price)
		paramCount++
	}
	if req.Disabled != nil {
		updateFields = append(updateFields, fmt.Sprintf("disabled = $%d", paramCount))
		updateValues = append(updateValues, *req.Disabled)
		paramCount++
	}
	if req.Alcohol != nil {
		updateFields = append(updateFields, fmt.Sprintf("alcohol = $%d", paramCount))
		updateValues = append(updateValues, *req.Alcohol)
		paramCount++
	}
	if req.OutOfStock != nil {
		updateFields = append(updateFields, fmt.Sprintf("out_of_stock = $%d", paramCount))
		updateValues = append(updateValues, *req.OutOfStock)
		paramCount++
	}
	if req.PriceHalf != nil {
		updateFields = append(updateFields, fmt.Sprintf("price_half = $%d", paramCount))
		updateValues = append(updateValues, *req.PriceHalf)
		paramCount++
	}

	if len(updateFields) == 0 {
		utils.LogMessage(utils.LevelWarn, "No fields to update")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No fields to update"})
	}

	updateValues = append(updateValues, id)
	query := fmt.Sprintf(`UPDATE traq SET %s WHERE id_traq = $%d`, strings.Join(updateFields, ", "), paramCount)

	result, err := h.DB.Exec(query, updateValues...)
	if err != nil {
		if req.TraqType != nil && isMissingTraqType(err) {
			utils.LogMessage(utils.LevelWarn, "Failed to update article: TraqType not found")
			utils.LogLineKeyValue(utils.LevelWarn, "Type Name", *req.TraqType)
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Traq type '%s' not found", *req.TraqType),
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to update Traq article")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update Traq article"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "Traq article not found for update")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Traq article not found"})
	}

	utils.LogMessage(utils.LevelInfo, "Traq article updated successfully")
	utils.LogFooter()
	return c.JSON(fiber.Map{"message": "Article updated successfully"})
}

func (h *TraqHandler) DeleteTraqArticle(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Delete Traq Article")

	id, err := parseIDParam(c, "id")
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid Traq article ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Traq article ID"})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Article ID", id)

	result, err := h.DB.Exec(`DELETE FROM traq WHERE id_traq = $1`, id)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete Traq article")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete Traq article"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "Traq article not found for deletion")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Traq article not found"})
	}

	utils.LogMessage(utils.LevelInfo, "Traq article deleted successfully")
	utils.LogFooter()
	return c.SendStatus(fiber.StatusOK)
}
