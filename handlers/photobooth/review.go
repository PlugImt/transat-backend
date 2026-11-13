package photobooth

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Review of the booth's pictures before they reach the app, by administrators and the Club Photo.

const thumbPixels = 480

func reviewer(c *fiber.Ctx) string {
	email, _ := c.Locals("email").(string)
	return email
}

// ListPending answers the events whose photobooth pictures await approval.
func (h *Handler) ListPending(c *fiber.Ctx) error {
	list, err := h.svc.ListPending(c.UserContext())
	if err != nil {
		return fail(c, err)
	}
	return c.JSON(list)
}

// PendingPictures answers the pictures of one pending folder.
func (h *Handler) PendingPictures(c *fiber.Ctx) error {
	pics, err := h.svc.PendingPictures(c.UserContext(), c.Params("folder"))
	if err != nil {
		return fail(c, err)
	}
	return c.JSON(pics)
}

// Approve publishes a pending folder: the gallery indexes it on its next sync.
func (h *Handler) Approve(c *fiber.Ctx) error {
	res, err := h.svc.Approve(c.UserContext(), c.Params("folder"), reviewer(c))
	if err != nil {
		return fail(c, err)
	}
	return c.JSON(res)
}

// DeletePendingPicture removes one picture before approval.
func (h *Handler) DeletePendingPicture(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	if err := h.svc.DeletePendingPicture(c.UserContext(), id, reviewer(c)); err != nil {
		return fail(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// PendingPicture streams a picture under review: ?size=thumb, ?size=preview or the original.
func (h *Handler) PendingPicture(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	size := 0
	switch c.Query("size") {
	case "thumb":
		size = thumbPixels
	case "preview":
		size = previewPixels
	}
	ctx, cancel := context.WithTimeout(context.Background(), streamTimeout)
	file, err := h.svc.OpenPendingPicture(ctx, id, size)
	if err != nil {
		cancel()
		return fail(c, err)
	}
	c.Set(fiber.HeaderContentType, file.ContentType)
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf(`inline; filename="%s"; filename*=UTF-8''%s`,
		asciiName(file.Name), url.PathEscape(file.Name)))
	c.Set(fiber.HeaderCacheControl, "private, max-age=3600")
	c.Set("X-Content-Type-Options", "nosniff")
	body := &cancelOnClose{ReadCloser: file.Body, cancel: cancel}
	if file.Size >= 0 {
		return c.SendStream(body, int(file.Size))
	}
	return c.SendStream(body)
}
