package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	photosHandler "github.com/plugimt/transat-backend/handlers/photos"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services/photos"
	"github.com/plugimt/transat-backend/utils"
)

// SetupPhotosRoutes exposes the Drive-backed photo gallery. sync may be nil when Drive is not configured.
func SetupPhotosRoutes(router fiber.Router, db *sql.DB, sync *photos.SyncService) {
	handler := photosHandler.NewPhotosHandler(db, sync)

	// Compression first, ETag inside it: unchanged pages cost a 304 with no body.
	group := router.Group("/photos",
		middlewares.JWTMiddleware, middlewares.NewfAuthMiddleware(db), utils.EnhanceSentryEventWithEmail,
		compress.New(compress.Config{Level: compress.LevelBestSpeed}), etag.New(etag.Config{Weak: true}))

	group.Get("", handler.GetGallery)
	group.Get("/", handler.GetGallery)
	// Static segments before "/:id".
	group.Get("/albums", handler.GetAlbums)
	group.Get("/albums/:id", handler.GetAlbum)
	group.Get("/albums/:id/photos", handler.GetAlbumPhotos)
	group.Get("/:id", handler.GetPhoto)

	admin := router.Group("/admin/photos",
		middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail, middlewares.AdminAuthMiddleware(db))
	admin.Get("/sync", handler.GetSyncStatus)
	admin.Post("/sync", handler.TriggerSync)
}
