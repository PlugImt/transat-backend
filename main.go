package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	// handlers "Transat_2.0_Backend/handlers" // Import handlers if needed directly (usually not)
	restaurantHandler "Transat_2.0_Backend/handlers/restaurant" // Import restaurant handler explicitly
	"Transat_2.0_Backend/i18n"

	"Transat_2.0_Backend/routes"
	"Transat_2.0_Backend/services"
	"Transat_2.0_Backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	_ "github.com/lib/pq"
	_ "github.com/nicksnyder/go-i18n/v2/i18n" // Keep for i18n initialization if still needed here
	"github.com/pressly/goose/v3"
)

var db *sql.DB
var jwtSecret []byte

// Menu check state - Consider moving this into RestaurantHandler or a dedicated service
var menuCheckedToday bool
var menuCheckMutex sync.Mutex

func init() {
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

	// Ensure data folder exists
	if err := utils.EnsureDataFolder(); err != nil {
		log.Fatalf("Failed to ensure data folder: %v", err)
	}
}

func main() {
	app := fiber.New()

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

	// Initialize Handlers that need explicit instantiation (e.g., for Cron)
	restHandler := restaurantHandler.NewRestaurantHandler(db, translationService, notificationService)

	// Cron Jobs - Requires access to handlers/services
	c := cron.New()

	// Reset menuCheckedToday flag at midnight
	_, err = c.AddFunc("0 0 * * *", func() {
		menuCheckMutex.Lock()
		menuCheckedToday = false
		menuCheckMutex.Unlock()
		log.Println("Reset menu check flag for new day")
	})
	if err != nil {
		log.Fatalf("Error scheduling midnight reset: %v", err)
	}

	// Check for menu updates every 10 minutes
	_, err = c.AddFunc("*/10 * * * *", func() {
		menuCheckMutex.Lock()
		shouldCheck := !menuCheckedToday
		menuCheckMutex.Unlock()

		if !shouldCheck {
			return
		}

		updated, err := restHandler.CheckAndUpdateMenuCron()

		if err != nil {
			log.Printf("Error checking menu via cron: %v", err)
			return
		}

		if updated {
			log.Println("Menu updated via cron.")
			menuCheckMutex.Lock()
			if time.Now().Hour() > 10 {
				menuCheckedToday = true
				log.Println("Marking menu as checked for today (update occurred after 10 AM).")
			} else {
				log.Println("Menu updated before 10 AM, will allow further checks today.")
			}
			menuCheckMutex.Unlock()
		} else {
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling menu check: %v", err)
	}

	c.Start()
	defer c.Stop()

	// Middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))

	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02 15:04:05",
	}))

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
	routes.SetupRealCampusRoutes(api, db) // Existing RealCampus routes
	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
