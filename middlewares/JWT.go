package middlewares

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/plugimt/transat-backend/utils"
)

// The mobile client matches these on 401 responses to decide to log the user out.
const (
	MsgAccountDeleted  = "Account deleted"
	MsgPasswordChanged = "Password changed, please log in again"
	MsgTokenExpired    = "Token expired, please log in again"
)

var errUserNotFound = errors.New("user not found")

// password_updated_date has no time zone, so a DST change can shift it by an hour; the extra minute covers clock drift.
const passwordChangeLeeway = time.Hour + time.Minute

var sessionDB *sql.DB

// InitJWTMiddleware sets the database used to check sessions. Call once at startup.
func InitJWTMiddleware(db *sql.DB) {
	sessionDB = db
}

// passwordChangedAt returns the stored email and the Unix time (seconds) of the user's last password change.
// Both the token's email and its lowercase form are tried so legacy mixed-case accounts still match.
func passwordChangedAt(ctx context.Context, email, lowerEmail string) (string, int64, error) {
	// The column is a timestamp without time zone written with NOW(), so cast using the session time zone.
	const q = `SELECT email, FLOOR(EXTRACT(EPOCH FROM password_updated_date::timestamptz))::bigint FROM newf WHERE email IN ($1, $2) ORDER BY (email = $1) DESC LIMIT 1`
	var storedEmail string
	var ts int64
	if err := sessionDB.QueryRowContext(ctx, q, email, lowerEmail).Scan(&storedEmail, &ts); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", 0, errUserNotFound
		}
		return "", 0, err
	}
	return storedEmail, ts, nil
}

// JWTMiddleware rejects a token if it is invalid or expired, the account no longer exists,
// or the password changed after the token was issued.
func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	utils.LogHeader("📧 JWT Middleware")
	defer utils.LogFooter()

	if authHeader == "" {
		utils.LogMessage(utils.LevelError, "Missing token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := utils.ValidateJWT(tokenString)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Invalid token")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		// An expired token can never work again, so the client must log the user out.
		if errors.Is(err, jwt.ErrTokenExpired) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": MsgTokenExpired})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.LogMessage(utils.LevelError, "Invalid claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
	}

	rawEmail, _ := claims["email"].(string)
	tokenEmail := strings.TrimSpace(rawEmail)
	iat, hasIat := claims["iat"].(float64)
	if tokenEmail == "" || !hasIat {
		utils.LogMessage(utils.LevelError, "Token is missing email or iat")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
	}

	email, changedAt, err := passwordChangedAt(c.UserContext(), tokenEmail, strings.ToLower(tokenEmail))
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			utils.LogMessage(utils.LevelWarn, "Token for deleted account")
			utils.LogLineKeyValue(utils.LevelWarn, "Email", tokenEmail)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": MsgAccountDeleted})
		}
		// Not a 401: a database hiccup must not log the user out.
		utils.LogMessage(utils.LevelError, "Failed to check session validity")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Service temporarily unavailable"})
	}

	// A token is rejected only if it predates the password change by more than the leeway.
	if int64(math.Floor(iat))+int64(passwordChangeLeeway/time.Second) < changedAt {
		utils.LogMessage(utils.LevelWarn, "Token issued before password change")
		utils.LogLineKeyValue(utils.LevelWarn, "Email", email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": MsgPasswordChanged})
	}

	c.Locals("email", email)

	utils.LogMessage(utils.LevelInfo, "Token is valid")
	utils.LogLineKeyValue(utils.LevelInfo, "Email", email)

	return c.Next()
}
