// Package photobooth exposes the booth's signed upload API and the public share pages.
package photobooth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	pb "github.com/plugimt/transat-backend/services/photobooth"
)

const (
	localDevice   = "photobooth_device"
	maxJSONBytes  = 4 << 10
	previewPixels = 1600
	streamTimeout = 10 * time.Minute
)

type Handler struct {
	svc      *pb.Service
	verifier *pb.Verifier
}

func NewHandler(svc *pb.Service, verifier *pb.Verifier) *Handler {
	return &Handler{svc: svc, verifier: verifier}
}

// ---- Booth authentication -------------------------------------------------------------------------------

// Authenticate accepts only requests signed with a booth key (see services/photobooth/auth.go).
func (h *Handler) Authenticate(c *fiber.Ctx) error {
	path := c.OriginalURL()
	if i := strings.IndexByte(path, '?'); i >= 0 {
		path = path[:i]
	}
	device, err := h.verifier.Verify(c.Method(), path, c.Get(pb.HeaderKey), c.Get(pb.HeaderTimestamp),
		c.Get(pb.HeaderNonce), c.Get(pb.HeaderSignature), c.Body())
	if err != nil {
		// The booth corrects its clock from this header and signs again.
		c.Set("X-Photobooth-Server-Time", strconv.FormatInt(time.Now().Unix(), 10))
		code := "unauthorized"
		if errors.Is(err, pb.ErrClockSkew) {
			code = "clock_skew"
		}
		log.Printf("Photobooth: rejected %s %s (%s) key=%q", c.Method(), path, code, truncate(c.Get(pb.HeaderKey), 32))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": code})
	}
	c.Locals(localDevice, device)
	return c.Next()
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func device(c *fiber.Ctx) string {
	d, _ := c.Locals(localDevice).(string)
	return d
}

// DeviceKey keys the per-booth rate limiter (only used after Authenticate).
func DeviceKey(c *fiber.Ctx) string { return "photobooth-" + device(c) }

// ---- Booth routes ---------------------------------------------------------------------------------------

// Ping lets the booth's admin panel check its key and the clock.
func (h *Handler) Ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"ok": true, "device": device(c), "server_time": time.Now().Unix()})
}

type startBody struct {
	Kind       string `json:"kind"`
	Number     int    `json:"number"`
	EventName  string `json:"event_name"`
	EventDate  string `json:"event_date"`
	LayoutType string `json:"layout_type"`
	Size       int64  `json:"size"`
	MD5        string `json:"md5"`
	MimeType   string `json:"mime_type"`
}

func (h *Handler) StartUpload(c *fiber.Ctx) error {
	var b startBody
	if err := parse(c, &b); err != nil {
		return err
	}
	res, err := h.svc.StartUpload(c.UserContext(), device(c), c.Params("session"), pb.StartRequest{
		Event: pb.Event{Name: b.EventName, Date: b.EventDate, LayoutType: b.LayoutType},
		Kind:  b.Kind, Number: b.Number, Size: b.Size, MD5: b.MD5, MimeType: b.MimeType,
	})
	if err != nil {
		return fail(c, err)
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *Handler) CompleteUpload(c *fiber.Ctx) error {
	var b struct {
		DriveFileID string `json:"drive_file_id"`
	}
	if err := parse(c, &b); err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid file id"})
	}
	if err := h.svc.CompleteUpload(c.UserContext(), device(c), c.Params("session"), id, b.DriveFileID); err != nil {
		return fail(c, err)
	}
	return c.JSON(fiber.Map{"ok": true})
}

func (h *Handler) DiscardSlot(c *fiber.Ctx) error {
	n, err := strconv.Atoi(c.Params("number"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid number"})
	}
	if err := h.svc.DiscardSlot(c.UserContext(), device(c), c.Params("session"), c.Params("kind"), n); err != nil {
		return fail(c, err)
	}
	return c.JSON(fiber.Map{"ok": true})
}

func (h *Handler) Cancel(c *fiber.Ctx) error {
	if err := h.svc.Cancel(c.UserContext(), device(c), c.Params("session")); err != nil {
		return fail(c, err)
	}
	return c.JSON(fiber.Map{"ok": true})
}

func (h *Handler) Finalize(c *fiber.Ctx) error {
	var b struct {
		PhotoCount int `json:"photo_count"`
	}
	if err := parse(c, &b); err != nil {
		return err
	}
	link, err := h.svc.Finalize(c.UserContext(), device(c), c.Params("session"), b.PhotoCount)
	if err != nil {
		return fail(c, err)
	}
	return c.JSON(fiber.Map{"url": link})
}

// ---- Public share ---------------------------------------------------------------------------------------

func (h *Handler) GetShare(c *fiber.Ctx) error {
	share, err := h.svc.GetShare(c.UserContext(), c.Params("token"))
	if err != nil {
		return fail(c, err)
	}
	c.Set("X-Robots-Tag", "noindex, nofollow")
	c.Set(fiber.HeaderCacheControl, "private, max-age=60")
	return c.JSON(share)
}

// GetSharedFile streams one picture: ?size=preview for a light version, ?download=1 to save the original.
func (h *Handler) GetSharedFile(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	preview := 0
	if c.Query("size") == "preview" {
		preview = previewPixels
	}
	// The body is read after this handler returns: its context lives until the stream is closed.
	ctx, cancel := context.WithTimeout(context.Background(), streamTimeout)
	file, err := h.svc.OpenSharedFile(ctx, c.Params("token"), id, preview)
	if err != nil {
		cancel()
		return fail(c, err)
	}
	disposition := "inline"
	if c.Query("download") == "1" {
		disposition = "attachment"
	}
	c.Set(fiber.HeaderContentType, file.ContentType)
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf(`%s; filename="%s"; filename*=UTF-8''%s`,
		disposition, asciiName(file.Name), url.PathEscape(file.Name)))
	c.Set(fiber.HeaderCacheControl, "public, max-age=604800, immutable")
	c.Set("X-Robots-Tag", "noindex, nofollow")
	c.Set("X-Content-Type-Options", "nosniff")
	body := &cancelOnClose{ReadCloser: file.Body, cancel: cancel}
	if file.Size >= 0 {
		return c.SendStream(body, int(file.Size))
	}
	return c.SendStream(body)
}

type cancelOnClose struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (r *cancelOnClose) Close() error {
	err := r.ReadCloser.Close()
	r.cancel()
	return err
}

func asciiName(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r < 0x20 || r > 0x7e || r == '"' || r == '\\' {
			b.WriteByte('_')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ---- helpers --------------------------------------------------------------------------------------------

func parse(c *fiber.Ctx, v any) error {
	if len(c.Body()) > maxJSONBytes {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{"error": "body too large"})
	}
	if err := c.BodyParser(v); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
	}
	return nil
}

func fail(c *fiber.Ctx, err error) error {
	var e *pb.Error
	if errors.As(err, &e) {
		return c.Status(e.Status).JSON(fiber.Map{"error": e.Message})
	}
	log.Printf("Photobooth: %s %s failed: %v", c.Method(), c.Path(), err)
	return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "storage temporarily unavailable"})
}
