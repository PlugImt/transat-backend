package services

import (
	"archive/zip"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
	"golang.org/x/sync/singleflight"
)

const (
	gtfsRefreshInterval       = 24 * time.Hour
	gtfsColdStartBackoffStart = 30 * time.Second
	gtfsColdStartBackoffCap   = 15 * time.Minute
	gtfsStaticMaxBody         = 80 << 20
	departureDatetimeLayout   = "2006-01-02 15:04:05"
	defaultMaxDepartures      = 3
	// Keep recently passed trips so delay-only / cancel updates can still match.
	scheduledLookback = 2 * time.Hour
)

type lineConfig struct {
	name           string
	routeShortName string
	directionID    string
	stopIDs        map[string]struct{}
}

type GTFSOptions struct {
	URL           string
	RealtimeURL   string
	Lines         []models.GTFSLineConfig
	MaxDepartures int
}

type calendarService struct {
	weekdays  [7]bool
	startDate time.Time
	endDate   time.Time
}

type calendarData struct {
	services   map[string]calendarService
	exceptions map[string]map[string]int // serviceID -> YYYYMMDD -> exception_type
}

type lineDeparture struct {
	tripID              string
	serviceID           string
	secondsFromMidnight int
}

type scheduledDeparture struct {
	at         time.Time
	tripID     string
	serviceDay time.Time // GTFS service day (not the calendar day of `at` for overnight trips)
}

type gtfsData struct {
	lineDepartures map[string][]lineDeparture
	calendar       calendarData
	routeToLine    map[string]string        // routes.route_id → line name
	tripToLine     map[string]string        // trips.trip_id → line name
	tripTimes      map[string]lineDeparture // trip_id → configured-stop departure
}

type GTFSService struct {
	gtfsURL       string
	realtimeURL   string
	lines         []lineConfig
	lineByName    map[string]lineConfig
	maxDepartures int
	mu            sync.RWMutex
	data          *gtfsData
	lastRefresh   time.Time
	refreshing    bool
	rtMu          sync.Mutex
	rtCache       *realtimeSnapshot
	rtGroup       singleflight.Group
}

func NewGTFSService(opts GTFSOptions) *GTFSService {
	lines, err := resolveLineConfigs(opts.Lines)
	if err != nil {
		log.Fatalf("💥 Invalid GTFS line config: %v", err)
	}

	maxDepartures := opts.MaxDepartures
	if maxDepartures <= 0 {
		maxDepartures = defaultMaxDepartures
	}

	lineByName := make(map[string]lineConfig, len(lines))
	for _, line := range lines {
		lineByName[line.name] = line
	}

	svc := &GTFSService{
		gtfsURL:       opts.URL,
		realtimeURL:   opts.RealtimeURL,
		lines:         lines,
		lineByName:    lineByName,
		maxDepartures: maxDepartures,
	}
	go svc.refreshLoop()
	return svc
}

func resolveLineConfigs(lines []models.GTFSLineConfig) ([]lineConfig, error) {
	if len(lines) == 0 {
		return defaultLineConfigs(), nil
	}
	parsed := lineConfigsFromModel(lines)
	if len(parsed) == 0 {
		return nil, fmt.Errorf("no valid lines in provided configuration")
	}
	return parsed, nil
}

func (s *GTFSService) refreshLoop() {
	failures := 0
	for {
		err := s.Refresh()
		ready := s.IsReady()
		if err != nil {
			failures++
		} else {
			failures = 0
		}
		time.Sleep(nextRefreshDelay(ready, failures))
	}
}

func nextRefreshDelay(ready bool, failures int) time.Duration {
	if ready {
		return gtfsRefreshInterval
	}
	if failures <= 0 {
		failures = 1
	}
	delay := gtfsColdStartBackoffStart
	for i := 1; i < failures; i++ {
		if delay >= gtfsColdStartBackoffCap {
			return gtfsColdStartBackoffCap
		}
		delay *= 2
	}
	if delay > gtfsColdStartBackoffCap {
		return gtfsColdStartBackoffCap
	}
	return delay
}

