package admin

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

// AdminHandler handles admin-specific operations
type AdminHandler struct {
	db                *sql.DB
	statisticsService *services.StatisticsService
}

// NewAdminHandler creates a new instance of AdminHandler
func NewAdminHandler(db *sql.DB, statisticsService *services.StatisticsService) *AdminHandler {
	return &AdminHandler{
		db:                db,
		statisticsService: statisticsService,
	}
}

// DashboardStats represents the admin dashboard statistics
type DashboardStats struct {
	TotalUsers          int                `json:"total_users"`
	ActiveUsers         int                `json:"active_users"`
	NewUsersToday       int                `json:"new_users_today"`
	NewUsersThisWeek    int                `json:"new_users_this_week"`
	TotalRequests       int64              `json:"total_requests"`
	RequestsToday       int64              `json:"requests_today"`
	AverageResponseTime float64            `json:"average_response_time"`
	TopEndpoints        []EndpointStat     `json:"top_endpoints"`
	UserGrowth          []UserGrowthStat   `json:"user_growth"`
	ServiceUsage        []ServiceUsageStat `json:"service_usage"`
}

type EndpointStat struct {
	Endpoint     string  `json:"endpoint"`
	RequestCount int64   `json:"request_count"`
	AvgDuration  float64 `json:"avg_duration"`
}

type UserGrowthStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ServiceUsageStat struct {
	Service string `json:"service"`
	Count   int64  `json:"count"`
}

// GetDashboardStats returns comprehensive dashboard statistics
func (h *AdminHandler) GetDashboardStats(c *fiber.Ctx) error {
	utils.LogHeader("📊 Admin Dashboard Stats")

	var stats DashboardStats

	// Get total users
	err := h.db.QueryRow("SELECT COUNT(*) FROM newf").Scan(&stats.TotalUsers)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get total users count")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user statistics",
		})
	}

	// Get active users (users who logged in within last 30 days)
	err = h.db.QueryRow(`
		SELECT COUNT(DISTINCT email) 
		FROM statistics 
		WHERE endpoint = '/api/auth/login' 
		AND timestamp > NOW() - INTERVAL '30 days'
	`).Scan(&stats.ActiveUsers)
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Failed to get active users count")
		stats.ActiveUsers = 0
	}

	// Get new users today
	err = h.db.QueryRow(`
		SELECT COUNT(*) 
		FROM newf 
		WHERE DATE(creation_date) = CURRENT_DATE
	`).Scan(&stats.NewUsersToday)
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Failed to get new users today")
		stats.NewUsersToday = 0
	}

	// Get new users this week
	err = h.db.QueryRow(`
		SELECT COUNT(*) 
		FROM newf 
		WHERE creation_date >= DATE_TRUNC('week', CURRENT_DATE)
	`).Scan(&stats.NewUsersThisWeek)
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Failed to get new users this week")
		stats.NewUsersThisWeek = 0
	}

	// Get statistics from statistics service
	globalStats, err := h.statisticsService.GetGlobalStatistics()
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Failed to get global statistics")
		stats.TotalRequests = 0
		stats.AverageResponseTime = 0
	} else {
		stats.TotalRequests = int64(globalStats.TotalRequestCount)
		stats.AverageResponseTime = globalStats.GlobalAvgDurationMs
	}

	// Get requests today
	err = h.db.QueryRow(`
		SELECT COUNT(*) 
		FROM statistics 
		WHERE DATE(timestamp) = CURRENT_DATE
	`).Scan(&stats.RequestsToday)
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Failed to get requests today")
		stats.RequestsToday = 0
	}

	// Get top endpoints
	endpointRows, err := h.db.Query(`
		SELECT endpoint, COUNT(*) as request_count, AVG(duration_ms) as avg_duration
		FROM statistics 
		WHERE timestamp > NOW() - INTERVAL '7 days'
		GROUP BY endpoint 
		ORDER BY request_count DESC 
		LIMIT 10
	`)
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Failed to get top endpoints")
		stats.TopEndpoints = []EndpointStat{}
	} else {
		defer endpointRows.Close()
		for endpointRows.Next() {
			var endpoint EndpointStat
			err := endpointRows.Scan(&endpoint.Endpoint, &endpoint.RequestCount, &endpoint.AvgDuration)
			if err != nil {
				continue
			}
			stats.TopEndpoints = append(stats.TopEndpoints, endpoint)
		}
	}

	// Get user growth over last 30 days
	growthRows, err := h.db.Query(`
		SELECT DATE(creation_date) as date, COUNT(*) as count
		FROM newf 
		WHERE creation_date > NOW() - INTERVAL '30 days'
		GROUP BY DATE(creation_date)
		ORDER BY date DESC
	`)
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Failed to get user growth data")
		stats.UserGrowth = []UserGrowthStat{}
	} else {
		defer growthRows.Close()
		for growthRows.Next() {
			var growth UserGrowthStat
			var date time.Time
			err := growthRows.Scan(&date, &growth.Count)
			if err != nil {
				continue
			}
			growth.Date = date.Format("2006-01-02")
			stats.UserGrowth = append(stats.UserGrowth, growth)
		}
	}

	// Get service usage statistics
	serviceRows, err := h.db.Query(`
		SELECT 
			CASE 
				WHEN endpoint LIKE '%restaurant%' THEN 'Restaurant'
				WHEN endpoint LIKE '%laundry%' OR endpoint LIKE '%washing%' THEN 'Laundry'
				WHEN endpoint LIKE '%statistics%' THEN 'Statistics'
				WHEN endpoint LIKE '%auth%' THEN 'Authentication'
				WHEN endpoint LIKE '%user%' THEN 'User Management'
				WHEN endpoint LIKE '%traq%' THEN 'Traq'
				WHEN endpoint LIKE '%realcampus%' THEN 'RealCampus'
				ELSE 'Other'
			END as service,
			COUNT(*) as count
		FROM statistics 
		WHERE timestamp > NOW() - INTERVAL '7 days'
		GROUP BY service
		ORDER BY count DESC
	`)
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Failed to get service usage data")
		stats.ServiceUsage = []ServiceUsageStat{}
	} else {
		defer serviceRows.Close()
		for serviceRows.Next() {
			var service ServiceUsageStat
			err := serviceRows.Scan(&service.Service, &service.Count)
			if err != nil {
				continue
			}
			stats.ServiceUsage = append(stats.ServiceUsage, service)
		}
	}

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved admin dashboard statistics")
	utils.LogLineKeyValue(utils.LevelInfo, "Total Users", stats.TotalUsers)
	utils.LogLineKeyValue(utils.LevelInfo, "Active Users", stats.ActiveUsers)
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// UserInfo represents user information for admin view
type UserInfo struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Role           string `json:"role"`
	CreationDate   string `json:"creation_date"`
	LastLoginDate  string `json:"last_login_date"`
	Campus         string `json:"campus"`
	GraduationYear int    `json:"graduation_year"`
	Language       string `json:"language"`
}

