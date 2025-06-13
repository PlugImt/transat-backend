package models

// Response représente une réponse API standard
type Response struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"Opération réussie"`
	Data    interface{} `json:"data,omitempty"`
} // @name Response

// ErrorResponse représente une réponse d'erreur API
type ErrorResponse struct {
	Error   string `json:"error" example:"Message d'erreur"`
	Message string `json:"message,omitempty" example:"Détail de l'erreur"`
} // @name ErrorResponse

// AuthResponse représente une réponse d'authentification
type AuthResponse struct {
	Success bool   `json:"success" example:"true"`
	Token   string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Message string `json:"message,omitempty" example:"Connexion réussie"`
} // @name AuthResponse

// StatusResponse représente une réponse de statut simple
type StatusResponse struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message,omitempty" example:"API fonctionnelle"`
} // @name StatusResponse
