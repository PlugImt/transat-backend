package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	if err == nil {
		log.Printf("Brevo: sent | to=%s subject=%q message=%+v", mailDetails.Recipient, subject, res)
		return nil
	}

	// Brevo failed — capture error details and notify Discord
	status := ""
	if httpResp != nil {
		status = httpResp.Status
	}
	var respBody string
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
		if b, rErr := io.ReadAll(httpResp.Body); rErr == nil {
			respBody = string(b)
		}
	}
	log.Printf("Brevo: send failed | status=%s to=%s subject=%q err=%v body=%s", status, mailDetails.Recipient, subject, err, truncate(respBody, 500))
	es.notifyDiscordEmailFailure("Brevo", mailDetails.Recipient, subject, fmt.Sprintf("status=%s err=%v body=%s", status, err, truncate(respBody, 2000)))

	// Fallback 1: Primary SMTP (Gmail)
	smtp1 := smtpConfig{
		Host:       firstNonEmpty(os.Getenv("EMAIL_HOST_GMAIL_1"), "smtp.gmail.com"),
		Port:       firstNonEmpty(os.Getenv("EMAIL_PORT_GMAIL_1"), "587"),
		Username:   os.Getenv("EMAIL_SENDER_NAME_GMAIL_1"),
		Password:   os.Getenv("EMAIL_PASSWORD_GMAIL_1"),
		SenderName: firstNonEmpty(os.Getenv("EMAIL_SENDER_NAME_GMAIL_1"), es.senderName),
		Sender:     firstNonEmpty(os.Getenv("EMAIL_SENDER_GMAIL_1"), es.senderEmail),
	}
	var err1 error
	err1 = es.sendSMTP(smtp1, mailDetails.Recipient, subject, body.String())
	if err1 == nil {
		log.Printf("SMTP1: sent | to=%s subject=%q via %s:%s as %s", mailDetails.Recipient, subject, smtp1.Host, smtp1.Port, smtp1.Sender)
		return nil
	}
	log.Printf("SMTP1: send failed | host=%s port=%s user=%s to=%s subject=%q err=%v", smtp1.Host, smtp1.Port, smtp1.Username, mailDetails.Recipient, subject, err1)
	es.notifyDiscordEmailFailure("SMTP1 - GOOGLE TRANSAT", mailDetails.Recipient, subject, err1.Error())

	// Fallback 2: Secondary SMTP (Second Gmail)
	smtp2 := smtpConfig{
		Host:       firstNonEmpty(os.Getenv("EMAIL_HOST_GMAIL_2"), "smtp.gmail.com"),
		Port:       firstNonEmpty(os.Getenv("EMAIL_PORT_GMAIL_2"), "587"),
		Username:   os.Getenv("EMAIL_SENDER_NAME_GMAIL_2"),
		Password:   os.Getenv("EMAIL_PASSWORD_GMAIL_2"),
		SenderName: firstNonEmpty(os.Getenv("EMAIL_SENDER_NAME_GMAIL_2"), es.senderName),
		Sender:     os.Getenv("EMAIL_SENDER_GMAIL_2"),
	}
	var err2 error
	err2 = es.sendSMTP(smtp2, mailDetails.Recipient, subject, body.String())
	if err2 == nil {
		log.Printf("SMTP2: sent | to=%s subject=%q via %s:%s as %s", mailDetails.Recipient, subject, smtp2.Host, smtp2.Port, smtp2.Sender)
		return nil
	}
	log.Printf("SMTP2: send failed | host=%s port=%s user=%s to=%s subject=%q err=%v", smtp2.Host, smtp2.Port, smtp2.Username, mailDetails.Recipient, subject, err2)
	es.notifyDiscordEmailFailure("SMTP2 - GOOGLE DESTIMT", mailDetails.Recipient, subject, err2.Error())
	return fmt.Errorf("all providers failed: brevo=%v; smtp1=%v; smtp2=%v", err, err1, err2)
}

// smtpConfig holds SMTP credentials
type smtpConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	SenderName string
	Sender     string
}

// sendSMTP sends HTML email via SMTP using net/smtp (STARTTLS on :587 if supported)
func (es *EmailService) sendSMTP(cfg smtpConfig, toEmail, subject, htmlBody string) error {
	if cfg.Host == "" || cfg.Port == "" || cfg.Username == "" || cfg.Password == "" || cfg.Sender == "" {
		return fmt.Errorf("smtp config incomplete")
	}
	addr := cfg.Host + ":" + cfg.Port
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	// Build MIME message
	headers := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", cfg.SenderName, cfg.Sender),
		"To":           toEmail,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}
	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(k)
		msg.WriteString(": ")
		msg.WriteString(v)
		msg.WriteString("\r\n")
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	// Note: net/smtp.SendMail will use STARTTLS if server supports it
	if err := smtp.SendMail(addr, auth, cfg.Sender, []string{toEmail}, msg.Bytes()); err != nil {
		return err
	}
	return nil
}

// notifyDiscordEmailFailure sends a discord embed to a dedicated webhook for email failures
func (es *EmailService) notifyDiscordEmailFailure(provider, recipient, subject, errMsg string) {
	webhook := firstNonEmpty(os.Getenv("DISCORD_EMAIL_ALERT_WEBHOOK"))
	ds := NewDiscordService(webhook)
	embed := discordEmbed{
		Title:     "Email Send Failure",
		Color:     0xE53935,
		Timestamp: timeNowUTC(),
		Fields: []discordEmbedField{
			{Name: "Provider", Value: provider, Inline: true},
			{Name: "Recipient", Value: recipient, Inline: true},
			{Name: "Subject", Value: safeStr(subject), Inline: false},
			{Name: "Error", Value: truncate(errMsg, 1000), Inline: false},
		},
	}
	if err := ds.sendEmbed(embed); err != nil {
		log.Printf("Discord alert failed: %v", err)
	}
}

// helpers
func truncate(s string, n int) string {
	if n <= 0 || len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}

func safeStr(s string) string {
	if strings.TrimSpace(s) == "" {
		return "(empty)"
	}
	return s
}

func timeNowUTC() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