// GetUsers returns a list of all users with their information
func (h *AdminHandler) GetUsers(c *fiber.Ctx) error {
	utils.LogHeader("👥 Admin Get Users")

	// Parse query parameters for pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	// Build query with optional search
	baseQuery := `
		SELECT n.id_newf, n.email, n.first_name, n.last_name, 
			   COALESCE(r.name, 'UNKNOWN') as role, n.creation_date,
			   n.campus, n.graduation_year, COALESCE(l.code, 'fr') as language
		FROM newf n
		LEFT JOIN newf_roles nr ON n.email = nr.email
		LEFT JOIN roles r ON nr.id_roles = r.id_roles
		LEFT JOIN languages l ON n.language = l.id_languages
	`

	var args []interface{}
	argIndex := 1

	if search != "" {
		baseQuery += ` WHERE (n.email ILIKE $` + string(rune(argIndex)) + ` OR n.first_name ILIKE $` + string(rune(argIndex)) + ` OR n.last_name ILIKE $` + string(rune(argIndex)) + `)`
		args = append(args, "%"+search+"%")
		argIndex++
	}

	baseQuery += ` ORDER BY n.creation_date DESC LIMIT $` + string(rune(argIndex)) + ` OFFSET $` + string(rune(argIndex+1))
	args = append(args, limit, offset)

	rows, err := h.db.Query(baseQuery, args...)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get users")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}
	defer rows.Close()

	var users []UserInfo
	for rows.Next() {
		var user UserInfo
		var creationDate time.Time
		err := rows.Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName,
			&user.Role, &creationDate, &user.Campus, &user.GraduationYear,
			&user.Language,
		)
		if err != nil {
			continue
		}
		user.CreationDate = creationDate.Format("2006-01-02 15:04:05")
		users = append(users, user)
	}

	// Get total count for pagination
	countQuery := `SELECT COUNT(*) FROM newf n`
	if search != "" {
		countQuery += ` WHERE (n.email ILIKE $1 OR n.first_name ILIKE $1 OR n.last_name ILIKE $1)`
	}

	var totalCount int
	if search != "" {
		err = h.db.QueryRow(countQuery, "%"+search+"%").Scan(&totalCount)
	} else {
		err = h.db.QueryRow(countQuery).Scan(&totalCount)
	}
	if err != nil {
		totalCount = len(users)
	}

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved users")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(users))
	utils.LogLineKeyValue(utils.LevelInfo, "Page", page)
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"users":       users,
			"total_count": totalCount,
			"page":        page,
			"limit":       limit,
			"total_pages": (totalCount + limit - 1) / limit,
		},
	})
}

// GetCurrentUser returns information about the currently logged-in user
func (h *AdminHandler) GetCurrentUser(c *fiber.Ctx) error {
	utils.LogHeader("👤 Get Current User")

	// Get user email from context (set by middleware)
	email, ok := c.Locals("email").(string)
	if !ok {
		utils.LogMessage(utils.LevelError, "User email not found in context")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Get user roles
	roles, err := middlewares.GetUserRoles(h.db, email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get user roles")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user information",
		})
	}

	// Get user details
	var user UserInfo
	var creationDate time.Time
	query := `
		SELECT n.id_newf, n.email, n.first_name, n.last_name, 
			   n.creation_date, n.campus, n.graduation_year, 
			   COALESCE(l.code, 'fr') as language
		FROM newf n
		LEFT JOIN languages l ON n.language = l.id_languages
		WHERE n.email = $1
	`

	err = h.db.QueryRow(query, strings.ToLower(email)).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&creationDate, &user.Campus, &user.GraduationYear, &user.Language,
	)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get user details")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user details",
		})
	}

	user.CreationDate = creationDate.Format("2006-01-02 15:04:05")
	user.Role = strings.Join(roles, ", ") // Join all roles

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved current user")
	utils.LogLineKeyValue(utils.LevelInfo, "Email", email)
	utils.LogLineKeyValue(utils.LevelInfo, "Roles", user.Role)
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user":  user,
			"roles": roles,
		},
	})
}
