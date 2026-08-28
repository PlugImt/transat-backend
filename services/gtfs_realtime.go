package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	gtfsrt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
	"google.golang.org/protobuf/proto"
)

const (
	gtfsRealtimeTimeout  = 5 * time.Second
	gtfsRealtimeTTL      = 20 * time.Second
	gtfsRealtimeMaxStale = 5 * time.Minute
	gtfsRealtimeMaxBody  = 8 << 20
	realtimeMatchWindow  = 20 * time.Minute
)

type liveStatus int

const (
	liveStatusPredicted liveStatus = iota
	liveStatusCanceled
	liveStatusSkipped
)

type liveDeparture struct {
	at         time.Time
	tripID     string
	status     liveStatus
	delay      *int
	serviceDay time.Time // TripDescriptor.start_date; zero if absent
}

type realtimeSnapshot struct {
	byLine    map[string][]liveDeparture
	fetchedAt time.Time // wall-clock fetch completion; drives TTL / max-stale
	feedAt    time.Time // FeedHeader.timestamp; exposed as lastRealtime
}

func (s *GTFSService) getRealtimeDepartures(now time.Time, data *gtfsData) (map[string][]liveDeparture, time.Time) {
	if s.realtimeURL == "" {
		return nil, time.Time{}
	}

	s.rtMu.Lock()
	if s.rtCache != nil && now.Sub(s.rtCache.fetchedAt) < gtfsRealtimeTTL {
		snapshot := s.rtCache
		s.rtMu.Unlock()
		return filterFutureLiveDepartures(snapshot.byLine, now), snapshot.feedAt
	}
	s.rtMu.Unlock()

	v, err, _ := s.rtGroup.Do("rt", func() (any, error) {
		// Re-check cache inside the singleflight critical section so waiters
		// share a successful refresh instead of launching another fetch.
		s.rtMu.Lock()
		if s.rtCache != nil && utils.Now().Sub(s.rtCache.fetchedAt) < gtfsRealtimeTTL {
			snapshot := s.rtCache
			s.rtMu.Unlock()
			return snapshot, nil
		}
		s.rtMu.Unlock()

		byLine, feedAt, err := s.fetchRealtimeDepartures(utils.Now(), data)
		if err != nil {
			return nil, err
		}
		snapshot := &realtimeSnapshot{
			byLine:    byLine,
			fetchedAt: utils.Now(),
			feedAt:    feedAt,
		}
		s.rtMu.Lock()
		s.rtCache = snapshot
		s.rtMu.Unlock()
		return snapshot, nil
	})
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Failed to refresh GTFS realtime data, using schedule only")
		utils.LogLineKeyValue(utils.LevelWarn, "Error", err)

		s.rtMu.Lock()
		defer s.rtMu.Unlock()
		if s.rtCache != nil && now.Sub(s.rtCache.fetchedAt) < gtfsRealtimeMaxStale {
			return filterFutureLiveDepartures(s.rtCache.byLine, now), s.rtCache.feedAt
		}
		return nil, time.Time{}
	}

	snapshot := v.(*realtimeSnapshot)
	return filterFutureLiveDepartures(snapshot.byLine, now), snapshot.feedAt
}

func filterFutureLiveDepartures(byLine map[string][]liveDeparture, now time.Time) map[string][]liveDeparture {
	if byLine == nil {
		return nil
	}
	out := make(map[string][]liveDeparture, len(byLine))
	for name, deps := range byLine {
		filtered := make([]liveDeparture, 0, len(deps))
		for _, dep := range deps {
			// Delay-only predictions (zero at) are kept until merge applies the delay.
			if dep.status == liveStatusPredicted && !dep.at.IsZero() && !dep.at.After(now) {
				continue
			}
			if (dep.status == liveStatusCanceled || dep.status == liveStatusSkipped) &&
				isExpiredServiceDay(dep.serviceDay, now) {
				continue
			}
			filtered = append(filtered, dep)
		}
		out[name] = filtered
	}
	return out
}

