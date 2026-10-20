package reservations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
)

// The read side and the manual actions of the administration panel.

var (
	ErrNotLoanItem = errors.New("item has slots: cancel its reservations instead")
	ErrNotFree     = errors.New("item or slot is not free")
)

// AdminHolder is somebody holding an item, or a slot, with the time that has gone by for a loan.
type AdminHolder struct {
	Email          string     `json:"email"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	ElapsedSeconds int64      `json:"elapsed_seconds"`
}

// AdminItem is a reservable element with what is happening to it right now.
type AdminItem struct {
	ID                  int     `json:"id"`
	Name                string  `json:"name"`
	Slot                bool    `json:"slot"`
	Description         *string `json:"description,omitempty"`
	Location            *string `json:"location,omitempty"`
	Status              string  `json:"status"`
	StatusNote          *string `json:"status_note,omitempty"`
	ClubID              int     `json:"club_id"`
	ClubName            string  `json:"club_name"`
	CategoryID          *int    `json:"category_id,omitempty"`
	CategoryPath        string  `json:"category_path"`
	WarningMessage      *string `json:"warning_message,omitempty"`
	ConfirmationMessage *string `json:"confirmation_message,omitempty"`
	// Availability: available, borrowed (loan), occupied (slot being used now), maintenance or inactive.
	Availability  string       `json:"availability"`
	Current       *AdminHolder `json:"current,omitempty"`
	Next          *AdminHolder `json:"next,omitempty"`
	UpcomingCount int          `json:"upcoming_count"`
	WaitlistCount int          `json:"waitlist_count"`
}

const adminItemsQuery = `
WITH RECURSIVE cat(id, path) AS (
	SELECT id_reservation_category, name::text FROM reservation_category WHERE id_parent_category IS NULL
	UNION ALL
	SELECT c.id_reservation_category, cat.path || ' / ' || c.name
	FROM reservation_category c JOIN cat ON c.id_parent_category = cat.id
)
SELECT re.id_reservation_element, COALESCE(re.name, ''), COALESCE(re.slot, FALSE), re.description, re.location,
       re.status, re.status_note, re.id_clubs, COALESCE(cl.name, ''), re.id_reservation_category, COALESCE(cat.path, ''),
       re.warning_message, re.confirmation_message,
       cur.email, cur.first_name, cur.last_name, cur.start_date, cur.end_date,
       nxt.email, nxt.first_name, nxt.last_name, nxt.start_date, nxt.end_date,
       (SELECT count(*) FROM reservation u WHERE u.id_reservation_element = re.id_reservation_element
              AND COALESCE(re.slot, FALSE) AND u.start_date > $1::timestamp),
       (SELECT count(*) FROM reservation_waitlist w WHERE w.id_reservation_element = re.id_reservation_element
              AND (w.slot_start IS NULL OR w.slot_start > $1::timestamp))
FROM reservation_element re
JOIN clubs cl ON cl.id_clubs = re.id_clubs
LEFT JOIN cat ON cat.id = re.id_reservation_category
LEFT JOIN LATERAL (
	SELECT r.email, n.first_name, n.last_name, r.start_date, r.end_date
	FROM reservation r JOIN newf n ON n.email = r.email
	WHERE r.id_reservation_element = re.id_reservation_element
	  AND ((NOT COALESCE(re.slot, FALSE) AND r.end_date IS NULL)
	       OR (COALESCE(re.slot, FALSE) AND r.start_date <= $1::timestamp AND r.end_date > $1::timestamp))
	ORDER BY r.start_date DESC LIMIT 1) cur ON TRUE
LEFT JOIN LATERAL (
	SELECT r.email, n.first_name, n.last_name, r.start_date, r.end_date
	FROM reservation r JOIN newf n ON n.email = r.email
	WHERE r.id_reservation_element = re.id_reservation_element AND COALESCE(re.slot, FALSE) AND r.start_date > $1::timestamp
	ORDER BY r.start_date LIMIT 1) nxt ON TRUE
ORDER BY cl.name, COALESCE(cat.path, ''), re.name`

type holderScan struct {
	email, first, last sql.NullString
	start, end         sql.NullTime
}

func (h holderScan) holder(now time.Time, loan bool) *AdminHolder {
	if !h.email.Valid {
		return nil
	}
	out := &AdminHolder{Email: h.email.String, FirstName: h.first.String, LastName: h.last.String, StartDate: h.start.Time.UTC()}
	if h.end.Valid {
		end := h.end.Time.UTC()
		out.EndDate = &end
	}
	if loan {
		out.ElapsedSeconds = int64(now.Sub(out.StartDate).Seconds())
	}
	return out
}

func optional(s sql.NullString) *string {
	if !s.Valid || s.String == "" {
		return nil
	}
	return &s.String
}

// ListItems returns every reservable element, with its category, its holder or next reservation, and the length of its queue.
func (s *Service) ListItems(ctx context.Context) ([]AdminItem, error) {
	now := s.now().UTC()
	rows, err := s.db.QueryContext(ctx, adminItemsQuery, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []AdminItem{}
	for rows.Next() {
		var it AdminItem
		var desc, loc, note, warn, confirm sql.NullString
		var category sql.NullInt64
		var cur, nxt holderScan
		if err := rows.Scan(&it.ID, &it.Name, &it.Slot, &desc, &loc, &it.Status, &note, &it.ClubID, &it.ClubName, &category, &it.CategoryPath,
			&warn, &confirm,
			&cur.email, &cur.first, &cur.last, &cur.start, &cur.end,
			&nxt.email, &nxt.first, &nxt.last, &nxt.start, &nxt.end,
			&it.UpcomingCount, &it.WaitlistCount); err != nil {
			return nil, err
		}
		it.Description, it.Location, it.StatusNote = optional(desc), optional(loc), optional(note)
		it.WarningMessage, it.ConfirmationMessage = optional(warn), optional(confirm)
		if category.Valid {
			id := int(category.Int64)
			it.CategoryID = &id
		}
		it.Current, it.Next = cur.holder(now, !it.Slot), nxt.holder(now, false)

		switch {
		case it.Status != "active":
			it.Availability = it.Status
		case it.Current != nil && it.Slot:
			it.Availability = "occupied"
		case it.Current != nil:
			it.Availability = "borrowed"
		default:
			it.Availability = "available"
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

// ReservationFilter narrows the administration list of reservations.
type ReservationFilter struct {
	// Status: active (a loan in progress or a slot being used), upcoming, past, cancelled, overdue (a loan older than a day); empty for all.
	Status string
	// Kind: slot or loan; empty for both.
	Kind   string
	ItemID int
	// Query matches the user's name or email.
	Query string
	// From and To bound the start of the reservation.
	From, To *time.Time
	// OverdueDays keeps loans in progress for at least that many days.
	OverdueDays int
	Ascending   bool
	Limit       int
	Offset      int
}

// AdminReservation is one line of the reservations list: a reservation, or a cancelled one.
type AdminReservation struct {
	Key       string     `json:"key"`
	ItemID    int        `json:"item_id"`
	ItemName  string     `json:"item_name"`
	Slot      bool       `json:"slot"`
	ClubName  string     `json:"club_name"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	// Status: upcoming, ongoing (slot being used), active (loan in progress), past or cancelled.
	Status string `json:"status"`
	// ElapsedSeconds is how long a loan has lasted (or lasted).
	ElapsedSeconds *int64 `json:"elapsed_seconds,omitempty"`
	// OverdueLevel is 1, 5 or 7 for a loan in progress that is that many days old, 0 otherwise.
	OverdueLevel int `json:"overdue_level"`
	// Reminders lists the reminders already sent (slot_1h, loan_d1, loan_d5, loan_d7).
	Reminders   []string   `json:"reminders"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
	CancelledBy *string    `json:"cancelled_by,omitempty"`
	// CancelKind: cancelled, cancelled_by_admin or returned_by_admin.
	CancelKind *string `json:"cancel_kind,omitempty"`
	// Cancellable: an administrator may cancel it (active, ongoing or upcoming).
	Cancellable bool `json:"cancellable"`
}

type ReservationPage struct {
	Items  []AdminReservation `json:"items"`
	Total  int                `json:"total"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
}

