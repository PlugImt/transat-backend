package main

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