func (s *GTFSService) fetchRealtimeDepartures(now time.Time, data *gtfsData) (map[string][]liveDeparture, time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gtfsRealtimeTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.realtimeURL, nil)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("create GTFS-RT request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("download GTFS-RT: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, time.Time{}, fmt.Errorf("download GTFS-RT: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, gtfsRealtimeMaxBody+1))
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("read GTFS-RT: %w", err)
	}
	if len(body) > gtfsRealtimeMaxBody {
		return nil, time.Time{}, fmt.Errorf("GTFS-RT payload exceeds %d bytes", gtfsRealtimeMaxBody)
	}

	var feed gtfsrt.FeedMessage
	if err := proto.Unmarshal(body, &feed); err != nil {
		return nil, time.Time{}, fmt.Errorf("parse GTFS-RT: %w", err)
	}

	feedAt, err := feedHeaderTime(feed.GetHeader(), now)
	if err != nil {
		return nil, time.Time{}, err
	}

	byLine := make(map[string][]liveDeparture, len(s.lines))
	for _, line := range s.lines {
		byLine[line.name] = nil
	}

	for _, entity := range feed.GetEntity() {
		tripUpdate := entity.GetTripUpdate()
		if tripUpdate == nil {
			continue
		}
		trip := tripUpdate.GetTrip()
		lineName := s.resolveLineName(trip, data)
		if lineName == "" {
			continue
		}
		line := s.lineByName[lineName]
		tripID := strings.TrimSpace(trip.GetTripId())
		serviceDay := tripServiceDay(trip)

		if trip.GetScheduleRelationship() == gtfsrt.TripDescriptor_CANCELED {
			byLine[lineName] = append(byLine[lineName], liveDeparture{
				tripID:     tripID,
				status:     liveStatusCanceled,
				serviceDay: serviceDay,
			})
			continue
		}

		if dep, ok := predictConfiguredStop(tripUpdate, line.stopIDs, tripID, now); ok {
			dep.serviceDay = serviceDay
			byLine[lineName] = append(byLine[lineName], dep)
		}
	}

	for name, deps := range byLine {
		sort.Slice(deps, func(i, j int) bool {
			if deps[i].status != liveStatusPredicted && deps[j].status == liveStatusPredicted {
				return true
			}
			if deps[i].status == liveStatusPredicted && deps[j].status != liveStatusPredicted {
				return false
			}
			if deps[i].at.IsZero() && !deps[j].at.IsZero() {
				return false
			}
			if !deps[i].at.IsZero() && deps[j].at.IsZero() {
				return true
			}
			return deps[i].at.Before(deps[j].at)
		})
		byLine[name] = deps
	}

	return byLine, feedAt, nil
}

// predictConfiguredStop walks StopTimeUpdates in stop_sequence order and
// builds a liveDeparture for the configured Chantrerie stop, including
// delay-only events and delays inherited from upstream stops.
func predictConfiguredStop(
	tripUpdate *gtfsrt.TripUpdate,
	stopIDs map[string]struct{},
	tripID string,
	now time.Time,
) (liveDeparture, bool) {
	updates := append([]*gtfsrt.TripUpdate_StopTimeUpdate(nil), tripUpdate.GetStopTimeUpdate()...)
	sort.SliceStable(updates, func(i, j int) bool {
		return updates[i].GetStopSequence() < updates[j].GetStopSequence()
	})

	var tripDelay *int
	if tripUpdate.Delay != nil {
		d := int(tripUpdate.GetDelay())
		tripDelay = &d
	}

	var inheritedDelay *int
	sawConfigured := false
	var result liveDeparture
	hasResult := false

	for _, stopUpdate := range updates {
		rel := stopUpdate.GetScheduleRelationship()
		isConfigured := false
		if stopID := stopUpdate.GetStopId(); stopID != "" {
			_, isConfigured = stopIDs[stopID]
		}

		if rel == gtfsrt.TripUpdate_StopTimeUpdate_NO_DATA {
			inheritedDelay = nil
			if isConfigured {
				sawConfigured = true
			}
			continue
		}

		if isConfigured && rel == gtfsrt.TripUpdate_StopTimeUpdate_SKIPPED {
			return liveDeparture{tripID: tripID, status: liveStatusSkipped}, true
		}

		event := stopUpdate.GetDeparture()
		if event == nil {
			event = stopUpdate.GetArrival()
		}

		var eventDelay *int
		if event != nil && event.Delay != nil {
			d := int(event.GetDelay())
			eventDelay = &d
		}

		if rel == gtfsrt.TripUpdate_StopTimeUpdate_SKIPPED {
			// Spec: previous delay propagates over skipped stops.
			continue
		}

		if event != nil && event.GetTime() != 0 {
			at := time.Unix(event.GetTime(), 0).In(utils.ParisLocation)
			if isConfigured {
				if !at.After(now) {
					return liveDeparture{}, false
				}
				dep := liveDeparture{
					at:     at,
					tripID: tripID,
					status: liveStatusPredicted,
					delay:  eventDelay,
				}
				return dep, true
			}
			if eventDelay != nil {
				inheritedDelay = eventDelay
			} else {
				inheritedDelay = nil
			}
			continue
		}

		if eventDelay != nil {
			inheritedDelay = eventDelay
		}

		if isConfigured {
			sawConfigured = true
			delay := eventDelay
			if delay == nil {
				delay = inheritedDelay
			}
			if delay == nil {
				delay = tripDelay
			}
			if delay == nil {
				continue
			}
			result = liveDeparture{
				tripID: tripID,
				status: liveStatusPredicted,
				delay:  delay,
			}
			hasResult = true
			break
		}
	}

	if hasResult {
		return result, true
	}

	// Configured stop had no StopTimeUpdate: inherit last delay or trip delay.
	if !sawConfigured {
		delay := inheritedDelay
		if delay == nil {
			delay = tripDelay
		}
		if delay != nil {
			return liveDeparture{
				tripID: tripID,
				status: liveStatusPredicted,
				delay:  delay,
			}, true
		}
	}

	return liveDeparture{}, false
}

