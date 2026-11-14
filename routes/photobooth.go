package routes

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	photoboothHandler "github.com/plugimt/transat-backend/handlers/photobooth"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services/photoreport"
	"github.com/plugimt/transat-backend/utils"
)

// SetupPhotoboothRoutes exposes the booth's private API (signed with a per-booth key, never a user token) and
// the public share links. h is nil when the photobooth is not configured: the routes then do not exist.
func SetupPhotoboothRoutes(router fiber.Router, db *sql.DB, h *photoboothHandler.Handler) {
	if h == nil {
		return
	}
	// The booth's pictures reach the app only once an administrator or the Club Photo approves them. Declared
	// before the admin-only "/admin" group, whose middleware would otherwise answer first.
	review := router.Group("/admin/photobooth",
		middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail, middlewares.AdminOrRolesMiddleware(db, photoreport.PhotoManagerRoles...))
	review.Get("/pending", h.ListPending)
	review.Get("/pending/:folder/pictures", h.PendingPictures)
	review.Post("/pending/:folder/approve", h.Approve)
	review.Get("/pictures/:id", h.PendingPicture)
	review.Delete("/pictures/:id", h.DeletePendingPicture)
	// Events prepared for the booths: settings, layouts switched on, frames.
	review.Get("/layouts", h.LayoutCatalog)
	review.Get("/events", h.ListEvents)
	review.Post("/events", h.CreateEvent)
	review.Get("/events/:id", h.GetEvent)
	review.Put("/events/:id", h.UpdateEvent)
	review.Delete("/events/:id", h.DeleteEvent)
	review.Get("/events/:id/frames/:type", h.EventFrame)
	review.Put("/events/:id/frames/:type", h.UploadEventFrame)
	review.Delete("/events/:id/frames/:type", h.DeleteEventFrame)
	// Each evening's audience: sessions, QR codes scanned, album visitors.
	review.Get("/stats", h.Stats)

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
	// Events prepared in the admin panel, and the booth's layout catalog the panel draws templates from.
	booth.Get("/events", h.BoothEvents)
	booth.Get("/events/:id/frames/:type", h.EventFrame)
	booth.Put("/layouts", h.PublishLayouts)

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
	// The site's pages report each opening of a share link or a public album, for the statistics.
	router.Post("/photobooth/shares/:token/visits", perShare, h.ShareVisit)
	perAlbum := limiter.New(limiter.Config{
		Max:          600,
		Expiration:   time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string { return "album-visit-" + c.Params("id") },
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Too many requests."})
		},
	})
	router.Post("/public/albums/:id/visits", perAlbum, h.AlbumVisit)
}
