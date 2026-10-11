package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	reportHandler "github.com/plugimt/transat-backend/handlers/photoreport"
	photosHandler "github.com/plugimt/transat-backend/handlers/photos"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services/photos"
	"github.com/plugimt/transat-backend/utils"
)

// SetupPhotosRoutes exposes the Drive-backed photo gallery. sync may be nil when Drive is not configured.
func SetupPhotosRoutes(router fiber.Router, db *sql.DB, sync *photos.SyncService, reports *reportHandler.Handler) {
	handler := photosHandler.NewPhotosHandler(db, sync)

	// Compression first, ETag inside it: unchanged pages cost a 304 with no body.
	group := router.Group("/photos",
		middlewares.JWTMiddleware, middlewares.NewfAuthMiddleware(db), utils.EnhanceSentryEventWithEmail,
		compress.New(compress.Config{Level: compress.LevelBestSpeed}), etag.New(etag.Config{Weak: true}))

	group.Get("", handler.GetGallery)
	group.Get("/", handler.GetGallery)
	// Static segments before "/:id".
	group.Get("/timeline", handler.GetTimeline)
	group.Get("/albums", handler.GetAlbums)
	group.Get("/albums/:id", handler.GetAlbum)
	group.Get("/albums/:id/photos", handler.GetAlbumPhotos)
	// Reporting needs the same student access as viewing; the reporter is the token's user.
	group.Post("/reports", reports.Submit)
	group.Get("/:id", handler.GetPhoto)

	admin := router.Group("/admin/photos",
		middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail, middlewares.AdminAuthMiddleware(db))
	admin.Get("/sync", handler.GetSyncStatus)
	admin.Post("/sync", handler.TriggerSync)
	admin.Get("/reports", reports.List)
	admin.Get("/reports/:id", reports.Get)
	admin.Post("/reports/:id/retry", reports.Retry)
	admin.Put("/albums/:id/club", reports.SetAlbumClub)
}
