package routes

import (
	"github.com/plugimt/transat-backend/handlers/event" // Import the event handlers
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupEventRoutes(router fiber.Router, eventHandler *event.EventHandler) {

	eventGroup := router.Group("/event", middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail)

	eventGroup.Get("", eventHandler.GetEvent)
	eventGroup.Get("/", eventHandler.GetEvent)
	eventGroup.Get("/:id", eventHandler.GetEventByID)
	eventGroup.Get("/:id/members", eventHandler.GetEventMembers)
	eventGroup.Post("", eventHandler.CreateEvent)
	eventGroup.Post("/", eventHandler.CreateEvent)
	eventGroup.Patch("/:id", eventHandler.UpdateEvent)

	eventGroup.Post("/:id/join", eventHandler.JoinEvent)
	eventGroup.Post("/:id/leave", eventHandler.LeaveEvent)
}
