package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	mailgun "github.com/mailgun/mailgun-go/v4"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	i18n "github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"
)

// EmailService handles sending emails using Mailgun API.
type EmailService struct {
	client      mailgun.Mailgun
	senderName  string
	senderEmail string
	domain      string
}

// NewEmailService creates a new EmailService instance using Mailgun API.
// It requires Mailgun API key, domain, sender email, and sender name.
func NewEmailService(apiKey, domain, senderEmail, senderName string) *EmailService {
	client := mailgun.NewMailgun(domain, apiKey)
	client.SetAPIBase("https://api.eu.mailgun.net/v3")

	return &EmailService{
		client:      client,
		senderName:  senderName,
		senderEmail: senderEmail,
		domain:      domain,
	}
}

// SendEmail sends an email using a template and data via Mailgun API.
// It utilizes the i18n setup for subject and template content localization.
func (es *EmailService) SendEmail(mailDetails models.Email, emailData interface{}) error {
	// Read the HTML template file
	htmlTemplate, err := os.ReadFile(mailDetails.Template)
	if err != nil {
		log.Printf("Error reading template file '%s': %v", mailDetails.Template, err)
		return fmt.Errorf("error reading template file '%s': %w", mailDetails.Template, err)
	}

	// Get localizer for the specified language
	localizer := i18n.GetLocalizer(mailDetails.Language)
	if localizer == nil {
		log.Printf("Error: Could not get localizer for language '%s'. Using default.", mailDetails.Language)
		localizer = i18n.GetLocalizer("fr") // Fallback
	}

	// Create the translation function for the template
	translateFunc := func(key string, data ...interface{}) template.HTML {
		config := &goi18n.LocalizeConfig{
			MessageID: key,
		}
		// Pass template data if provided (e.g., for pluralization or variable substitution)
		if len(data) > 0 {
			// Check if the data is a map or struct suitable for template data
			if templateMap, ok := data[0].(map[string]interface{}); ok {
				config.TemplateData = templateMap
			} else {
				// Otherwise, assume it's the primary data object itself
				config.TemplateData = data[0]
			}
		} else {
			// Pass the main emailData if no specific data is given to T
			config.TemplateData = emailData
		}

		message, err := localizer.Localize(config)
		if err != nil {
			log.Printf("Error translating key '%s' for language '%s': %v", key, mailDetails.Language, err)
			//gosec:disable G203 -- Faux positif
			return template.HTML("!!!" + template.HTMLEscapeString(key) + "!!!") // Return key clearly marked as untranslated
		}
		//gosec:disable G203 -- Faux positif
		return template.HTML(message)
	}

	// Prepare data for the template, including the translation function
	templateData := struct {
		Data interface{}
		T    func(string, ...interface{}) template.HTML // Make T available in the template
	}{
		Data: emailData,
		T:    translateFunc,
	}

	// Translate the subject line separately
	// Use GetLocalizer to ensure fallback logic is consistent
	subjectLocalizer := i18n.GetLocalizer(mailDetails.Language)
	subject, err := subjectLocalizer.Localize(&goi18n.LocalizeConfig{
		MessageID:    mailDetails.SubjectKey,
		TemplateData: emailData, // Pass data to subject translation as well
	})
	if err != nil {
		log.Printf("Error translating subject key '%s' for language '%s': %v", mailDetails.SubjectKey, mailDetails.Language, err)
		subject = mailDetails.SubjectKey // Fallback to the key itself
	}

	// Parse the HTML template
	tmpl, err := template.New(filepath.Base(mailDetails.Template)).Parse(string(htmlTemplate))
	if err != nil {
		log.Printf("Error parsing template '%s': %v", mailDetails.Template, err)
		return fmt.Errorf("error parsing template '%s': %w", mailDetails.Template, err)
	}

	// Execute the template to generate the email body
	var body bytes.Buffer
	if err := tmpl.Execute(&body, templateData); err != nil {
		log.Printf("Error executing template '%s': %v", mailDetails.Template, err)
		return fmt.Errorf("error executing template '%s': %w", mailDetails.Template, err)
	}

	// Create the email message using Mailgun
	message := mailgun.NewMessage(
		fmt.Sprintf("%s <%s>", es.senderName, es.senderEmail), // From
		subject,               // Subject
		"",                    // Text body (empty for HTML-only)
		mailDetails.Recipient, // To
	)
	message.SetHTML(body.String())

	// Send the email using Mailgun API
	ctx := context.Background()
	_, _, err = es.client.Send(ctx, message)
	if err != nil {
		log.Printf("Error sending email to %s: %v", mailDetails.Recipient, err)
		return fmt.Errorf("error sending email: %w", err)
	}

	log.Printf("Email sent successfully to %s (Subject: %s)", mailDetails.Recipient, subject)
	return nil
}
