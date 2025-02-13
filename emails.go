package main

import (
	"bytes"
	"fmt"
	gomail "gopkg.in/mail.v2"
	"html/template"
	"log"
	"os"
)

type Email struct {
	Recipient string      `json:"recipient" bson:"recipient"`
	Subject   string      `json:"subject" bson:"subject"`
	Template  string      `json:"template" bson:"template"`
	Sender    EmailSender `json:"sender" bson:"sender"`
}

type EmailSender struct {
	Email string `json:"email" bson:"email"`
	Name  string `json:"name" bson:"name"`
}

func sendEmail(mailDetails Email, emailData interface{}) error {
	// Read the HTML template file
	htmlTemplate, err := os.ReadFile(mailDetails.Template)
	if err != nil {
		log.Printf("error reading template file: %v", err)
		return fmt.Errorf("error reading template file: %v", err)
	}

	// Parse the template
	tmpl, err := template.New("emailTemplate").Parse(string(htmlTemplate))
	if err != nil {
		log.Printf("error parsing template: %v", err)
		return fmt.Errorf("error parsing template: %v", err)
	}

	// Execute the template with the data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, emailData); err != nil {
		log.Printf("error executing template: %v", err)
		return fmt.Errorf("error executing template: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", mailDetails.Sender.Name+" <"+mailDetails.Sender.Email+">")
	m.SetHeader("To", mailDetails.Recipient)
	m.SetHeader("Subject", mailDetails.Subject)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, mailDetails.Sender.Email, os.Getenv("EMAIL_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		log.Printf("error sending email: %v", err)
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}
