package routes

import (
	"github.com/plugimt/transat-backend/handlers/club" // Import the club handlers
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupClubRoutes(router fiber.Router, clubHandler *club.ClubHandler) {

	clubGroup := router.Group("/club", middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail)

	clubGroup.Get("", clubHandler.GetClub)
	clubGroup.Get("/", clubHandler.GetClub)
	clubGroup.Get("/:id", clubHandler.GetClubByID)
	clubGroup.Get("/:id/members", clubHandler.GetClubMembers)
	clubGroup.Post("", clubHandler.CreateClub)
	clubGroup.Post("/", clubHandler.CreateClub)
	clubGroup.Patch("/:id", clubHandler.UpdateClub)
	clubGroup.Post("/:id/respo", clubHandler.AddClubRespo)

	clubGroup.Post("/:id/join", clubHandler.JoinClub)
	clubGroup.Post("/:id/leave", clubHandler.LeaveClub)

}