func (s *GTFSService) Refresh() error {
	s.mu.Lock()
	if s.refreshing {
		s.mu.Unlock()
		return nil
	}
	s.refreshing = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.refreshing = false
		s.mu.Unlock()
	}()

	data, err := s.downloadAndParse(context.Background())
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to refresh GTFS data")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return err
	}

	s.mu.Lock()
	s.data = data
	s.lastRefresh = utils.Now()
	s.mu.Unlock()

	utils.LogMessage(utils.LevelInfo, "GTFS data refreshed successfully")
	return nil
}

func (s *GTFSService) IsReady() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data != nil
}

func (s *GTFSService) GetChantrerieDepartures() (*models.BusDeparturesResponse, error) {
	s.mu.RLock()
	data := s.data
	lastRefresh := s.lastRefresh
	s.mu.RUnlock()

	if data == nil {
		return nil, fmt.Errorf("GTFS data not loaded yet")
	}

	now := utils.Now()
	liveByLine, lastRealtime := s.getRealtimeDepartures(now, data)

	response := &models.BusDeparturesResponse{
		LastRefresh: lastRefresh.Format(departureDatetimeLayout),
		Lines:       make([]models.BusLineDepartures, 0, len(s.lines)),
	}
	if !lastRealtime.IsZero() {
		response.LastRealtime = lastRealtime.Format(departureDatetimeLayout)
	}

	for _, line := range s.lines {
		scheduled := findNextDepartures(data.lineDepartures[line.name], data.calendar, now, s.maxDepartures*3)
		departures := mergeDepartures(scheduled, liveByLine[line.name], s.maxDepartures, now)
		response.Lines = append(response.Lines, models.BusLineDepartures{
			Name:       line.name,
			Departures: departures,
		})
	}

	return response, nil
}

func (s *GTFSService) downloadAndParse(ctx context.Context) (*gtfsData, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.gtfsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create GTFS request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("download GTFS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download GTFS: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, gtfsStaticMaxBody+1))
	if err != nil {
		return nil, fmt.Errorf("read GTFS body: %w", err)
	}
	if len(body) > gtfsStaticMaxBody {
		return nil, fmt.Errorf("GTFS payload exceeds %d bytes", gtfsStaticMaxBody)
	}

	zipReader, err := zip.NewReader(bytesReaderAt(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("open GTFS zip: %w", err)
	}

	return parseGTFSZip(zipReader, s.lines)
}

func parseGTFSZip(zipReader *zip.Reader, lines []lineConfig) (*gtfsData, error) {
	calendar, err := parseCalendar(zipFileByName(zipReader, "calendar.txt"), zipFileByName(zipReader, "calendar_dates.txt"))
	if err != nil {
		return nil, err
	}

	routeToLine, err := parseRoutes(zipFileByName(zipReader, "routes.txt"), lines)
	if err != nil {
		return nil, err
	}

	targetTrips, err := parseTrips(zipFileByName(zipReader, "trips.txt"), routeToLine, lines)
	if err != nil {
		return nil, err
	}

	lineDepartures, err := parseStopTimes(zipFileByName(zipReader, "stop_times.txt"), targetTrips, lines)
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		if len(lineDepartures[line.name]) == 0 {
			return nil, fmt.Errorf("configured line %q matched no stop times", line.name)
		}
	}

	tripToLine := make(map[string]string, len(targetTrips))
	tripTimes := make(map[string]lineDeparture, len(targetTrips))
	for tripID, trip := range targetTrips {
		tripToLine[tripID] = trip.lineName
	}
	for _, deps := range lineDepartures {
		for _, dep := range deps {
			tripTimes[dep.tripID] = dep
		}
	}

	return &gtfsData{
		lineDepartures: lineDepartures,
		calendar:       calendar,
		routeToLine:    routeToLine,
		tripToLine:     tripToLine,
		tripTimes:      tripTimes,
	}, nil
}

func parseRoutes(file *zip.File, lines []lineConfig) (map[string]string, error) {
	rc, csvReader, cols, err := openGTFSCSV(file, "routes.txt", "route_id", "route_short_name")
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	shortNameToLine := make(map[string]string, len(lines))
	for _, line := range lines {
		shortNameToLine[line.routeShortName] = line.name
	}

	routeToLine := make(map[string]string)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse routes.txt: %w", err)
		}

		shortName, ok := cols.get(record, "route_short_name")
		if !ok {
			continue
		}
		lineName, ok := shortNameToLine[shortName]
		if !ok {
			continue
		}
		routeID, ok := cols.get(record, "route_id")
		if !ok || routeID == "" {
			continue
		}
		routeToLine[routeID] = lineName
	}

	return routeToLine, nil
}

