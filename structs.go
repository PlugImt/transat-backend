package main

import (
	"database/sql"
	"github.com/go-resty/resty/v2"
)

type Newf struct {
	ID                      int    `json:"id_newf"`
	Email                   string `json:"email"`
	Password                string `json:"password"`
	NewPassword             string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
	PasswordUpdatedDate     string `json:"password_updated_date"`
	verificationCodeData
	CreationDate      string `json:"creation_date"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	PhoneNumber       string `json:"phone_number"`
	ProfilePicture    string `json:"profile_picture"`
	NotificationToken string `json:"notification_token"`
	GraduationYear    int    `json:"graduation_year"`
	Campus            string `json:"campus"`
	TotalUsers        int    `json:"total_newf"`
}

type verificationCodeData struct {
	VerificationCode           string `json:"verification_code"`
	VerificationCodeExpiration string `json:"verification_code_expiration"`
}

type NotificationPayload struct {
	UserEmails []string `json:"userEmails,omitempty"`
	Groups     []string `json:"groups,omitempty"`
	Title      string   `json:"title"`
	Message    string   `json:"message,omitempty"`
	ImageURL   string   `json:"imageUrl,omitempty"`
	Screen     string   `json:"screen,omitempty"`
}

type NotificationTarget struct {
	Email             string
	NotificationToken string
}

type NotificationService struct {
	db     *sql.DB
	client *resty.Client
}
