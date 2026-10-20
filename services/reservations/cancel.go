package reservations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	appI18n "github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
)

// CancelRequest asks to cancel a reservation of an item.
type CancelRequest struct {
	ItemID int
	// Start identifies the slot to cancel; required for a slot item, ignored for a loan item.
	Start      *time.Time
	ActorEmail string
	// AsAdmin is set by the administration routes only. Without it the actor can cancel nothing but their own
	// reservation, whatever their role: the app's route never sets it.
	AsAdmin bool
}

// CancelResult describes what was cancelled.
type CancelResult struct {
	Item Item
	// Owner is the user whose reservation was cancelled.
	Owner string
	// Start is the cancelled slot, or the start of the cancelled loan.
	Start time.Time
	// ByAdmin: an administrator cancelled the reservation of somebody else.
	ByAdmin bool
	// NextNotified: somebody from the waiting list was told that the resource is free.
	NextNotified bool
}

// Cancel cancels a reservation. Without AsAdmin the actor may cancel a slot of their own and nothing else. With
// AsAdmin, which the administration routes set and which still requires the ADMIN role, any slot may be cancelled
// and the loan of an item without a slot may be ended (the borrower returns it themselves otherwise).
// When an administrator cancels the reservation of somebody else, that person gets a push notification and an
// email naming the administrator. Whoever cancels, the head of the waiting list is told that the resource is free.
func (s *Service) Cancel(ctx context.Context, req CancelRequest) (CancelResult, error) {
	item, err := s.loadItem(ctx, req.ItemID)
	if err != nil {
		return CancelResult{}, err
	}
	admin := false
	if req.AsAdmin {
		if admin, err = s.IsAdmin(ctx, req.ActorEmail); err != nil {
			return CancelResult{}, fmt.Errorf("failed to check the admin role: %w", err)
		}
		if !admin {
			return CancelResult{}, ErrNotAdmin
		}
	}

	var rows *sql.Rows
	res := CancelResult{Item: item}
	if item.Slot {
		if req.Start == nil {
			return res, ErrStartRequired
		}
		res.Start = req.Start.UTC()
		query := `DELETE FROM reservation WHERE id_reservation_element = $1 AND start_date = $2`
		args := []any{item.ID, res.Start}
		if !admin {
			query += ` AND email = $3`
			args = append(args, req.ActorEmail)
		}
		rows, err = s.db.QueryContext(ctx, query+` RETURNING email`, args...)
	} else {
		if !admin {
			return res, ErrNotSlotItem
		}
		// The loan ends now; the end must stay after the start for the table's range check.
		rows, err = s.db.QueryContext(ctx, `
			UPDATE reservation
			SET end_date = GREATEST(now() AT TIME ZONE 'UTC', start_date + interval '1 second')
			WHERE id_reservation_element = $1 AND end_date IS NULL
			RETURNING email, start_date`, item.ID)
	}
	if err != nil {
		return res, err
	}
	var owners []string
	for rows.Next() {
		var owner string
		if item.Slot {
			err = rows.Scan(&owner)
		} else {
			err = rows.Scan(&owner, &res.Start)
		}
		if err != nil {
			rows.Close()
			return res, err
		}
		owners = append(owners, owner)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return res, err
	}
	if len(owners) == 0 {
		return res, ErrNotFound
	}
	res.Owner = owners[0]

	end := s.now().UTC()
	if item.Slot {
		end = res.Start.Add(time.Hour)
	}
	for _, owner := range owners {
		kind := "cancelled"
		if owner != req.ActorEmail {
			kind = "cancelled_by_admin"
		}
		s.logCancellation(ctx, item.ID, owner, res.Start, end, kind, req.ActorEmail)
	}

	// The reservation is gone. What follows must not fail the request: log and move on.
	var start *time.Time
	if item.Slot {
		start = &res.Start
	}
	adminName := ""
	for _, owner := range owners {
		if owner == req.ActorEmail || !admin {
			continue
		}
		res.ByAdmin = true
		if adminName == "" {
			adminName = s.displayName(ctx, req.ActorEmail)
		}
		if err := s.notifyAdminCancellation(ctx, item, start, owner, adminName); err != nil {
			log.Printf("Admin cancellation of item %d: notification to the user failed: %v", item.ID, err)
		}
	}
	notified, err := s.notifyNext(ctx, item, start)
	if err != nil {
		log.Printf("Waiting list of item %d: could not notify the next user: %v", item.ID, err)
	}
	res.NextNotified = notified
	return res, nil
}

// ReleaseOf tells the waiting list that an item or slot was released by other means than Cancel (the owner
// returned a loan). It returns whether somebody was notified.
func (s *Service) ReleaseOf(ctx context.Context, itemID int, start *time.Time) (bool, error) {
	item, err := s.loadItem(ctx, itemID)
	if err != nil {
		return false, err
	}
	if !item.Slot {
		start = nil
	}
	return s.notifyNext(ctx, item, start)
}

// CancelEmailData feeds the cancellation email template. Names are rendered by html/template, which escapes
// them; they never go through a translation message, which would not.
type CancelEmailData struct {
	ItemName    string
	HasSlot     bool
	Slot        string
	AdminName   string
	CancelledAt string
}

func (s *Service) notifyAdminCancellation(ctx context.Context, item Item, start *time.Time, owner, adminName string) error {
	if adminName == "" {
		adminName = "Admin"
	}
	lang := s.userLanguage(ctx, owner)
	var errs []error

	if s.notifier != nil {
		now := s.now()
		var slotStart *time.Time
		if item.Slot {
			slotStart = start
		}
		_, err := s.notifier.Announce(ctx, services.Announcement{
			Kind:      "reservation_cancelled",
			Category:  models.NotificationReservations,
			DedupeKey: dedupe("reservation_cancelled", item.ID, owner, now.UnixNano()),
			Audience:  usersAudience(owner),
			Compose: func(l string, loc *i18n.Localizer) (string, string) {
				data := map[string]interface{}{"AdminName": adminName, "ItemName": item.Name}
				title := appI18n.Localize(loc, "reservation_notification.admin_cancelled_title", "Reservation cancelled", nil)
				if slotStart == nil {
					return title, appI18n.Localize(loc, "reservation_notification.admin_cancelled_message",
						"{{.AdminName}} cancelled your reservation of {{.ItemName}}.", data)
				}
				data["Slot"] = FormatSlot(l, *slotStart)
				return title, appI18n.Localize(loc, "reservation_notification.admin_cancelled_message_slot",
					"{{.AdminName}} cancelled your reservation of {{.ItemName}} on {{.Slot}}.", data)
			},
			Navigation: navigation(item, slotStart, nil),
			ExpiresAt:  now.Add(24 * time.Hour),
		})
		if err != nil {
			errs = append(errs, fmt.Errorf("push: %w", err))
		}
	}

	if s.mailer != nil {
		data := CancelEmailData{
			ItemName:    item.Name,
			HasSlot:     item.Slot && start != nil,
			AdminName:   adminName,
			CancelledAt: FormatSlot(lang, s.now()),
		}
		if data.HasSlot {
			data.Slot = FormatSlot(lang, *start)
		}
		mail := models.Email{Recipient: owner, Template: emailTemplate, SubjectKey: emailSubjectKey, Language: lang}
		s.background("cancellation email", func(context.Context) error {
			return s.mailer.SendEmail(mail, data)
		})
	}
	return errors.Join(errs...)
}
