package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ensureDataFolder checks if the /data folder exists, creates it if not
func ensureDataFolder() error {
	path := "/data"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("/data folder does not exist. Creating...")
		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Println("Failed to create /data folder:", err)
			return err
		}
		fmt.Println("/data folder created successfully.")
	} else {
		fmt.Println("/data folder exists.")
	}
	return nil
}

// uploadImage handles file uploads and saves them to the /data directory
func uploadImage(c *fiber.Ctx) error {
	fmt.Println("Received file upload request")

	// Get uploaded file from the request
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println("No image uploaded or invalid file:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No image uploaded or invalid file",
		})
	}
	fmt.Println("Uploaded file name:", file.Filename)

	// Validate file is an image
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		fmt.Println("Uploaded file is not an image")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Uploaded file is not an image",
		})
	}

	// Ensure the /data directory exists
	if err := ensureDataFolder(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to ensure storage directory",
		})
	}

	// Set destination path
	filename := filepath.Base(file.Filename) // Get filename without path
	dst := fmt.Sprintf("/data/%s", filename)
	fmt.Println("Destination path:", dst)

	// Save file to /data directory
	fileReader, err := file.Open()
	if err != nil {
		fmt.Println("Failed to read uploaded file:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read uploaded file",
		})
	}
	defer fileReader.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Println("Failed to create destination file:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create destination file",
		})
	}
	defer dstFile.Close()

	// Copy file contents
	if _, err := io.Copy(dstFile, fileReader); err != nil {
		fmt.Println("Failed to save file:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}
	fmt.Println("File uploaded successfully:", dst)

	// Return success response with image URL
	return c.JSON(fiber.Map{
		"success": true,
		"url":     "/data/" + filename,
	})
}

// serveImage serves images from the /data directory
func serveImage(c *fiber.Ctx) error {
	filename := c.Params("filename")
	filepath := fmt.Sprintf("/data/%s", filename)

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Println("Image not found:", filepath)
		return c.Status(fiber.StatusNotFound).SendString("Image not found")
	}

	fmt.Println("Serving image:", filepath)
	// Send the file
	return c.SendFile(filepath)
}
