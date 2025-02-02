package main

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/http"
)

func registerToken(c *fiber.Ctx) error {
	var token PushToken
	if err := c.BodyParser(&token); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Store token
	tokenStore.Tokens[token.Token] = token

	// Update groups
	for _, group := range token.Groups {
		if _, exists := tokenStore.Groups[group]; !exists {
			tokenStore.Groups[group] = make([]string, 0)
		}
		tokenStore.Groups[group] = append(tokenStore.Groups[group], token.Token)
	}

	return c.JSON(fiber.Map{
		"message": "Token registered successfully",
	})
}

func sendNotification(c *fiber.Ctx) error {
	type SendRequest struct {
		Group string      `json:"group"`
		Title string      `json:"title"`
		Body  string      `json:"body"`
		Data  interface{} `json:"data,omitempty"`
	}

	var req SendRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get tokens for the group
	tokens, exists := tokenStore.Groups[req.Group]
	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "Group not found",
		})
	}

	// Send notification to each token in the group
	for _, token := range tokens {
		err := sendExpoNotification(NotificationPayload{
			To:    token,
			Title: req.Title,
			Body:  req.Body,
			Data:  req.Data,
		})
		if err != nil {
			log.Printf("Error sending notification to token %s: %v", token, err)
		}
	}

	return c.JSON(fiber.Map{
		"message": "Notifications sent successfully",
	})
}

func sendExpoNotification(payload NotificationPayload) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		"https://exp.host/--/api/v2/push/send",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	return nil
}
