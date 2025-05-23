package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

type NaolibService struct {
	mu         sync.Mutex
	cache      map[string][]models.Departure
	cacheTime  time.Duration
	lastUpdate map[string]time.Time
}

func NewNaolibService(refreshTime time.Duration) *NaolibService {
	return &NaolibService{
		cache:      make(map[string][]models.Departure),
		cacheTime:  refreshTime,
		lastUpdate: make(map[string]time.Time),
	}
}

func (s *NaolibService) GetNextDepartures(stopID string) ([]models.Departure, error) {
	s.mu.Lock()
	if lastUpdateTime, updateExists := s.lastUpdate[stopID]; updateExists {
		if departures, ok := s.cache[stopID]; ok && time.Since(lastUpdateTime) < s.cacheTime {
			s.mu.Unlock()
			utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Using cached departures for %s", stopID))
			return departures, nil
		}
	}

	url := fmt.Sprintf("https://open.tan.fr/ewp/tempsattentelieu.json/%s/2", stopID)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la requête HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur HTTP: code %d", resp.StatusCode)
	}

	var departuresData []models.Departure
	if err := json.NewDecoder(resp.Body).Decode(&departuresData); err != nil {
		return nil, fmt.Errorf("erreur lors du décodage JSON: %v", err)
	}

	s.cache[stopID] = departuresData
	s.lastUpdate[stopID] = time.Now()
	s.mu.Unlock()

	return departuresData, nil
}
