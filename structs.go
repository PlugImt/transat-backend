package main

import (
	"database/sql"
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
	UserEmails []string    `json:"userEmails,omitempty"`
	Groups     []string    `json:"groups,omitempty"`
	Title      string      `json:"title"`
	Message    string      `json:"body,omitempty"`
	ImageURL   string      `json:"imageUrl,omitempty"`
	TTL        int         `json:"ttl,omitempty"`
	Subtitle   string      `json:"subtitle,omitempty"`
	Sound      string      `json:"sound,omitempty"`
	ChannelID  string      `json:"channelId,omitempty"`
	Badge      int         `json:"badge,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

type NotificationTarget struct {
	Email             string
	NotificationToken string
}

type NotificationService struct {
	db *sql.DB
}

type TraqArticle struct {
	ID           int     `json:"id_traq"`
	Name         string  `json:"name"`
	Disabled     bool    `json:"disabled"`
	Limited      bool    `json:"limited"`
	Alcohol      float32 `json:"alcohol"`
	OutOfStock   bool    `json:"out_of_stock"`
	CreationDate string  `json:"creation_date"`
	Picture      string  `json:"picture"`
	Description  string  `json:"description"`
	Price        float32 `json:"price"`
	PriceHalf    float32 `json:"price_half"`
	TraqType     string  `json:"traq_type"`
}

type TraqType struct {
	IDType int    `json:"id_traq_types"`
	Name   string `json:"name"`
}
