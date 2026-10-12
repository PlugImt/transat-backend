package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	devicesHandler "github.com/plugimt/transat-backend/handlers/devices"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"
)

// signOutLimiter limits one installation, not one address: behind the proxy every user shares an address, and the
// account limiter that logins use would let a burst of logins keep devices signed in.
var signOutLimiter = limiter.New(limiter.Config{
	Max:          20,
	Expiration:   time.Minute,
	KeyGenerator: func(c *fiber.Ctx) string { return "sign-out-" + c.Params("installation_id") },
	LimitReached: func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Too many requests for this device."})
	},
})

// SetupDeviceRoutes exposes device registration, presence and sign-out.
func SetupDeviceRoutes(router fiber.Router, h *devicesHandler.Handler) {
	auth := []fiber.Handler{middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail}

	router.Put("/devices/:installation_id", append(auth, h.Register)...)
	router.Get("/devices", append(auth, h.List)...)
	// No JWT on purpose: the app signs the device out when its session may already be over.
	router.Post("/devices/:installation_id/sign-out", signOutLimiter, h.SignOut)
}
