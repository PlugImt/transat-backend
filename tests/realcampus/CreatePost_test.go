package realcampus

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plugimt/transat-backend/realcampus/posts"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupApp(db *sql.DB) *fiber.App {
	app := fiber.New()
	app.Post("/posts", func(c *fiber.Ctx) error {
		c.Locals("email", "user@example.com")
		return posts.CreatePost(db)(c)
	})
	return app
}

func TestCreatePost_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT COUNT\\(\\*\\)").
		WithArgs(1, 2, "user@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery("INSERT INTO realcampus_posts").
		WithArgs(1, 2, "user@example.com", "Paris", "PUBLIC").
		WillReturnRows(sqlmock.NewRows([]string{"id_post"}).AddRow(42))

	app := setupApp(db)

	reqBody := map[string]interface{}{
		"file_id_1": 1,
		"file_id_2": 2,
		"location":  "Paris",
		"privacy":   "PUBLIC",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCreatePost_InvalidPrivacy(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT COUNT\\(\\*\\)").
		WithArgs(1, 2, "user@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery("INSERT INTO realcampus_posts").
		WithArgs(1, 2, "user@example.com", "Paris", "PRIVATE"). // Defaults to PRIVATE.
		WillReturnRows(sqlmock.NewRows([]string{"id_post"}).AddRow(100))

	app := setupApp(db)

	reqBody := map[string]interface{}{
		"file_id_1": 1,
		"file_id_2": 2,
		"location":  "Paris",
		"privacy":   "FRIENDS", // Invalid value.
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCreatePost_UnauthorizedFile(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT COUNT\\(\\*\\)").
		WithArgs(1, 2, "user@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1)) // Only 1 file owned.

	app := setupApp(db)

	reqBody := map[string]interface{}{
		"file_id_1": 1,
		"file_id_2": 2,
		"location":  "Paris",
		"privacy":   "PUBLIC",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}
