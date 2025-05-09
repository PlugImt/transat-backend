package utils

import (
	"os"
)

func GetEnvName() string {
	env := os.Getenv("RAILWAY_ENVIRONMENT_NAME")
	if env == "" {
		env = "development"
	}
	return env
}

func GetEnvCommitSHA() string {
	commitSHA := os.Getenv("RAILWAY_GIT_COMMIT_SHA")
	if commitSHA == "" {
		commitSHA = "unknown"
	}
	return commitSHA
}

func GetEnvHost() string {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
		LogMessage(LevelError, "Error getting hostname: "+err.Error())
	}
	return host
}
