// Package faces exposes the face recognition opt-in to the app: consent, reference selfies, "not me" and
// withdrawal. The album itself is GET /photos/me.
package faces

import (
	"errors"
	"io"

	"github.com/gofiber/fiber/v2"
	svc "github.com/plugimt/transat-backend/services/faces"
	"github.com/plugimt/transat-backend/utils"
)

// MaxSelfieBytes bounds a reference selfie. The request body limit keeps it in memory: it never reaches the disk.
const MaxSelfieBytes = 10 << 20

type Handler struct {
	svc *svc.Service
}

func NewHandler(s *svc.Service) *Handler { return &Handler{svc: s} }

func email(c *fiber.Ctx) string {
	e, _ := c.Locals("email").(string)
	return e
}

// errorBody gives the app a stable code to pick its own translated message.
func fail(c *fiber.Ctx, err error) error {
	type out struct {
		status int
		code   string
	}
	known := map[error]out{
		svc.ErrUnavailable:    {fiber.StatusServiceUnavailable, "unavailable"},
		svc.ErrNoConsent:      {fiber.StatusForbidden, "consent_required"},
		svc.ErrConsentVersion: {fiber.StatusConflict, "consent_outdated"},
		svc.ErrNoFace:         {fiber.StatusUnprocessableEntity, "no_face"},
		svc.ErrManyFaces:      {fiber.StatusUnprocessableEntity, "multiple_faces"},
		svc.ErrFaceTooSmall:   {fiber.StatusUnprocessableEntity, "face_too_small"},
		svc.ErrTooManyRefs:    {fiber.StatusConflict, "too_many_references"},
		svc.ErrNotFound:       {fiber.StatusNotFound, "not_found"},
	}
	for e, o := range known {
		if errors.Is(err, e) {
			return c.Status(o.status).JSON(fiber.Map{"error": err.Error(), "code": o.code})
		}
	}
	utils.LogMessage(utils.LevelError, "Face recognition request failed")
	utils.LogLineKeyValue(utils.LevelError, "Error", err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error", "code": "internal"})
}

// GetStatus handles GET /faces.
func (h *Handler) GetStatus(c *fiber.Ctx) error {
	st, err := h.svc.Status(c.UserContext(), email(c))
	if err != nil {
		return fail(c, err)
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.JSON(st)
}

type consentBody struct {
	Version  string `json:"version"`
	Accepted bool   `json:"accepted"`
}

// PostConsent handles POST /faces/consent {"version": "...", "accepted": true}: the explicit opt-in.
func (h *Handler) PostConsent(c *fiber.Ctx) error {
	var b consentBody
	if err := c.BodyParser(&b); err != nil || !b.Accepted {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "explicit consent is required", "code": "consent_required"})
	}
	if err := h.svc.Consent(c.UserContext(), email(c), b.Version); err != nil {
		return fail(c, err)
	}
	return h.GetStatus(c)
}

// PostReference handles POST /faces/references (multipart, field "photo"): a selfie, analysed then discarded.
func (h *Handler) PostReference(c *fiber.Ctx) error {
	header, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": `send the selfie in the "photo" field`, "code": "bad_request"})
	}
	if header.Size > MaxSelfieBytes {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{"error": "image too large", "code": "too_large"})
	}
	f, err := header.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unreadable file", "code": "bad_request"})
	}
	data, err := io.ReadAll(io.LimitReader(f, MaxSelfieBytes+1))
	_ = f.Close()
	if err != nil || len(data) == 0 || len(data) > MaxSelfieBytes {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unreadable file", "code": "bad_request"})
	}
	if err := h.svc.AddReference(c.UserContext(), email(c), data); err != nil {
		return fail(c, err)
	}
	c.Status(fiber.StatusCreated)
	return h.GetStatus(c)
}

// Delete handles DELETE /faces: withdraws consent and erases every face data of the user.
func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.svc.Withdraw(c.UserContext(), email(c)); err != nil {
		return fail(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// RejectPhoto handles POST /faces/photos/:id/reject ("this is not me").
func (h *Handler) RejectPhoto(c *fiber.Ctx) error {
	if err := h.svc.Reject(c.UserContext(), email(c), c.Params("id")); err != nil {
		return fail(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
