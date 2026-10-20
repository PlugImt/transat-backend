package services

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services/push"
)

func reservationAnnouncement(emails ...string) Announcement {
	return Announcement{
		Kind:           "reservation_waitlist",
		Category:       models.NotificationReservations,
		DedupeKey:      "k",
		Audience:       usersOf(emails...),
		ActionCategory: "reservation_waitlist",
		Navigation:     &models.NavigationTarget{Type: models.NavigationReservation, ID: "7", Params: map[string]string{"waitlist": "1"}},
		Compose: func(lang string, _ *i18n.Localizer) (string, string) {
			return "title " + lang, "body " + lang
		},
	}
}

func usersOf(emails ...string) push.Audience {
	return push.Audience{NotificationAudience: models.NotificationAudience{Type: models.AudienceUsers, Emails: emails}}
}

func TestWithoutThePipelineAnAnnouncementGoesStraightToExpoWithItsActionsAndLanguages(t *testing.T) {
	f := newFixture(t)
	f.user("fr@imt.fr", 1, false, tok("fr"))
	f.user("en@imt.fr", 2, false, tok("en"))
	f.user("off@imt.fr", 1, false, tok("off"))
	f.user("nodevice@imt.fr", 1, false)
	if err := f.ns.SetPreference("off@imt.fr", models.NotificationReservations, false); err != nil {
		t.Fatal(err)
	}
	f.user("other@imt.fr", 1, false, tok("other"))
	if f.ns.OutboxEnabled() {
		t.Fatal("the fixture must run without the pipeline")
	}

	if _, err := f.ns.Announce(context.Background(), reservationAnnouncement("fr@imt.fr", "en@imt.fr", "off@imt.fr", "nodevice@imt.fr")); err != nil {
		t.Fatal(err)
	}

	if got, want := strings.Join(f.ex.tokens(), ","), tok("en")+","+tok("fr"); got != want {
		t.Fatalf("addressed %s, want %s: only the chosen users who kept the category and have a device", got, want)
	}
	for _, msg := range f.ex.messages() {
		lang := "fr"
		if msg["to"] == tok("en") {
			lang = "en"
		}
		if msg["title"] != "title "+lang || msg["body"] != "body "+lang {
			t.Fatalf("%v: not in the user's language", msg)
		}
		if msg["categoryId"] != "reservation_waitlist" {
			t.Fatalf("%v: the action buttons are missing", msg)
		}
		data, _ := msg["data"].(map[string]any)
		nav, _ := data["navigation"].(map[string]any)
		if data["type"] != "reservation_waitlist" || nav["type"] != "reservation" || data["screen"] != "Reservation" {
			t.Fatalf("data = %v: a tap must still open the offer", data)
		}
	}
}

func TestWithoutThePipelineCanReachFollowsTheDevicesAndThePreference(t *testing.T) {
	f := newFixture(t)
	f.user("ok@imt.fr", 1, false, tok("ok"))
	f.user("off@imt.fr", 1, false, tok("off"))
	f.user("nodevice@imt.fr", 1, false)
	if err := f.ns.SetPreference("off@imt.fr", models.NotificationReservations, false); err != nil {
		t.Fatal(err)
	}
	for email, want := range map[string]bool{"ok@imt.fr": true, "off@imt.fr": false, "nodevice@imt.fr": false, "nobody@imt.fr": false} {
		got, err := f.ns.CanReach(context.Background(), email, models.NotificationReservations)
		if err != nil || got != want {
			t.Fatalf("CanReach(%s) = %v, %v; want %v", email, got, err, want)
		}
	}
}

func TestWithoutThePipelineAFailedSendIsAnErrorSoTheCallerCanRetryOrPassTheTurn(t *testing.T) {
	f := newFixture(t)
	f.user("fr@imt.fr", 1, false, tok("fr"))
	f.ex.statuses = []int{http.StatusInternalServerError, http.StatusInternalServerError, http.StatusInternalServerError}

	if _, err := f.ns.Announce(context.Background(), reservationAnnouncement("fr@imt.fr")); err == nil {
		t.Fatal("an announcement Expo never accepted reported success")
	}
	// Nobody to send to is not a failure.
	if _, err := f.ns.Announce(context.Background(), reservationAnnouncement("nobody@imt.fr")); err != nil {
		t.Fatalf("no recipient: %v", err)
	}
}
