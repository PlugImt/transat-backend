package routes

import (
	"github.com/gofiber/fiber/v2"
	clubsHandler "github.com/plugimt/transat-backend/handlers/clubs"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"
)

func SetupClubsRoutes(router fiber.Router, h *clubsHandler.ClubsHandler) {
	// Protected routes - require JWT authentication
	clubs := router.Group("/clubs", middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail)

	// Admin-only routes
	clubs.Post("/", h.CreateClub)

	// Admin or club respo routes
	clubs.Put("/:id", h.UpdateClub)
	clubs.Post("/:id/respos", h.AddRespo)
	clubs.Delete("/:id/respos", h.RemoveRespo)

	// Member routes
	clubs.Post("/:id/join", h.JoinClub)
	clubs.Post("/:id/leave", h.LeaveClub)

	clubs.Get("/my", h.GetUserClubs)
	clubs.Get("/respos", h.GetAllClubRespos)
	clubs.Get("/:id/members", h.GetClubMembers)
	clubs.Get("/", h.GetAllClubs)
	clubs.Get("/:id", h.GetClubDetails)
}
