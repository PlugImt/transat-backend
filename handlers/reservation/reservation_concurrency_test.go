package reservation_test

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/plugimt/transat-backend/services/reservations"
)

// firstComeFirstServed runs the same request for several users at once and returns the status codes.
func (a *api) race(method, path string, body any, users []string) []int {
	statuses := make([]int, len(users))
	var wg sync.WaitGroup
	for i, user := range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			statuses[i] = a.call(method, path, user, body).status
		}()
	}
	wg.Wait()
	return statuses
}

func count(statuses []int, want int) int {
	n := 0
	for _, s := range statuses {
		if s == want {
			n++
		}
	}
	return n
}

func TestConcurrentBookingsOfTheSameSlotOrItemLeaveExactlyOneReservation(t *testing.T) {
	a := newAPI(t)
	var users []string
	for i := 0; i < 8; i++ {
		email := fmt.Sprintf("u%d@x.test", i)
		a.user(email, "U", fmt.Sprint(i))
		users = append(users, email)
	}
	room, cable := a.item("Salle", true), a.item("Câble", false)
	_, start := slot()

	got := a.race(http.MethodPatch, fmt.Sprintf("/reservation/item/%d", room), map[string]string{"start_date": start}, users)
	if count(got, 200) != 1 || count(got, 409) != 7 || a.rows(room) != 1 {
		t.Fatalf("slot: statuses %v, %d rows", got, a.rows(room))
	}
	now := time.Now().UTC().Add(-time.Minute).Format(reservations.APITimeLayout)
	got = a.race(http.MethodPatch, fmt.Sprintf("/reservation/item/%d", cable), map[string]string{"start_date": now}, users)
	if count(got, 200) != 1 || count(got, 409) != 7 || a.rows(cable) != 1 {
		t.Fatalf("loan: statuses %v, %d rows", got, a.rows(cable))
	}
}

func TestADirectBookingAndTheReserveButtonNeverBothWin(t *testing.T) {
	a := newAPI(t)
	a.user("owner@x.test", "O", "X")
	a.user("q1@x.test", "Q", "One")
	a.user("fast@x.test", "F", "X")
	room := a.item("Salle", true)
	_, start := slot()
	item := fmt.Sprintf("/reservation/item/%d", room)
	body := map[string]string{"start_date": start}

	a.call(http.MethodPatch, item, "owner@x.test", body)
	a.call(http.MethodPost, item+"/waitlist", "q1@x.test", body)
	if r := a.call(http.MethodDelete, item, "owner@x.test", body); r.status != 200 {
		t.Fatalf("release: %d", r.status)
	}

	var wg sync.WaitGroup
	var reserve, direct int
	wg.Add(2)
	go func() {
		defer wg.Done()
		reserve = a.call(http.MethodPost, item+"/waitlist/reserve", "q1@x.test", body).status
	}()
	go func() { defer wg.Done(); direct = a.call(http.MethodPatch, item, "fast@x.test", body).status }()
	wg.Wait()

	if (reserve == 200) == (direct == 200) || a.rows(room) != 1 {
		t.Fatalf("reserve %d, direct %d, %d reservations", reserve, direct, a.rows(room))
	}
	if reserve != 200 && reserve != 409 || direct != 200 && direct != 409 {
		t.Fatalf("unexpected statuses: reserve %d, direct %d", reserve, direct)
	}
}
