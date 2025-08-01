package utils

import (
	"time"
)

var (
	ParisLocation *time.Location
)

func init() {
	var err error
	ParisLocation, err = time.LoadLocation("Europe/Paris")
	if err != nil {
		// Fallback to UTC if Paris timezone can't be loaded
		LogMessage(LevelError, "Failed to load Paris timezone, falling back to UTC")
		ParisLocation = time.UTC
	}
}

// Now returns the current time in Paris timezone
// This should replace all time.Now() calls in the application
func Now() time.Time {
	return time.Now().In(ParisLocation)
}

// ParseInParis parses a time string and converts it to Paris timezone
func ParseInParis(layout, value string) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(ParisLocation), nil
}

// ParseTimeInParis parses a time with a specific location and converts to Paris timezone
func ParseTimeInParis(layout, value string, loc *time.Location) (time.Time, error) {
	t, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(ParisLocation), nil
}

// ToParisTime converts any time to Paris timezone
func ToParisTime(t time.Time) time.Time {
	return t.In(ParisLocation)
}

// StartOfDayParis returns the start of the day (00:00:00) in Paris timezone
func StartOfDayParis(t time.Time) time.Time {
	parisTime := t.In(ParisLocation)
	return time.Date(parisTime.Year(), parisTime.Month(), parisTime.Day(), 0, 0, 0, 0, ParisLocation)
}

// EndOfDayParis returns the end of the day (23:59:59.999999999) in Paris timezone
func EndOfDayParis(t time.Time) time.Time {
	parisTime := t.In(ParisLocation)
	return time.Date(parisTime.Year(), parisTime.Month(), parisTime.Day(), 23, 59, 59, 999999999, ParisLocation)
}

// TodayInParis returns today's date range in Paris timezone
func TodayInParis() (start, end time.Time) {
	now := Now()
	start = StartOfDayParis(now)
	end = start.Add(24 * time.Hour)
	return start, end
}

// FormatParis formats time in Paris timezone with the given layout
func FormatParis(t time.Time, layout string) string {
	return t.In(ParisLocation).Format(layout)
}

// UnixMilliParis returns Unix timestamp in milliseconds for Paris time
func UnixMilliParis(t time.Time) int64 {
	return t.In(ParisLocation).UnixMilli()
}

// UnixNanoParis returns Unix timestamp in nanoseconds for Paris time
func UnixNanoParis(t time.Time) int64 {
	return t.In(ParisLocation).UnixNano()
}

// IsWeekdayParis checks if the given time is a weekday in Paris timezone
func IsWeekdayParis(t time.Time) bool {
	parisTime := t.In(ParisLocation)
	weekday := parisTime.Weekday()
	return weekday != time.Saturday && weekday != time.Sunday
}

// GetHourParis returns the hour in Paris timezone
func GetHourParis(t time.Time) int {
	return t.In(ParisLocation).Hour()
}

func GetTimeInParis(t time.Time) time.Time {
	return t.In(ParisLocation)
}

// GetWeekdayParis returns the weekday in Paris timezone
func GetWeekdayParis(t time.Time) time.Weekday {
	return t.In(ParisLocation).Weekday()
}

// GetYearParis returns the year in Paris timezone
func GetYearParis(t time.Time) int {
	return t.In(ParisLocation).Year()
}

// SinceInParis returns the duration since the given time, using Paris timezone for calculations
func SinceInParis(t time.Time) time.Duration {
	return Now().Sub(t.In(ParisLocation))
}

// AddInParis adds duration to time and ensures result is in Paris timezone
func AddInParis(t time.Time, d time.Duration) time.Time {
	return t.In(ParisLocation).Add(d)
}

func ParseDateTimeInParis(layout, value string) (time.Time, error) {
	t, err := time.ParseInLocation(layout, value, ParisLocation)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(ParisLocation), nil
}
