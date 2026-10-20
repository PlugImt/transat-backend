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

// The waiting list is a first come first served queue of people interested in something that is reserved: an
// item that is on loan, or one slot of a slot item. It never blocks a booking: anybody may still book the
// resource the moment it is free, queue or not. Its only effect is that, each time the resource is released,
// the first person of the queue is told, with two actions: "Réserver" (Reserve) and "Annuler" (Remove), which
// passes the baton to the next person.

// WaitlistEntry is one place in a queue.
type WaitlistEntry struct {
	ItemID   int        `json:"item_id"`
	ItemName string     `json:"item_name"`
	Slot     bool       `json:"slot"`
	Start    *time.Time `json:"-"`
	// StartDate is the slot in the API's UTC format; absent for an item without slots.
	StartDate *string `json:"start_date,omitempty"`
	// Position is 1 for the first of the queue.
	Position int `json:"position"`
	// Notified: this person was told that the resource is free and has not answered yet.
	Notified bool      `json:"notified"`
	JoinedAt time.Time `json:"joined_at"`
}

// WaitlistStatus is what a user sees of the queue of one item or slot.
type WaitlistStatus struct {
	Count    int  `json:"count"`
	Joined   bool `json:"joined"`
	Position int  `json:"position,omitempty"`
	Notified bool `json:"notified,omitempty"`
}

// Join puts the user in the queue of an item or slot that is reserved. It returns the user's position.
func (s *Service) Join(ctx context.Context, email string, itemID int, start *time.Time) (int, error) {
	item, err := s.loadItem(ctx, itemID)
	if err != nil {
		return 0, err
	}
	slot, err := s.target(item, start)
	if err != nil {
		return 0, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	if err := lockItem(ctx, tx, item.ID); err != nil {
		return 0, err
	}

	free, err := isFree(ctx, tx, item, slot)
	if err != nil {
		return 0, err
	}
	if free {
		return 0, ErrStillAvailable
	}
	var own bool
	if item.Slot {
		err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM reservation
			WHERE id_reservation_element = $1 AND start_date = $2 AND email = $3)`, item.ID, slot.UTC(), email).Scan(&own)
	} else {
		err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM reservation
			WHERE id_reservation_element = $1 AND end_date IS NULL AND email = $2)`, item.ID, email).Scan(&own)
	}
	if err != nil {
		return 0, err
	}
	if own {
		return 0, ErrOwnReservation
	}
	var entries int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM reservation_waitlist WHERE email = $1`, email).Scan(&entries); err != nil {
		return 0, err
	}
	if entries >= maxQueueEntries {
		return 0, ErrQueueFull
	}

	var id int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO reservation_waitlist (email, id_reservation_element, slot_start)
		VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id`, email, item.ID, startArg(slot)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrAlreadyQueued
	}
	if err != nil {
		return 0, err
	}
	var position int
	if err := tx.QueryRowContext(ctx, `
		SELECT count(*) FROM reservation_waitlist
		WHERE id_reservation_element = $1 AND slot_start IS NOT DISTINCT FROM $2::timestamp AND id <= $3`,
		item.ID, startArg(slot), id).Scan(&position); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return position, nil
}

// Remove takes the user out of a queue. When the user was the one being waited for (told that the resource is
// free), the resource is still free and somebody else is queued, that next person is told straight away; the
// returned bool says so. This is the "Annuler" action of the notification, and also how a user leaves a queue.
func (s *Service) Remove(ctx context.Context, email string, itemID int, start *time.Time) (bool, error) {
	item, err := s.loadItem(ctx, itemID)
	if err != nil {
		return false, err
	}
	slot := start
	if !item.Slot {
		slot = nil
	} else if slot == nil {
		return false, ErrStartRequired
	} else {
		u := slot.UTC()
		slot = &u
	}
	var notified sql.NullTime
	err = s.db.QueryRowContext(ctx, `
		DELETE FROM reservation_waitlist
		WHERE email = $1 AND id_reservation_element = $2 AND slot_start IS NOT DISTINCT FROM $3::timestamp
		RETURNING notified_at`, email, item.ID, startArg(slot)).Scan(&notified)
	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrNotQueued
	}
	if err != nil {
		return false, err
	}
	if !notified.Valid {
		return false, nil
	}
	next, err := s.notifyNext(ctx, item, slot)
	if err != nil {
		log.Printf("Waiting list of item %d: could not notify the next user: %v", item.ID, err)
	}
	return next, nil
}

// ReserveResult is a reservation made from the waiting list.
type ReserveResult struct {
	Item  Item
	Start time.Time
	End   *time.Time
	User  models.ReservationUser
}

