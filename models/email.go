package models

import "html/template"

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
