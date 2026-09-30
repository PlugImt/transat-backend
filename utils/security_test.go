package utils

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGeneratedTokenHasNoExpiryAndSurvivesRestart(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_EXPIRATION_HOURS", "1") // must be ignored
	jwtSecret = nil

	tok, err := GenerateJWT("a@b.c", []string{"NEWF"}, "")
	if err != nil {
		t.Fatal(err)
	}

	// Simulate a restart: the secret is re-read from the environment.
	jwtSecret = nil
	parsed, err := ValidateJWT(tok)
	if err != nil {
		t.Fatalf("token rejected after restart: %v", err)
	}
	claims := parsed.Claims.(jwt.MapClaims)
	for _, k := range []string{"exp", "nbf"} {
		if _, ok := claims[k]; ok {
			t.Fatalf("token unexpectedly carries %q", k)
		}
	}
	if _, ok := claims["iat"]; !ok {
		t.Fatal("token is missing iat")
	}
}