// Reserve books the resource for a user of its queue if it is still free, and takes them out of the queue. It
// fails with ErrNotAvailable when somebody else got it first, in which case the user keeps their place.
func (s *Service) Reserve(ctx context.Context, email string, itemID int, start *time.Time) (ReserveResult, error) {
	item, err := s.loadItem(ctx, itemID)
	if err != nil {
		return ReserveResult{}, err
	}
	slot, err := s.target(item, start)
	if err != nil {
		return ReserveResult{}, err
	}
	res := ReserveResult{Item: item}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return res, err
	}
	defer tx.Rollback()
	if err := lockItem(ctx, tx, item.ID); err != nil {
		return res, err
	}

	var queued bool
	if err := tx.QueryRowContext(ctx, `
		SELECT EXISTS (SELECT 1 FROM reservation_waitlist
		WHERE email = $1 AND id_reservation_element = $2 AND slot_start IS NOT DISTINCT FROM $3::timestamp)`,
		email, item.ID, startArg(slot)).Scan(&queued); err != nil {
		return res, err
	}
	if !queued {
		return res, ErrNotQueued
	}
	free, err := isFree(ctx, tx, item, slot)
	if err != nil {
		return res, err
	}
	if !free {
		return res, ErrNotAvailable
	}

	if item.Slot {
		end := slot.Add(time.Hour)
		res.Start, res.End = *slot, &end
		_, err = tx.ExecContext(ctx, `
			INSERT INTO reservation (email, id_reservation_element, start_date, end_date) VALUES ($1, $2, $3, $4)`,
			email, item.ID, res.Start, end)
	} else {
		err = tx.QueryRowContext(ctx, `
			INSERT INTO reservation (email, id_reservation_element, start_date)
			VALUES ($1, $2, date_trunc('second', now() AT TIME ZONE 'UTC')) RETURNING start_date`,
			email, item.ID).Scan(&res.Start)
	}
	if err != nil {
		return res, err
	}
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM reservation_waitlist
		WHERE email = $1 AND id_reservation_element = $2 AND slot_start IS NOT DISTINCT FROM $3::timestamp`,
		email, item.ID, startArg(slot)); err != nil {
		return res, err
	}
	if err := tx.Commit(); err != nil {
		return res, err
	}

	var picture sql.NullString
	res.User.Email = email
	if err := s.db.QueryRowContext(ctx, `SELECT first_name, last_name, profile_picture FROM newf WHERE email = $1`, email).
		Scan(&res.User.FirstName, &res.User.LastName, &picture); err == nil {
		res.User.ProfilePicture = picture.String
	}
	return res, nil
}

// ForgetInterest removes the user's place in a queue; it is called when they book the resource directly.
func (s *Service) ForgetInterest(ctx context.Context, email string, itemID int, start *time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM reservation_waitlist
		WHERE email = $1 AND id_reservation_element = $2 AND slot_start IS NOT DISTINCT FROM $3::timestamp`,
		email, itemID, startArg(start))
	return err
}

// Status says how long the queue of an item or slot is and where the user stands in it.
func (s *Service) Status(ctx context.Context, email string, itemID int, start *time.Time) (WaitlistStatus, error) {
	item, err := s.loadItem(ctx, itemID)
	if err != nil {
		return WaitlistStatus{}, err
	}
	var slot *time.Time
	if item.Slot {
		if start == nil {
			return WaitlistStatus{}, ErrStartRequired
		}
		u := start.UTC()
		slot = &u
	}
	var st WaitlistStatus
	var notified sql.NullTime
	err = s.db.QueryRowContext(ctx, `
		SELECT (SELECT count(*) FROM reservation_waitlist q
		        WHERE q.id_reservation_element = $1 AND q.slot_start IS NOT DISTINCT FROM $3::timestamp),
		       COALESCE((SELECT count(*) FROM reservation_waitlist q
		        WHERE q.id_reservation_element = $1 AND q.slot_start IS NOT DISTINCT FROM $3::timestamp
		          AND q.id <= (SELECT w.id FROM reservation_waitlist w
		                       WHERE w.email = $2 AND w.id_reservation_element = $1
		                         AND w.slot_start IS NOT DISTINCT FROM $3::timestamp)), 0),
		       (SELECT w.notified_at FROM reservation_waitlist w
		        WHERE w.email = $2 AND w.id_reservation_element = $1 AND w.slot_start IS NOT DISTINCT FROM $3::timestamp)`,
		item.ID, email, startArg(slot)).Scan(&st.Count, &st.Position, &notified)
	if err != nil {
		return st, err
	}
	st.Joined = st.Position > 0
	st.Notified = notified.Valid
	return st, nil
}