const adminReservationsCTE = `
WITH lines AS (
	SELECT 'r:' || r.id_reservation_element || ':' || r.email || ':' || to_char(r.start_date, 'YYYY-MM-DD"T"HH24:MI:SS') AS key,
	       r.id_reservation_element AS item_id, COALESCE(re.name, '') AS item_name, COALESCE(re.slot, FALSE) AS slot,
	       COALESCE(cl.name, '') AS club_name, r.email, COALESCE(n.first_name, '') AS first_name, COALESCE(n.last_name, '') AS last_name,
	       r.start_date, r.end_date,
	       CASE WHEN COALESCE(re.slot, FALSE) THEN
	                CASE WHEN r.start_date > $1::timestamp THEN 'upcoming' WHEN r.end_date > $1::timestamp THEN 'ongoing' ELSE 'past' END
	            WHEN r.end_date IS NULL THEN 'active'
	            WHEN cx.kind = 'cancelled_by_admin' THEN 'cancelled'
	            ELSE 'past' END AS status,
	       CASE WHEN NOT COALESCE(re.slot, FALSE) THEN EXTRACT(EPOCH FROM (COALESCE(r.end_date, $1::timestamp) - r.start_date))::bigint END AS elapsed,
	       cx.cancelled_at, cx.cancelled_by, cx.kind AS cancel_kind,
	       COALESCE((SELECT array_agg(rr.milestone ORDER BY rr.sent_at) FROM reservation_reminders rr
	                 WHERE rr.email = r.email AND rr.id_reservation_element = r.id_reservation_element AND rr.start_date = r.start_date),
	                '{}'::varchar[]) AS reminders
	FROM reservation r
	JOIN reservation_element re ON re.id_reservation_element = r.id_reservation_element
	LEFT JOIN clubs cl ON cl.id_clubs = re.id_clubs
	LEFT JOIN newf n ON n.email = r.email
	LEFT JOIN LATERAL (
		SELECT c.kind, c.cancelled_by, c.cancelled_at FROM reservation_cancellations c
		WHERE c.id_reservation_element = r.id_reservation_element AND c.email = r.email AND c.start_date = r.start_date
		  AND NOT COALESCE(re.slot, FALSE)
		ORDER BY c.id DESC LIMIT 1) cx ON TRUE
	UNION ALL
	SELECT 'c:' || c.id, c.id_reservation_element, COALESCE(re.name, ''), TRUE, COALESCE(cl.name, ''), c.email,
	       COALESCE(n.first_name, ''), COALESCE(n.last_name, ''), c.start_date, c.end_date, 'cancelled', NULL::bigint,
	       c.cancelled_at, c.cancelled_by, c.kind, '{}'::varchar[]
	FROM reservation_cancellations c
	JOIN reservation_element re ON re.id_reservation_element = c.id_reservation_element
	LEFT JOIN clubs cl ON cl.id_clubs = re.id_clubs
	LEFT JOIN newf n ON n.email = c.email
	WHERE COALESCE(re.slot, FALSE) AND c.kind IN ('cancelled', 'cancelled_by_admin')
)`

