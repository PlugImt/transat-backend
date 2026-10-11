package photos

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata" // Europe/Paris must resolve in minimal containers
)

type TimeSource string

const (
	TimeExif         TimeSource = "exif"
	TimeFilename     TimeSource = "filename"
	TimeEventDate    TimeSource = "event_date"
	TimeDriveCreated TimeSource = "drive_created"
	TimeDriveUpdated TimeSource = "drive_modified"
)

// Window around the event date in which a capture time is trusted (night events end after midnight,
// photographers sometimes deliver days later).
const (
	suspectBefore = 2 * 24 * time.Hour
	suspectAfter  = 30 * 24 * time.Hour
)

var (
	parisLoc, _ = time.LoadLocation("Europe/Paris")

	exifLayout   = "2006:01:02 15:04:05"
	nameFullRe   = regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2}) at (\d{2})\.(\d{2})\.(\d{2})`)
	nameCompact  = regexp.MustCompile(`(?:^|[^\d])(\d{4})(\d{2})(\d{2})[_-](\d{2})(\d{2})(\d{2})(?:[^\d]|$)`)
	nameDateOnly = regexp.MustCompile(`^(?:IMG|VID)-(\d{4})(\d{2})(\d{2})-WA\d+`)
)

// PhotoTimes are the raw date signals known for one file.
type PhotoTimes struct {
	ExifTime      string // Drive imageMediaMetadata.time, camera wall clock without zone
	Name          string
	EventDate     *time.Time // from the album, when known
	DriveCreated  time.Time
	DriveModified time.Time
}

type ResolvedTimes struct {
	TakenAt      *time.Time // capture time (EXIF or filename); never derived from Drive/event dates
	TakenSource  TimeSource
	TakenSuspect bool // far from the event date, kept for the record but not used for sorting
	SortAt       time.Time
	SortAtSource TimeSource
}

// ResolveTimes applies the fallback chain: EXIF, filename timestamp, event date, Drive creation, Drive modification.
// Camera clocks carry no zone, so they are read as Europe/Paris wall time.
func ResolveTimes(in PhotoTimes, now time.Time) ResolvedTimes {
	var res ResolvedTimes
	if t, ok := parseExif(in.ExifTime, now); ok {
		res.TakenAt, res.TakenSource = &t, TimeExif
	} else if t, ok := parseFilenameTime(in.Name, now); ok {
		res.TakenAt, res.TakenSource = &t, TimeFilename
	}

	if res.TakenAt != nil && in.EventDate != nil {
		start := time.Date(in.EventDate.Year(), in.EventDate.Month(), in.EventDate.Day(), 0, 0, 0, 0, parisLoc)
		if res.TakenAt.Before(start.Add(-suspectBefore)) || res.TakenAt.After(start.Add(suspectAfter)) {
			res.TakenSuspect = true
		}
	}

	switch {
	case res.TakenAt != nil && !res.TakenSuspect:
		res.SortAt, res.SortAtSource = *res.TakenAt, res.TakenSource
	case in.EventDate != nil:
		d := in.EventDate
		res.SortAt, res.SortAtSource = time.Date(d.Year(), d.Month(), d.Day(), 12, 0, 0, 0, parisLoc), TimeEventDate
	case res.TakenAt != nil:
		res.SortAt, res.SortAtSource = *res.TakenAt, res.TakenSource
	case !in.DriveCreated.IsZero():
		res.SortAt, res.SortAtSource = in.DriveCreated, TimeDriveCreated
	default:
		res.SortAt, res.SortAtSource = in.DriveModified, TimeDriveUpdated
	}
	return res
}

func plausible(t, now time.Time) bool {
	return t.Year() >= 2010 && t.Before(now.Add(24*time.Hour))
}

func parseExif(s string, now time.Time) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	t, err := time.ParseInLocation(exifLayout, s, parisLoc)
	if err != nil || !plausible(t, now) {
		return time.Time{}, false
	}
	return t, true
}

func parseFilenameTime(name string, now time.Time) (time.Time, bool) {
	var m []string
	if m = nameFullRe.FindStringSubmatch(name); m == nil {
		m = nameCompact.FindStringSubmatch(name)
	}
	if m != nil {
		n := make([]int, 6)
		for i := range n {
			n[i], _ = strconv.Atoi(m[i+1])
		}
		if _, ok := validDate(n[0], n[1], n[2]); ok {
			t := time.Date(n[0], time.Month(n[1]), n[2], n[3], n[4], n[5], 0, parisLoc)
			if t.Hour() == n[3] && plausible(t, now) {
				return t, true
			}
		}
		return time.Time{}, false
	}
	if m = nameDateOnly.FindStringSubmatch(name); m != nil {
		y, _ := strconv.Atoi(m[1])
		mo, _ := strconv.Atoi(m[2])
		d, _ := strconv.Atoi(m[3])
		if _, ok := validDate(y, mo, d); ok {
			t := time.Date(y, time.Month(mo), d, 12, 0, 0, 0, parisLoc)
			if plausible(t, now) {
				return t, true
			}
		}
	}
	return time.Time{}, false
}
