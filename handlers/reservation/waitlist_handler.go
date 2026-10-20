package reservation

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services/reservations"
	"github.com/plugimt/transat-backend/utils"
)

var errBadWaitlistRequest = errors.New("invalid request")

// waitlistTarget reads the item id of the path and the optional slot (start_date, in the query or the body).
func waitlistTarget(c *fiber.Ctx) (int, *time.Time, error) {
	itemID, err := strconv.Atoi(c.Params("id"))
	if err != nil || itemID <= 0 {
		return 0, nil, fmt.Errorf("%w: invalid item ID", errBadWaitlistRequest)
	}
	raw := c.Query("start_date")
	if raw == "" && len(c.Body()) > 0 {
		var body models.ReservationManagementRequest
		if err := c.BodyParser(&body); err != nil {
			return 0, nil, fmt.Errorf("%w: invalid request format", errBadWaitlistRequest)
		}
		if body.StartDate != nil {
			raw = *body.StartDate
		}
	}
	if raw == "" {
		return itemID, nil, nil
	}
	start, err := time.Parse(reservations.APITimeLayout, raw)
	if err != nil {
		return 0, nil, fmt.Errorf("%w: invalid start_date format", errBadWaitlistRequest)
	}
	return itemID, &start, nil
}

// waitlistFailure answers a failed waiting list call with a status, a message and a stable code for the app.
func waitlistFailure(c *fiber.Ctx, err error) error {
	status, code, message := fiber.StatusInternalServerError, "internal", "Request failed"
	switch {
	case errors.Is(err, errBadWaitlistRequest):
		status, code, message = fiber.StatusBadRequest, "bad_request", err.Error()
	case errors.Is(err, reservations.ErrItemNotFound):
		status, code, message = fiber.StatusNotFound, "item_not_found", "The item does not exist"
	case errors.Is(err, reservations.ErrStartRequired):
		status, code, message = fiber.StatusBadRequest, "start_required", "start_date is required for a slot item"
	case errors.Is(err, reservations.ErrInvalidStart):
		status, code, message = fiber.StatusBadRequest, "invalid_start", "start_date must be a future slot on a plain hour"
	case errors.Is(err, reservations.ErrStillAvailable):
		status, code, message = fiber.StatusConflict, "still_available", "This item is available: reserve it directly"
	case errors.Is(err, reservations.ErrOwnReservation):
		status, code, message = fiber.StatusConflict, "own_reservation", "You already hold this reservation"
	case errors.Is(err, reservations.ErrAlreadyQueued):
		status, code, message = fiber.StatusConflict, "already_queued", "You are already on the waiting list"
	case errors.Is(err, reservations.ErrQueueFull):
		status, code, message = fiber.StatusConflict, "queue_limit", "You are on too many waiting lists"
	case errors.Is(err, reservations.ErrNotQueued):
		status, code, message = fiber.StatusNotFound, "not_queued", "You are not on the waiting list of this item"
	case errors.Is(err, reservations.ErrNotAvailable):
		status, code, message = fiber.StatusConflict, "not_available", "Too late: this item has just been reserved by someone else"
	default:
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Waiting list request failed: %v", err))
	}
	return c.Status(status).JSON(fiber.Map{"error": message, "code": code})
}

// JoinWaitlist handles POST /reservation/item/:id/waitlist: show interest in an item or slot that is reserved.
func (h *ReservationHandler) JoinWaitlist(c *fiber.Ctx) error {
	itemID, start, err := waitlistTarget(c)
	if err != nil {
		return waitlistFailure(c, err)
	}
	position, err := h.Reservations.Join(c.UserContext(), c.Locals("email").(string), itemID, start)
	if err != nil {
		return waitlistFailure(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "You are on the waiting list",
		"position": position,
	})
}

// LeaveWaitlist handles DELETE /reservation/item/:id/waitlist and POST /reservation/item/:id/waitlist/decline
// (the "Annuler" action of the notification): the user leaves the queue, and when they were the one being waited
// for, the next person is told that the item is free.
func (h *ReservationHandler) LeaveWaitlist(c *fiber.Ctx) error {
	itemID, start, err := waitlistTarget(c)
	if err != nil {
		return waitlistFailure(c, err)
	}
	next, err := h.Reservations.Remove(c.UserContext(), c.Locals("email").(string), itemID, start)
	if err != nil {
		return waitlistFailure(c, err)
	}
	return c.JSON(fiber.Map{
		"message":       "You left the waiting list",
		"next_notified": next,
	})
}

// ReserveFromWaitlist handles POST /reservation/item/:id/waitlist/reserve (the "Réserver" action of the
// notification): the item or slot is booked at once if it is still free, otherwise the answer is a clear 409.
func (h *ReservationHandler) ReserveFromWaitlist(c *fiber.Ctx) error {
	itemID, start, err := waitlistTarget(c)
	if err != nil {
		return waitlistFailure(c, err)
	}
	res, err := h.Reservations.Reserve(c.UserContext(), c.Locals("email").(string), itemID, start)
	if err != nil {
		return waitlistFailure(c, err)
	}

	if h.DiscordService != nil {
		var endDate string
		if res.End != nil {
			endDate = res.End.Format(reservations.APITimeLayout)
		}
		if err := h.DiscordService.SendReservationCreated(res.Item.Name, res.Start.Format(reservations.APITimeLayout), endDate, res.Item.Slot, res.User); err != nil {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Failed to send Discord notification for reservation creation: %v", err))
		}
	}

	return c.JSON(fiber.Map{
		"message": "Reservation created successfully",
		"reservation": models.ReservationItem{
			ID:   res.Item.ID,
			Name: res.Item.Name,
			Slot: res.Item.Slot,
			User: &res.User,
		},
	})
}

// GetItemWaitlist handles GET /reservation/item/:id/waitlist: the length of the queue and the caller's place.
func (h *ReservationHandler) GetItemWaitlist(c *fiber.Ctx) error {
	itemID, start, err := waitlistTarget(c)
	if err != nil {
		return waitlistFailure(c, err)
	}
	status, err := h.Reservations.Status(c.UserContext(), c.Locals("email").(string), itemID, start)
	if err != nil {
		return waitlistFailure(c, err)
	}
	return c.JSON(status)
}

// GetMyWaitlist handles GET /reservation/waitlist: every queue the caller is in.
func (h *ReservationHandler) GetMyWaitlist(c *fiber.Ctx) error {
	entries, err := h.Reservations.List(c.UserContext(), c.Locals("email").(string))
	if err != nil {
		return waitlistFailure(c, err)
	}
	return c.JSON(fiber.Map{"waitlist": entries})
}