// ListReservations returns the reservations (and cancelled ones) matching the filter, most recent first.
func (s *Service) ListReservations(ctx context.Context, f ReservationFilter) (ReservationPage, error) {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	args := []any{s.now().UTC()}
	arg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}
	var where []string

	switch f.Status {
	case "":
	case "active":
		where = append(where, "status IN ('active', 'ongoing')")
	case "upcoming", "past", "cancelled":
		where = append(where, "status = "+arg(f.Status))
	case "overdue":
		f.OverdueDays = max(f.OverdueDays, 1)
	default:
		return ReservationPage{}, fmt.Errorf("unknown status %q", f.Status)
	}
	switch f.Kind {
	case "":
	case "slot":
		where = append(where, "slot")
	case "loan":
		where = append(where, "NOT slot")
	default:
		return ReservationPage{}, fmt.Errorf("unknown kind %q", f.Kind)
	}
	if f.OverdueDays > 0 {
		where = append(where, "status = 'active' AND elapsed >= "+arg(f.OverdueDays*86400))
	}
	if f.ItemID > 0 {
		where = append(where, "item_id = "+arg(f.ItemID))
	}
	if q := strings.TrimSpace(f.Query); q != "" {
		p := arg("%" + strings.ToLower(q) + "%")
		where = append(where, fmt.Sprintf("(LOWER(email) LIKE %[1]s OR LOWER(first_name || ' ' || last_name) LIKE %[1]s OR LOWER(last_name || ' ' || first_name) LIKE %[1]s)", p))
	}
	if f.From != nil {
		where = append(where, "start_date >= "+arg(f.From.UTC()))
	}
	if f.To != nil {
		where = append(where, "start_date < "+arg(f.To.UTC()))
	}

	query := adminReservationsCTE + `
	SELECT key, item_id, item_name, slot, club_name, email, first_name, last_name, start_date, end_date, status, elapsed,
	       cancelled_at, cancelled_by, cancel_kind, reminders, count(*) OVER ()
	FROM lines`
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	direction := "DESC"
	if f.Ascending {
		direction = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY start_date %s, key LIMIT %s OFFSET %s", direction, arg(f.Limit), arg(f.Offset))

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return ReservationPage{}, err
	}
	defer rows.Close()

	page := ReservationPage{Items: []AdminReservation{}, Limit: f.Limit, Offset: f.Offset}
	for rows.Next() {
		var r AdminReservation
		var end, cancelledAt sql.NullTime
		var elapsed sql.NullInt64
		var by, kind sql.NullString
		var reminders pq.StringArray
		if err := rows.Scan(&r.Key, &r.ItemID, &r.ItemName, &r.Slot, &r.ClubName, &r.Email, &r.FirstName, &r.LastName, &r.StartDate, &end,
			&r.Status, &elapsed, &cancelledAt, &by, &kind, &reminders, &page.Total); err != nil {
			return ReservationPage{}, err
		}
		r.StartDate = r.StartDate.UTC()
		if end.Valid {
			t := end.Time.UTC()
			r.EndDate = &t
		}
		if elapsed.Valid {
			r.ElapsedSeconds = &elapsed.Int64
			if r.Status == "active" {
				switch days := elapsed.Int64 / 86400; {
				case days >= 7:
					r.OverdueLevel = 7
				case days >= 5:
					r.OverdueLevel = 5
				case days >= 1:
					r.OverdueLevel = 1
				}
			}
		}
		if cancelledAt.Valid {
			t := cancelledAt.Time.UTC()
			r.CancelledAt = &t
		}
		r.CancelledBy, r.CancelKind = optional(by), optional(kind)
		r.Reminders = []string(reminders)
		if r.Reminders == nil {
			r.Reminders = []string{}
		}
		r.Cancellable = r.Status == "active" || r.Status == "ongoing" || r.Status == "upcoming"
		page.Items = append(page.Items, r)
	}
	return page, rows.Err()
}

