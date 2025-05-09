package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type CachedWeather struct {
	Data      *models.WeatherData
	Timestamp time.Time
}

const cacheDuration = 5 * time.Minute

type WeatherService struct {
	apiKey    string
	lat       string
	lon       string
	cache     map[string]CachedWeather
	cacheLock sync.Mutex
}

func NewWeatherService() (*WeatherService, error) {
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENWEATHERMAP_API_KEY environment variable is not set")
	}
	lat := "47.218371"
	lon := "-1.553621"

	return &WeatherService{
		apiKey:    apiKey,
		lat:       lat,
		lon:       lon,
		cache:     make(map[string]CachedWeather),
		cacheLock: sync.Mutex{},
	}, nil
}

func (s *WeatherService) GetWeather(lang string) (*models.WeatherData, error) {
	if lang == "" {
		lang = "fr"
	}

	s.cacheLock.Lock()
	cachedEntry, found := s.cache[lang]
	if found && time.Since(cachedEntry.Timestamp) < cacheDuration {
		s.cacheLock.Unlock()
		return cachedEntry.Data, nil
	}
	s.cacheLock.Unlock()

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&units=metric&appid=%s&lang=%s", s.lat, s.lon, s.apiKey, lang)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get weather data: status code %d", resp.StatusCode)
	}

	var data models.WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode weather data: %w", err)
	}

	// Translate the Weather[].Main field using TranslationService
	localizer := i18n.GetLocalizer(lang)

	if localizer == nil {
		// This case should ideally be handled within GetLocalizer itself
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Error: Could not get localizer for language '%s'. Using default.", lang))
		localizer = i18n.GetLocalizer("fr") // Fallback
	}

	if len(data.Weather) > 0 {
		for i := range data.Weather {
			englishMain := data.Weather[i].Main
			conditionKey := strings.ToLower(englishMain)
			messageID := "weather." + conditionKey

			config := &goi18n.LocalizeConfig{
				MessageID: messageID,
			}

			message, err := localizer.Localize(config)
			if err != nil {
				log.Printf("Error translating key '%s' for language '%s': %v", messageID, lang, err)
				return nil, fmt.Errorf("error translating key '%s' for language '%s': %v", messageID, lang, err)
			}

			if message != "" && message != messageID {
				data.Weather[i].Main = message
			} else {
				data.Weather[i].Main = englishMain
			}
		}
	}

	// Store in cache
	s.cacheLock.Lock()
	s.cache[lang] = CachedWeather{
		Data:      &data,
		Timestamp: time.Now(),
	}
	s.cacheLock.Unlock()

	return &data, nil
}
