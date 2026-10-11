// Package photoreport exposes the photo-report endpoints: submission for students, review for administrators.
package photoreport

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	service "github.com/plugimt/transat-backend/services/photoreport"
	"github.com/plugimt/transat-backend/utils"
)

type Handler struct{ svc *service.Service }

func NewHandler(svc *service.Service) *Handler { return &Handler{svc: svc} }

type submitRequest struct {
	PhotoID string `json:"photo_id"`
	Reason  string `json:"reason"`
	Comment string `json:"comment"`
}

// Submit handles POST /photos/reports. The reporter is the authenticated user; the club and the recipients are
// worked out from the stored photo, so the body carries nothing but the photo, the reason and the comment.
func (h *Handler) Submit(c *fiber.Ctx) error {
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		return fail(c, fiber.StatusUnauthorized, "Authentication required", nil)
	}
	var body submitRequest
	if err := c.BodyParser(&body); err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	reason, err := service.ParseReason(body.Reason)
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid report reason", nil)
	}

	receipt, err := h.svc.Submit(c.Context(), service.Submission{
		PhotoID: body.PhotoID, Reason: reason, Comment: body.Comment, ReporterEmail: email,
	})
	var limited *service.RateLimitedError
	switch {
	case errors.Is(err, service.ErrInvalidPhotoID):
		return fail(c, fiber.StatusBadRequest, "Invalid photo id", nil)
	case errors.Is(err, service.ErrInvalidReason):
		return fail(c, fiber.StatusBadRequest, "Invalid report reason", nil)
	case errors.Is(err, service.ErrCommentTooLong):
		return fail(c, fiber.StatusBadRequest, "Comment too long", nil)
	case errors.Is(err, service.ErrCommentRequired):
		return fail(c, fiber.StatusBadRequest, "A comment is required for this reason", nil)
	case errors.Is(err, service.ErrPhotoNotFound):
		return fail(c, fiber.StatusNotFound, "Photo not found", nil)
	case errors.As(err, &limited):
		c.Set(fiber.HeaderRetryAfter, strconv.Itoa(int(limited.RetryAfter/time.Second)+1))
		return fail(c, fiber.StatusTooManyRequests, "Too many reports, try again later", nil)
	case err != nil:
		// The comment is user text and stays out of logs; the error does not contain it.
		return fail(c, fiber.StatusInternalServerError, "Could not record the report", err)
	}

	if !receipt.Duplicate {
		h.deliverInBackground()
	}
	status := fiber.StatusCreated
	if receipt.Duplicate {
		status = fiber.StatusOK
	}
	return c.Status(status).JSON(fiber.Map{"id": receipt.ID, "duplicate": receipt.Duplicate})
}

// deliverInBackground sends right away so the reporter gets an immediate answer; the scheduled job covers failures.
func (h *Handler) deliverInBackground() {
	go func() {
		defer func() { _ = recover() }()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if _, err := h.svc.DeliverDue(ctx); err != nil {
			utils.LogMessage(utils.LevelError, "Photo report delivery failed")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
		}
	}()
}

// List handles GET /admin/photos/reports.
func (h *Handler) List(c *fiber.Ctx) error {
	filter := service.ListFilter{
		Reason:   c.Query("reason"),
		Delivery: c.Query("delivery"),
		Limit:    c.QueryInt("limit", 25),
		Offset:   max(c.QueryInt("offset", 0), 0),
	}
	if filter.Reason != "" {
		if _, err := service.ParseReason(filter.Reason); err != nil {
			return fail(c, fiber.StatusBadRequest, "Invalid report reason", nil)
		}
	}
	switch filter.Delivery {
	case "", "sent", "pending", "failed", "no_recipients":
	default:
		return fail(c, fiber.StatusBadRequest, "Invalid delivery status", nil)
	}
	reports, total, err := h.svc.List(c.Context(), filter)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Could not load reports", err)
	}
	return c.JSON(fiber.Map{"reports": reports, "total": total})
}

// Get handles GET /admin/photos/reports/:id.
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id < 1 {
		return fail(c, fiber.StatusBadRequest, "Invalid report id", nil)
	}
	report, err := h.svc.Get(c.Context(), int64(id))
	if errors.Is(err, service.ErrReportNotFound) {
		return fail(c, fiber.StatusNotFound, "Report not found", nil)
	}
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Could not load the report", err)
	}
	return c.JSON(report)
}

// Retry handles POST /admin/photos/reports/:id/retry.
func (h *Handler) Retry(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id < 1 {
		return fail(c, fiber.StatusBadRequest, "Invalid report id", nil)
	}
	if _, err := h.svc.Get(c.Context(), int64(id)); errors.Is(err, service.ErrReportNotFound) {
		return fail(c, fiber.StatusNotFound, "Report not found", nil)
	}
	n, err := h.svc.RetryFailed(c.Context(), int64(id))
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Could not retry the report", err)
	}
	return c.JSON(fiber.Map{"retried": n})
}

type albumClubRequest struct {
	ClubID *int64 `json:"club_id"`
}

// SetAlbumClub handles PUT /admin/photos/albums/:id/club (club_id null removes the link).
func (h *Handler) SetAlbumClub(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id < 1 {
		return fail(c, fiber.StatusBadRequest, "Invalid album id", nil)
	}
	var body albumClubRequest
	if err := c.BodyParser(&body); err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	switch err := h.svc.SetAlbumClub(c.Context(), int64(id), body.ClubID); {
	case errors.Is(err, service.ErrAlbumNotFound):
		return fail(c, fiber.StatusNotFound, "Album not found", nil)
	case errors.Is(err, service.ErrClubNotFound):
		return fail(c, fiber.StatusBadRequest, "Club not found", nil)
	case err != nil:
		return fail(c, fiber.StatusInternalServerError, "Could not update the album", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func fail(c *fiber.Ctx, status int, msg string, err error) error {
	if err != nil {
		utils.LogMessage(utils.LevelError, msg)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}
	return c.Status(status).JSON(fiber.Map{"error": msg})
}
