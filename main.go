package main

import (
	"database/sql"
	"encoding/json"
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
}

func main() {
	app := fiber.New()
	notificationService := NewNotificationService(db)
	c := cron.New()

	// Reset cache at midnight
	_, err := c.AddFunc("0 0 * * *", func() {
		resetCache()
		log.Println("Cache reset at midnight")
	})

	if err != nil {
		log.Fatalf("Error scheduling midnight cache reset: %v", err)
	}

	// Check for menu updates every 10 minutes
	_, err = c.AddFunc("*/10 * * * *", func() {
		log.Println("Checking for menu updates...")
		
		// Get the latest menu from database
		request := `
			SELECT restaurant.articles, restaurant.updated_date 
			FROM restaurant 
			WHERE restaurant.language = 1 
			ORDER BY updated_date DESC 
			LIMIT 1;
		`

		stmt, err := db.Prepare(request)
		if err != nil {
			log.Printf("Error preparing statement: %v", err)
			return
		}
		defer stmt.Close()

		var articles string
		var updatedDate string
		err = stmt.QueryRow().Scan(&articles, &updatedDate)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error getting latest menu: %v", err)
			return
		}

		// If no menu in database, fetch from API
		if err == sql.ErrNoRows {
			menuData, err := fetchMenuFromAPI()
			if err != nil {
				log.Printf("Error fetching menu from API: %v", err)
				return
			}

			// Update database with new menu
			menuJSON, err := json.Marshal(menuData)
			if err != nil {
				log.Printf("Error marshaling menu data: %v", err)
				return
			}

			updateQuery := `
				INSERT INTO restaurant (articles, language) 
				VALUES ($1, $2)
			`

			updateStmt, err := db.Prepare(updateQuery)
			if err != nil {
				log.Printf("Error preparing update statement: %v", err)
				return
			}
			defer updateStmt.Close()

			_, err = updateStmt.Exec(string(menuJSON), 1)
			if err != nil {
				log.Printf("Error updating menu in database: %v", err)
				return
			}

			log.Println("Menu updated successfully")
			err = notificationService.SendDailyMenuNotification()
			if err != nil {
				log.Printf("Error sending notification: %v", err)
			}
		}
	})

	if err != nil {
		log.Fatalf("Error scheduling menu check: %v", err)
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
