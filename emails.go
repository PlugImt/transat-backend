package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"

	"Transat_2.0_Backend/models"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	gomail "gopkg.in/mail.v2"
)

type Email struct {
	Recipient string      `json:"recipient" bson:"recipient"`
	Template  string      `json:"template" bson:"template"`
	Sender    EmailSender `json:"sender" bson:"sender"`
	Language  string      `json:"language" bson:"language"`
}

type EmailSender struct {
	Email string `json:"email" bson:"email"`
	Name  string `json:"name" bson:"name"`
}

type EmailData struct {
	Data interface{}
	T    func(string, ...interface{}) template.HTML
}

func sendEmail(mailDetails models.Email, emailData interface{}) error {
	// Read the HTML template file
	htmlTemplate, err := os.ReadFile(mailDetails.Template)
	if err != nil {
		log.Printf("error reading template file: %v", err)
		return fmt.Errorf("error reading template file: %v", err)
	}

	// Créer le localizer pour la langue spécifiée
	localizer := GetLocalizer(mailDetails.Language)
	if localizer == nil {
		log.Printf("error getting localizer: %v", err)
		localizer = GetLocalizer("fr") // Fallback sur le français
	}

	// Créer la fonction de traduction
	translate := func(key string, data ...interface{}) template.HTML {
		config := &i18n.LocalizeConfig{
			MessageID: key,
		}
		if len(data) > 0 {
			config.TemplateData = data[0]
		}
		message, err := localizer.Localize(config)
		if err != nil {
			log.Printf("error translating key %s: %v", key, err)
			return template.HTML(key)
		}
		return template.HTML(message)
	}

	// Traduire l'objet de l'email avec les données de template
	subject := translate(mailDetails.SubjectKey, emailData)

	// Créer les données pour le template
	templateData := struct {
		Data interface{}
		T    func(string, ...interface{}) template.HTML
	}{
		Data: emailData,
		T:    translate,
	}

	// Parse the template
	tmpl, err := template.New("emailTemplate").Parse(string(htmlTemplate))
	if err != nil {
		log.Printf("error parsing template: %v", err)
		return fmt.Errorf("error parsing template: %v", err)
	}

	// Execute the template with the data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, templateData); err != nil {
		log.Printf("error executing template: %v", err)
		return fmt.Errorf("error executing template: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", mailDetails.Sender.Name+" <"+mailDetails.Sender.Email+">")
	m.SetHeader("To", mailDetails.Recipient)
	m.SetHeader("Subject", string(subject))
	m.SetBody("text/html", body.String())

	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), port, mailDetails.Sender.Email, os.Getenv("EMAIL_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		log.Printf("error sending email: %v", err)
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}
