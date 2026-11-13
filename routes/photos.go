package routes

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	reportHandler "github.com/plugimt/transat-backend/handlers/photoreport"
	photosHandler "github.com/plugimt/transat-backend/handlers/photos"
	"github.com/plugimt/transat-backend/middlewares"
	reportService "github.com/plugimt/transat-backend/services/photoreport"
	"github.com/plugimt/transat-backend/services/photos"
	"github.com/plugimt/transat-backend/utils"
)

// SetupPhotosRoutes exposes the Drive-backed photo gallery. sync may be nil when Drive is not configured.
// media may be nil too: the public album pages then answer 503 for every image.
func SetupPhotosRoutes(router fiber.Router, db *sql.DB, sync *photos.SyncService, media photos.MediaSource, reports *reportHandler.Handler) {
	handler := photosHandler.NewPhotosHandler(db, sync)
	if media != nil {
		handler.WithMedia(media)
	}

	// Public album pages on transat.dev: no account, only shown albums and their images. The site's server proxies
	// every call, so all visitors share its address: limit per album, not per address.
	perAlbum := limiter.New(limiter.Config{
		Max:          3000,
		Expiration:   time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string { return "public-album-" + c.Params("id") },
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Too many requests."})
		},
	})
	public := router.Group("/public/albums", compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	public.Get("/:id", perAlbum, etag.New(etag.Config{Weak: true}), handler.GetPublicAlbum)
	public.Get("/:id/photos", perAlbum, etag.New(etag.Config{Weak: true}), handler.GetPublicAlbumPhotos)
	public.Get("/:id/files", perAlbum, etag.New(etag.Config{Weak: true}), handler.GetPublicAlbumFiles)
	public.Get("/:id/cover", perAlbum, handler.GetPublicCover)
	public.Get("/:id/photos/:photo/:variant", perAlbum, handler.GetPublicPhoto)

	// Album management: administrators and the Club Photo. Declared before the admin-only group below.
	albums := router.Group("/admin/photos/albums",
		middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail, middlewares.AdminOrRolesMiddleware(db, reportService.PhotoManagerRoles...))
	albums.Get("", handler.ListAdminAlbums)
	albums.Put("/:id/hidden", handler.SetAlbumHidden)
	albums.Get("/:id/cover", handler.GetAdminAlbumCover)

	// Compression first, ETag inside it: unchanged pages cost a 304 with no body.
	group := router.Group("/photos",
		middlewares.JWTMiddleware, middlewares.NewfAuthMiddleware(db), utils.EnhanceSentryEventWithEmail,
		compress.New(compress.Config{Level: compress.LevelBestSpeed}), etag.New(etag.Config{Weak: true}))

	group.Get("", handler.GetGallery)
	group.Get("/", handler.GetGallery)
	// Static segments before "/:id".
	group.Get("/me", handler.GetMyPhotos)
	group.Get("/albums", handler.GetAlbums)
	group.Get("/albums/:id", handler.GetAlbum)
	group.Get("/albums/:id/photos", handler.GetAlbumPhotos)
	// Reporting needs the same student access as viewing; the reporter is the token's user.
	group.Post("/reports", reports.Submit)
	group.Get("/:id", handler.GetPhoto)

	// Reports are reviewed by administrators and by the Club Photo's owners. Declared before the admin-only
	// group, whose middleware would otherwise answer first.
	review := router.Group("/admin/photos/reports",
		middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail, middlewares.AdminOrRolesMiddleware(db, reportService.PhotoManagerRoles...))
	review.Get("", reports.List)
	review.Post("/approve-all", reports.ApproveAll)
	review.Get("/:id", reports.Get)
	review.Post("/:id/decision", reports.Decide)
	review.Post("/:id/retry", reports.Retry)

	admin := router.Group("/admin/photos",
		middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail, middlewares.AdminAuthMiddleware(db))
	admin.Get("/sync", handler.GetSyncStatus)
	admin.Post("/sync", handler.TriggerSync)
	admin.Put("/albums/:id/club", reports.SetAlbumClub)
}
