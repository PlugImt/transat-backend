package photobooth

import (
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	pb "github.com/plugimt/transat-backend/services/photobooth"
)

// maxFrameSessionJSON fits the geometry of every layout of an event.
const maxFrameSessionJSON = 64 << 10

type frameSessionBody struct {
	EventName string `json:"event_name"`
	// EventID is the admin-panel event the booth shows, when it came from there.
	EventID int64            `json:"event_id"`
	Layouts []pb.FrameLayout `json:"layouts"`
}

// CreateFrameSession handles POST /photobooth/booth/frame-sessions: a link for the QR code of the admin panel.
func (h *Handler) CreateFrameSession(c *fiber.Ctx) error {
	if len(c.Body()) > maxFrameSessionJSON {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{"error": "body too large"})
	}
	var b frameSessionBody
	if err := json.Unmarshal(c.Body(), &b); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
	}
	session, err := h.svc.CreateFrameSession(c.UserContext(), device(c), b.EventName, b.EventID, b.Layouts)
	if err != nil {
		return fail(c, err)
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.Status(fiber.StatusCreated).JSON(session)
}

func frameSessionID(c *fiber.Ctx) (int64, bool) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	return id, err == nil && id > 0
}

// FrameSessionStatus handles GET /photobooth/booth/frame-sessions/:id: what was uploaded so far.
func (h *Handler) FrameSessionStatus(c *fiber.Ctx) error {
	id, ok := frameSessionID(c)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	uploads, expires, err := h.svc.FrameUploads(c.UserContext(), device(c), id)
	if err != nil {
		return fail(c, err)
	}
	return c.JSON(fiber.Map{"uploads": uploads, "expires_at": expires, "expired": time.Now().After(expires)})
}

// FrameImage handles GET /photobooth/booth/frame-sessions/:id/frames/:type: the PNG to apply.
func (h *Handler) FrameImage(c *fiber.Ctx) error {
	id, ok := frameSessionID(c)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	png, err := h.svc.FrameImage(c.UserContext(), device(c), id, c.Params("type"))
	if err != nil {
		return fail(c, err)
	}
	c.Set(fiber.HeaderContentType, "image/png")
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.Send(png)
}

// CloseFrameSession handles DELETE /photobooth/booth/frame-sessions/:id.
func (h *Handler) CloseFrameSession(c *fiber.Ctx) error {
	id, ok := frameSessionID(c)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	if err := h.svc.CloseFrameSession(c.UserContext(), device(c), id); err != nil {
		return fail(c, err)
	}
	return c.JSON(fiber.Map{"ok": true})
}

// GetFrameSession handles GET /photobooth/frames/:token: the upload page's layouts and what was already sent.
func (h *Handler) GetFrameSession(c *fiber.Ctx) error {
	s, err := h.svc.PublicFrameSession(c.UserContext(), c.Params("token"))
	if err != nil {
		return fail(c, err)
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	c.Set("X-Robots-Tag", "noindex, nofollow")
	return c.JSON(s)
}

// UploadFrame handles POST /photobooth/frames/:token/:type (multipart, field "file"): the frame of one layout.
func (h *Handler) UploadFrame(c *fiber.Ctx) error {
	header, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "send the PNG in the \"file\" field"})
	}
	if header.Size > pb.MaxFrameBytes {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{"error": "image too large"})
	}
	f, err := header.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unreadable file"})
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, pb.MaxFrameBytes+1))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unreadable file"})
	}
	u, err := h.svc.UploadFrame(c.UserContext(), c.Params("token"), c.Params("type"), data)
	if err != nil {
		return fail(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(u)
}
