package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

var configuredDataFolder string

// InitFilePaths initializes the data folder path and ensures it exists.
func InitFilePaths(dataFolderFromConfig string) error {
	path := dataFolderFromConfig
	if path == "" {
		path = "./data" // Default to a relative path
		LogMessage(LevelInfo, "Data folder not specified in config, using default './data'")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		LogMessage(LevelError, fmt.Sprintf("Error converting data folder path '%s' to absolute: %v", path, err))
		return fmt.Errorf("error converting data folder path to absolute: %w", err)
	}
	configuredDataFolder = absPath
	LogMessage(LevelInfo, fmt.Sprintf("Data folder path configured to: %s", configuredDataFolder))

	// Ensure the directory exists
	if _, err := os.Stat(configuredDataFolder); os.IsNotExist(err) {
		LogMessage(LevelInfo, fmt.Sprintf("Data folder '%s' does not exist. Creating...", configuredDataFolder))
		if errMkdir := os.MkdirAll(configuredDataFolder, 0750); errMkdir != nil {
			LogMessage(LevelError, fmt.Sprintf("Failed to create data folder '%s': %v", configuredDataFolder, errMkdir))
			return fmt.Errorf("failed to create data folder '%s': %w", configuredDataFolder, errMkdir)
		}
		LogMessage(LevelInfo, fmt.Sprintf("Data folder '%s' created successfully.", configuredDataFolder))
	} else if err != nil {
		LogMessage(LevelError, fmt.Sprintf("Error checking data folder '%s': %v", configuredDataFolder, err))
		return fmt.Errorf("error checking data folder '%s': %w", configuredDataFolder, err)
	} else {
		LogMessage(LevelInfo, fmt.Sprintf("Data folder '%s' already exists.", configuredDataFolder))
	}
	return nil
}

// EnsureDataFolder checks if the configured data folder exists.
// It is recommended to call InitFilePaths at startup instead of this function directly.
// If InitFilePaths was not called, this function attempts a fallback initialization.
func EnsureDataFolder() error {
	if configuredDataFolder == "" {
		LogMessage(LevelWarn, "EnsureDataFolder called before InitFilePaths. Attempting default initialization for safety.")
		if err := InitFilePaths("./data"); err != nil { // Use default path for fallback
			return fmt.Errorf("fallback InitFilePaths failed in EnsureDataFolder: %w", err)
		}
	} else {
		// If already configured, just verify existence as an additional check, though InitFilePaths should have handled it.
		if _, err := os.Stat(configuredDataFolder); os.IsNotExist(err) {
			LogMessage(LevelError, fmt.Sprintf("Configured data folder '%s' does not exist. This should have been created by InitFilePaths.", configuredDataFolder))
			return fmt.Errorf("configured data folder '%s' not found after InitFilePaths was expected to run", configuredDataFolder)
		} else if err != nil {
			LogMessage(LevelError, fmt.Sprintf("Error re-checking data folder '%s': %v", configuredDataFolder, err))
			return fmt.Errorf("error re-checking data folder: %w", err)
		}
	}
	LogMessage(LevelDebug, "EnsureDataFolder check completed.")
	return nil
}

// GetDataFolderPath returns the configured data folder path.
// Relies on InitFilePaths being called for proper initialization.
// Returns a fallback path if not initialized, with a warning.
func GetDataFolderPath() string {
	if configuredDataFolder == "" {
		LogMessage(LevelWarn, "GetDataFolderPath called before InitFilePaths. Returning fallback default './data'.")
		absPath, err := filepath.Abs("./data") // Fallback default
		if err != nil {
			LogMessage(LevelError, "Could not get absolute path for default ./data in GetDataFolderPath fallback.")
			return "./data"
		}
		return absPath
	}
	return configuredDataFolder
}

// Note: Other file-related utility functions could be added here later.
