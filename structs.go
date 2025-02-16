package main

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
	GraduationYear    string `json:"graduation_year"`
	Campus            string `json:"campus"`
	TotalUsers        int    `json:"total_users"`
}

type verificationCodeData struct {
	VerificationCode           string `json:"verification_code"`
	VerificationCodeExpiration string `json:"verification_code_expiration"`
}

type PushToken struct {
	Token  string   `json:"token"`
	UserID string   `json:"userId"`
	Groups []string `json:"groups"`
}

type TokenStore struct {
	Tokens map[string]PushToken // Key: token
	Groups map[string][]string  // Key: group name, Value: array of tokens
}

var tokenStore = TokenStore{
	Tokens: make(map[string]PushToken),
	Groups: make(map[string][]string),
}

type NotificationPayload struct {
	To    string      `json:"to"`
	Title string      `json:"title"`
	Body  string      `json:"body"`
	Data  interface{} `json:"data,omitempty"`
}
