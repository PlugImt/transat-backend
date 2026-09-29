package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/plugimt/transat-backend/utils"
)

// StatisticsService handles logging of request statistics
type StatisticsService struct {
	db *sql.DB
}

// NewStatisticsService creates a new instance of StatisticsService
func NewStatisticsService(db *sql.DB) *StatisticsService {
	return &StatisticsService{
		db: db,
	}
}

// EndpointStatistic represents aggregated statistics for an endpoint
type EndpointStatistic struct {
	Endpoint           string    `json:"endpoint"`
	Method             string    `json:"method"`
	RequestCount       int       `json:"request_count"`
	AvgDurationMs      float64   `json:"avg_duration_ms"`
	MinDurationMs      int       `json:"min_duration_ms"`
	MaxDurationMs      int       `json:"max_duration_ms"`
	SuccessRatePercent float64   `json:"success_rate_percent"`
	FirstRequest       time.Time `json:"first_request"`
	LastRequest        time.Time `json:"last_request"`
	SuccessCount       int       `json:"success_count"`
	ErrorCount         int       `json:"error_count"`
}

// UserStatistic represents aggregated statistics for a user
type UserStatistic struct {
	Email              string    `json:"email"`
	RequestCount       int       `json:"request_count"`
	AvgDurationMs      float64   `json:"avg_duration_ms"`
	MinDurationMs      int       `json:"min_duration_ms"`
	MaxDurationMs      int       `json:"max_duration_ms"`
	SuccessRatePercent float64   `json:"success_rate_percent"`
	FirstRequest       time.Time `json:"first_request"`
	LastRequest        time.Time `json:"last_request"`
	SuccessCount       int       `json:"success_count"`
	ErrorCount         int       `json:"error_count"`
}

// GlobalStatistic represents aggregated statistics across all endpoints
type GlobalStatistic struct {
	TotalRequestCount        int       `json:"total_request_count"`
	GlobalAvgDurationMs      float64   `json:"global_avg_duration_ms"`
	GlobalMinDurationMs      int       `json:"global_min_duration_ms"`
	GlobalMaxDurationMs      int       `json:"global_max_duration_ms"`
	GlobalSuccessRatePercent float64   `json:"global_success_rate_percent"`
	FirstRequest             time.Time `json:"first_request"`
	LastRequest              time.Time `json:"last_request"`
	SuccessCount             int       `json:"success_count"`
	ErrorCount               int       `json:"error_count"`
}

