package services

import (
	"archive/zip"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

const (
	defaultGTFSURL          = "https://data.nantesmetropole.fr/explore/dataset/244400404_tan-arrets-horaires-circuits/files/16a1a0af5946619af621baa4ad9ee662/download/"
	gtfsRefreshInterval     = 24 * time.Hour
	departureDatetimeLayout = "2006-01-02 15:04:05"
)

type lineConfig struct {
	name           string
	routeShortName string
	directionID    string
	stopIDs        map[string]struct{}
}

var chantrerieLeavingLines = []lineConfig{
	{name: "C6", routeShortName: "C6", directionID: "1", stopIDs: map[string]struct{}{"CTRE2": {}}},
	{name: "75", routeShortName: "75", directionID: "1", stopIDs: map[string]struct{}{"CTRE4": {}}},
	{name: "E5", routeShortName: "E5", directionID: "1", stopIDs: map[string]struct{}{"PTEC2": {}}},
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
	serviceID           string
	secondsFromMidnight int
}

type gtfsData struct {
	lineDepartures map[string][]lineDeparture
	calendar       calendarData
}

type GTFSService struct {
	gtfsURL     string
	mu          sync.RWMutex
	data        *gtfsData
	lastRefresh time.Time
	refreshing  bool
}

func NewGTFSService() *GTFSService {
	svc := &GTFSService{
		gtfsURL: defaultGTFSURL,
	}
	go svc.refreshLoop()
	return svc
}

func (s *GTFSService) refreshLoop() {
	s.Refresh()

	ticker := time.NewTicker(gtfsRefreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.Refresh()
	}
}

func (s *GTFSService) Refresh() {
	s.mu.Lock()
	if s.refreshing {
		s.mu.Unlock()
		return
	}
	s.refreshing = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.refreshing = false
		s.mu.Unlock()
	}()

	data, err := s.downloadAndParse()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to refresh GTFS data")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return
	}

	s.mu.Lock()
	s.data = data
	s.lastRefresh = utils.Now()
	s.mu.Unlock()

	utils.LogMessage(utils.LevelInfo, "GTFS data refreshed successfully")
}

func (s *GTFSService) IsReady() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data != nil
}

func (s *GTFSService) GetChantrerieDepartures() ([]models.BusDeparture, error) {
	s.mu.RLock()
	data := s.data
	s.mu.RUnlock()

	if data == nil {
		return nil, fmt.Errorf("GTFS data not loaded yet")
	}

	now := utils.Now()
	result := make([]models.BusDeparture, 0, len(chantrerieLeavingLines))

	for _, line := range chantrerieLeavingLines {
		departures := data.lineDepartures[line.name]
		next := findNextDepartures(departures, data.calendar, now, 2)
		item := models.BusDeparture{Name: line.name}
		if len(next) > 0 {
			item.NextDeparture = next[0].Format(departureDatetimeLayout)
		}
		if len(next) > 1 {
			item.NextDeparture2 = next[1].Format(departureDatetimeLayout)
		}
		result = append(result, item)
	}

	return result, nil
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read GTFS body: %w", err)
	}

	zipReader, err := zip.NewReader(bytesReaderAt(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("open GTFS zip: %w", err)
	}

	return parseGTFSZip(zipReader)
}

func parseGTFSZip(zipReader *zip.Reader) (*gtfsData, error) {
	files := map[string]*zip.File{}
	for _, f := range zipReader.File {
		files[f.Name] = f
	}

	calendar, err := parseCalendar(files["calendar.txt"], files["calendar_dates.txt"])
	if err != nil {
		return nil, err
	}

	routeToLine, err := parseRoutes(files["routes.txt"])
	if err != nil {
		return nil, err
	}

	targetTrips, err := parseTrips(files["trips.txt"], routeToLine)
	if err != nil {
		return nil, err
	}

	lineDepartures, err := parseStopTimes(files["stop_times.txt"], targetTrips)
	if err != nil {
		return nil, err
	}

	return &gtfsData{
		lineDepartures: lineDepartures,
		calendar:       calendar,
	}, nil
}

func parseRoutes(file *zip.File) (map[string]string, error) {
	if file == nil {
		return nil, fmt.Errorf("routes.txt not found in GTFS archive")
	}

	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	routeToLine := make(map[string]string)
	csvReader := csv.NewReader(reader)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse routes.txt: %w", err)
	}

	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 3 {
			continue
		}
		for _, line := range chantrerieLeavingLines {
			if record[2] == line.routeShortName {
				routeToLine[record[0]] = line.name
			}
		}
	}

	return routeToLine, nil
}

type tripInfo struct {
	lineName  string
	serviceID string
}