type tripInfo struct {
	lineName  string
	serviceID string
}

func parseTrips(file *zip.File, routeToLine map[string]string, lines []lineConfig) (map[string]tripInfo, error) {
	rc, csvReader, cols, err := openGTFSCSV(file, "trips.txt", "route_id", "service_id", "trip_id")
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	needsDirection := false
	lineByName := make(map[string]lineConfig, len(lines))
	for _, line := range lines {
		lineByName[line.name] = line
		if line.directionID != "" {
			needsDirection = true
		}
	}
	if needsDirection {
		if _, ok := cols["direction_id"]; !ok {
			return nil, fmt.Errorf("trips.txt: missing direction_id column required by line config")
		}
	}

	targetTrips := make(map[string]tripInfo)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse trips.txt: %w", err)
		}

		routeID, ok := cols.get(record, "route_id")
		if !ok {
			continue
		}
		lineName, ok := routeToLine[routeID]
		if !ok {
			continue
		}

		line := lineByName[lineName]
		if line.directionID != "" {
			directionID, ok := cols.get(record, "direction_id")
			if !ok || directionID != line.directionID {
				continue
			}
		}

		tripID, ok := cols.get(record, "trip_id")
		if !ok || tripID == "" {
			continue
		}
		serviceID, ok := cols.get(record, "service_id")
		if !ok || serviceID == "" {
			continue
		}

		targetTrips[tripID] = tripInfo{
			lineName:  lineName,
			serviceID: serviceID,
		}
	}

	return targetTrips, nil
}

func parseStopTimes(file *zip.File, targetTrips map[string]tripInfo, lines []lineConfig) (map[string][]lineDeparture, error) {
	rc, csvReader, cols, err := openGTFSCSV(file, "stop_times.txt", "trip_id", "departure_time", "stop_id")
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	stopIDsByLine := make(map[string]map[string]struct{}, len(lines))
	for _, line := range lines {
		stopIDsByLine[line.name] = line.stopIDs
	}

	lineDepartures := make(map[string]map[string]lineDeparture, len(lines))
	for _, line := range lines {
		lineDepartures[line.name] = make(map[string]lineDeparture)
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse stop_times.txt: %w", err)
		}

		tripID, ok := cols.get(record, "trip_id")
		if !ok {
			continue
		}
		trip, ok := targetTrips[tripID]
		if !ok {
			continue
		}

		stopID, ok := cols.get(record, "stop_id")
		if !ok {
			continue
		}
		if _, ok := stopIDsByLine[trip.lineName][stopID]; !ok {
			continue
		}

		departureTime, ok := cols.get(record, "departure_time")
		if !ok {
			continue
		}
		seconds, err := parseGTFSTime(departureTime)
		if err != nil {
			continue
		}

		lineDepartures[trip.lineName][tripID] = lineDeparture{
			tripID:              tripID,
			serviceID:           trip.serviceID,
			secondsFromMidnight: seconds,
		}
	}

	result := make(map[string][]lineDeparture, len(lines))
	for lineName, departures := range lineDepartures {
		list := make([]lineDeparture, 0, len(departures))
		for _, dep := range departures {
			list = append(list, dep)
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].secondsFromMidnight < list[j].secondsFromMidnight
		})
		result[lineName] = list
	}

	return result, nil
}

