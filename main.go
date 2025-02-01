package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatalf("Error loading .env file: %v", err)
	//}
	var err error
	//user := os.Getenv("DB_USER")
	//pass := os.Getenv("DB_PASS")
	//host := os.Getenv("DB_HOST")
	//port := os.Getenv("DB_PORT")
	//name := os.Getenv("DB_NAME")

	connStr := os.Getenv("Postgres.DATABASE_URL")

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Successfully connected to the database!")
}

func main() {
	app := fiber.New()

	// Update CORS configuration to allow WebSocket
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Turnstile-Token",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	loginRegisterLimiter := limiter.New(limiter.Config{
		Max:        3,
		Expiration: 300 * time.Second,
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

	app.Get("/login", loginRegisterLimiter, login)
	app.Get("/traq", traq)

	log.Fatal(app.Listen(":3000"))
}
