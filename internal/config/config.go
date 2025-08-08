package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

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
	}

	return cfg
}

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName,
	)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