func parseTrips(file *zip.File, routeToLine map[string]string) (map[string]tripInfo, error) {
	if file == nil {
		return nil, fmt.Errorf("trips.txt not found in GTFS archive")
	}

	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	lineByDirection := make(map[string]lineConfig)
	for _, line := range chantrerieLeavingLines {
		lineByDirection[line.name] = line
	}

	targetTrips := make(map[string]tripInfo)
	csvReader := csv.NewReader(reader)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse trips.txt: %w", err)
		}
		if len(record) < 6 {
			continue
		}

		lineName, ok := routeToLine[record[0]]
		if !ok {
			continue
		}

		line := lineByDirection[lineName]
		if record[5] != line.directionID {
			continue
		}

		targetTrips[record[2]] = tripInfo{
			lineName:  lineName,
			serviceID: record[1],
		}
	}

	return targetTrips, nil
}

func parseStopTimes(file *zip.File, targetTrips map[string]tripInfo) (map[string][]lineDeparture, error) {
	if file == nil {
		return nil, fmt.Errorf("stop_times.txt not found in GTFS archive")
	}

	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	stopIDsByLine := make(map[string]map[string]struct{})
	for _, line := range chantrerieLeavingLines {
		stopIDsByLine[line.name] = line.stopIDs
	}

	lineDepartures := make(map[string]map[string]lineDeparture)
	for _, line := range chantrerieLeavingLines {
		lineDepartures[line.name] = make(map[string]lineDeparture)
	}

	csvReader := csv.NewReader(reader)
	isHeader := true

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse stop_times.txt: %w", err)
		}
		if isHeader {
			isHeader = false
			continue
		}
		if len(record) < 4 {
			continue
		}

		trip, ok := targetTrips[record[0]]
		if !ok {
			continue
		}

		stopIDs := stopIDsByLine[trip.lineName]
		if _, ok := stopIDs[record[3]]; !ok {
			continue
		}

		seconds, err := parseGTFSTime(record[2])
		if err != nil {
			continue
		}

		key := fmt.Sprintf("%s|%d", trip.serviceID, seconds)
		lineDepartures[trip.lineName][key] = lineDeparture{
			serviceID:           trip.serviceID,
			secondsFromMidnight: seconds,
		}
	}

	result := make(map[string][]lineDeparture, len(chantrerieLeavingLines))
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

	if calendarFile == nil {
		return data, fmt.Errorf("calendar.txt not found in GTFS archive")
	}

	reader, err := calendarFile.Open()
	if err != nil {
		return data, err
	}
	defer reader.Close()

	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return data, fmt.Errorf("parse calendar.txt: %w", err)
	}

	for i, record := range records {
		if i == 0 || len(record) < 10 {
			continue
		}

		startDate, err := parseGTFSDate(record[8])
		if err != nil {
			continue
		}
		endDate, err := parseGTFSDate(record[9])
		if err != nil {
			continue
		}

		data.services[record[0]] = calendarService{
			weekdays: [7]bool{
				record[7] == "1",
				record[1] == "1",
				record[2] == "1",
				record[3] == "1",
				record[4] == "1",
				record[5] == "1",
				record[6] == "1",
			},
			startDate: startDate,
			endDate:   endDate,
		}
	}

	if exceptionsFile != nil {
		excReader, err := exceptionsFile.Open()
		if err != nil {
			return data, err
		}
		defer excReader.Close()

		excCSV := csv.NewReader(excReader)
		excRecords, err := excCSV.ReadAll()
		if err != nil {
			return data, fmt.Errorf("parse calendar_dates.txt: %w", err)
		}

		for i, record := range excRecords {
			if i == 0 || len(record) < 3 {
				continue
			}
			exceptionType, err := strconv.Atoi(record[2])
			if err != nil {
				continue
			}
			if data.exceptions[record[0]] == nil {
				data.exceptions[record[0]] = make(map[string]int)
			}
			data.exceptions[record[0]][record[1]] = exceptionType
		}
	}

	return data, nil
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

func findNextDepartures(departures []lineDeparture, calendar calendarData, now time.Time, count int) []time.Time {
	if len(departures) == 0 || count == 0 {
		return nil
	}

	result := make([]time.Time, 0, count)
	now = now.In(utils.ParisLocation)

	for dayOffset := 0; dayOffset < 8 && len(result) < count; dayOffset++ {
		day := utils.StartOfDayParis(now).AddDate(0, 0, dayOffset)

		for _, dep := range departures {
			if !calendar.isServiceActive(dep.serviceID, day) {
				continue
			}

			departureTime := day.Add(time.Duration(dep.secondsFromMidnight) * time.Second)
			if !departureTime.After(now) {
				continue
			}

			result = append(result, departureTime)
			if len(result) == count {
				return result
			}
		}
	}

	return result
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
