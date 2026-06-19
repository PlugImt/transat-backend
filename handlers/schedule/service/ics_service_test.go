package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const sampleICS = `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//EN
BEGIN:VEVENT
UID:PASS-322849@pass.example.com
SUMMARY:PARCOURS FISA  - Cours1
DTSTART;TZID=Europe/Paris:20250903T082000
DTEND;TZID=Europe/Paris:20250903T120000
LOCATION:
END:VEVENT
BEGIN:VEVENT
UID:PASS-315288@pass.example.com
SUMMARY:Evenements FISE A1S5 - N  - Theatre
DTSTART;TZID=Europe/Paris:20250903T150000
DTEND;TZID=Europe/Paris:20250903T180000
LOCATION:NA-Espace Manifestation du Forum - (A015)
END:VEVENT
BEGIN:VEVENT
UID:PASS-999999@pass.example.com
SUMMARY:Cancelled Event
STATUS:CANCELLED
DTSTART:20250903T100000Z
DTEND:20250903T110000Z
END:VEVENT
END:VCALENDAR`

func TestParseICS(t *testing.T) {
	data, err := ParseICS([]byte(sampleICS))
	if err != nil {
		t.Fatalf("ParseICS() error = %v", err)
	}

	events, ok := data["2025-09-03"]
	if !ok {
		t.Fatalf("expected date key 2025-09-03, got keys: %v", data)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	if events[0].ID != "PASS-322849" {
		t.Errorf("events[0].ID = %q, want PASS-322849", events[0].ID)
	}
	if events[0].Name != "PARCOURS FISA  - Cours1" {
		t.Errorf("events[0].Name = %q", events[0].Name)
	}
	if events[0].StartTime != "08:20" {
		t.Errorf("events[0].StartTime = %q, want 08:20", events[0].StartTime)
	}
	if events[0].EndTime != "12:00" {
		t.Errorf("events[0].EndTime = %q, want 12:00", events[0].EndTime)
	}
	if events[0].Location != "-" {
		t.Errorf("events[0].Location = %q, want -", events[0].Location)
	}

	if events[1].ID != "PASS-315288" {
		t.Errorf("events[1].ID = %q, want PASS-315288", events[1].ID)
	}
	if events[1].Location != "NA-Espace Manifestation du Forum - (A015)" {
		t.Errorf("events[1].Location = %q", events[1].Location)
	}
}

func TestExtractEventID(t *testing.T) {
	tests := []struct {
		uid  string
		want string
	}{
		{"PASS-322849@pass.example.com", "PASS-322849"},
		{"some-id@example.com", "some-id"},
		{"plain-id", "plain-id"},
	}

	for _, tc := range tests {
		if got := extractEventID(tc.uid); got != tc.want {
			t.Errorf("extractEventID(%q) = %q, want %q", tc.uid, got, tc.want)
		}
	}
}

func TestFetchICS_RetryOnServerError(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("BEGIN:VCALENDAR\nEND:VCALENDAR"))
	}))
	defer server.Close()

	svc := &IcsService{
		client: server.Client(),
	}

	data, err := svc.FetchICS(server.URL)
	if err != nil {
		t.Fatalf("FetchICS() error = %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty ICS data")
	}
}

func TestValidateICSURL(t *testing.T) {
	if err := validateICSURL("ftp://example.com/cal.ics"); err == nil {
		t.Fatal("expected error for ftp scheme")
	}
	if err := validateICSURL("not-a-url"); err == nil {
		t.Fatal("expected error for invalid url")
	}
	if err := validateICSURL("https://example.com/cal.ics"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
