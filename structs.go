package main

type TurnstileResponse struct {
	Success bool     `json:"success"`
	Error   []string `json:"error-codes"`
}
