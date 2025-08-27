package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"

	brevo "github.com/getbrevo/brevo-go/lib"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	i18n "github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"
)

// EmailService handles sending emails using Brevo (Sendinblue) API.
type EmailService struct {
	client      *brevo.APIClient
	apiKey      string
	senderName  string
	senderEmail string
}

// NewEmailService creates a new EmailService instance using Brevo API.
// It requires Brevo API key, sender email, and sender name.
func NewEmailService(apiKey, senderEmail, senderName string) *EmailService {
	cfg := brevo.NewConfiguration()
	client := brevo.NewAPIClient(cfg)

	return &EmailService{
		client:      client,
		apiKey:      apiKey,
		senderName:  senderName,
		senderEmail: senderEmail,
	}
}

// SendEmail sends an email using a template and data via Brevo API.
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

	// Debug information before sending
	log.Printf("Brevo: preparing send | to=%s from=%s <%s> subject=%q template=%s lang=%s",
		mailDetails.Recipient, es.senderName, es.senderEmail, subject, mailDetails.Template, mailDetails.Language,
	)
	// Log safe API key fingerprint and length (do not log full secret)
	ak := es.apiKey
	mask := "(empty)"
	if ak != "" {
		if len(ak) >= 8 {
			mask = ak[:4] + "..." + ak[len(ak)-4:]
		} else {
			mask = "***"
		}
	}
	log.Printf("Brevo: api key fingerprint=%s len=%d", mask, len(ak))

	// Build Brevo payload
	send := brevo.SendSmtpEmail{
		To: []brevo.SendSmtpEmailTo{
			{Email: mailDetails.Recipient},
		},
		Sender: &brevo.SendSmtpEmailSender{
			Email: es.senderEmail,
			Name:  es.senderName,
		},
		Subject:     subject,
		HtmlContent: body.String(),
	}

	// Send the email using Brevo API (set API key in context)
	ctx := context.WithValue(context.Background(), brevo.ContextAPIKey, brevo.APIKey{Key: es.apiKey})
	res, httpResp, err := es.client.TransactionalEmailsApi.SendTransacEmail(ctx, send)
	if err != nil {
		status := ""
		if httpResp != nil {
			status = httpResp.Status
		}
		log.Printf("Brevo: send failed | status=%s to=%s subject=%q", status, mailDetails.Recipient, subject)
		// Try to log API error body when available
		if apiErr, ok := err.(interface{ Error() string }); ok {
			log.Printf("Brevo: error=%s", apiErr.Error())
		}
		if httpResp != nil {
			log.Printf("Brevo: response headers=%v", httpResp.Header)
			if httpResp.Request != nil && httpResp.Request.URL != nil {
				log.Printf("Brevo: request url=%s", httpResp.Request.URL.String())
			}
			// Try to read response body for detailed error message
			if httpResp.Body != nil {
				defer httpResp.Body.Close()
				if b, rErr := io.ReadAll(httpResp.Body); rErr == nil {
					log.Printf("Brevo: response body=%s", string(b))
				}
			}
		}
		return fmt.Errorf("error sending email: %w", err)
	}

	// Success logging (include message ID if available)
	log.Printf("Brevo: sent | to=%s subject=%q message=%+v", mailDetails.Recipient, subject, res)
	return nil
}
