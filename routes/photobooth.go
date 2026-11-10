package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	photoboothHandler "github.com/plugimt/transat-backend/handlers/photobooth"
)

// SetupPhotoboothRoutes exposes the booth's private API (signed with a per-booth key, never a user token) and
// the public share links. h is nil when the photobooth is not configured: the routes then do not exist.
func SetupPhotoboothRoutes(router fiber.Router, h *photoboothHandler.Handler) {
	if h == nil {
		return
	}
	perBooth := limiter.New(limiter.Config{
		Max:          600,
		Expiration:   time.Minute,
		KeyGenerator: photoboothHandler.DeviceKey,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Too many requests from this booth."})
		},
	})
	booth := router.Group("/photobooth/booth", h.Authenticate, perBooth)
	booth.Get("/ping", h.Ping)
	booth.Post("/sessions/:session/files", h.StartUpload)
	booth.Post("/sessions/:session/files/:id/complete", h.CompleteUpload)
	booth.Delete("/sessions/:session/slots/:kind/:number", h.DiscardSlot)
	booth.Post("/sessions/:session/finalize", h.Finalize)
	booth.Delete("/sessions/:session", h.Cancel)
	// Frames sent from a phone through a temporary link shown as a QR code in the booth's admin panel.
	booth.Post("/frame-sessions", h.CreateFrameSession)
	booth.Get("/frame-sessions/:id", h.FrameSessionStatus)
	booth.Get("/frame-sessions/:id/frames/:type", h.FrameImage)
	booth.Delete("/frame-sessions/:id", h.CloseFrameSession)

	// Every visitor may come through the same proxy address: limit per link, not per address.
	perShare := limiter.New(limiter.Config{
		Max:          600,
		Expiration:   time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string { return "photobooth-share-" + c.Params("token") },
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Too many requests."})
		},
	})
	router.Get("/photobooth/shares/:token", perShare, h.GetShare)
	router.Get("/photobooth/frames/:token", perShare, h.GetFrameSession)
	router.Post("/photobooth/frames/:token/:type", perShare, h.UploadFrame)
	router.Get("/photobooth/shares/:token/files/:id", perShare, h.GetSharedFile)
}
