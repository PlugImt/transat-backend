package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// EnsureDataFolder checks if the data folder exists, creates it if not.
// Reads the target path from the DATA_FOLDER environment variable or defaults to "./data".
func EnsureDataFolder() error {
	path := os.Getenv("DATA_FOLDER")
	if path == "" {
		// Default to a relative path in the current working directory
		path = "./data"
		log.Println("DATA_FOLDER environment variable not set, using default './data'")
	}

	// Convert to absolute path if relative
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			log.Printf("Error converting relative path to absolute: %v", err)
			return fmt.Errorf("error converting relative path to absolute: %w", err)
		}
		path = absPath
		log.Printf("Using absolute path: %s", path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Data folder '%s' does not exist. Creating...", path)
		// Create the directory and any necessary parents
		if err := os.MkdirAll(path, 0750); err != nil { // 0750 allows owner rwx, group rx
			log.Printf("Failed to create data folder '%s': %v", path, err)
			return fmt.Errorf("failed to create data folder '%s': %w", path, err)
		}
		log.Printf("Data folder '%s' created successfully.", path)
	} else if err != nil {
		// Another error occurred (e.g., permission denied)
		log.Printf("Error checking data folder '%s': %v", path, err)
		return fmt.Errorf("error checking data folder '%s': %w", path, err)
	} else {
		// Folder exists
		log.Printf("Data folder '%s' already exists.", path)
	}
	return nil
}

// GetDataFolderPath returns the configured data folder path
func GetDataFolderPath() string {
	path := os.Getenv("DATA_FOLDER")
	if path == "" {
		path = "./data"
	}

	// Convert to absolute path if relative
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err == nil {
			path = absPath
		}
	}

	return path
}

// Note: Other file-related utility functions could be added here later.
