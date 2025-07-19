package main

import (
	"database/sql"
	"fmt"
	"github.com/plugimt/transat-backend/handlers/event"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/plugimt/transat-backend/handlers"
	"github.com/plugimt/transat-backend/handlers/club"
	restaurantHandler "github.com/plugimt/transat-backend/handlers/restaurant"
	"github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/routes"
	"github.com/plugimt/transat-backend/scheduler"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
	"github.com/robfig/cron/v3"

	_ "github.com/lib/pq"
	_ "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pressly/goose/v3"
)

var db *sql.DB
var jwtSecret []byte

func init() {
	utils.SentryInit()

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
		StrictRouting:           true,
		EnableTrustedProxyCheck: true,
		BodyLimit:               15 * 1024 * 1024, // 15MB
		ReadTimeout:             10 * time.Second,
		WriteTimeout:            15 * time.Second,
		IdleTimeout:             120 * time.Second,
	})

	app.Use(utils.SentryHandler)

	notificationService := services.NewNotificationService(db)
	translationService, err := services.NewTranslationService()
	if err != nil {
		log.Fatalf("💥 Failed to create Translation Service: %v", err)
	}

	emailService := services.NewEmailService(
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_PORT"),
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_SENDER_NAME"),
	)

	statisticsService := services.NewStatisticsService(db)

	restHandler := restaurantHandler.NewRestaurantHandler(db, translationService, notificationService)

	weatherService, err := services.NewWeatherService()
	if err != nil {
		log.Fatalf("💥 Failed to create Weather Service: %v", err)
	}
	weatherHandler := handlers.NewWeatherHandler(weatherService)

	r2Service, err := services.NewR2Service()
	if err != nil {
		log.Fatalf("💥 Failed to create R2 Service: %v", err)
	}

	clubsHandler := club.NewclubHandler(db)

	eventHandler := event.NeweventHandler(db)

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
		corsConfig.AllowOrigins = os.Getenv("ALLOWED_ORIGINS")
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
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// 8. Add statistics middleware to capture all requests
	app.Use(middlewares.StatisticsMiddleware(statisticsService))

	// Status Route
	app.Get("/status", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})

	// API Group --- NEW ROUTES
	routes.SetupAuthRoutes(app, db, jwtSecret, notificationService, emailService)
	routes.SetupUserRoutes(app, db, notificationService)
	routes.SetupTraqRoutes(app, db)
	routes.SetupFileRoutes(app, db, r2Service)
	routes.SetupRestaurantRoutes(app, restHandler)
	routes.SetupClubRoutes(app, clubsHandler)
	routes.SetupRealCampusRoutes(app, db)
	routes.SetupPlanningRoutes(app, db)
	routes.SetupNotificationRoutes(app, db, notificationService)
	routes.SetupStatisticsRoutes(app, db, statisticsService)
	routes.SetupWashingMachineRoutes(app)
	routes.SetupWeatherRoutes(app, weatherHandler)
	routes.SetupNotificationRoutes(app, db, notificationService)
	routes.SetupEventRoutes(app, eventHandler)

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
