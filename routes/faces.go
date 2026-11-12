package routes

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	facesHandler "github.com/plugimt/transat-backend/handlers/faces"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services/faces"
	"github.com/plugimt/transat-backend/utils"
)

// SetupFacesRoutes exposes the face recognition opt-in. The album is GET /photos/me.
func SetupFacesRoutes(router fiber.Router, db *sql.DB, svc *faces.Service) {
	h := facesHandler.NewHandler(svc)
	group := router.Group("/faces", middlewares.JWTMiddleware, middlewares.NewfAuthMiddleware(db), utils.EnhanceSentryEventWithEmail)

	group.Get("", h.GetStatus)
	group.Post("/consent", h.PostConsent)
	// Each selfie costs a model run: a few tries per hour are plenty.
	group.Post("/references", limiter.New(limiter.Config{
		Max:          10,
		Expiration:   time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string { e, _ := c.Locals("email").(string); return "faces-ref:" + e },
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "too many attempts", "code": "rate_limited"})
		},
	}), h.PostReference)
	group.Delete("", h.Delete)
	group.Post("/photos/:id/reject", h.RejectPhoto)
}