func parseCalendar(calendarFile, exceptionsFile *zip.File) (calendarData, error) {
	data := calendarData{
		services:   make(map[string]calendarService),
		exceptions: make(map[string]map[string]int),
	}

	if calendarFile == nil && exceptionsFile == nil {
		return data, fmt.Errorf("calendar.txt and calendar_dates.txt both missing from GTFS archive")
	}

	if calendarFile != nil {
		rc, csvReader, cols, err := openGTFSCSV(calendarFile, "calendar.txt",
			"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date")
		if err != nil {
			return data, err
		}
		defer rc.Close()

		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return data, fmt.Errorf("parse calendar.txt: %w", err)
			}

			serviceID, ok := cols.get(record, "service_id")
			if !ok || serviceID == "" {
				continue
			}
			startDateRaw, ok := cols.get(record, "start_date")
			if !ok {
				continue
			}
			endDateRaw, ok := cols.get(record, "end_date")
			if !ok {
				continue
			}
			startDate, err := parseGTFSDate(startDateRaw)
			if err != nil {
				continue
			}
			endDate, err := parseGTFSDate(endDateRaw)
			if err != nil {
				continue
			}

			data.services[serviceID] = calendarService{
				weekdays: [7]bool{
					cols.equals(record, "sunday", "1"),
					cols.equals(record, "monday", "1"),
					cols.equals(record, "tuesday", "1"),
					cols.equals(record, "wednesday", "1"),
					cols.equals(record, "thursday", "1"),
					cols.equals(record, "friday", "1"),
					cols.equals(record, "saturday", "1"),
				},
				startDate: startDate,
				endDate:   endDate,
			}
		}
	}

	if exceptionsFile != nil {
		excRC, excCSV, excCols, err := openGTFSCSV(exceptionsFile, "calendar_dates.txt", "service_id", "date", "exception_type")
		if err != nil {
			return data, err
		}
		defer excRC.Close()

		for {
			record, err := excCSV.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return data, fmt.Errorf("parse calendar_dates.txt: %w", err)
			}

			serviceID, ok := excCols.get(record, "service_id")
			if !ok || serviceID == "" {
				continue
			}
			dateKey, ok := excCols.get(record, "date")
			if !ok || dateKey == "" {
				continue
			}
			exceptionRaw, ok := excCols.get(record, "exception_type")
			if !ok {
				continue
			}
			exceptionType, err := strconv.Atoi(exceptionRaw)
			if err != nil {
				continue
			}
			if data.exceptions[serviceID] == nil {
				data.exceptions[serviceID] = make(map[string]int)
			}
			data.exceptions[serviceID][dateKey] = exceptionType
		}
	}

	return data, nil
}

type csvCols map[string]int

func zipFileByName(zipReader *zip.Reader, name string) *zip.File {
	for _, f := range zipReader.File {
		if path.Base(f.Name) == name {
			return f
		}
	}
	return nil
}

func openGTFSCSV(file *zip.File, name string, required ...string) (io.ReadCloser, *csv.Reader, csvCols, error) {
	if file == nil {
		return nil, nil, nil, fmt.Errorf("%s not found in GTFS archive", name)
	}

	reader, err := file.Open()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("open %s: %w", name, err)
	}

	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1
	csvReader.LazyQuotes = true

	cols, err := readCSVHeader(csvReader)
	if err != nil {
		reader.Close()
		return nil, nil, nil, fmt.Errorf("parse %s header: %w", name, err)
	}
	if err := cols.require(required...); err != nil {
		reader.Close()
		return nil, nil, nil, fmt.Errorf("%s: %w", name, err)
	}

	return reader, csvReader, cols, nil
}

func readCSVHeader(r *csv.Reader) (csvCols, error) {
	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	cols := make(csvCols, len(header))
	for i, name := range header {
		name = strings.TrimSpace(name)
		if i == 0 {
			name = strings.TrimPrefix(name, "\ufeff")
		}
		cols[name] = i
	}
	return cols, nil
}

func (c csvCols) get(record []string, name string) (string, bool) {
	i, ok := c[name]
	if !ok || i < 0 || i >= len(record) {
		return "", false
	}
	return strings.TrimSpace(record[i]), true
}

func (c csvCols) equals(record []string, name, want string) bool {
	value, ok := c.get(record, name)
	return ok && value == want
}

