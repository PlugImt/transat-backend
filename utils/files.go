package utils

import (
	"fmt"
	"log"
	"os"
)

// EnsureDataFolder checks if the /data folder exists, creates it if not.
// Reads the target path from the DATA_FOLDER environment variable or defaults to "/data".
func EnsureDataFolder() error {
	path := os.Getenv("DATA_FOLDER")
	if path == "" {
		path = "/data" // Default path
		log.Println("DATA_FOLDER environment variable not set, using default '/data'")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Data folder '%s' does not exist. Creating...", path)
		// Create the directory and any necessary parents
		if err := os.MkdirAll(path, 0755); err != nil { // 0755 allows owner rwx, group rx, others rx
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

// Note: Other file-related utility functions could be added here later.
