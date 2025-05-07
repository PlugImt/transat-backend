package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// weatherTranslations map[englishCondition]map[lang]translation
var weatherTranslations = map[string]map[string]string{
	"Thunderstorm": {
		"fr": "Orage",
		"es": "Tormenta",
		"de": "Gewitter",
	},
	"Drizzle": {
		"fr": "Bruine",
		"es": "Llovizna",
		"de": "Nieselregen",
	},
	"Rain": {
		"fr": "Pluie",
		"es": "Lluvia",
		"de": "Regen",
	},
	"Snow": {
		"fr": "Neige",
		"es": "Nieve",
		"de": "Schnee",
	},
	"Mist": {
		"fr": "Brume",
		"es": "Neblina",
		"de": "Nebel",
	},
	"Smoke": {
		"fr": "Fumée",
		"es": "Humo",
		"de": "Rauch",
	},
	"Haze": {
		"fr": "Brume sèche",
		"es": "Calima",
		"de": "Dunst",
	},
	"Dust": {
		"fr": "Poussière",
		"es": "Polvo",
		"de": "Staub",
	},
	"Fog": {
		"fr": "Brouillard",
		"es": "Niebla",
		"de": "Nebel",
	},
	"Sand": {
		"fr": "Sable",
		"es": "Arena",
		"de": "Sand",
	},
	"Ash": {
		"fr": "Cendre",
		"es": "Ceniza",
		"de": "Asche",
	},
	"Squall": {
		"fr": "Rafale",
		"es": "Chubasco",
		"de": "Böe",
	},
	"Tornado": {
		"fr": "Tornade",
		"es": "Tornado",
		"de": "Tornado",
	},
	"Clear": {
		"fr": "Dégagé",
		"es": "Despejado",
		"de": "Klar",
	},
	"Clouds": {
		"fr": "Nuageux",
		"es": "Nublado",
		"de": "Bewölkt",
	},
}

type WeatherData struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

type CachedWeather struct {
	Data      *WeatherData
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

	// Coordonnées pour Nantes
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

func (s *WeatherService) GetWeather(lang string) (*WeatherData, error) {
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

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode weather data: %w", err)
	}

	// Traduire le champ Weather[].Main et le stocker directement dans Main
	if len(data.Weather) > 0 {
		for i := range data.Weather {
			englishMain := data.Weather[i].Main
			if translations, ok := weatherTranslations[englishMain]; ok {
				if translatedText, langOk := translations[lang]; langOk {
					data.Weather[i].Main = translatedText
				}
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
