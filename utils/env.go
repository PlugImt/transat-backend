package utils

import (
	"os"
)

// GetEnvOrDefault retrieves an environment variable or returns a default value.
func GetEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetEnvName() string {
	// Uses RAILWAY_ENVIRONMENT_NAME first, then falls back to "development"
	return GetEnvOrDefault("RAILWAY_ENVIRONMENT_NAME", "development")
}

func GetEnvCommitSHA() string {
	// Uses RAILWAY_GIT_COMMIT_SHA first, then falls back to "unknown"
	return GetEnvOrDefault("RAILWAY_GIT_COMMIT_SHA", "unknown")
}

func GetEnvHost() string {
	host, err := os.Hostname()
	if err != nil {
		LogMessage(LevelError, "Error getting hostname: "+err.Error())
		return "unknown"
	}
	return host
}
