package models

import (
	"encoding/json"
	"testing"
)

func TestNavigationTargetValidate(t *testing.T) {
	cases := []struct {
		name    string
		target  *NavigationTarget
		wantErr bool
	}{
		{"event with id", EventNavigation(42), false},
		{"club with id", ClubNavigation(3), false},
		{"restaurant without id", RestaurantNavigation(), false},
		{"reservation with id", ReservationNavigation(7, map[string]string{"date": "2026-10-06"}), false},
		{"service with key", ServiceNavigation("laundry"), false},
		{"event without id", &NavigationTarget{Type: NavigationEvent}, true},
		{"unknown type", &NavigationTarget{Type: "nope", ID: "1"}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.target.Validate(); (err != nil) != tc.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestExpoDataWithoutNavigationKeepsData(t *testing.T) {
	p := NotificationPayload{Data: map[string]interface{}{"foo": "bar"}}
	data := p.ExpoData()
	if data["foo"] != "bar" || len(data) != 1 {
		t.Fatalf("unexpected data: %v", data)
	}
	if (NotificationPayload{}).ExpoData() != nil {
		t.Fatal("expected nil data for empty payload")
	}
}

func TestExpoDataEmbedsNavigation(t *testing.T) {
	p := NotificationPayload{
		Data:       map[string]interface{}{"foo": "bar"},
		Navigation: EventNavigation(42),
	}

	data := p.ExpoData()
	if data["foo"] != "bar" {
		t.Fatalf("custom data lost: %v", data)
	}
	if data["screen"] != "Events" {
		t.Fatalf("legacy screen = %v, want Events", data["screen"])
	}
	if _, mutated := p.Data["navigation"]; mutated {
		t.Fatal("ExpoData must not mutate the payload Data map")
	}

	raw, err := json.Marshal(data["navigation"])
	if err != nil {
		t.Fatal(err)
	}
	const want = `{"v":1,"type":"event","id":"42"}`
	if string(raw) != want {
		t.Fatalf("navigation JSON = %s, want %s", raw, want)
	}
}

func TestExpoDataServiceHasNoLegacyScreen(t *testing.T) {
	data := NotificationPayload{Navigation: ServiceNavigation("laundry")}.ExpoData()
	if _, ok := data["screen"]; ok {
		t.Fatalf("unexpected legacy screen: %v", data["screen"])
	}
}