func tripServiceDay(trip *gtfsrt.TripDescriptor) time.Time {
	if trip == nil {
		return time.Time{}
	}
	raw := strings.TrimSpace(trip.GetStartDate())
	if raw == "" {
		return time.Time{}
	}
	day, err := parseGTFSDate(raw)
	if err != nil {
		return time.Time{}
	}
	return utils.StartOfDayParis(day)
}

func isExpiredServiceDay(serviceDay, now time.Time) bool {
	if serviceDay.IsZero() {
		return false
	}
	earliest := utils.StartOfDayParis(now).AddDate(0, 0, -1)
	return utils.StartOfDayParis(serviceDay).Before(earliest)
}

func sameServiceDay(a, b time.Time) bool {
	if a.IsZero() || b.IsZero() {
		return false
	}
	return utils.StartOfDayParis(a).Equal(utils.StartOfDayParis(b))
}

func feedHeaderTime(header *gtfsrt.FeedHeader, now time.Time) (time.Time, error) {
	if header == nil || header.GetTimestamp() == 0 {
		return time.Time{}, nil
	}
	feedAt := time.Unix(int64(header.GetTimestamp()), 0).In(utils.ParisLocation)
	if now.Sub(feedAt) > gtfsRealtimeMaxStale {
		return time.Time{}, fmt.Errorf("GTFS-RT feed timestamp is stale (%s)", feedAt.Format(time.RFC3339))
	}
	return feedAt, nil
}

func (s *GTFSService) resolveLineName(trip *gtfsrt.TripDescriptor, data *gtfsData) string {
	if trip == nil {
		return ""
	}
	tripID := strings.TrimSpace(trip.GetTripId())
	if data != nil && tripID != "" {
		if lineName, ok := data.tripToLine[tripID]; ok {
			return lineName
		}
	}
	routeID := strings.TrimSpace(trip.GetRouteId())
	if data != nil && routeID != "" {
		if lineName, ok := data.routeToLine[routeID]; ok {
			return lineName
		}
	}
	return s.lineNameForRoute(routeID)
}

func (s *GTFSService) lineNameForRoute(routeID string) string {
	for _, line := range s.lines {
		if routeMatchesLine(routeID, line.routeShortName) {
			return line.name
		}
	}
	return ""
}

func routeMatchesLine(routeID, shortName string) bool {
	if routeID == "" || shortName == "" {
		return false
	}
	return routeID == shortName ||
		strings.HasSuffix(routeID, ":"+shortName) ||
		strings.HasSuffix(routeID, ":Line:"+shortName)
}

