package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	_ "github.com/lib/pq"
	_ "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pressly/goose/v3"

	"Transat_2.0_Backend/routes"
)

var db *sql.DB
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var menuCheckedToday bool
var menuCheckMutex sync.Mutex

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Warning: 💥 Error loading .env file: ", err)
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

	// If there is an error connecting to the database, exit the program
	if err != nil {
		log.Fatalf("💥 Error connecting to the database : %v", err)
	} else {
		log.Println("Connected to the database")
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
}

func main() {
	app := fiber.New()
	notificationService := NewNotificationService(db)
	c := cron.New()

	// Reset menuCheckedToday flag at midnight
	_, err := c.AddFunc("0 0 * * *", func() {
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
		fmt.Println("Checking for menu updates...")
		menuCheckMutex.Lock()
		if menuCheckedToday {
			menuCheckMutex.Unlock()
			log.Println("Menu already updated today, skipping check")
			return
		}
		menuCheckMutex.Unlock()

		log.Println("Checking for menu updates...")
		updated, err := checkAndUpdateMenu(notificationService)
		if err != nil {
			log.Printf("Error checking menu: %v", err)
			return
		}

		if updated {
			menuCheckMutex.Lock()
			if time.Now().Hour() > 12 {
				menuCheckedToday = true
			}
			menuCheckMutex.Unlock()
			log.Println("Menu updated and notifications sent, won't check again today")
		}
	})

	if err != nil {
		log.Fatalf("Error scheduling daily menu notification: %v", err)
	}

	c.Start()
	defer c.Stop()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))

	loginRegisterLimiter := limiter.New(limiter.Config{
		Max:        3,
		Expiration: 200 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	})

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})

	api := app.Group("/api")

	// User routes
	newf := api.Group("/newf")
	newf.Post("/", loginRegisterLimiter, register)
	newf.Delete("/", jwtMiddleware, deleteNewf)
	newf.Get("/me", jwtMiddleware, getNewf)
	newf.Patch("/me", jwtMiddleware, updateNewf)
	newf.Post("/notification", jwtMiddleware, addNotification)
	newf.Get("/notification", jwtMiddleware, getNotification)
	//newf.Get("/", getAllNewfs)
	//newf.Get(":id", getNewf)
	//newf.Put(":id", updateNewf)
	//newf.Delete(":id", deleteNewf)
	newf.Post("/send-notification", sendNotification)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", loginRegisterLimiter, login)
	auth.Post("/verify-account", verifyAccount)
	auth.Post("/verification-code", verificationCode)
	auth.Patch("/change-password", changePassword)

	// Traq routes
	traq := api.Group("/traq")
	traq.Post("/", createTraqArticle)
	traq.Get("/", getAllTraqArticles)
	//traq.Get(":id", getTraqArticle)
	//traq.Put(":id", updateTraqArticle)
	//traq.Delete(":id", deleteTraqArticle)

	traqTypes := traq.Group("/types")
	traqTypes.Post("/", createTraqTypes)
	traqTypes.Get("/", getAllTraqTypes)

	// restaurant routes
	restaurant := api.Group("/restaurant")
	restaurant.Get("/", getRestaurant)

	api.Post("/upload", jwtMiddleware, uploadImage)
	api.Get("/data/:filename", serveImage)

	api.Get("/files", jwtMiddleware, listUserFiles)
	api.Delete("/files/:filename", jwtMiddleware, deleteFile)

	api.Get("/all-files", jwtMiddleware, listAllFiles)

	// Setup RealCampus routes.
	routes.SetupRealCampusRoutes(api, db)

	// Initialiser i18n
	if err := initI18n(); err != nil {
		log.Fatalf("Failed to initialize i18n: %v", err)
	}

	log.Fatal(app.Listen(":3000"))
}
