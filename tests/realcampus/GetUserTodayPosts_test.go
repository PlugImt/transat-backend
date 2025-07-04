package realcampus

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"github.com/plugimt/transat-backend/realcampus/posts"
	"github.com/plugimt/transat-backend/utils"
)

// TestGetUserTodayPosts_Success tests a successful response from GetUserTodayPosts.
func TestGetUserTodayPosts_Success(t *testing.T) {
	// Create a sqlmock database connection and a mock instance.
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Calculate today's start and end time in Paris timezone.
	startOfDay, endOfDay := utils.TodayInParis()

	// Prepare an expected query (using regexp for flexible whitespace matching).
	expectedQuery := regexp.QuoteMeta(`
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
			ORDER BY p.creation_date DESC`)

	// Set up the expected row(s).
	columns := []string{"id_post", "id_file_1", "id_file_2", "location", "privacy", "creation_date", "path1", "name1", "path2", "name2"}
	sampleTime := utils.Now()
	mockRows := sqlmock.NewRows(columns).AddRow(
		1,                 // id_post
		101,               // id_file_1
		102,               // id_file_2
		"TestPlace",       // location
		"public",          // privacy
		sampleTime,        // creation_date
		"/data/test1.jpg", // path1
		"Test File 1",     // name1
		"/data/test2.jpg", // path2
		"Test File 2",     // name2
	)

	// Expect the query with the right arguments.
	mock.ExpectQuery(expectedQuery).
		WithArgs("test@example.com", startOfDay, endOfDay).
		WillReturnRows(mockRows)

	// Create a new Fiber app.
	app := fiber.New()

	// Add a simple middleware to set the "email" in locals.
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("email", "test@example.com")
		return c.Next()
	})

	// Register the route using the handler.
	app.Get("/posts/today", posts.GetUserTodayPosts(db))

	// Create a dummy HTTP request to the route.
	req := httptest.NewRequest("GET", "/posts/today", nil)
	// Execute the request using app.Test() which handles the conversion to fasthttp internally.
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Read the response body.
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	// Parse and validate the JSON response.
	var response struct {
		Success bool `json:"success"`
		Count   int  `json:"count"`
		Posts   []struct {
			ID       interface{} `json:"id"`
			Location string      `json:"location"`
			Privacy  string      `json:"privacy"`
			Created  string      `json:"created"`
			File1    struct {
				ID   interface{} `json:"id"`
				Name string      `json:"name"`
				URL  string      `json:"url"`
			} `json:"file1"`
			File2 struct {
				ID   interface{} `json:"id"`
				Name string      `json:"name"`
				URL  string      `json:"url"`
			} `json:"file2"`
		} `json:"posts"`
	}

	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	// Validate that the response indicates success and contains one post.
	require.True(t, response.Success)
	require.Equal(t, 1, response.Count)
	require.Len(t, response.Posts, 1)
	require.Contains(t, response.Posts[0].File1.URL, "/api/data/")
	require.Contains(t, response.Posts[0].File2.URL, "/api/data/")

	// Ensure all expectations were met.
	require.NoError(t, mock.ExpectationsWereMet())
}

// TestGetUserTodayPosts_DBError tests the response when the database query fails.
func TestGetUserTodayPosts_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	startOfDay, endOfDay := utils.TodayInParis()

	expectedQuery := regexp.QuoteMeta(`
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
			ORDER BY p.creation_date DESC`)

	// Simulate a query error.
	mock.ExpectQuery(expectedQuery).
		WithArgs("test@example.com", startOfDay, endOfDay).
		WillReturnError(sql.ErrConnDone)

	app := fiber.New()

	// Set local "email" using middleware.
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("email", "test@example.com")
		return c.Next()
	})
	app.Get("/posts/today", posts.GetUserTodayPosts(db))

	req := httptest.NewRequest("GET", "/posts/today", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// In a failure case, the handler should return an error response (or set an error status).
	require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
