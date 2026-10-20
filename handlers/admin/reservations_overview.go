package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/reservations"
	"github.com/plugimt/transat-backend/utils"
)

// ReservationsAdminHandler serves the reservations part of the administration panel: the catalogue with what is
// happening to each item, every reservation, the manual cancellation and return, and the waiting lists.
// Its routes are mounted behind the administrator check.
type ReservationsAdminHandler struct {
	db      *sql.DB
	service *reservations.Service
	discord *services.DiscordService
}

func NewReservationsAdminHandler(db *sql.DB, service *reservations.Service, discord *services.DiscordService) *ReservationsAdminHandler {
	return &ReservationsAdminHandler{db: db, service: service, discord: discord}
}

func badRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": message})
}

func internalError(c *fiber.Ctx, what string, err error) error {
	utils.LogMessage(utils.LevelError, fmt.Sprintf("%s: %v", what, err))
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": what})
}

// parseInstant reads an instant from the panel: a day, an RFC 3339 time, or the API's UTC wall time.
func parseInstant(raw string) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	for _, layout := range []string{time.RFC3339, reservations.APITimeLayout, "2006-01-02"} {
		if t, err := time.Parse(layout, raw); err == nil {
			t = t.UTC()
			return &t, nil
		}
	}
	return nil, fmt.Errorf("invalid date %q", raw)
}

func reservationError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, reservations.ErrItemNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	case errors.Is(err, reservations.ErrStartRequired), errors.Is(err, reservations.ErrInvalidStart):
		return badRequest(c, err.Error())
	case errors.Is(err, reservations.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No such reservation in progress"})
	case errors.Is(err, reservations.ErrNotQueued):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "This waiting list entry no longer exists"})
	case errors.Is(err, reservations.ErrNotLoanItem):
		return badRequest(c, err.Error())
	case errors.Is(err, reservations.ErrNotFree):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "The item or slot is not free: nobody can be told yet"})
	case errors.Is(err, reservations.ErrNotAdmin):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Administrator role required"})
	}
	return internalError(c, "Reservation request failed", err)
}

// GetItems handles GET /admin/reservations/items.
func (h *ReservationsAdminHandler) GetItems(c *fiber.Ctx) error {
	items, err := h.service.ListItems(c.UserContext())
	if err != nil {
		return internalError(c, "Failed to list the reservable items", err)
	}
	return c.JSON(items)
}

// GetStats handles GET /admin/reservations/stats.
func (h *ReservationsAdminHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.service.Stats(c.UserContext())
	if err != nil {
		return internalError(c, "Failed to compute the reservation statistics", err)
	}
	return c.JSON(stats)
}

// GetReservations handles GET /admin/reservations with the filters status, kind, item_id, q, from, to,
// overdue_days, order (asc or desc), limit and offset.
func (h *ReservationsAdminHandler) GetReservations(c *fiber.Ctx) error {
	filter := reservations.ReservationFilter{
		Status:    c.Query("status"),
		Kind:      c.Query("kind"),
		Query:     c.Query("q"),
		Ascending: c.Query("order") == "asc",
	}
	var err error
	if filter.ItemID, err = queryInt(c, "item_id"); err != nil {
		return badRequest(c, err.Error())
	}
	if filter.OverdueDays, err = queryInt(c, "overdue_days"); err != nil {
		return badRequest(c, err.Error())
	}
	if filter.Limit, err = queryInt(c, "limit"); err != nil {
		return badRequest(c, err.Error())
	}
	if filter.Offset, err = queryInt(c, "offset"); err != nil {
		return badRequest(c, err.Error())
	}
	if filter.From, err = parseInstant(c.Query("from")); err != nil {
		return badRequest(c, err.Error())
	}
	if filter.To, err = parseInstant(c.Query("to")); err != nil {
		return badRequest(c, err.Error())
	}

	page, err := h.service.ListReservations(c.UserContext(), filter)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unknown ") {
			return badRequest(c, err.Error())
		}
		return internalError(c, "Failed to list the reservations", err)
	}
	return c.JSON(page)
}

func queryInt(c *fiber.Ctx, name string) (int, error) {
	raw := c.Query(name)
	if raw == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("invalid %s", name)
	}
	return n, nil
}

type cancelRequest struct {
	ItemID    int     `json:"item_id"`
	StartDate *string `json:"start_date"`
}

