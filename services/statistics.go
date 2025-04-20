package services

import (
	"database/sql"
	"log"
	"time"

	"Transat_2.0_Backend/utils"
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
	Endpoint           string     `json:"endpoint"`
	Method             string     `json:"method"`
	RequestCount       int        `json:"request_count"`
	AvgDurationMs      float64    `json:"avg_duration_ms"`
	MinDurationMs      int        `json:"min_duration_ms"`
	MaxDurationMs      int        `json:"max_duration_ms"`
	SuccessRatePercent float64    `json:"success_rate_percent"`
	FirstRequest       time.Time  `json:"first_request"`
	LastRequest        time.Time  `json:"last_request"`
	SuccessCount       int        `json:"success_count"`
	ErrorCount         int        `json:"error_count"`
}

// GlobalStatistic represents aggregated statistics across all endpoints
type GlobalStatistic struct {
	TotalRequestCount      int       `json:"total_request_count"`
	GlobalAvgDurationMs    float64   `json:"global_avg_duration_ms"`
	GlobalMinDurationMs    int       `json:"global_min_duration_ms"`
	GlobalMaxDurationMs    int       `json:"global_max_duration_ms"`
	GlobalSuccessRatePercent float64 `json:"global_success_rate_percent"`
	FirstRequest           time.Time `json:"first_request"`
	LastRequest            time.Time `json:"last_request"`
	SuccessCount           int       `json:"success_count"`
	ErrorCount             int       `json:"error_count"`
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