package services

import (
	"database/sql"
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

	utils.LogHeader("📊 Statistics Service")
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)
	utils.LogLineKeyValue(utils.LevelInfo, "Endpoint", endpoint)
	utils.LogLineKeyValue(utils.LevelInfo, "Method", method)
	utils.LogLineKeyValue(utils.LevelInfo, "Status", statusCode)
	utils.LogLineKeyValue(utils.LevelInfo, "Duration", durationMs)

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