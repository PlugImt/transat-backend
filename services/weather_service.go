package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

const openWeatherMapURL = "https://api.openweathermap.org/data/2.5/weather"

// WeatherService is responsible for fetching weather data.
type WeatherService struct {
	apiKey     string
	httpClient *http.Client
	lat        string
	lon        string
	cache      map[string]CachedWeather
	cacheLock  sync.Mutex
}

// NewWeatherService creates a new WeatherService.
// It requires the OpenWeatherMap API key to be passed as a parameter.
func NewWeatherService(apiKey string) (*WeatherService, error) {
	if apiKey == "" {
		utils.LogMessage(utils.LevelWarn, "OPENWEATHERMAP_API_KEY is not set. Weather service will be disabled.")
		// Return a service instance that will gracefully fail or indicate it's disabled.
		// Depending on requirements, this might return an error instead.
		return &WeatherService{apiKey: "", httpClient: &http.Client{}}, nil // Client is not nil, but API calls will fail
	}
	return &WeatherService{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		lat:        "47.218371",
		lon:        "-1.553621",
		cache:      make(map[string]CachedWeather),
		cacheLock:  sync.Mutex{},
	}, nil
}

func (s *WeatherService) GetWeather(lang string) (*models.WeatherData, error) {
	if s.apiKey == "" { // Check if service is disabled
		return nil, fmt.Errorf("weather service is disabled (API key not configured)")
	}

	if lang == "" {
		lang = "fr" // Default language
	}

	s.cacheLock.Lock()
	cachedEntry, found := s.cache[lang]
	if found && time.Since(cachedEntry.Timestamp) < cacheDuration {
		s.cacheLock.Unlock()
		return cachedEntry.Data, nil
	}
	s.cacheLock.Unlock()

	url := fmt.Sprintf("%s?lat=%s&lon=%s&units=metric&appid=%s&lang=%s", openWeatherMapURL, s.lat, s.lon, s.apiKey, lang)

	if !strings.HasPrefix(url, "https://api.openweathermap.org/") {
		return nil, fmt.Errorf("invalid weather API URL")
	}

	resp, err := s.httpClient.Get(url) // #nosec G107
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
				return nil, fmt.Errorf("error translating weather condition '%s' for language '%s': %v", messageID, lang, err)
			}

			if message != "" && message != messageID { // Check if translation was successful and different
				data.Weather[i].Main = message
			} else {
				// If translation is same as ID or empty, keep original (or handle as untranslated)
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

// GetWeatherByCity fetches the current weather for a given city.
func (s *WeatherService) GetWeatherByCity(city string) (*models.WeatherData, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("weather service is not configured (OpenWeatherMap API key missing)")
	}

	url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", openWeatherMapURL, city, s.apiKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
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

	return &data, nil
}
