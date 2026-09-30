package middlewares

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strings"

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

// ErrUserNotFound is returned by a UserSessionStore when the account no longer exists.
var ErrUserNotFound = errors.New("user not found")

// UserSessionStore provides the data needed to decide whether a token is still valid.
type UserSessionStore interface {
	// PasswordChangedAt returns the Unix time (seconds) of the last password change.
	PasswordChangedAt(ctx context.Context, email string) (int64, error)
}

type sqlSessionStore struct{ db *sql.DB }

// NewSQLSessionStore returns a UserSessionStore backed by the newf table.
func NewSQLSessionStore(db *sql.DB) UserSessionStore { return sqlSessionStore{db: db} }

func (s sqlSessionStore) PasswordChangedAt(ctx context.Context, email string) (int64, error) {
	// The column is a timestamp without time zone written with NOW(), so cast using the session time zone.
	const q = `SELECT FLOOR(EXTRACT(EPOCH FROM password_updated_date::timestamptz))::bigint FROM newf WHERE email = $1`
	var ts int64
	if err := s.db.QueryRowContext(ctx, q, email).Scan(&ts); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		return 0, err
	}
	return ts, nil
}

var defaultSessionStore UserSessionStore

// InitJWTMiddleware configures the store used by JWTMiddleware. Call once at startup.
func InitJWTMiddleware(db *sql.DB) {
	defaultSessionStore = NewSQLSessionStore(db)
}

// JWTMiddleware authenticates the request with the configured session store.
func JWTMiddleware(c *fiber.Ctx) error {
	return NewJWTMiddleware(defaultSessionStore)(c)
}

// NewJWTMiddleware rejects a token only if it is malformed, the account no longer exists,
// or the password changed after the token was issued.
func NewJWTMiddleware(store UserSessionStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
			// Only legacy tokens carry an exp; they can never work again, so the client must re-login.
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
		email := strings.ToLower(strings.TrimSpace(rawEmail))
		iat, hasIat := claims["iat"].(float64)
		if email == "" || !hasIat {
			utils.LogMessage(utils.LevelError, "Token is missing email or iat")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
		}

		if store == nil {
			utils.LogMessage(utils.LevelError, "JWT middleware session store is not configured")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Authentication error"})
		}

		changedAt, err := store.PasswordChangedAt(c.UserContext(), email)
		if err != nil {
			if errors.Is(err, ErrUserNotFound) {
				utils.LogMessage(utils.LevelWarn, "Token for deleted account")
				utils.LogLineKeyValue(utils.LevelWarn, "Email", email)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": MsgAccountDeleted})
			}
			utils.LogMessage(utils.LevelError, "Failed to check session validity")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Service temporarily unavailable"})
		}

		// iat has second precision, so a token issued in the same second as the change is still accepted.
		if int64(math.Floor(iat)) < changedAt {
			utils.LogMessage(utils.LevelWarn, "Token issued before password change")
			utils.LogLineKeyValue(utils.LevelWarn, "Email", email)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": MsgPasswordChanged})
		}

		c.Locals("email", email)

		utils.LogMessage(utils.LevelInfo, "Token is valid")
		utils.LogLineKeyValue(utils.LevelInfo, "Email", email)

		return c.Next()
	}
}