// GetEndpointStatistics retrieves all endpoint statistics from the view
func (s *StatisticsService) GetEndpointStatistics() ([]EndpointStatistic, error) {
	rows, err := s.db.Query(`SELECT 
		endpoint, method, request_count, avg_duration_ms, 
		min_duration_ms, max_duration_ms, success_rate_percent,
		first_request, last_request, success_count, error_count
		FROM endpoint_statistics`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []EndpointStatistic
	for rows.Next() {
		var stat EndpointStatistic
		if err := rows.Scan(
			&stat.Endpoint, &stat.Method, &stat.RequestCount, &stat.AvgDurationMs,
			&stat.MinDurationMs, &stat.MaxDurationMs, &stat.SuccessRatePercent,
			&stat.FirstRequest, &stat.LastRequest, &stat.SuccessCount, &stat.ErrorCount,
		); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

// GetTopUserStatistics retrieves the top 10 users by request count
func (s *StatisticsService) GetTopUserStatistics() ([]UserStatistic, error) {
	rows, err := s.db.Query(`SELECT 
		email, request_count, avg_duration_ms, 
		min_duration_ms, max_duration_ms, success_rate_percent,
		first_request, last_request, success_count, error_count
		FROM top_users_statistics`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []UserStatistic
	for rows.Next() {
		var stat UserStatistic
		if err := rows.Scan(
			&stat.Email, &stat.RequestCount, &stat.AvgDurationMs,
			&stat.MinDurationMs, &stat.MaxDurationMs, &stat.SuccessRatePercent,
			&stat.FirstRequest, &stat.LastRequest, &stat.SuccessCount, &stat.ErrorCount,
		); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

// GetGlobalStatistics retrieves global statistics from the view
func (s *StatisticsService) GetGlobalStatistics() (*GlobalStatistic, error) {
	var stat GlobalStatistic

	err := s.db.QueryRow(`SELECT 
		total_request_count, global_avg_duration_ms, 
		global_min_duration_ms, global_max_duration_ms, 
		global_success_rate_percent, first_request, last_request,
		success_count, error_count
		FROM global_statistics`).Scan(
		&stat.TotalRequestCount, &stat.GlobalAvgDurationMs,
		&stat.GlobalMinDurationMs, &stat.GlobalMaxDurationMs,
		&stat.GlobalSuccessRatePercent, &stat.FirstRequest, &stat.LastRequest,
		&stat.SuccessCount, &stat.ErrorCount,
	)
	if err != nil {
		return nil, err
	}

	return &stat, nil
}

// ActiveUsersPoint represents the number of distinct active users in a time bucket
type ActiveUsersPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// activeUsersPeriodConfig describes how a bucket period maps to SQL truncation/step/lookback
type activeUsersPeriodConfig struct {
	unit     string // date_trunc unit
	step     string // generate_series step interval
	lookback string // how far back to look from now
}

// ActiveUsersPeriods lists the allowed granularities for GetActiveUsersOverTime
var ActiveUsersPeriods = map[string]activeUsersPeriodConfig{
	"day":     {unit: "day", step: "1 day", lookback: "90 days"},
	"week":    {unit: "week", step: "1 week", lookback: "364 days"},
	"month":   {unit: "month", step: "1 month", lookback: "24 months"},
	"quarter": {unit: "quarter", step: "3 months", lookback: "24 months"},
	"year":    {unit: "year", step: "1 year", lookback: "5 years"},
}

// GetActiveUsersOverTime retrieves the number of distinct active users bucketed by the given period
func (s *StatisticsService) GetActiveUsersOverTime(period string) ([]ActiveUsersPoint, error) {
	cfg, ok := ActiveUsersPeriods[period]
	if !ok {
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	// cfg fields come exclusively from the ActiveUsersPeriods whitelist above, never from raw user input
	query := fmt.Sprintf(`
		SELECT d.bucket::date, COUNT(DISTINCT r.user_email)
		FROM generate_series(
			date_trunc('%s', NOW() - INTERVAL '%s'),
			date_trunc('%s', NOW()),
			INTERVAL '%s'
		) AS d(bucket)
		LEFT JOIN request_statistics r
			ON date_trunc('%s', r.request_received) = d.bucket
			AND r.user_email IS NOT NULL
		GROUP BY d.bucket
		ORDER BY d.bucket ASC`,
		cfg.unit, cfg.lookback, cfg.unit, cfg.step, cfg.unit,
	)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []ActiveUsersPoint
	for rows.Next() {
		var bucket time.Time
		var count int
		if err := rows.Scan(&bucket, &count); err != nil {
			return nil, err
		}
		points = append(points, ActiveUsersPoint{
			Date:  bucket.Format("2006-01-02"),
			Count: count,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return points, nil
}

// ActivityHourPoint represents the number of distinct active users for a given weekday/hour slot
type ActivityHourPoint struct {
	DayOfWeek int `json:"dayOfWeek"` // ISO day of week: 1 = Monday ... 7 = Sunday
	Hour      int `json:"hour"`      // 0-23
	Count     int `json:"count"`
}

// ActivityHeatmapRanges lists the allowed lookback windows for GetActivityHeatmap, mapped to a
// SQL expression evaluating to the earliest timestamp to include. Values are hardcoded, never user input.
var ActivityHeatmapRanges = map[string]string{
	"week":  "date_trunc('week', NOW())",
	"month": "date_trunc('month', NOW())",
	"year":  "date_trunc('year', NOW())",
	"all":   "'-infinity'::timestamp",
}

// GetActivityHeatmap retrieves the number of distinct active users for each (day of week, hour)
// slot within the given range, to visualize when users are typically active during the week
func (s *StatisticsService) GetActivityHeatmap(rangeParam string) ([]ActivityHourPoint, error) {
	sinceExpr, ok := ActivityHeatmapRanges[rangeParam]
	if !ok {
		return nil, fmt.Errorf("invalid range: %s", rangeParam)
	}

	query := fmt.Sprintf(`
		SELECT d.dow, h.hour, COUNT(DISTINCT r.user_email)
		FROM generate_series(1, 7) AS d(dow)
		CROSS JOIN generate_series(0, 23) AS h(hour)
		LEFT JOIN request_statistics r
			ON EXTRACT(ISODOW FROM r.request_received)::int = d.dow
			AND EXTRACT(HOUR FROM r.request_received)::int = h.hour
			AND r.user_email IS NOT NULL
			AND r.request_received >= %s
		GROUP BY d.dow, h.hour
		ORDER BY d.dow, h.hour`,
		sinceExpr,
	)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []ActivityHourPoint
	for rows.Next() {
		var point ActivityHourPoint
		if err := rows.Scan(&point.DayOfWeek, &point.Hour, &point.Count); err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return points, nil
}

// LogRequest records a request statistics entry in the database
func (s *StatisticsService) LogRequest(
	userEmail string,
	endpoint string,
	method string,
	requestReceived time.Time,
	responseSent time.Time,
	statusCode int,
) {
	// Calculate duration in milliseconds
	durationMs := int(responseSent.Sub(requestReceived).Milliseconds())

	// Add more detailed logging for troubleshooting
	log.Printf("Logging request: %s %s - Status: %d", method, endpoint, statusCode)

	utils.LogHeader("📊 Statistics Service")
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)
	utils.LogLineKeyValue(utils.LevelInfo, "Endpoint", endpoint)
	utils.LogLineKeyValue(utils.LevelInfo, "Method", method)
	utils.LogLineKeyValue(utils.LevelInfo, "Status", statusCode)
	utils.LogLineKeyValue(utils.LevelInfo, "Duration", durationMs)

	// Verify the status code is appropriate
	if statusCode < 100 || statusCode > 599 {
		log.Printf("WARNING: Invalid status code %d for %s %s - correcting to 500",
			statusCode, method, endpoint)
		statusCode = 500 // Default to 500 if the status code is invalid
	}

	// Set userEmail to nil for empty strings to handle the foreign key constraint properly
	var userEmailOrNil interface{}
	if userEmail == "" {
		userEmailOrNil = nil
	} else {
		userEmailOrNil = userEmail
	}

	// Insert into database
	_, err := s.db.Exec(
		`INSERT INTO request_statistics 
		(user_email, endpoint, method, request_received, response_sent, status_code, duration_ms) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userEmailOrNil, endpoint, method, requestReceived, responseSent, statusCode, durationMs,
	)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to log request statistics")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	} else {
		utils.LogMessage(utils.LevelInfo, "Successfully logged request statistics")
	}
	utils.LogFooter()
}
