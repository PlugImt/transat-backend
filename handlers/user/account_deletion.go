package user

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/services/accounts"
	"github.com/plugimt/transat-backend/utils"
)

// Accounts deletes accounts; set by the routes. Nil keeps the endpoint answering 503.
var Accounts *accounts.Service

// DeleteNewf handles DELETE /newf/me: the signed-in user deletes their own account and personal data.
func (h *UserHandler) DeleteNewf(c *fiber.Ctx) error {
	email, _ := c.Locals("email").(string)
	if Accounts == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Account deletion unavailable"})
	}
	switch err := Accounts.Delete(c.UserContext(), email); {
	case errors.Is(err, accounts.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	case errors.Is(err, accounts.ErrProtected):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "This account cannot be deleted"})
	case err != nil:
		utils.LogMessage(utils.LevelError, "Failed to delete user account")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete account"})
	}
	return c.SendStatus(fiber.StatusOK)
}
