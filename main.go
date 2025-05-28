package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/plugimt/transat-backend/handlers"
	restaurantHandler "github.com/plugimt/transat-backend/handlers/restaurant" // Import restaurant handler explicitly
	"github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/routes"
	"github.com/plugimt/transat-backend/scheduler" // Import our scheduler package
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
	"github.com/robfig/cron/v3"

	_ "github.com/lib/pq"
	_ "github.com/nicksnyder/go-i18n/v2/i18n" // Keep for i18n initialization if still needed here
	"github.com/pressly/goose/v3"
)

// Config holds all configuration for the application
type Config struct {
	JWTSecret              string
	DBUser                 string
	DBPass                 string
	DBHost                 string
	DBPort                 string
	DBName                 string
	EmailHost              string
	EmailPort              string
	EmailSender            string
	EmailPassword          string
	EmailSenderName        string
	AllowedOrigins         string
	Env                    string
	Port                   string
	SentryDSN              string
	SentryTracesSampleRate float64
	LogLevel               string
	LogFormat              string
	GoogleTranslateAPIKey  string
	OpenWeatherMapAPIKey   string
	JWTExpirationHours     string
	DataFolder             string
	ExpoPushURL            string
}

// loadConfig loads configuration from environment variables
func loadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: ℹ️ Error loading .env file, relying on environment variables: ", err)
	}

	sentryDSN := os.Getenv("SENTRY_DSN")

	sentryTracesSampleRateStr := utils.GetEnvOrDefault("SENTRY_TRACES_SAMPLE_RATE", "1.0")
	sentryTracesSampleRate, err := strconv.ParseFloat(sentryTracesSampleRateStr, 64)
	if err != nil {
		log.Printf("Info: ℹ️ SENTRY_TRACES_SAMPLE_RATE ('%s') is not a valid float, defaulting to 1.0. Error: %v", sentryTracesSampleRateStr, err)
		sentryTracesSampleRate = 1.0
	}

	cfg := &Config{
		JWTSecret:              os.Getenv("JWT_SECRET"),
		DBUser:                 utils.GetEnvOrDefault("DB_USER", "postgres"),
		DBPass:                 os.Getenv("DB_PASS"),
		DBHost:                 utils.GetEnvOrDefault("DB_HOST", "localhost"),
		DBPort:                 utils.GetEnvOrDefault("DB_PORT", "5432"),
		DBName:                 utils.GetEnvOrDefault("DB_NAME", "transat_db"),
		EmailHost:              utils.GetEnvOrDefault("EMAIL_HOST", "localhost"),
		EmailPort:              utils.GetEnvOrDefault("EMAIL_PORT", "1025"),
		EmailSender:            utils.GetEnvOrDefault("EMAIL_SENDER", "noreply@example.com"),
		EmailPassword:          os.Getenv("EMAIL_PASSWORD"),
		EmailSenderName:        utils.GetEnvOrDefault("EMAIL_SENDER_NAME", "Transat App"),
		AllowedOrigins:         utils.GetEnvOrDefault("ALLOWED_ORIGINS", "*"),
		Env:                    utils.GetEnvOrDefault("ENV", "development"),
		Port:                   utils.GetEnvOrDefault("PORT", "3000"),
		SentryDSN:              sentryDSN,
		SentryTracesSampleRate: sentryTracesSampleRate,
		LogLevel:               utils.GetEnvOrDefault("LOG_LEVEL", "INFO"),
		LogFormat:              utils.GetEnvOrDefault("LOG_FORMAT", "text"),
		GoogleTranslateAPIKey:  os.Getenv("GOOGLE_TRANSLATE_API_KEY"),
		OpenWeatherMapAPIKey:   os.Getenv("OPENWEATHERMAP_API_KEY"),
		JWTExpirationHours:     utils.GetEnvOrDefault("JWT_EXPIRATION_HOURS", "24"),
		DataFolder:             utils.GetEnvOrDefault("DATA_FOLDER", "./data"),
		ExpoPushURL:            utils.GetEnvOrDefault("EXPO_PUSH_URL", "https://api.expo.dev/v2/push/send"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("💥 JWT_SECRET environment variable is not set")
	}
	if cfg.Port == "" {
		cfg.Port = "3000" // Default port
	}

	return cfg, nil
}

func initGlobal(cfg *Config) {
	utils.SentryInit(cfg.SentryDSN, utils.GetEnvName(), utils.GetEnvCommitSHA(), utils.GetEnvHost(), cfg.SentryTracesSampleRate)
	utils.InitLogger(cfg.Env, cfg.LogLevel)

	if err := utils.InitJWT([]byte(cfg.JWTSecret), cfg.JWTExpirationHours); err != nil {
		log.Fatalf("💥 Failed to initialize JWT: %v", err)
	}

	if err := utils.InitFilePaths(cfg.DataFolder); err != nil {
		log.Fatalf("💥 Failed to initialize file paths: %v", err)
	}

	if err := i18n.Init(); err != nil {
		log.Fatalf("Failed to initialize i18n: %v", err)
	}
}