// ReservationStats is the summary shown above the reservations.
type ReservationStats struct {
	ActiveLoans      int `json:"active_loans"`
	OngoingSlots     int `json:"ongoing_slots"`
	UpcomingSlots    int `json:"upcoming_slots"`
	OverdueDay1      int `json:"overdue_1d"`
	OverdueDay5      int `json:"overdue_5d"`
	OverdueDay7      int `json:"overdue_7d"`
	WaitlistEntries  int `json:"waitlist_entries"`
	AwaitingResponse int `json:"waitlist_awaiting_response"`
	Cancelled30Days  int `json:"cancelled_30d"`
	Items            int `json:"items"`
	UnavailableItems int `json:"unavailable_items"`
}

func (s *Service) Stats(ctx context.Context) (ReservationStats, error) {
	now := s.now().UTC()
	var st ReservationStats
	err := s.db.QueryRowContext(ctx, `
		SELECT count(*) FILTER (WHERE NOT COALESCE(re.slot, FALSE) AND r.end_date IS NULL),
		       count(*) FILTER (WHERE COALESCE(re.slot, FALSE) AND r.start_date <= $1::timestamp AND r.end_date > $1::timestamp),
		       count(*) FILTER (WHERE COALESCE(re.slot, FALSE) AND r.start_date > $1::timestamp),
		       count(*) FILTER (WHERE NOT COALESCE(re.slot, FALSE) AND r.end_date IS NULL AND r.start_date <= $1::timestamp - interval '1 day'),
		       count(*) FILTER (WHERE NOT COALESCE(re.slot, FALSE) AND r.end_date IS NULL AND r.start_date <= $1::timestamp - interval '5 days'),
		       count(*) FILTER (WHERE NOT COALESCE(re.slot, FALSE) AND r.end_date IS NULL AND r.start_date <= $1::timestamp - interval '7 days')
		FROM reservation r JOIN reservation_element re ON re.id_reservation_element = r.id_reservation_element`, now).
		Scan(&st.ActiveLoans, &st.OngoingSlots, &st.UpcomingSlots, &st.OverdueDay1, &st.OverdueDay5, &st.OverdueDay7)
	if err != nil {
		return st, err
	}
	err = s.db.QueryRowContext(ctx, `
		SELECT (SELECT count(*) FROM reservation_waitlist WHERE slot_start IS NULL OR slot_start > $1::timestamp),
		       (SELECT count(*) FROM reservation_waitlist WHERE notified_at IS NOT NULL AND (slot_start IS NULL OR slot_start > $1::timestamp)),
		       (SELECT count(*) FROM reservation_cancellations WHERE kind IN ('cancelled', 'cancelled_by_admin')
		                                                      AND cancelled_at > now() - interval '30 days'),
		       (SELECT count(*) FROM reservation_element),
		       (SELECT count(*) FROM reservation_element WHERE status <> 'active')`, now).
		Scan(&st.WaitlistEntries, &st.AwaitingResponse, &st.Cancelled30Days, &st.Items, &st.UnavailableItems)
	return st, err
}

