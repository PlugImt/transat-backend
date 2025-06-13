package models

type Newf struct {
	ID                      int    `json:"id_newf"`
	Email                   string `json:"email"`
	Password                string `json:"password"`
	NewPassword             string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
	PasswordUpdatedDate     string `json:"password_updated_date"`
	VerificationCodeData
	CreationDate      string `json:"creation_date"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	PhoneNumber       string `json:"phone_number"`
	ProfilePicture    string `json:"profile_picture"`
	NotificationToken string `json:"notification_token"`
	GraduationYear    int    `json:"graduation_year"`
	Campus            string `json:"campus"`
	TotalUsers        int    `json:"total_newf"`
	Language          string `json:"language"`
}

// VerificationCodeData holds verification code details.
type VerificationCodeData struct {
	VerificationCode           string `json:"verification_code"`
	VerificationCodeExpiration string `json:"verification_code_expiration"` // Formatted time (e.g., "15h04")
}

// --- Request Structs ---

// LoginRequest is used for user authentication.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"motdepasse123"`
} // @name LoginRequest

// VerificationRequest is used for verifying an account.
type VerificationRequest struct {
	Email            string `json:"email"`
	VerificationCode string `json:"verification_code"`
}

// EmailRequest is used when only an email is needed, e.g., resending a code.
type EmailRequest struct {
	Email string `json:"email"`
}

// ChangePasswordRequest is used for changing a user's password.
type ChangePasswordRequest struct {
	Email                   string `json:"email"`
	OldPassword             string `json:"old_password,omitempty"` // Optional: Used when logged in
	NewPassword             string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
	VerificationCode        string `json:"verification_code,omitempty"` // Optional: Used for password reset
}
