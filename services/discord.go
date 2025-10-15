package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/plugimt/transat-backend/models"
)

type DiscordService struct {
	WebhookURL string
	client     *http.Client
}

func NewDiscordService(webhookURL string) *DiscordService {
	if webhookURL == "" {
		webhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
	}
	return &DiscordService{
		WebhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

type discordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type discordEmbed struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Color       int                    `json:"color,omitempty"`
	Fields      []discordEmbedField    `json:"fields,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Thumbnail   *discordEmbedThumbnail `json:"thumbnail,omitempty"`
	Image       *discordEmbedImage     `json:"image,omitempty"`
}

type discordEmbedThumbnail struct {
	URL string `json:"url"`
}

type discordEmbedImage struct {
	URL string `json:"url"`
}

type discordPayload struct {
	Content string         `json:"content,omitempty"`
	Embeds  []discordEmbed `json:"embeds,omitempty"`
}

func (ds *DiscordService) sendEmbed(e discordEmbed) error {
	if ds == nil || ds.WebhookURL == "" {
		return fmt.Errorf("discord webhook URL not configured")
	}
	payload := discordPayload{Embeds: []discordEmbed{e}}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", ds.WebhookURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := ds.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook returned status %d", resp.StatusCode)
	}
	return nil
}

func safe(v string, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func (ds *DiscordService) SendUserRegistered(user models.Newf, numberOfAccounts int) error {
	embed := discordEmbed{
		Title:     "New Account Created (Unverified)",
		Color:     0xFFA000, // amber
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Fields: []discordEmbedField{
			{Name: "Email", Value: safe(user.Email, "N/A"), Inline: true},
			{Name: "First Name", Value: safe(user.FirstName, "N/A"), Inline: true},
			{Name: "Last Name", Value: safe(user.LastName, "N/A"), Inline: true},
			{Name: "Language", Value: safe(user.Language, "N/A"), Inline: true},
			{Name: "Created At", Value: safe(user.CreationDate, "N/A"), Inline: true},
			{Name: "Number of Accounts", Value: strconv.Itoa(numberOfAccounts), Inline: true},
		},
	}
	return ds.sendEmbed(embed)
}

func (ds *DiscordService) SendUserVerified(user models.Newf, numberOfAccounts int) error {
	embed := discordEmbed{
		Title:     "Account Verified",
		Color:     0x00C853, // green
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Fields: []discordEmbedField{
			{Name: "Email", Value: safe(user.Email, "N/A"), Inline: true},
			{Name: "First Name", Value: safe(user.FirstName, "N/A"), Inline: true},
			{Name: "Last Name", Value: safe(user.LastName, "N/A"), Inline: true},
			{Name: "Language", Value: safe(user.Language, "N/A"), Inline: true},
			{Name: "Created At", Value: safe(user.CreationDate, "N/A"), Inline: true},
			{Name: "Number of Accounts", Value: strconv.Itoa(numberOfAccounts), Inline: true},
		},
	}
	return ds.sendEmbed(embed)
}

func (ds *DiscordService) SendReservationCreated(itemName, startDate, endDate string, isSlot bool, user models.ReservationUser) error {
	var slotType string
	if isSlot {
		slotType = "Slot-based"
	} else {
		slotType = "Open-ended"
	}

	fields := []discordEmbedField{
		{Name: "Item", Value: safe(itemName, "N/A"), Inline: false},
		{Name: "Type", Value: slotType, Inline: true},
		{Name: "Start Date", Value: safe(startDate, "N/A"), Inline: true},
		{Name: "User", Value: fmt.Sprintf("%s %s", safe(user.FirstName, ""), safe(user.LastName, "")), Inline: true},
		{Name: "Email", Value: safe(user.Email, "N/A"), Inline: true},
	}

	if endDate != "" {
		fields = append(fields, discordEmbedField{Name: "End Date", Value: endDate, Inline: true})
	}

	embed := discordEmbed{
		Title:     "New Reservation Created",
		Color:     0x4CAF50, // green
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Fields:    fields,
	}
	if user.ProfilePicture != "" {
		embed.Thumbnail = &discordEmbedThumbnail{URL: user.ProfilePicture}
	}
	return ds.sendEmbed(embed)
}

func (ds *DiscordService) SendReservationCancelled(itemName, startDate, endDate string, isSlot bool, user models.ReservationUser) error {
	var slotType string
	if isSlot {
		slotType = "Slot-based"
	} else {
		slotType = "Open-ended"
	}

	fields := []discordEmbedField{
		{Name: "Item", Value: safe(itemName, "N/A"), Inline: false},
		{Name: "Type", Value: slotType, Inline: true},
		{Name: "Start Date", Value: safe(startDate, "N/A"), Inline: true},
		{Name: "User", Value: fmt.Sprintf("%s %s", safe(user.FirstName, ""), safe(user.LastName, "")), Inline: true},
		{Name: "Email", Value: safe(user.Email, "N/A"), Inline: true},
	}

	if endDate != "" {
		fields = append(fields, discordEmbedField{Name: "End Date", Value: endDate, Inline: true})
	}

	embed := discordEmbed{
		Title:     "Reservation Cancelled",
		Color:     0xF44336, // red
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Fields:    fields,
	}
	if user.ProfilePicture != "" {
		embed.Thumbnail = &discordEmbedThumbnail{URL: user.ProfilePicture}
	}
	return ds.sendEmbed(embed)
}