func (c csvCols) require(names ...string) error {
	var missing []string
	for _, name := range names {
		if _, ok := c[name]; !ok {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing columns: %s", strings.Join(missing, ", "))
	}
	return nil
}

func parseGTFSDate(value string) (time.Time, error) {
	return time.ParseInLocation("20060102", value, utils.ParisLocation)
}

func parseGTFSTime(value string) (int, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid GTFS time: %s", value)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}

	return hours*3600 + minutes*60 + seconds, nil
}

func (c calendarData) isServiceActive(serviceID string, day time.Time) bool {
	day = utils.StartOfDayParis(day)
	dateKey := day.Format("20060102")

	if serviceExceptions, ok := c.exceptions[serviceID]; ok {
		if exceptionType, ok := serviceExceptions[dateKey]; ok {
			return exceptionType == 1
		}
	}

	service, ok := c.services[serviceID]
	if !ok {
		return false
	}

	if day.Before(service.startDate) || day.After(service.endDate) {
		return false
	}

	weekday := int(day.Weekday())
	return service.weekdays[weekday]
}

// gtfsServiceTime converts a GTFS Time (seconds since "noon minus 12h" of the
// service day) into a wall-clock instant in Europe/Paris. Using noon as the
// anchor keeps times correct across DST transitions.
func gtfsServiceTime(serviceDay time.Time, secondsFromMidnight int) time.Time {
	day := serviceDay.In(utils.ParisLocation)
	noon := time.Date(day.Year(), day.Month(), day.Day(), 12, 0, 0, 0, utils.ParisLocation)
	return noon.Add(time.Duration(secondsFromMidnight-12*3600) * time.Second)
}

func findNextDepartures(departures []lineDeparture, calendar calendarData, now time.Time, count int) []scheduledDeparture {
	if len(departures) == 0 || count <= 0 {
		return nil
	}

	now = now.In(utils.ParisLocation)
	recent := make([]scheduledDeparture, 0, count)
	upcoming := make([]scheduledDeparture, 0, count*2)

	for dayOffset := -1; dayOffset < 8; dayOffset++ {
		day := utils.StartOfDayParis(now).AddDate(0, 0, dayOffset)

		for _, dep := range departures {
			if !calendar.isServiceActive(dep.serviceID, day) {
				continue
			}

			departureTime := gtfsServiceTime(day, dep.secondsFromMidnight)
			item := scheduledDeparture{
				at:         departureTime,
				tripID:     dep.tripID,
				serviceDay: day,
			}
			if departureTime.After(now) {
				upcoming = append(upcoming, item)
				continue
			}
			if now.Sub(departureTime) <= scheduledLookback {
				recent = append(recent, item)
			}
		}
	}

	sort.Slice(recent, func(i, j int) bool {
		return recent[i].at.Before(recent[j].at)
	})
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].at.Before(upcoming[j].at)
	})
	if len(upcoming) > count {
		upcoming = upcoming[:count]
	}
	return append(recent, upcoming...)
}

func defaultLineConfigs() []lineConfig {
	return lineConfigsFromModel([]models.GTFSLineConfig{
		{Name: "C6", RouteShortName: "C6", DirectionID: "1", StopIDs: []string{"FR_NAOLIB:Quay:634"}},
		{Name: "75", RouteShortName: "75", DirectionID: "1", StopIDs: []string{"FR_NAOLIB:Quay:632"}},
		{Name: "E5", RouteShortName: "E5", DirectionID: "1", StopIDs: []string{"FR_NAOLIB:Quay:2521"}},
	})
}

func lineConfigsFromModel(lines []models.GTFSLineConfig) []lineConfig {
	out := make([]lineConfig, 0, len(lines))
	for _, line := range lines {
		if line.Name == "" || line.RouteShortName == "" {
			continue
		}
		stopIDs := make(map[string]struct{}, len(line.StopIDs))
		for _, stopID := range line.StopIDs {
			if stopID != "" {
				stopIDs[stopID] = struct{}{}
			}
		}
		if len(stopIDs) == 0 {
			continue
		}
		out = append(out, lineConfig{
			name:           line.Name,
			routeShortName: line.RouteShortName,
			directionID:    line.DirectionID,
			stopIDs:        stopIDs,
		})
	}
	return out
}

// bytesReaderAt adapts a byte slice for zip.NewReader.
type bytesReaderAt []byte

func (b bytesReaderAt) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(b)) {
		return 0, io.EOF
	}
	n := copy(p, b[off:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}
