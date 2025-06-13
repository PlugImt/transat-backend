package planning

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type PlanningHandler struct {
	db *sql.DB
}

func NewPlanningHandler(db *sql.DB) *PlanningHandler {
	return &PlanningHandler{db: db}
}

// GET /planning/users
// @Summary		Obtenir les utilisateurs avec Pass ID
// @Description	Récupère la liste de tous les utilisateurs avec leurs informations de Pass ID pour la planification
// @Tags			Planning
// @Produce		json
// @Security		BearerAuth
// @Success		200	{array}		models.UserWithPassID	"Liste des utilisateurs récupérée avec succès"
// @Failure		401	{object}	models.ErrorResponse	"Non autorisé"
// @Failure		500	{object}	models.ErrorResponse	"Erreur serveur"
// @Router			/planning/users [get]
func (h *PlanningHandler) GetUsersWithPassID(c *fiber.Ctx) error {
	utils.LogHeader("👥 Get Users With Pass ID")

	rows, err := h.db.Query(`SELECT id_newf, first_name, last_name, pass_id FROM newf`)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch users")
		utils.LogFooter()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var firstName, lastName sql.NullString
		var passID sql.NullInt64
		if err := rows.Scan(&id, &firstName, &lastName, &passID); err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan user")
			utils.LogFooter()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan user"})
		}
		var passIDValue interface{}
		if passID.Valid {
			passIDValue = passID.Int64
		} else {
			passIDValue = nil
		}
		users = append(users, fiber.Map{
			"id":         id,
			"first_name": firstName.String,
			"last_name":  lastName.String,
			"pass_id":    passIDValue,
		})
	}
	utils.LogMessage(utils.LevelInfo, "Successfully fetched users")
	utils.LogFooter()
	return c.JSON(users)
}

// PATCH /planning/users/:id_newf/passid
// @Summary		Mettre à jour le Pass ID d'un utilisateur
// @Description	Met à jour le Pass ID d'un utilisateur spécifique s'il est actuellement nul
// @Tags			Planning
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id_newf	path	int								true	"ID de l'utilisateur"
// @Param			passid	body	models.UpdatePassIDRequest	true	"Nouveau Pass ID"
// @Success		200		{object}	models.Response			"Pass ID mis à jour avec succès"
// @Failure		400		{object}	models.ErrorResponse	"Données invalides"
// @Failure		401		{object}	models.ErrorResponse	"Non autorisé"
// @Failure		409		{object}	models.ErrorResponse	"Pass ID déjà défini"
// @Failure		500		{object}	models.ErrorResponse	"Erreur serveur"
// @Router			/planning/users/{id_newf}/passid [patch]
func (h *PlanningHandler) UpdateUserPassIDIfNull(c *fiber.Ctx) error {
	utils.LogHeader("🔄 Update User Pass ID If Null")

	idStr := c.Params("id_newf")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Invalid id_newf parameter")
		utils.LogFooter()
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid id_newf parameter"})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Newf ID", id)

	var body models.UpdatePassIDRequest
	if err := c.BodyParser(&body); err != nil {
		utils.LogMessage(utils.LevelError, "Invalid request body")
		utils.LogFooter()
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Requested Pass ID", body.PassID)

	// Only update if pass_id is NULL or 0
	res, err := h.db.Exec(`UPDATE newf SET pass_id = $1 WHERE id_newf = $2 AND (pass_id IS NULL OR pass_id = 0)`, body.PassID, id)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to update pass_id")
		utils.LogFooter()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update pass_id"})
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		utils.LogMessage(utils.LevelWarn, "PassID already set or user not found")
		utils.LogFooter()
		return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "PassID already set or user not found"})
	}
	utils.LogMessage(utils.LevelInfo, "Successfully updated PassID")
	utils.LogFooter()
	return c.JSON(fiber.Map{"success": true})
}