// List returns the queues the user is in, with their position, skipping slots that are over.
func (s *Service) List(ctx context.Context, email string) ([]WaitlistEntry, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT w.id_reservation_element, COALESCE(re.name, ''), COALESCE(re.slot, FALSE), w.slot_start,
		       (SELECT count(*) FROM reservation_waitlist q
		        WHERE q.id_reservation_element = w.id_reservation_element
		          AND q.slot_start IS NOT DISTINCT FROM w.slot_start AND q.id <= w.id),
		       w.notified_at IS NOT NULL, w.created_at
		FROM reservation_waitlist w
		JOIN reservation_element re ON re.id_reservation_element = w.id_reservation_element
		WHERE w.email = $1 AND (w.slot_start IS NULL OR w.slot_start > $2)
		ORDER BY w.id`, email, s.now().UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []WaitlistEntry{}
	for rows.Next() {
		var e WaitlistEntry
		var start sql.NullTime
		if err := rows.Scan(&e.ItemID, &e.ItemName, &e.Slot, &start, &e.Position, &e.Notified, &e.JoinedAt); err != nil {
			return nil, err
		}
		if start.Valid {
			t := start.Time.UTC()
			e.Start = &t
			formatted := t.Format(APITimeLayout)
			e.StartDate = &formatted
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// notifyNext tells the first person of the queue that the resource is free, and nobody else. It does nothing when
// the resource is not free (somebody booked it first), when the queue is empty, or when the push pipeline is
// off. People who could not receive the push (no device, notifications of reservations switched off) are passed
// over, so that the baton does not stall with somebody who will never see it. It reports whether somebody was
// notified.
func (s *Service) notifyNext(ctx context.Context, item Item, start *time.Time) (bool, error) {
	if s.notifier == nil || !s.notifier.OutboxEnabled() {
		return false, nil
	}
	if item.Slot && (start == nil || !start.After(s.now())) {
		return false, nil
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	if err := lockItem(ctx, tx, item.ID); err != nil {
		return false, err
	}
	free, err := isFree(ctx, tx, item, start)
	if err != nil || !free {
		return false, err
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT id, email FROM reservation_waitlist
		WHERE id_reservation_element = $1 AND slot_start IS NOT DISTINCT FROM $2::timestamp
		ORDER BY id FOR UPDATE`, item.ID, startArg(start))
	if err != nil {
		return false, err
	}
	type queued struct {
		id    int64
		email string
	}
	var queue []queued
	for rows.Next() {
		var q queued
		if err := rows.Scan(&q.id, &q.email); err != nil {
			rows.Close()
			return false, err
		}
		queue = append(queue, q)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return false, err
	}

	for _, q := range queue {
		reach, err := s.notifier.CanReach(ctx, q.email, models.NotificationReservations)
		if err != nil {
			return false, err
		}
		if !reach {
			continue
		}
		var stamp int64
		if err := tx.QueryRowContext(ctx, `
			UPDATE reservation_waitlist SET notified_at = now() WHERE id = $1
			RETURNING (extract(epoch FROM notified_at) * 1000000)::bigint`, q.id).Scan(&stamp); err != nil {
			return false, err
		}
		if err := tx.Commit(); err != nil {
			return false, err
		}
		if err := s.announceAvailable(ctx, item, start, q.email, q.id, stamp); err != nil {
			// Not queued: this person never heard of it, so they must not keep the baton.
			if _, rerr := s.db.ExecContext(ctx, `UPDATE reservation_waitlist SET notified_at = NULL WHERE id = $1`, q.id); rerr != nil {
				log.Printf("Waiting list entry %d: baton not released: %v", q.id, rerr)
			}
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (s *Service) announceAvailable(ctx context.Context, item Item, start *time.Time, email string, entryID, stamp int64) error {
	now := s.now()
	extra := map[string]string{"waitlist": "1"}
	expires := now.Add(12 * time.Hour)
	if start != nil {
		extra["start"] = start.UTC().Format(APITimeLayout)
		if start.Before(expires) {
			expires = *start
		}
	}
	_, err := s.notifier.Announce(ctx, services.Announcement{
		Kind:           "reservation_waitlist",
		Category:       models.NotificationReservations,
		DedupeKey:      dedupe("reservation_waitlist", entryID, stamp),
		Audience:       usersAudience(email),
		ActionCategory: WaitlistActionCategory,
		Navigation:     navigation(item, start, extra),
		ExpiresAt:      expires,
		Compose: func(l string, loc *i18n.Localizer) (string, string) {
			data := map[string]interface{}{"ItemName": item.Name}
			title := appI18n.Localize(loc, "reservation_notification.waitlist_available_title", "{{.ItemName}} is available", data)
			if start == nil {
				return title, appI18n.Localize(loc, "reservation_notification.waitlist_available_message",
					"You were on the waiting list: book it before someone else does.", nil)
			}
			data["Slot"] = FormatSlot(l, *start)
			return title, appI18n.Localize(loc, "reservation_notification.waitlist_available_message_slot",
				"The {{.Slot}} slot just became free. Book it before someone else does.", data)
		},
	})
	if err != nil {
		return fmt.Errorf("waiting list notification: %w", err)
	}
	return nil
}
