package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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

var db *sql.DB
var jwtSecret []byte

func init() {
	utils.SentryInit()
	utils.InitLogger()

	err := godotenv.Load()
	if err != nil {
		log.Println("Info: ℹ️ Error loading .env file: ", err)
	}

	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Fatal("💥 JWT_SECRET environment variable is not set")
	}

	// Connect to the database
	db, err = sql.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		),
	)
	if err != nil {
		log.Fatalf("💥 Error connecting to the database : %v", err)
	} else {
		log.Println("Successfully connected to the database")
	}

	// Ping the database to verify connection
	if err = db.Ping(); err != nil {
		log.Fatalf("💥 Error pinging the database : %v", err)
	} else {
		log.Println("Database connection verified")
	}

	// Run migrations
	log.Println("Running database migrations...")
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("💥 Failed to set goose dialect: %v", err)
	}
	if err := goose.Up(db, "db/migrations"); err != nil {
		log.Fatalf("💥 Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	if err := i18n.Init(); err != nil {
		log.Fatalf("Failed to initialize i18n: %v", err)
	}
}

func main() {
	app := fiber.New(fiber.Config{
		// Set strict routing
		StrictRouting: true,
		// Add security-focused settings
		EnableTrustedProxyCheck: true,
		// Limit the size of request body
		BodyLimit: 15 * 1024 * 1024, // 15MB
		// Set read timeout to avoid slow request attacks
		ReadTimeout: 10 * time.Second,
		// Set write timeout
		WriteTimeout: 15 * time.Second,
		// Set idle timeout
		IdleTimeout: 120 * time.Second,
	})

	app.Use(utils.SentryHandler)

	// Initialize Services
	notificationService := services.NewNotificationService(db)
	translationService, err := services.NewTranslationService()
	if err != nil {
		log.Fatalf("💥 Failed to create Translation Service: %v", err)
	}
	// Initialize Email Service
	emailService := services.NewEmailService(
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_PORT"),
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_SENDER_NAME"),
	)
	// Initialize Statistics Service
	statisticsService := services.NewStatisticsService(db)

	// Initialize Handlers that need explicit instantiation (e.g., for Cron)
	restHandler := restaurantHandler.NewRestaurantHandler(db, translationService, notificationService)

	naolibService := services.NewNaolibService(30 * time.Second)

	// Initialize Weather Service and Handler
	weatherService, err := services.NewWeatherService()
	if err != nil {
		log.Fatalf("💥 Failed to create Weather Service: %v", err)
	}
	weatherHandler := handlers.NewWeatherHandler(weatherService)

	// Initialize and start schedulers
	appScheduler := scheduler.NewScheduler(restHandler)
	appScheduler.StartAll()
	defer appScheduler.StopAll()

	// Cron Jobs - Requires access to handlers/services
	c := cron.New()

	c.Start()
	defer c.Stop()

	// ---- SECURITY MIDDLEWARES ----

	// 1. Add security headers to all responses
	app.Use(middlewares.SecurityHeadersMiddleware())

	// 2. CORS configuration - restrict to allowed origins in production
	corsConfig := cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}

	// Use proper CORS origins in production, or * in development
	if os.Getenv("ENV") == "production" {
		corsConfig.AllowOrigins = os.Getenv("ALLOWED_ORIGINS") // Comma-separated list from env var
	} else {
		corsConfig.AllowOrigins = "*" // Allow all origins in development
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
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// 8. Add statistics middleware to capture all requests
	app.Use(middlewares.StatisticsMiddleware(statisticsService))

	// Status Route
	app.Get("/status", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})

	// API Group
	api := app.Group("/api")

	// Setup Routes using the new structure
	// Pass necessary dependencies (db, jwtSecret, services) to setup functions
	routes.SetupAuthRoutes(api, db, jwtSecret, notificationService, emailService)
	routes.SetupUserRoutes(api, db, notificationService) // Pass notificationService if needed by user handlers
	routes.SetupTraqRoutes(api, db)
	routes.SetupFileRoutes(api, db)
	routes.SetupRestaurantRoutes(api, restHandler)
	routes.SetupRealCampusRoutes(api, db)
	routes.SetupPlanningRoutes(api, db)
	routes.SetupStatisticsRoutes(api, db, statisticsService)     // Setup statistics routes
	routes.SetupWashingMachineRoutes(api)                        // Setup washing machine routes
	routes.SetupWeatherRoutes(api, weatherHandler)               // Setup weather routes
	routes.SetupNotificationRoutes(api, db, notificationService) // Setup notification test routes
	routes.SetupNaolibRoutes(api, naolibService)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
