package models

import (
	"encoding/json"
	"strings"
	"testing"
)

// The app's tests (transat-app/tests/notifications/navigation.test.mjs, "the app understands what the backend
// sends") parse these exact strings: change one side and the other must follow.
func TestTheWireFormatOfEveryDestination(t *testing.T) {
	cases := map[string]struct {
		payload NotificationPayload
		want    string
	}{
		"event": {
			NotificationPayload{Kind: "event_created", Navigation: EventNavigation(7)},
			`{"navigation":{"v":1,"type":"event","id":"7"},"screen":"Events","type":"event_created"}`,
		},
		"club": {
			NotificationPayload{Kind: "admin", Navigation: ClubNavigation(3)},
			`{"navigation":{"v":1,"type":"club","id":"3"},"screen":"Clubs","type":"admin"}`,
		},
		"photo album": {
			NotificationPayload{Kind: "photo_album", Navigation: PhotoAlbumNavigation(42, "Gala")},
			`{"navigation":{"v":1,"type":"photo_album","id":"42","params":{"title":"Gala"}},"type":"photo_album"}`,
		},
		"restaurant": {
			NotificationPayload{Kind: "restaurant_menu", Navigation: RestaurantNavigation()},
			`{"navigation":{"v":1,"type":"restaurant"},"screen":"Restaurant","type":"restaurant_menu"}`,
		},
		"reservation": {
			NotificationPayload{Kind: "admin", Navigation: ReservationNavigation(11, map[string]string{"date": "2026-10-12"})},
			`{"navigation":{"v":1,"type":"reservation","id":"11","params":{"date":"2026-10-12"}},"screen":"Reservation","type":"admin"}`,
		},
		"service": {
			NotificationPayload{Kind: "admin", Navigation: ServiceNavigation("laundry")},
			`{"navigation":{"v":1,"type":"service","id":"laundry"},"type":"admin"}`,
		},
	}
	for name, c := range cases {
		raw, err := json.Marshal(c.payload.ExpoData())
		if err != nil || string(raw) != c.want {
			t.Errorf("%s:\n got %s\nwant %s", name, raw, c.want)
		}
	}
}

func TestImageURLMustBeAnAbsoluteHTTPSLink(t *testing.T) {
	valid := []string{
		"https://cdn.example.org/a/b.jpg",
		"https://cdn.example.org/a.jpg?w=400&q=80",
	}
	for _, raw := range valid {
		if err := ValidateImageURL(raw); err != nil {
			t.Errorf("%q refused: %v", raw, err)
		}
	}

	invalid := map[string]string{
		"plain http":    "http://cdn.example.org/a.jpg",
		"script":        "javascript:alert(1)",
		"file":          "file:///etc/passwd",
		"data":          "data:image/png;base64,AAAA",
		"relative":      "/a.jpg",
		"no host":       "https:///a.jpg",
		"credentials":   "https://user:pass@cdn.example.org/a.jpg",
		"whitespace":    "https://cdn.example.org/a b.jpg",
		"empty":         "",
		"too long":      "https://cdn.example.org/" + strings.Repeat("a", MaxNotificationImageURLLength),
		"newline split": "https://cdn.example.org/a.jpg\nhttps://evil.example.org",
	}
	for name, raw := range invalid {
		if err := ValidateImageURL(raw); err == nil {
			t.Errorf("%s: %q accepted", name, raw)
		}
	}
}

func TestOnlyImagesOnOurOwnStorageAreTrusted(t *testing.T) {
	const base = "https://cdn.example.org/"
	if got := OwnedImageURL("  https://cdn.example.org/events/a.jpg ", base); got != "https://cdn.example.org/events/a.jpg" {
		t.Fatalf("got %q", got)
	}
	for name, raw := range map[string]string{
		"empty":       "",
		"other host":  "https://evil.example.net/a.jpg",
		"look-alike":  "https://cdn.example.org.evil.net/a.jpg",
		"sub domain":  "https://x.cdn.example.org/a.jpg",
		"credentials": "https://cdn.example.org@evil.net/a.jpg",
		"http":        "http://cdn.example.org/a.jpg",
		"other port":  "https://cdn.example.org:8443/a.jpg",
		"not a link":  "poster",
		"script":      "javascript:alert(1)",
		"relative":    "/events/a.jpg",
	} {
		if got := OwnedImageURL(raw, base); got != "" {
			t.Errorf("%s: %q was trusted", name, raw)
		}
	}
	// Without a configured storage domain nothing is trusted.
	for _, base := range []string{"", "  ", "not a url"} {
		if got := OwnedImageURL("https://cdn.example.org/a.jpg", base); got != "" {
			t.Errorf("base %q trusted %q", base, got)
		}
	}
}

func TestExpoDataIsTheContractWithTheApp(t *testing.T) {
	if data := (NotificationPayload{}).ExpoData(); data != nil {
		t.Fatalf("an empty payload carries data: %v", data)
	}

	data := NotificationPayload{Kind: "photo_album", Navigation: PhotoAlbumNavigation(42, "Gala")}.ExpoData()
	nav, ok := data["navigation"].(*NavigationTarget)
	if data["type"] != "photo_album" || !ok || nav.Type != NavigationPhotoAlbum || nav.ID != "42" || nav.Version != NavigationVersion {
		t.Fatalf("data = %v", data)
	}
	if _, legacy := data["screen"]; legacy {
		t.Fatal("a destination without a legacy tab sent one")
	}

	event := NotificationPayload{Kind: "event_created", Navigation: EventNavigation(9)}.ExpoData()
	if event["screen"] != "Events" {
		t.Fatalf("legacy screen = %v", event["screen"])
	}

	// The type alone is still worth sending: it does not depend on a destination.
	if only := (NotificationPayload{Kind: "admin"}).ExpoData(); only["type"] != "admin" || only["navigation"] != nil {
		t.Fatalf("data = %v", only)
	}
}

func TestNavigationTargetsUseStableIDs(t *testing.T) {
	cases := map[string]*NavigationTarget{
		"event":       EventNavigation(5),
		"club":        ClubNavigation(3),
		"reservation": ReservationNavigation(11, map[string]string{"date": "2026-10-12"}),
		"photo_album": PhotoAlbumNavigation(8, "Title"),
	}
	for name, target := range cases {
		if err := target.Validate(); err != nil {
			t.Errorf("%s: %v", name, err)
		}
		if target.ID == "" {
			t.Errorf("%s has no id", name)
		}
	}
	if err := (&NavigationTarget{Type: NavigationPhotoAlbum, ID: "my-album"}).Validate(); err == nil {
		t.Error("an album addressed by name was accepted")
	}
	if err := (&NavigationTarget{Type: NavigationEvent}).Validate(); err == nil {
		t.Error("an event without id was accepted")
	}
}
