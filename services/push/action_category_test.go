package push

import (
	"testing"

	"github.com/plugimt/transat-backend/models"
)

func TestAnActionCategoryReachesTheProviderOnlyForTheMessagesThatHaveOne(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESERVATIONS")
	token := f.device(1, "alice@x.test", "ios", true)

	plain := f.message(models.NotificationReservations, "plain")
	f.enqueue(plain)
	f.run()
	if got := f.fp.sent(); len(got) != 1 || got[0].ActionCategory != "" {
		t.Fatalf("plain message: %+v", got)
	}
	if _, has := expoMessage(f.fp.sent()[0])["categoryId"]; has {
		t.Fatal("a plain notification asks the phone for action buttons")
	}

	withActions := f.message(models.NotificationReservations, "actions")
	withActions.ActionCategory = "reservation_waitlist"
	id := f.enqueue(withActions)
	f.run()
	sent := f.fp.sent()
	if len(sent) != 2 || sent[1].Token != token || sent[1].ActionCategory != "reservation_waitlist" {
		t.Fatalf("sent = %+v", sent)
	}
	if got := expoMessage(sent[1])["categoryId"]; got != "reservation_waitlist" {
		t.Fatalf("categoryId = %v", got)
	}
	if n := f.count(`SELECT count(*) FROM push_messages WHERE id = $1 AND action_category = 'reservation_waitlist'`, id); n != 1 {
		t.Fatal("the category is not recorded with the message")
	}
}

func TestAnInvalidActionCategoryIsRefused(t *testing.T) {
	f := newFx(t)
	for _, bad := range []string{"has space", "semi;colon", "a/b", string(make([]byte, 65))} {
		m := f.message(models.NotificationReservations, "k")
		m.ActionCategory = bad
		if _, err := f.p.Enqueue(f.ctx, m); err == nil {
			t.Fatalf("%q was accepted", bad)
		}
	}
	if n := f.count(`SELECT count(*) FROM push_messages`); n != 0 {
		t.Fatalf("%d invalid messages were stored", n)
	}
}
