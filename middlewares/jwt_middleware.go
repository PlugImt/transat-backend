package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/plugimt/transat-backend/utils"
)

func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	utils.LogHeader("📧 JWT Middleware")

	if authHeader == "" {
		utils.LogMessage(utils.LevelError, "Missing token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenString := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = authHeader[7:]
	}

	token, err := utils.ValidateJWT(tokenString)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Invalid token")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.LogMessage(utils.LevelError, "Invalid claims")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
	}

	c.Locals("email", claims["email"])

	utils.LogMessage(utils.LevelInfo, "Token is valid")
	utils.LogLineKeyValue(utils.LevelInfo, "Email", claims["email"])
	utils.LogFooter()

	return c.Next()
}
