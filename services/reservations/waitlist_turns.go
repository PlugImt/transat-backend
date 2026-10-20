package reservations

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

const (
	// turnTimeoutNear is how long the person told about a slot that starts soon has to answer.
	turnTimeoutNear = 15 * time.Minute
	// turnTimeoutFar is how long they have for an item without a slot, or a slot that is still far away.
	turnTimeoutFar = 30 * time.Minute
	// nearSlotWindow: a slot that starts within this delay of the notification is "near".
	nearSlotWindow = 2 * time.Hour
)

// turnDeadline is when the person told now that the resource is free loses their turn.
func (s *Service) turnDeadline(start *time.Time) time.Time {
	now := s.now()
	if start != nil && start.Sub(now) < nearSlotWindow {
		return now.Add(turnTimeoutNear).UTC()
	}
	return now.Add(turnTimeoutFar).UTC()
}

// ExpireWaitlistTurns passes the turn on for everybody who was told that a resource is free and did not answer in
// time. It returns how many people were passed over. Safe to call from several processes: each entry is decided
// under the item's lock.
func (s *Service) ExpireWaitlistTurns(ctx context.Context) (int, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id FROM reservation_waitlist
		WHERE notified_at IS NOT NULL AND turn_expires_at IS NOT NULL AND turn_expires_at <= $1
		ORDER BY turn_expires_at, id`, s.now().UTC())
	if err != nil {
		return 0, err
	}
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return 0, err
		}
		ids = append(ids, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}

	passed := 0
	for _, id := range ids {
		ok, err := s.expireTurn(ctx, id)
		if err != nil {
			log.Printf("Waiting list entry %d: turn not passed on: %v", id, err)
			continue
		}
		if ok {
			passed++
		}
	}
	return passed, nil
}

// expireTurn drops the entry if its person still has not answered and the resource is still free, then tells the
// next person. If somebody else got the resource first, the queue is kept for the next release and only the
// baton is taken back.
func (s *Service) expireTurn(ctx context.Context, entryID int64) (bool, error) {
	var itemID int
	var slotStart sql.NullTime
	err := s.db.QueryRowContext(ctx, `SELECT id_reservation_element, slot_start FROM reservation_waitlist WHERE id = $1`, entryID).
		Scan(&itemID, &slotStart)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	item, err := s.loadItem(ctx, itemID)
	if err != nil {
		return false, err
	}
	var start *time.Time
	if slotStart.Valid {
		t := slotStart.Time.UTC()
		start = &t
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	if err := lockItem(ctx, tx, item.ID); err != nil {
		return false, err
	}
	// Under the lock the answer may have come in: the person booked or left, or the baton changed hands.
	var held bool
	if err := tx.QueryRowContext(ctx, `
		SELECT EXISTS (SELECT 1 FROM reservation_waitlist
		WHERE id = $1 AND notified_at IS NOT NULL AND turn_expires_at <= $2 FOR UPDATE)`, entryID, s.now().UTC()).Scan(&held); err != nil {
		return false, err
	}
	if !held {
		return false, nil
	}

	free := false
	if !item.Slot || start != nil {
		if free, err = isFree(ctx, tx, item, start); err != nil {
			return false, err
		}
	}
	stillOffered := free && item.Status == "active" && (start == nil || start.After(s.now()))
	if !stillOffered {
		if _, err := tx.ExecContext(ctx, `UPDATE reservation_waitlist SET notified_at = NULL, turn_expires_at = NULL WHERE id = $1`, entryID); err != nil {
			return false, err
		}
		return false, tx.Commit()
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM reservation_waitlist WHERE id = $1`, entryID); err != nil {
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	if _, err := s.notifyNext(ctx, item, start); err != nil {
		log.Printf("Waiting list of item %d: could not notify the next user: %v", item.ID, err)
	}
	return true, nil
}
