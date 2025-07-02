package planning

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/plugimt/transat-backend/utils"
)

type PlanningHandler struct {
	db *sql.DB
}

func NewPlanningHandler(db *sql.DB) *PlanningHandler {
	return &PlanningHandler{db: db}
}

// GET /planning/users
func (h *PlanningHandler) GetUsersWithPassID(c *fiber.Ctx) error {
	utils.LogHeader("👥 Get Users With Pass ID")

	rows, err := h.db.Query(`SELECT id_newf, first_name, last_name, pass_id, email FROM newf`)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch users")
		utils.LogFooter()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var firstName, lastName, email sql.NullString
		var passID sql.NullInt64
		if err := rows.Scan(&id, &firstName, &lastName, &passID, &email); err != nil {
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
			"email":      email.String,
			"pass_id":    passIDValue,
		})
	}
	utils.LogMessage(utils.LevelInfo, "Successfully fetched users")
	utils.LogFooter()
	return c.JSON(users)
}

// PATCH /planning/users/:id_newf/passid
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

	type reqBody struct {
		PassID int `json:"pass_id"`
	}
	var body reqBody
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
			if t, ok := vals[i].(time.Time); ok {
				course[col] = utils.FormatParis(t, time.RFC3339)
			} else {
				course[col] = vals[i]
			}
		}
		courses = append(courses, course)
	}
	utils.LogMessage(utils.LevelInfo, "Successfully fetched courses")
	utils.LogFooter()
	return c.JSON(courses)
}

// GET /planning/users/:email/courses/today
func (h *PlanningHandler) GetUserCoursesToday(c *fiber.Ctx) error {
	utils.LogHeader("📅 Get User Courses Today")

	email := c.Params("email")
	if email == "" {
		utils.LogMessage(utils.LevelError, "Missing email parameter")
		utils.LogFooter()
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing email parameter"})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "User Email", email)

	today := utils.FormatParis(utils.Now(), "2006-01-02")
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
			if t, ok := vals[i].(time.Time); ok {
				course[col] = utils.FormatParis(t, time.RFC3339)
			} else {
				course[col] = vals[i]
			}
		}
		courses = append(courses, course)
	}
	utils.LogMessage(utils.LevelInfo, "Successfully fetched today's courses")
	utils.LogFooter()
	return c.JSON(courses)
}

// POST /planning/courses
func (h *PlanningHandler) CreateCourse(c *fiber.Ctx) error {
	utils.LogHeader("📚 Process Course (Create/Replace)")

	type Course struct {
		Date      string `json:"date"`
		Title     string `json:"title"`
		StartTime string `json:"start_time"` // Assumes ISO 8601 format: "2023-10-27T10:00:00Z"
		EndTime   string `json:"end_time"`
		Teacher   string `json:"teacher"`
		Room      string `json:"room"`
		Group     string `json:"group"`
		UserEmail string `json:"user_email"`
	}
	var course Course
	if err := c.BodyParser(&course); err != nil {
		utils.LogMessage(utils.LevelError, "Invalid request body")
		utils.LogFooter()
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Processing course for", course.UserEmail)
	utils.LogLineKeyValue(utils.LevelInfo, "Course Data", course)

	// Begin a new database transaction. Ensure atomicity.
	tx, err := h.db.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction")
		utils.LogFooter()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Database transaction could not be started"})
	}

	// Defer a rollback. If anything goes wrong, the transaction will be cancelled.
	// If tx.Commit() succeeds, this defer does nothing.
	defer tx.Rollback()

	// --- Step 1: Delete any existing courses that conflict with the new time slot ---
	// The OVERLAPS operator is perfect for this. It checks if two time periods intersect.
	// We also ensure we don't delete the exact same course if its title is the same,
	// as that's an update, not a replacement of a different course. We will handle
	// its "update" by deleting and re-inserting it.
	deleteQuery := `
		DELETE FROM courses
		WHERE user_email = $1
		  AND (start_time, end_time) OVERLAPS ($2, $3)
	`
	res, err := tx.Exec(deleteQuery, course.UserEmail, course.StartTime, course.EndTime)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete conflicting courses during transaction")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to clear conflicting time slot"})
	}

	deletedCount, _ := res.RowsAffected()
	if deletedCount > 0 {
		utils.LogLineKeyValue(utils.LevelInfo, "Deleted conflicting courses", deletedCount)
	}

	// --- Step 2: Insert the new course into the now-clear time slot ---
	// We use RETURNING id to get the ID of the newly created course.
	insertQuery := `
		INSERT INTO courses (date, title, start_time, end_time, teacher, room, "group", user_email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	var newCourseID int
	err = tx.QueryRow(
		insertQuery,
		course.Date, course.Title, course.StartTime, course.EndTime,
		course.Teacher, course.Room, course.Group, course.UserEmail,
	).Scan(&newCourseID)

	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to insert new course during transaction")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert new course"})
	}

	// --- Step 3: Commit the transaction ---
	// If both the delete and insert were successful, we commit the changes to the database.
	if err := tx.Commit(); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction")
		utils.LogFooter()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to finalize course creation"})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Successfully created/replaced course with ID", newCourseID)
	utils.LogFooter()
	// Using 201 Created is appropriate because a new resource representation has been created.
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"status":  "created_or_replaced",
		"id":      newCourseID,
	})
}
