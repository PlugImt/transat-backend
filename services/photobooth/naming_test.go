package photobooth

import (
	"testing"
	"time"
)

func day(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSchoolYearRunsFromAugustToJuly(t *testing.T) {
	cases := map[string]string{
		"2026-03-27": "2025-2026",
		"2025-07-31": "2024-2025",
		"2025-08-01": "2025-2026",
		"2026-10-27": "2026-2027",
		"2027-01-01": "2026-2027",
	}
	for date, want := range cases {
		if got := SchoolYear(day(date)); got != want {
			t.Errorf("SchoolYear(%s) = %s, want %s", date, got, want)
		}
	}
}

func TestEventFolderNameIsTheStrictClassifierForm(t *testing.T) {
	if got := EventFolderName(day("2026-03-07"), "Gala"); got != "2026.03.07 - Gala" {
		t.Fatalf("got %q", got)
	}
}

func TestMatchesEventFolderAcceptsExistingSpellings(t *testing.T) {
	d := day("2026-03-07")
	for _, name := range []string{"2026.03.07 - Soirée Intégration", "2026.3.7 - soiree   integration", "2026-03-07 Soiree Integration"} {
		if !MatchesEventFolder(name, d, "Soirée intégration") {
			t.Errorf("%q should match", name)
		}
	}
	for _, name := range []string{"2026.03.08 - Soirée Intégration", "2026.03.07 - Gala", "Soirée Intégration"} {
		if MatchesEventFolder(name, d, "Soirée intégration") {
			t.Errorf("%q should not match", name)
		}
	}
}

func TestIsSchoolYearFolder(t *testing.T) {
	if !IsSchoolYearFolder("2025 - 2026", "2025-2026") || !IsSchoolYearFolder("2025-2026", "2025-2026") {
		t.Fatal("year folder not recognised")
	}
	if IsSchoolYearFolder("2024-2025", "2025-2026") || IsSchoolYearFolder("Photos 2025-2026", "2025-2026") {
		t.Fatal("wrong year folder accepted")
	}
}

func TestCleanEventName(t *testing.T) {
	cases := map[string]string{
		"  Gala\t2026 ":        "Gala 2026",
		"AG/BDE":               "AG-BDE",
		"a\x00b":               "a b",
		"":                     "",
		"   ":                  "",
		"Soirée d'intégration": "Soirée d'intégration",
	}
	for in, want := range cases {
		if got := CleanEventName(in); got != want {
			t.Errorf("CleanEventName(%q) = %q, want %q", in, got, want)
		}
	}
	long := ""
	for i := 0; i < 300; i++ {
		long += "é"
	}
	if got := []rune(CleanEventName(long)); len(got) != maxEventNameRunes {
		t.Errorf("long name kept %d runes", len(got))
	}
}

func TestFileNamesAreUniquePerAttempt(t *testing.T) {
	start := time.Date(2026, 11, 10, 20, 4, 33, 0, time.UTC)
	if got := FileName(start, "abcdef1234567", KindPhoto, 2, 1); got != "2026-11-10_21-04-33_abcdef12_photo-2.jpg" {
		t.Errorf("got %q", got)
	}
	if got := FileName(start, "abcdef12", KindPhoto, 2, 3); got != "2026-11-10_21-04-33_abcdef12_photo-2_v3.jpg" {
		t.Errorf("got %q", got)
	}
	if got := FileName(start, "abcdef12", KindCollage, 0, 1); got != "2026-11-10_21-04-33_abcdef12_collage.jpg" {
		t.Errorf("got %q", got)
	}
}