// AdminWaitlistEntry is one place in a queue, with the person.
type AdminWaitlistEntry struct {
	ID         int64      `json:"id"`
	Position   int        `json:"position"`
	Email      string     `json:"email"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	SlotStart  *time.Time `json:"slot_start,omitempty"`
	JoinedAt   time.Time  `json:"joined_at"`
	NotifiedAt *time.Time `json:"notified_at,omitempty"`
	// TurnExpiresAt: when the person told that the resource is free loses their turn.
	TurnExpiresAt *time.Time `json:"turn_expires_at,omitempty"`
}

// AdminWaitlist is the queue of an item (no slot given, or an item without slot) or of one slot.
type AdminWaitlist struct {
	Entries []AdminWaitlistEntry `json:"entries"`
	// Free: the item or slot is free right now, which makes a notified person's answer the only thing in the way.
	Free bool `json:"free"`
	// AwaitingResponse: the first of a queue was told that it is free and has not answered.
	AwaitingResponse bool `json:"awaiting_response"`
}

// AdminWaitlist lists a queue in order of arrival. For a slot item without a slot given it lists the queues of every
// upcoming slot, each with its own positions.
func (s *Service) AdminWaitlist(ctx context.Context, itemID int, start *time.Time) (AdminWaitlist, error) {
	item, err := s.loadItem(ctx, itemID)
	if err != nil {
		return AdminWaitlist{}, err
	}
	var slot *time.Time
	if item.Slot && start != nil {
		u := start.UTC()
		slot = &u
	}
	out := AdminWaitlist{Entries: []AdminWaitlistEntry{}}
	if !item.Slot || slot != nil {
		if out.Free, err = isFree(ctx, s.db, item, slot); err != nil {
			return out, err
		}
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT w.id, row_number() OVER (PARTITION BY w.slot_start ORDER BY w.id), w.email, COALESCE(n.first_name, ''), COALESCE(n.last_name, ''),
		       w.slot_start, w.created_at, w.notified_at, w.turn_expires_at
		FROM reservation_waitlist w LEFT JOIN newf n ON n.email = w.email
		WHERE w.id_reservation_element = $1
		  AND (($2::timestamp IS NOT NULL AND w.slot_start = $2::timestamp)
		       OR ($2::timestamp IS NULL AND (w.slot_start IS NULL OR w.slot_start > $3::timestamp)))
		ORDER BY w.slot_start NULLS FIRST, w.id`, item.ID, startArg(slot), s.now().UTC())
	if err != nil {
		return out, err
	}
	defer rows.Close()
	for rows.Next() {
		var e AdminWaitlistEntry
		var slotStart, notified, expires sql.NullTime
		if err := rows.Scan(&e.ID, &e.Position, &e.Email, &e.FirstName, &e.LastName, &slotStart, &e.JoinedAt, &notified, &expires); err != nil {
			return out, err
		}
		if slotStart.Valid {
			t := slotStart.Time.UTC()
			e.SlotStart = &t
		}
		if notified.Valid {
			t := notified.Time
			e.NotifiedAt = &t
			out.AwaitingResponse = true
		}
		if expires.Valid {
			t := expires.Time
			e.TurnExpiresAt = &t
		}
		out.Entries = append(out.Entries, e)
	}
	return out, rows.Err()
}

