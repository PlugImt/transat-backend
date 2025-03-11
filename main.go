package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var menuCheckedToday bool
var menuCheckMutex sync.Mutex

func init() {
	err := godotenv.Load()

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

	// Check for menu updates every 15 minutes
	_, err = c.AddFunc("*/2 * * * *", func() {
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
			menuCheckedToday = true
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

	log.Fatal(app.Listen(":3000"))
}
