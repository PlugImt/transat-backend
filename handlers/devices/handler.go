// Package devices exposes the device-registration endpoints used by the mobile app.
package devices

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	service "github.com/plugimt/transat-backend/services/devices"
	"github.com/plugimt/transat-backend/utils"
)

type Handler struct{ svc *service.Service }

func NewHandler(svc *service.Service) *Handler { return &Handler{svc: svc} }

type registerRequest struct {
	Platform string `json:"platform"`
	// PushToken is null or empty when the installation has none.
	PushToken        *string `json:"push_token"`
	PermissionStatus string  `json:"permission_status"`
	AppVersion       string  `json:"app_version"`
}

// Register handles PUT /devices/:installation_id. It is both the registration and the heartbeat: the owner is
// the authenticated user, whatever the installation was bound to before.
func (h *Handler) Register(c *fiber.Ctx) error {
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		return fail(c, fiber.StatusUnauthorized, "Authentication required", nil)
	}
	var body registerRequest
	if err := c.BodyParser(&body); err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	reg := service.Registration{
		InstallationID:   c.Params("installation_id"),
		Email:            email,
		Platform:         body.Platform,
		PermissionStatus: body.PermissionStatus,
		AppVersion:       body.AppVersion,
	}
	if body.PushToken != nil {
		reg.PushToken = *body.PushToken
	}

	_, err := h.svc.Register(c.Context(), reg)
	switch {
	case errors.Is(err, service.ErrInvalidInstallationID):
		return fail(c, fiber.StatusBadRequest, "Invalid installation id", nil)
	case errors.Is(err, service.ErrInvalidPlatform):
		return fail(c, fiber.StatusBadRequest, "Invalid platform", nil)
	case errors.Is(err, service.ErrInvalidPushToken):
		return fail(c, fiber.StatusBadRequest, "Invalid push token", nil)
	case errors.Is(err, service.ErrInvalidPermission):
		return fail(c, fiber.StatusBadRequest, "Invalid permission status", nil)
	case errors.Is(err, service.ErrInvalidAppVersion):
		return fail(c, fiber.StatusBadRequest, "Invalid app version", nil)
	case err != nil:
		// Never answer 200 for a registration that was not stored: the app would stop retrying.
		return fail(c, fiber.StatusInternalServerError, "Could not register the device", err)
	}

	cfg := h.svc.Config()
	return c.JSON(fiber.Map{
		"installation_id":            reg.InstallationID,
		"heartbeat_interval_seconds": int(cfg.HeartbeatInterval.Seconds()),
		"presence_window_seconds":    int(cfg.PresenceWindow.Seconds()),
	})
}

// List handles GET /devices: the devices of the authenticated user and whether the user is connected.
func (h *Handler) List(c *fiber.Ctx) error {
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		return fail(c, fiber.StatusUnauthorized, "Authentication required", nil)
	}
	devices, err := h.svc.List(c.Context(), email)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Could not load the devices", err)
	}
	return c.JSON(fiber.Map{
		"devices":            devices,
		"connected":          service.AnyConnected(devices),
		"has_active_session": service.AnyActiveSession(devices),
	})
}

// SignOut handles POST /devices/:installation_id/sign-out. It needs no token: it runs when the session is already
// gone (expired, revoked, offline retry), and an installation id is only known to its own app and its owner.
func (h *Handler) SignOut(c *fiber.Ctx) error {
	_, err := h.svc.SignOut(c.Context(), c.Params("installation_id"))
	switch {
	case errors.Is(err, service.ErrInvalidInstallationID):
		return fail(c, fiber.StatusBadRequest, "Invalid installation id", nil)
	case err != nil:
		return fail(c, fiber.StatusInternalServerError, "Could not sign the device out", err)
	}
	// The same answer for an unknown, already disabled or just disabled installation.
	return c.SendStatus(fiber.StatusNoContent)
}

func fail(c *fiber.Ctx, status int, msg string, err error) error {
	if err != nil {
		utils.LogMessage(utils.LevelError, msg)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}
	return c.Status(status).JSON(fiber.Map{"error": msg})
}
