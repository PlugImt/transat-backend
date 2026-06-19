package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/plugimt/transat-backend/handlers/schedule/repository"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
	"github.com/sethvargo/go-retry"
)

const (
	fetchTimeout    = 30 * time.Second
	maxFetchRetries = 2                // 3 total attempts (initial + 2 retries)
	maxICSSize      = 10 * 1024 * 1024 // 10MB
)

var passIDRegex = regexp.MustCompile(`PASS-\d+`)

type IcsService struct {
	repo   *repository.UserScheduleRepository
	client *http.Client
}

func NewIcsService(db *sql.DB) *IcsService {
	return &IcsService{
		repo: repository.NewUserScheduleRepository(db),
		client: &http.Client{
			Timeout: fetchTimeout,
		},
	}
}

func (s *IcsService) SyncUserSchedule(userID int, icsURL string) error {
	utils.LogHeader("📅 ICS Sync")
	utils.LogLineKeyValue(utils.LevelInfo, "UserID", userID)

	data, err := s.FetchICS(icsURL)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch ICS calendar")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return err
	}

	calendarData, err := ParseICS(data)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse ICS calendar")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return err
	}

	raw, err := json.Marshal(calendarData)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to marshal calendar data")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return err
	}

	if err := s.repo.UpdateCalendarData(userID, raw); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to save calendar data")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return err
	}

	utils.LogMessage(utils.LevelInfo, "ICS sync completed successfully")
	utils.LogFooter()
	return nil
}

func (s *IcsService) SyncAll() error {
	schedules, err := s.repo.ListAllWithIcsURL()
	if err != nil {
		return fmt.Errorf("list schedules with ics_url: %w", err)
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Starting ICS sync for %d users", len(schedules)))

	var syncErrs []error
	for _, schedule := range schedules {
		if err := s.SyncUserSchedule(schedule.UserID, schedule.IcsURL); err != nil {
			syncErrs = append(syncErrs, err)
		}
	}

	return errors.Join(syncErrs...)
}

func (s *IcsService) FetchICS(icsURL string) ([]byte, error) {
	if err := validateICSURL(icsURL); err != nil {
		return nil, err
	}

	var body []byte
	backoff := retry.WithMaxRetries(maxFetchRetries, retry.NewExponential(1*time.Second))

	err := retry.Do(context.Background(), backoff, func(ctx context.Context) error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, icsURL, nil)
		if err != nil {
			return err
		}

		resp, err := s.client.Do(req)
		if err != nil {
			return retry.RetryableError(fmt.Errorf("http request failed: %w", err))
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= http.StatusInternalServerError {
			io.Copy(io.Discard, resp.Body) //nolint:errcheck
			return retry.RetryableError(fmt.Errorf("unexpected status code %d", resp.StatusCode))
		}
		if resp.StatusCode != http.StatusOK {
			io.Copy(io.Discard, resp.Body) //nolint:errcheck
			return fmt.Errorf("unexpected status code %d", resp.StatusCode)
		}

		data, err := io.ReadAll(io.LimitReader(resp.Body, maxICSSize))
		if err != nil {
			return retry.RetryableError(fmt.Errorf("read response body: %w", err))
		}
		if int64(len(data)) == maxICSSize {
			return fmt.Errorf("ics response exceeds maximum size of %d bytes", maxICSSize)
		}
		body = data
		return nil
	})
	if err != nil {
		return nil, err
	}

	return body, nil
}

func validateICSURL(icsURL string) error {
	parsed, err := url.Parse(icsURL)
	if err != nil {
		return fmt.Errorf("invalid ics url: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("ics url must use http or https scheme")
	}
	if parsed.Host == "" {
		return fmt.Errorf("ics url must have a host")
	}
	return nil
}

func ParseICS(data []byte) (models.CalendarData, error) {
	cal, err := ics.ParseCalendar(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("parse calendar: %w", err)
	}

	result := make(models.CalendarData)

	for _, event := range cal.Events() {
		if isCancelled(event) {
			continue
		}

		start, err := event.GetStartAt()
		if err != nil {
			continue
		}
		end, err := event.GetEndAt()
		if err != nil {
			continue
		}

		startParis := utils.ToParisTime(start)
		endParis := utils.ToParisTime(end)

		dateKey := utils.FormatParis(startParis, "2006-01-02")
		name := propertyValue(event, ics.ComponentPropertySummary)
		location := propertyValue(event, ics.ComponentPropertyLocation)
		if location == "" {
			location = "-"
		}

		entry := models.CalendarEvent{
			ID:        extractEventID(event.Id()),
			Name:      name,
			StartTime: utils.FormatParis(startParis, "15:04"),
			EndTime:   utils.FormatParis(endParis, "15:04"),
			Location:  location,
		}

		result[dateKey] = append(result[dateKey], entry)
	}

	for dateKey, events := range result {
		sort.Slice(events, func(i, j int) bool {
			return events[i].StartTime < events[j].StartTime
		})
		result[dateKey] = events
	}

	return result, nil
}

func isCancelled(event *ics.VEvent) bool {
	status := propertyValue(event, ics.ComponentPropertyStatus)
	return strings.EqualFold(status, string(ics.ObjectStatusCancelled))
}

func propertyValue(event *ics.VEvent, prop ics.ComponentProperty) string {
	p := event.GetProperty(prop)
	if p == nil {
		return ""
	}
	return p.Value
}

func extractEventID(uid string) string {
	if match := passIDRegex.FindString(uid); match != "" {
		return match
	}
	if at := strings.Index(uid, "@"); at > 0 {
		return uid[:at]
	}
	return uid
}