// AdminRemoveWaitlistEntry takes somebody out of a queue. When they were the one being waited for, the next person is told.
func (s *Service) AdminRemoveWaitlistEntry(ctx context.Context, entryID int64) (bool, error) {
	var itemID int
	var slotStart, notified sql.NullTime
	err := s.db.QueryRowContext(ctx, `
		DELETE FROM reservation_waitlist WHERE id = $1 RETURNING id_reservation_element, slot_start, notified_at`, entryID).
		Scan(&itemID, &slotStart, &notified)
	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrNotQueued
	}
	if err != nil || !notified.Valid {
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
	return s.notifyNext(ctx, item, start)
}

// AdminAdvanceWaitlist passes the turn: a first person who was told and has not answered is dropped, then the next
// reachable person is told. It needs the item or slot to be free.
func (s *Service) AdminAdvanceWaitlist(ctx context.Context, itemID int, start *time.Time) (bool, error) {
	item, err := s.loadItem(ctx, itemID)
	if err != nil {
		return false, err
	}
	slot, err := s.target(item, start)
	if err != nil {
		return false, err
	}
	free, err := isFree(ctx, s.db, item, slot)
	if err != nil {
		return false, err
	}
	if !free {
		return false, ErrNotFree
	}
	if _, err := s.db.ExecContext(ctx, `
		DELETE FROM reservation_waitlist
		WHERE id = (SELECT id FROM reservation_waitlist
		            WHERE id_reservation_element = $1 AND slot_start IS NOT DISTINCT FROM $2::timestamp ORDER BY id LIMIT 1)
		  AND notified_at IS NOT NULL`, item.ID, startArg(slot)); err != nil {
		return false, err
	}
	return s.notifyNext(ctx, item, slot)
}

// ReturnResult describes a loan an administrator marked as returned.
type ReturnResult struct {
	Item         Item
	Owner        string
	Start        time.Time
	NextNotified bool
}

// MarkReturned ends the loan of an item without slot for somebody who gave it back without saying so in the app.
// The person is not notified (they did nothing wrong), the reminders stop with the loan, the item is free again and
// the first of the queue is told.
func (s *Service) MarkReturned(ctx context.Context, itemID int, actor string) (ReturnResult, error) {
	item, err := s.loadItem(ctx, itemID)
	if err != nil {
		return ReturnResult{}, err
	}
	res := ReturnResult{Item: item}
	if item.Slot {
		return res, ErrNotLoanItem
	}
	err = s.db.QueryRowContext(ctx, `
		UPDATE reservation
		SET end_date = GREATEST(now() AT TIME ZONE 'UTC', start_date + interval '1 second')
		WHERE id_reservation_element = $1 AND end_date IS NULL
		RETURNING email, start_date`, item.ID).Scan(&res.Owner, &res.Start)
	if errors.Is(err, sql.ErrNoRows) {
		return res, ErrNotFound
	}
	if err != nil {
		return res, err
	}
	res.Start = res.Start.UTC()
	s.logCancellation(ctx, item.ID, res.Owner, res.Start, s.now().UTC(), "returned_by_admin", actor)
	res.NextNotified, _ = s.notifyNext(ctx, item, nil)
	return res, nil
}

func (s *Service) logCancellation(ctx context.Context, itemID int, owner string, start, end time.Time, kind, by string) {
	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO reservation_cancellations (email, id_reservation_element, start_date, end_date, kind, cancelled_by)
		VALUES ($1, $2, $3, $4, $5, $6)`, owner, itemID, start.UTC(), end.UTC(), kind, by); err != nil {
		log.Printf("Reservation history of item %d could not be written: %v", itemID, err)
	}
}