func initDB(cfg *Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("💥 Error connecting to the database: %v", err)
	}

	if err = database.Ping(); err != nil {
		return nil, fmt.Errorf("💥 Error pinging the database: %v", err)
	}
	log.Println("Successfully connected to the database and connection verified")

	log.Println("Running database migrations...")
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("💥 Failed to set goose dialect: %v", err)
	}
	if err := goose.Up(database, "db/migrations"); err != nil {
		return nil, fmt.Errorf("💥 Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	return database, nil
}

var db *sql.DB
var jwtSecret []byte
var appConfig *Config

func init() {
	var err error
	appConfig, err = loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret = []byte(appConfig.JWTSecret) // Keep global jwtSecret for now, or pass config around

	initGlobal(appConfig)

	// Connect to the database
	db, err = initDB(appConfig)
	if err != nil {
		log.Fatal(err)
	}
}

// Services holds all initialized services
type Services struct {
	Notification *services.NotificationService
	Translation  *services.TranslationService
	Email        *services.EmailService
	Statistics   *services.StatisticsService
	Weather      *services.WeatherService
}

func initServices(cfg *Config, db *sql.DB) (*Services, error) {
	notificationService := services.NewNotificationService(db, cfg.ExpoPushURL)
	translationService, err := services.NewTranslationService(cfg.GoogleTranslateAPIKey)
	if err != nil {
		return nil, fmt.Errorf("💥 Failed to create Translation Service: %v", err)
	}
	emailService := services.NewEmailService(
		cfg.EmailHost,
		cfg.EmailPort,
		cfg.EmailSender,
		cfg.EmailPassword,
		cfg.EmailSenderName,
	)
	statisticsService := services.NewStatisticsService(db)
	weatherService, err := services.NewWeatherService(cfg.OpenWeatherMapAPIKey)
	if err != nil {
		return nil, fmt.Errorf("💥 Failed to create Weather Service: %v", err)
	}

	return &Services{
		Notification: notificationService,
		Translation:  translationService,
		Email:        emailService,
		Statistics:   statisticsService,
		Weather:      weatherService,
	}, nil
}

// Handlers holds all initialized handlers
type Handlers struct {
	Restaurant *restaurantHandler.RestaurantHandler
	Weather    *handlers.WeatherHandler
}

func initHandlers(db *sql.DB, s *Services) *Handlers {
	restHandler := restaurantHandler.NewRestaurantHandler(db, s.Translation, s.Notification)
	weatherHandler := handlers.NewWeatherHandler(s.Weather)

	return &Handlers{
		Restaurant: restHandler,
		Weather:    weatherHandler,
	}
}

func initScheduler(h *Handlers) *scheduler.Scheduler {
	appScheduler := scheduler.NewScheduler(h.Restaurant)
	appScheduler.StartAll()
	return appScheduler
}

func setupMiddlewares(app *fiber.App, cfg *Config, s *Services) {
	app.Use(utils.SentryHandler)

	// 1. Add security headers to all responses
	app.Use(middlewares.SecurityHeadersMiddleware())

	// 2. CORS configuration
	corsConfig := cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}
	if cfg.Env == "production" {
		corsConfig.AllowOrigins = cfg.AllowedOrigins
	} else {
		corsConfig.AllowOrigins = "*"
	}
	app.Use(cors.New(corsConfig))

	// 3. Add general rate limiting
	app.Use(middlewares.GeneralRateLimiter)

	// 4. Add input validation middleware
	app.Use(middlewares.ValidationMiddleware())

	// 5. Add SQL injection protection
	app.Use(middlewares.SQLInjectionProtectionMiddleware())

	// 6. Add XSS protection
	app.Use(middlewares.XSSProtectionMiddleware())

	// 7. Logger middleware
	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02 15:04:05", // Consider making this configurable
	}))

	// 8. Add statistics middleware to capture all requests
	app.Use(middlewares.StatisticsMiddleware(s.Statistics))
}

func setupRoutes(api fiber.Router, db *sql.DB, jwtSecret []byte, s *Services, h *Handlers) {
	routes.SetupAuthRoutes(api, db, jwtSecret, s.Notification, s.Email)
	routes.SetupUserRoutes(api, db, s.Notification)
	routes.SetupTraqRoutes(api, db)
	routes.SetupFileRoutes(api, db)
	routes.SetupRestaurantRoutes(api, h.Restaurant)
	routes.SetupRealCampusRoutes(api, db)
	routes.SetupStatisticsRoutes(api, db, s.Statistics)
	routes.SetupWashingMachineRoutes(api)
	routes.SetupWeatherRoutes(api, h.Weather)
	routes.SetupNotificationRoutes(api, db, s.Notification)
}

func main() {
	app := fiber.New(fiber.Config{
		StrictRouting:           true,
		EnableTrustedProxyCheck: true,
		BodyLimit:               15 * 1024 * 1024, // 15MB
		ReadTimeout:             10 * time.Second,
		WriteTimeout:            15 * time.Second,
		IdleTimeout:             120 * time.Second,
	})

	// Initialize Services
	appServices, err := initServices(appConfig, db)
	if err != nil {
		log.Fatalf("💥 Failed to initialize services: %v", err)
	}

	// Initialize Handlers
	appHandlers := initHandlers(db, appServices)

	// Initialize and start schedulers
	appScheduler := initScheduler(appHandlers)
	defer appScheduler.StopAll()

	// Cron Jobs
	c := cron.New()
	// Example: c.AddFunc("@daily", func() { appHandlers.Restaurant.SomeDailyTask() })
	c.Start()
	defer c.Stop()

	// Setup Middlewares
	setupMiddlewares(app, appConfig, appServices)

	// Status Route
	app.Get("/status", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})

	// API Group
	api := app.Group("/api")

	// Setup Routes
	setupRoutes(api, db, jwtSecret, appServices, appHandlers)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Start Server
	port := appConfig.Port
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
