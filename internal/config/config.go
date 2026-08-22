package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/plugimt/transat-backend/models"
)

const defaultGTFSURL = "https://transport.data.gouv.fr/resources/84101/download"
const defaultGTFSMaxDepartures = 3

//go:embed gtfs_lines.json
var embeddedGTFSLines []byte

type Config struct {
	// Server
	Port           string
	AllowedOrigins string

	// Database
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string

	// JWT
	JWTSecret []byte

	// Email
	EmailHost       string
	EmailPort       string
	EmailSender     string
	EmailPassword   string
	EmailSenderName string

	// GTFS
	GTFSURL           string
	GTFSLines         []models.GTFSLineConfig
	GTFSMaxDepartures int
}

func Load() *Config {
	// Load .env if available; keep logging on error as before
	if err := godotenv.Load(); err != nil {
		log.Println("Info: ℹ️ Error loading .env file: ", err)
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Fatal("💥 JWT_SECRET environment variable is not set")
	}

	cfg := &Config{
		Port:           firstNonEmpty(os.Getenv("PORT"), "3000"),
		AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),

		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),

		JWTSecret: jwtSecret,

		EmailHost:       os.Getenv("EMAIL_HOST"),
		EmailPort:       os.Getenv("EMAIL_PORT"),
		EmailSender:     os.Getenv("EMAIL_SENDER"),
		EmailPassword:   os.Getenv("EMAIL_PASSWORD"),
		EmailSenderName: os.Getenv("EMAIL_SENDER_NAME"),

		GTFSURL:           firstNonEmpty(os.Getenv("GTFS_URL"), defaultGTFSURL),
		GTFSLines:         loadGTFSLines(),
		GTFSMaxDepartures: parsePositiveInt(os.Getenv("GTFS_DEPARTURE_COUNT"), defaultGTFSMaxDepartures),
	}

	return cfg
}

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName,
	)
}

func loadGTFSLines() []models.GTFSLineConfig {
	if raw := strings.TrimSpace(os.Getenv("GTFS_LINES")); raw != "" {
		return mustParseGTFSLines([]byte(raw), "GTFS_LINES")
	}

	if path := strings.TrimSpace(os.Getenv("GTFS_LINES_FILE")); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("💥 Unable to read GTFS_LINES_FILE %s: %v", path, err)
		}
		return mustParseGTFSLines(data, path)
	}

	return mustParseGTFSLines(embeddedGTFSLines, "embedded gtfs_lines.json")
}

func mustParseGTFSLines(data []byte, source string) []models.GTFSLineConfig {
	var lines []models.GTFSLineConfig
	if err := json.Unmarshal(data, &lines); err != nil {
		log.Fatalf("💥 Invalid GTFS line config in %s: %v", source, err)
	}
	if len(lines) == 0 {
		log.Fatalf("💥 GTFS line config in %s is empty", source)
	}
	return lines
}

func parsePositiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