// CancelReservation handles POST /admin/reservations/cancel: the reservation is cancelled, its owner is told by push
// notification and email with the administrator's name, and the first of the waiting list is told.
func (h *ReservationsAdminHandler) CancelReservation(c *fiber.Ctx) error {
	var req cancelRequest
	if err := c.BodyParser(&req); err != nil || req.ItemID <= 0 {
		return badRequest(c, "item_id is required")
	}
	start, err := parseInstant(deref(req.StartDate))
	if err != nil {
		return badRequest(c, err.Error())
	}

	result, err := h.service.Cancel(c.UserContext(), reservations.CancelRequest{
		ItemID: req.ItemID, Start: start, ActorEmail: c.Locals("email").(string), AsAdmin: true,
	})
	if err != nil {
		return reservationError(c, err)
	}

	if h.discord != nil {
		user := models.ReservationUser{Email: result.Owner}
		_ = h.db.QueryRow(`SELECT first_name, last_name, COALESCE(profile_picture, '') FROM newf WHERE email = $1`, result.Owner).
			Scan(&user.FirstName, &user.LastName, &user.ProfilePicture)
		end := ""
		if result.Item.Slot {
			end = result.Start.Add(time.Hour).Format(reservations.APITimeLayout)
		}
		if err := h.discord.SendReservationCancelled(result.Item.Name, result.Start.Format(reservations.APITimeLayout), end, result.Item.Slot, user); err != nil {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Discord notice of the cancellation failed: %v", err))
		}
	}

	return c.JSON(fiber.Map{
		"message":            "Reservation cancelled",
		"owner":              result.Owner,
		"cancelled_by_admin": result.ByAdmin,
		"waitlist_notified":  result.NextNotified,
	})
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// MarkReturned handles POST /admin/reservations/items/:id/return: the loan of an item without slot ends because it
// was given back. Nobody is notified of it except the first of the waiting list; the reminders stop.
func (h *ReservationsAdminHandler) MarkReturned(c *fiber.Ctx) error {
	itemID, err := strconv.Atoi(c.Params("id"))
	if err != nil || itemID <= 0 {
		return badRequest(c, "Invalid item ID")
	}
	result, err := h.service.MarkReturned(c.UserContext(), itemID, c.Locals("email").(string))
	if err != nil {
		return reservationError(c, err)
	}
	return c.JSON(fiber.Map{
		"message":           "Item marked as returned",
		"owner":             result.Owner,
		"waitlist_notified": result.NextNotified,
	})
}

// GetWaitlist handles GET /admin/reservations/waitlist?item_id=&start_date=.
func (h *ReservationsAdminHandler) GetWaitlist(c *fiber.Ctx) error {
	itemID, err := queryInt(c, "item_id")
	if err != nil || itemID == 0 {
		return badRequest(c, "item_id is required")
	}
	start, err := parseInstant(c.Query("start_date"))
	if err != nil {
		return badRequest(c, err.Error())
	}
	list, err := h.service.AdminWaitlist(c.UserContext(), itemID, start)
	if err != nil {
		return reservationError(c, err)
	}
	return c.JSON(list)
}

// RemoveWaitlistEntry handles DELETE /admin/reservations/waitlist/:id.
func (h *ReservationsAdminHandler) RemoveWaitlistEntry(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return badRequest(c, "Invalid entry ID")
	}
	next, err := h.service.AdminRemoveWaitlistEntry(c.UserContext(), id)
	if err != nil {
		return reservationError(c, err)
	}
	return c.JSON(fiber.Map{"message": "Removed from the waiting list", "next_notified": next})
}

// AdvanceWaitlist handles POST /admin/reservations/waitlist/advance.
func (h *ReservationsAdminHandler) AdvanceWaitlist(c *fiber.Ctx) error {
	var req cancelRequest
	if err := c.BodyParser(&req); err != nil || req.ItemID <= 0 {
		return badRequest(c, "item_id is required")
	}
	start, err := parseInstant(deref(req.StartDate))
	if err != nil {
		return badRequest(c, err.Error())
	}
	next, err := h.service.AdminAdvanceWaitlist(c.UserContext(), req.ItemID, start)
	if err != nil {
		return reservationError(c, err)
	}
	return c.JSON(fiber.Map{"message": "The turn was passed on", "next_notified": next})
}
