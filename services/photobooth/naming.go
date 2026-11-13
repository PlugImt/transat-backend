// Package photobooth stores the photobooth's pictures in the shared Photos drive and issues the permanent share
// links shown as QR codes on the booth.
package photobooth

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	_ "time/tzdata" // the container image may ship without a zone database

	"github.com/plugimt/transat-backend/services/photos"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	// PendingFolderName is the sub-folder of an event folder that receives the booth's pictures. The gallery
	// never indexes it: the pictures reach the app once an administrator or the Club Photo approves them.
	PendingFolderName = photos.PendingPhotoboothFolder
	// PhotoboothFolderName is what the pending folder becomes once approved.
	PhotoboothFolderName = "Photobooth"
	maxEventNameRunes    = 100
)

var (
	yearFolderRe  = regexp.MustCompile(`^(\d{4})\s*[-–/]\s*(\d{4})$`)
	eventFolderRe = regexp.MustCompile(`^(\d{4})\s*[.\-/]\s*(\d{1,2})\s*[.\-/]\s*(\d{1,2})(?:\s*[-–—:]\s*|\s+|$)(.*)$`)
	spacesRe      = regexp.MustCompile(`\s+`)
	accents       = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
)

// SchoolYear names the school-year folder (August to July) holding an event: 2026-03-27 → "2025-2026".
func SchoolYear(date time.Time) string {
	start := date.Year()
	if date.Month() < time.August {
		start--
	}
	return fmt.Sprintf("%d-%d", start, start+1)
}

// EventFolderName is the name given to a new event folder, the strict "YYYY.MM.DD - Name" form the photo
// classifier ranks as certain.
func EventFolderName(date time.Time, name string) string {
	return fmt.Sprintf("%04d.%02d.%02d - %s", date.Year(), int(date.Month()), date.Day(), name)
}

// CleanEventName makes an event name safe as a Drive folder name: no control characters, no slash (the gallery
// uses "/" in its paths), single spaces, at most maxEventNameRunes characters. Empty means unusable.
func CleanEventName(name string) string {
	var b strings.Builder
	for _, r := range name {
		switch {
		case r == '/' || r == '\\':
			b.WriteRune('-')
		case unicode.IsControl(r) || r == unicode.ReplacementChar:
			b.WriteRune(' ')
		default:
			b.WriteRune(r)
		}
	}
	s := strings.TrimSpace(spacesRe.ReplaceAllString(b.String(), " "))
	if r := []rune(s); len(r) > maxEventNameRunes {
		s = strings.TrimSpace(string(r[:maxEventNameRunes]))
	}
	return s
}

// NameKey is the comparison form of a folder or event name: lower case, no accents, single spaces.
func NameKey(name string) string {
	s, _, err := transform.String(accents, name)
	if err != nil {
		s = name
	}
	return spacesRe.ReplaceAllString(strings.ToLower(strings.TrimSpace(s)), " ")
}

// IsSchoolYearFolder tells whether a folder name designates the given school year ("2025-2026", "2025 - 2026"…).
func IsSchoolYearFolder(folderName, year string) bool {
	m := yearFolderRe.FindStringSubmatch(strings.TrimSpace(folderName))
	return m != nil && m[1]+"-"+m[2] == year
}

// MatchesEventFolder tells whether an existing folder is the one of this event: same day (with or without zero
// padding, "." "-" or "/" separators) and same name ignoring case, accents and spacing.
func MatchesEventFolder(folderName string, date time.Time, eventName string) bool {
	m := eventFolderRe.FindStringSubmatch(strings.TrimSpace(folderName))
	if m == nil {
		return false
	}
	y, _ := strconv.Atoi(m[1])
	mo, _ := strconv.Atoi(m[2])
	d, _ := strconv.Atoi(m[3])
	if y != date.Year() || mo != int(date.Month()) || d != date.Day() {
		return false
	}
	return NameKey(m[4]) == NameKey(eventName)
}

// FileName is the Drive name of one picture: sortable, unique per shoot, readable by the organisers. attempt
// counts the uploads of the slot (retakes), so no two uploads ever share a name.
func FileName(shootStart time.Time, sessionID, kind string, number, attempt int) string {
	short := sessionID
	if len(short) > 8 {
		short = short[:8]
	}
	base := shootStart.In(paris()).Format("2006-01-02_15-04-05") + "_" + short + "_"
	if kind == KindCollage {
		base += "collage"
	} else {
		base += fmt.Sprintf("photo-%d", number)
	}
	if attempt > 1 {
		base += fmt.Sprintf("_v%d", attempt)
	}
	return base + ".jpg"
}

func paris() *time.Location {
	if loc, err := time.LoadLocation("Europe/Paris"); err == nil {
		return loc
	}
	return time.UTC
}
