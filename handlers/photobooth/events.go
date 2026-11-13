package photobooth

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/gofiber/fiber/v2"
	pb "github.com/plugimt/transat-backend/services/photobooth"
)

// Events prepared in the admin panel, and what the booths download of them.

// maxEventJSON fits an event's settings and the list of its layouts.
const maxEventJSON = 16 << 10

func eventID(c *fiber.Ctx) (int64, bool) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	return id, err == nil && id > 0
}

func notFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
}

// parseSettings reads the body over the booth's defaults; on false the error answer is already written.
func parseSettings(c *fiber.Ctx) (pb.EventSettings, bool) {
	in := pb.DefaultEventSettings()
	if len(c.Body()) > maxEventJSON {
		_ = c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{"error": "body too large"})
		return in, false
	}
	if err := json.Unmarshal(c.Body(), &in); err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
		return in, false
	}
	return in, true
}

// ---- Admin panel ----------------------------------------------------------------------------------------

// LayoutCatalog handles GET /admin/photobooth/layouts: the layouts the booths offer, per photo ratio.
func (h *Handler) LayoutCatalog(c *fiber.Ctx) error {
	catalog, err := h.svc.LayoutCatalog(c.UserContext())
	if err != nil {
		return fail(c, err)
	}
	return c.JSON(catalog)
}

// ListEvents handles GET /admin/photobooth/events.
func (h *Handler) ListEvents(c *fiber.Ctx) error {
	list, err := h.svc.Events(c.UserContext())
	if err != nil {
		return fail(c, err)
	}
	return c.JSON(list)
}

// GetEvent handles GET /admin/photobooth/events/:id.
func (h *Handler) GetEvent(c *fiber.Ctx) error {
	id, ok := eventID(c)
	if !ok {
		return notFound(c)
	}
	e, err := h.svc.Event(c.UserContext(), id)
	if err != nil {
		return fail(c, err)
	}
	return c.JSON(e)
}

// CreateEvent handles POST /admin/photobooth/events; omitted settings take the booth's defaults.
func (h *Handler) CreateEvent(c *fiber.Ctx) error {
	in, ok := parseSettings(c)
	if !ok {
		return nil
	}
	e, err := h.svc.CreateEvent(c.UserContext(), reviewer(c), in)
	if err != nil {
		return fail(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(e)
}

// UpdateEvent handles PUT /admin/photobooth/events/:id with every setting.
func (h *Handler) UpdateEvent(c *fiber.Ctx) error {
	id, ok := eventID(c)
	if !ok {
		return notFound(c)
	}
	in, ok := parseSettings(c)
	if !ok {
		return nil
	}
	e, err := h.svc.UpdateEvent(c.UserContext(), id, in)
	if err != nil {
		return fail(c, err)
	}
	return c.JSON(e)
}

// DeleteEvent handles DELETE /admin/photobooth/events/:id.
func (h *Handler) DeleteEvent(c *fiber.Ctx) error {
	id, ok := eventID(c)
	if !ok {
		return notFound(c)
	}
	if err := h.svc.DeleteEvent(c.UserContext(), id); err != nil {
		return fail(c, err)
	}
	return c.JSON(fiber.Map{"ok": true})
}

// UploadEventFrame handles PUT /admin/photobooth/events/:id/frames/:type (multipart, field "file").
func (h *Handler) UploadEventFrame(c *fiber.Ctx) error {
	id, ok := eventID(c)
	if !ok {
		return notFound(c)
	}
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
	frame, err := h.svc.SetEventFrame(c.UserContext(), id, c.Params("type"), data)
	if err != nil {
		return fail(c, err)
	}
	return c.JSON(frame)
}

// EventFrame handles GET /admin/photobooth/events/:id/frames/:type, and the booths' same route.
func (h *Handler) EventFrame(c *fiber.Ctx) error {
	id, ok := eventID(c)
	if !ok {
		return notFound(c)
	}
	png, err := h.svc.EventFrameImage(c.UserContext(), id, c.Params("type"))
	if err != nil {
		return fail(c, err)
	}
	c.Set(fiber.HeaderContentType, "image/png")
	c.Set(fiber.HeaderCacheControl, "private, no-cache")
	return c.Send(png)
}

// DeleteEventFrame handles DELETE /admin/photobooth/events/:id/frames/:type.
func (h *Handler) DeleteEventFrame(c *fiber.Ctx) error {
	id, ok := eventID(c)
	if !ok {
		return notFound(c)
	}
	if err := h.svc.DeleteEventFrame(c.UserContext(), id, c.Params("type")); err != nil {
		return fail(c, err)
	}
	return c.JSON(fiber.Map{"ok": true})
}

// ---- Booth ----------------------------------------------------------------------------------------------

// BoothEvents handles GET /photobooth/booth/events: recent and upcoming events with their settings and layouts.
func (h *Handler) BoothEvents(c *fiber.Ctx) error {
	list, err := h.svc.BoothEvents(c.UserContext())
	if err != nil {
		return fail(c, err)
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.JSON(fiber.Map{"events": list})
}

// PublishLayouts handles PUT /photobooth/booth/layouts: the booth's layout catalog, {"aspects": {"RATIO_4_3": [...]}}.
func (h *Handler) PublishLayouts(c *fiber.Ctx) error {
	if len(c.Body()) > pb.MaxCatalogJSON {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{"error": "body too large"})
	}
	var b struct {
		Aspects map[string][]pb.CatalogLayout `json:"aspects"`
	}
	if err := json.Unmarshal(c.Body(), &b); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
	}
	if err := h.svc.PublishLayoutCatalog(c.UserContext(), device(c), b.Aspects); err != nil {
		return fail(c, err)
	}
	return c.JSON(fiber.Map{"ok": true})
}
