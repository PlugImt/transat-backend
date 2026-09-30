package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/plugimt/transat-backend/utils"
)

const testSecret = "test-secret-that-is-stable-across-restarts"

type fakeStore struct {
	changedAt int64
	err       error
}

func (f fakeStore) PasswordChangedAt(context.Context, string) (int64, error) {
	return f.changedAt, f.err
}

func setup(t *testing.T, store UserSessionStore) *fiber.App {
	t.Helper()
	t.Setenv("JWT_SECRET", testSecret)
	app := fiber.New()
	app.Get("/me", NewJWTMiddleware(store), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"email": c.Locals("email")})
	})
	return app
}

// tokenIssuedAt signs a token as an earlier login would have (same secret, older iat).
func tokenIssuedAt(t *testing.T, email string, iat time.Time) string {
	t.Helper()
	claims := jwt.MapClaims{"email": email, "iss": "transat-backend", "iat": iat.Unix(), "nbf": iat.Unix()}
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func call(t *testing.T, app *fiber.App, token string) (int, string) {
	t.Helper()
	req := httptest.NewRequest("GET", "/me", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var body struct {
		Error string `json:"error"`
	}
	_ = json.Unmarshal(raw, &body)
	return resp.StatusCode, body.Error
}

func TestValidTokenStaysValidAcrossDaysRestartsAndLogins(t *testing.T) {
	store := fakeStore{changedAt: time.Now().AddDate(-1, 0, 0).Unix()}
	app := setup(t, store)

	// Tokens issued 1 day, 30 days and 2 years ago (multiple logins, several deploys) all remain valid.
	for _, age := range []time.Duration{24 * time.Hour, 30 * 24 * time.Hour, 2 * 365 * 24 * time.Hour} {
		iat := time.Now().Add(-age)
		if age > 365*24*time.Hour {
			store.changedAt = iat.Add(-time.Hour).Unix()
			app = setup(t, store)
		}
		if code, msg := call(t, app, tokenIssuedAt(t, "a@b.c", iat)); code != 200 {
			t.Fatalf("age %v: got %d %q, want 200", age, code, msg)
		}
	}

	// A token freshly issued by GenerateJWT has no exp and stays valid after a "restart" (new app instance).
	tok, err := utils.GenerateJWT("a@b.c", []string{"NEWF"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if code, msg := call(t, setup(t, fakeStore{changedAt: time.Now().Add(-time.Hour).Unix()}), tok); code != 200 {
		t.Fatalf("got %d %q, want 200", code, msg)
	}
	// A second login does not affect the first token.
	if _, err := utils.GenerateJWT("a@b.c", []string{"NEWF"}, ""); err != nil {
		t.Fatal(err)
	}
	if code, msg := call(t, setup(t, fakeStore{changedAt: time.Now().Add(-time.Hour).Unix()}), tok); code != 200 {
		t.Fatalf("got %d %q, want 200", code, msg)
	}
}

func TestTokenIssuedBeforePasswordChangeIsRejected(t *testing.T) {
	changed := time.Now().Add(-time.Hour)
	app := setup(t, fakeStore{changedAt: changed.Unix()})

	code, msg := call(t, app, tokenIssuedAt(t, "a@b.c", changed.Add(-time.Minute)))
	if code != 401 || msg != MsgPasswordChanged {
		t.Fatalf("got %d %q, want 401 %q", code, msg, MsgPasswordChanged)
	}

	// Token issued in the same second as the change (e.g. login right after) stays valid.
	if code, msg := call(t, app, tokenIssuedAt(t, "a@b.c", changed)); code != 200 {
		t.Fatalf("got %d %q, want 200", code, msg)
	}
}

func TestTokenForDeletedUserIsRejected(t *testing.T) {
	app := setup(t, fakeStore{err: ErrUserNotFound})

	code, msg := call(t, app, tokenIssuedAt(t, "a@b.c", time.Now()))
	if code != 401 || msg != MsgAccountDeleted {
		t.Fatalf("got %d %q, want 401 %q", code, msg, MsgAccountDeleted)
	}
}

func TestTransientDBErrorIsNot401(t *testing.T) {
	app := setup(t, fakeStore{err: errors.New("connection refused")})

	code, msg := call(t, app, tokenIssuedAt(t, "a@b.c", time.Now()))
	if code < 500 || code > 599 {
		t.Fatalf("got %d %q, want 5xx", code, msg)
	}
}

func TestMalformedOrMissingTokenIs401(t *testing.T) {
	app := setup(t, fakeStore{})

	if code, _ := call(t, app, ""); code != 401 {
		t.Fatalf("missing token: got %d, want 401", code)
	}
	if code, _ := call(t, app, "garbage"); code != 401 {
		t.Fatalf("garbage token: got %d, want 401", code)
	}
}

func TestLegacyExpiredTokenReturnsExpiredMessage(t *testing.T) {
	app := setup(t, fakeStore{})
	claims := jwt.MapClaims{
		"email": "a@b.c", "iss": "transat-backend",
		"iat": time.Now().Add(-48 * time.Hour).Unix(),
		"exp": time.Now().Add(-24 * time.Hour).Unix(),
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}
	if code, msg := call(t, app, tok); code != 401 || msg != MsgTokenExpired {
		t.Fatalf("got %d %q, want 401 %q", code, msg, MsgTokenExpired)
	}
}