func mergeDepartures(scheduled []scheduledDeparture, live []liveDeparture, count int, now time.Time) []models.BusDeparture {
	if count <= 0 {
		return []models.BusDeparture{}
	}

	usedScheduled := make([]bool, len(scheduled))
	merged := make([]models.BusDeparture, 0, len(scheduled)+len(live))

	// First pass: suppress canceled/skipped scheduled departures.
	for _, item := range live {
		if item.status != liveStatusCanceled && item.status != liveStatusSkipped {
			continue
		}
		if idx := matchScheduledIndex(scheduled, usedScheduled, item, now); idx >= 0 {
			usedScheduled[idx] = true
		}
	}

	// Second pass: predicted realtime departures.
	for _, item := range live {
		if item.status != liveStatusPredicted {
			continue
		}

		at := item.at
		idx := matchScheduledIndex(scheduled, usedScheduled, item, now)
		if idx >= 0 {
			usedScheduled[idx] = true
			// Delay-only: derive wall time from the matched schedule.
			if at.IsZero() {
				if item.delay == nil {
					continue
				}
				at = scheduled[idx].at.Add(time.Duration(*item.delay) * time.Second)
			}
		} else if at.IsZero() {
			// No schedule to apply the delay against.
			continue
		}

		if !at.After(now) {
			continue
		}

		merged = append(merged, models.BusDeparture{
			Time:         at.Format(departureDatetimeLayout),
			Realtime:     true,
			DelaySeconds: item.delay,
		})
	}

	for i, scheduledAt := range scheduled {
		if usedScheduled[i] {
			continue
		}
		if !scheduledAt.at.After(now) {
			continue
		}
		merged = append(merged, models.BusDeparture{
			Time:     scheduledAt.at.Format(departureDatetimeLayout),
			Realtime: false,
		})
	}

	sort.Slice(merged, func(i, j int) bool { return merged[i].Time < merged[j].Time })
	if len(merged) > count {
		merged = merged[:count]
	}
	if merged == nil {
		return []models.BusDeparture{}
	}
	return merged
}

func normalizeTripID(id string) string {
	return strings.TrimSpace(id)
}

func matchScheduledIndex(scheduled []scheduledDeparture, used []bool, item liveDeparture, now time.Time) int {
	itemTripID := normalizeTripID(item.tripID)
	if itemTripID != "" {
		if idx := matchTripIDIndex(scheduled, used, item, now); idx >= 0 {
			return idx
		}
	}

	if item.status != liveStatusPredicted {
		return -1
	}

	// Delay-only predictions without a trip_id match cannot use the time window.
	if item.at.IsZero() {
		return -1
	}

	closest := -1
	closestDelta := realtimeMatchWindow + 1
	for i, scheduledAt := range scheduled {
		if used[i] {
			continue
		}
		if !item.serviceDay.IsZero() && !scheduledAt.serviceDay.IsZero() &&
			!sameServiceDay(item.serviceDay, scheduledAt.serviceDay) {
			continue
		}
		// Prefer exact trip_id above; when IDs differ or one side is empty,
		// still allow a time-window match so delayed live rows de-dupe.
		delta := item.at.Sub(scheduledAt.at)
		if delta < 0 {
			delta = -delta
		}
		if delta <= realtimeMatchWindow && delta < closestDelta {
			closest = i
			closestDelta = delta
		}
	}
	return closest
}

func matchTripIDIndex(scheduled []scheduledDeparture, used []bool, item liveDeparture, now time.Time) int {
	itemTripID := normalizeTripID(item.tripID)
	best := -1
	var bestDelta time.Duration

	for i, scheduledAt := range scheduled {
		if used[i] {
			continue
		}
		if normalizeTripID(scheduledAt.tripID) != itemTripID {
			continue
		}
		if !item.serviceDay.IsZero() && !scheduledAt.serviceDay.IsZero() &&
			!sameServiceDay(item.serviceDay, scheduledAt.serviceDay) {
			continue
		}

		ref := now
		if !item.at.IsZero() {
			ref = item.at
		}
		delta := scheduledAt.at.Sub(ref)
		if delta < 0 {
			delta = -delta
		}
		if best < 0 || delta < bestDelta {
			best = i
			bestDelta = delta
		}
	}
	return best
}