// GET /planning/users/:email/courses?start=YYYY-MM-DD&end=YYYY-MM-DD
// @Summary		Obtenir les cours d'un utilisateur entre deux dates
// @Description	Récupère la liste des cours d'un utilisateur pour une période donnée
// @Tags			Planning
// @Produce		json
// @Security		BearerAuth
// @Param			email	path	string	true	"Email de l'utilisateur"
// @Param			start	query	string	true	"Date de début (YYYY-MM-DD)"	example(2024-01-01)
// @Param			end		query	string	true	"Date de fin (YYYY-MM-DD)"		example(2024-01-31)
// @Success		200		{array}		models.Course			"Liste des cours récupérée avec succès"
// @Failure		400		{object}	models.ErrorResponse	"Paramètres manquants ou invalides"
// @Failure		401		{object}	models.ErrorResponse	"Non autorisé"
// @Failure		500		{object}	models.ErrorResponse	"Erreur serveur"
// @Router			/planning/users/{email}/courses [get]
func (h *PlanningHandler) GetUserCoursesBetweenDates(c *fiber.Ctx) error {
	utils.LogHeader("📅 Get User Courses Between Dates")

	email := c.Params("email")
	start := c.Query("start")
	end := c.Query("end")

	utils.LogLineKeyValue(utils.LevelInfo, "User Email", email)
	utils.LogLineKeyValue(utils.LevelInfo, "Start Date", start)
	utils.LogLineKeyValue(utils.LevelInfo, "End Date", end)

	if start == "" || end == "" {
		utils.LogMessage(utils.LevelError, "Missing start or end date")
		utils.LogFooter()
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing start or end date"})
	}

	rows, err := h.db.Query(`SELECT * FROM courses WHERE user_email = $1 AND date BETWEEN $2 AND $3 ORDER BY start_time`, email, start, end)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch courses")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch courses", "details": err.Error()})
	}
	defer rows.Close()

	var courses []map[string]interface{}
	cols, _ := rows.Columns()
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		valPtrs := make([]interface{}, len(cols))
		for i := range vals {
			valPtrs[i] = &vals[i]
		}
		if err := rows.Scan(valPtrs...); err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan course")
			utils.LogFooter()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan course"})
		}
		course := make(map[string]interface{})
		for i, col := range cols {
			course[col] = vals[i]
		}
		courses = append(courses, course)
	}
	utils.LogMessage(utils.LevelInfo, "Successfully fetched courses")
	utils.LogFooter()
	return c.JSON(courses)
}

// GET /planning/users/:email/courses/today
// @Summary		Obtenir les cours d'aujourd'hui pour un utilisateur
// @Description	Récupère la liste des cours d'un utilisateur pour la journée actuelle
// @Tags			Planning
// @Produce		json
// @Security		BearerAuth
// @Param			email	path	string	true	"Email de l'utilisateur"
// @Success		200		{array}		models.Course			"Liste des cours d'aujourd'hui récupérée avec succès"
// @Failure		400		{object}	models.ErrorResponse	"Email manquant"
// @Failure		401		{object}	models.ErrorResponse	"Non autorisé"
// @Failure		500		{object}	models.ErrorResponse	"Erreur serveur"
// @Router			/planning/users/{email}/courses/today [get]
func (h *PlanningHandler) GetUserCoursesToday(c *fiber.Ctx) error {
	utils.LogHeader("📅 Get User Courses Today")

	email := c.Params("email")
	if email == "" {
		utils.LogMessage(utils.LevelError, "Missing email parameter")
		utils.LogFooter()
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing email parameter"})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "User Email", email)

	today := time.Now().Format("2006-01-02")
	utils.LogLineKeyValue(utils.LevelInfo, "Today", today)

	rows, err := h.db.Query(`SELECT * FROM courses WHERE user_email = $1 AND date = $2 ORDER BY start_time`, email, today)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch courses")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch courses", "details": err.Error()})
	}
	defer rows.Close()

	var courses []map[string]interface{}
	cols, _ := rows.Columns()
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		valPtrs := make([]interface{}, len(cols))
		for i := range vals {
			valPtrs[i] = &vals[i]
		}
		if err := rows.Scan(valPtrs...); err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan course")
			utils.LogFooter()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan course"})
		}
		course := make(map[string]interface{})
		for i, col := range cols {
			course[col] = vals[i]
		}
		courses = append(courses, course)
	}
	utils.LogMessage(utils.LevelInfo, "Successfully fetched today's courses")
	utils.LogFooter()
	return c.JSON(courses)
}

// POST /planning/courses
// @Summary		Créer un nouveau cours
// @Description	Ajoute un nouveau cours dans le planning
// @Tags			Planning
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			course	body	models.CreateCourseRequest	true	"Données du cours"
// @Success		201		{object}	models.Response			"Cours créé avec succès"
// @Failure		400		{object}	models.ErrorResponse	"Données invalides"
// @Failure		401		{object}	models.ErrorResponse	"Non autorisé"
// @Failure		500		{object}	models.ErrorResponse	"Erreur serveur"
// @Router			/planning/courses [post]
func (h *PlanningHandler) CreateCourse(c *fiber.Ctx) error {
	utils.LogHeader("📚 Create Course")

	var course models.CreateCourseRequest
	if err := c.BodyParser(&course); err != nil {
		utils.LogMessage(utils.LevelError, "Invalid request body")
		utils.LogFooter()
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Course Data", course)

	_, err := h.db.Exec(`INSERT INTO courses (date, title, start_time, end_time, teacher, room, "group", user_email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		course.Date, course.Title, course.StartTime, course.EndTime, course.Teacher, course.Room, course.Group, course.UserEmail)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to create course")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create course", "details": err.Error()})
	}
	utils.LogMessage(utils.LevelInfo, "Successfully created course")
	utils.LogFooter()
	return c.Status(http.StatusCreated).JSON(fiber.Map{"success": true})
}
