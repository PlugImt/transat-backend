package routes

import (
	"github.com/plugimt/transat-backend/handlers/association"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupAssociationRoutes(router fiber.Router, associationHandler *association.AssociationHandler) {

	associationGroup := router.Group("/association", middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail)

	associationGroup.Get("", associationHandler.GetAssociation)
	associationGroup.Get("/", associationHandler.GetAssociation)
	associationGroup.Get("/:id", associationHandler.GetAssociationByID)
	associationGroup.Get("/:id/members", associationHandler.GetAssociationMembers)
	associationGroup.Post("", associationHandler.CreateAssociation)
	associationGroup.Post("/", associationHandler.CreateAssociation)
	associationGroup.Patch("/:id", associationHandler.UpdateAssociation)
	associationGroup.Post("/:id/respo", associationHandler.AddAssociationRespo)

	associationGroup.Post("/:id/join", associationHandler.JoinAssociation)
	associationGroup.Post("/:id/leave", associationHandler.LeaveAssociation)

}
