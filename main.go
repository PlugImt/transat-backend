package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"github.com/plugimt/transat-backend/handlers/event"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	_ "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/plugimt/transat-backend/handlers"
	"github.com/plugimt/transat-backend/handlers/club"
	devicesHandler "github.com/plugimt/transat-backend/handlers/devices"
	photoboothHandler "github.com/plugimt/transat-backend/handlers/photobooth"
	reportHandler "github.com/plugimt/transat-backend/handlers/photoreport"
	restaurantHandler "github.com/plugimt/transat-backend/handlers/restaurant"
	userHandlers "github.com/plugimt/transat-backend/handlers/user"
	"github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/internal/config"
	"github.com/plugimt/transat-backend/internal/database"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/routes"
	"github.com/plugimt/transat-backend/scheduler"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/accounts"
	devicesService "github.com/plugimt/transat-backend/services/devices"
	photoboothService "github.com/plugimt/transat-backend/services/photobooth"
	"github.com/plugimt/transat-backend/services/photoreport"
	photosService "github.com/plugimt/transat-backend/services/photos"
	"github.com/plugimt/transat-backend/services/push"
	"github.com/plugimt/transat-backend/services/reservations"
	"github.com/plugimt/transat-backend/utils"
	"github.com/robfig/cron/v3"
)

var db *sql.DB
var jwtSecret []byte

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: ℹ️ Error loading .env file: ", err)
	}

	utils.SentryInit()

	cfg := config.Load()
	jwtSecret = cfg.JWTSecret

	var err error
	db, err = database.Open(cfg.DSN(), "db/migrations")
	if err != nil {
		log.Fatalf("💥 Error initializing database: %v", err)
	}
	middlewares.InitJWTMiddleware(db)

	if err := i18n.Init(); err != nil {
		log.Fatalf("Failed to initialize i18n: %v", err)
	}
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
	deviceConfig := devicesService.ConfigFromEnv()
	// PUSH_PIPELINE=legacy goes back to sending straight from the request or the scheduler, without the worker.
	if strings.ToLower(os.Getenv("PUSH_PIPELINE")) != "legacy" {
		pipeline := push.New(db, push.NewExpoProvider(push.ExpoConfigFromEnv()), push.ConfigFromEnv(), deviceConfig)
		notificationService.SetPipeline(pipeline)
		pushWorker := pipeline.Start()
		defer pushWorker.Stop()
	} else {
		log.Println("Info: PUSH_PIPELINE=legacy, push notifications are sent without the delivery pipeline")
	}
	translationService, err := services.NewTranslationService()
	if err != nil {
		log.Fatalf("💥 Failed to create Translation Service: %v", err)
	}

	// Initialize Brevo Email Service (API-based, not SMTP)
	brevoAPIKey := os.Getenv("BREVO_API_KEY")
	if brevoAPIKey == "" {
		log.Fatal("💥 BREVO_API_KEY environment variable is not set")
	}

	emailService := services.NewEmailService(
		brevoAPIKey,
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_SENDER_NAME"),
	)

	statisticsService := services.NewStatisticsService(db)

	// Discord webhook notifications
	discordService := services.NewDiscordService(os.Getenv("DISCORD_WEBHOOK_URL"))
	// Discord reservation-specific webhook (falls back to default if unset)
	reservationDiscordService := services.NewDiscordService(os.Getenv("DISCORD_RESERVATION_ALERT_WEBHOOK"))

	restHandler := restaurantHandler.NewRestaurantHandler(db, translationService, notificationService)

	reservationService := reservations.New(db, notificationService, emailService)

	weatherService, err := services.NewWeatherService()
	if err != nil {
		log.Fatalf("💥 Failed to create Weather Service: %v", err)
	}
	weatherHandler := handlers.NewWeatherHandler(weatherService)

	gtfsService := services.NewGTFSService(services.GTFSOptions{
		URL:           cfg.GTFSURL,
		RealtimeURL:   cfg.GTFSRealtimeURL,
		Lines:         cfg.GTFSLines,
		MaxDepartures: cfg.GTFSMaxDepartures,
	})

	r2Service, err := services.NewR2Service()
	if err != nil {
		log.Fatalf("💥 Failed to create R2 Service: %v", err)
	}

	clubsHandler := club.NewclubHandler(db)

	eventHandler := event.NewEventHandler(db, notificationService)

	appScheduler := scheduler.NewScheduler(restHandler, db, notificationService, reservationService)
	appScheduler.StartAll()
	defer appScheduler.StopAll()

	// Cron Jobs - Requires access to handlers/services
	c := cron.New()

	c.Start()
	defer c.Stop()

	photosSync, err := photosService.FromEnv(context.Background(), db)
	if err != nil {
		log.Printf("Warning: Photos sync disabled: %v", err)
	} else if photosSync == nil {
		log.Println("Info: ℹ️ Photos sync disabled: GOOGLE_DRIVE_CREDENTIALS_JSON and PHOTOS_DRIVE_ID are not set")
	} else {
		photosSync.EnableAlbumNotifications(func(_ context.Context, n photosService.AlbumNotice) (photosService.Delivery, error) {
			d, err := notificationService.NotifyNewPhotoAlbum(n.AlbumID, n.Title, n.EventDate)
			return photosService.Delivery{Users: d.Users, Devices: d.Devices}, err
		})
		photosSync.Schedule(c)
	}

	// Photobooth: booth uploads into the Photos drive and permanent share links.
	var photobooth *photoboothHandler.Handler
	if svc, verifier, err := photoboothService.FromEnv(context.Background(), db); err != nil {
		log.Printf("Warning: Photobooth disabled: %v", err)
	} else if svc == nil {
		log.Println("Info: ℹ️ Photobooth disabled: PHOTOBOOTH_DEVICE_KEYS is not set")
	} else {
		photobooth = photoboothHandler.NewHandler(svc, verifier)
		// Retakes and cancelled sessions are trashed right away; this retries what Drive refused then.
		if _, err := c.AddFunc("@every 15m", func() {
			if _, err := svc.PurgeFrameSessions(context.Background()); err != nil {
				log.Printf("Photobooth frame links cleanup failed: %v", err)
			}
			if _, err := svc.RunTrash(context.Background()); err != nil {
				log.Printf("Photobooth trash job failed: %v", err)
			}
		}); err != nil {
			log.Printf("Photobooth trash scheduling failed: %v", err)
		}
	}

	// Reports are stored before any email is sent; this job retries the ones the email service refused.
	photoReportService := photoreport.NewService(db, emailService, os.Getenv("ADMIN_PANEL_URL"))
	// Approving a report deletes the photo: that needs a Drive account with write access (the photo sync's is
	// read-only). PHOTOS_WRITE_CREDENTIALS_JSON, else the photobooth's account, which manages the same drive.
	if creds := firstEnv("PHOTOS_WRITE_CREDENTIALS_JSON", "PHOTOBOOTH_DRIVE_CREDENTIALS_JSON"); creds != "" {
		if writer, err := photoboothService.NewDriveClient(context.Background(), creds); err != nil {
			log.Printf("Warning: photo deletion from reports disabled: %v", err)
		} else {
			photoReportService.SetRemover(writer)
		}
	} else {
		log.Println("Info: ℹ️ Photo deletion from reports disabled: PHOTOS_WRITE_CREDENTIALS_JSON is not set")
	}
	photoReports := reportHandler.NewHandler(photoReportService)
	if _, err := c.AddFunc("@every 5m", func() {
		if _, err := photoReportService.DeliverDue(context.Background()); err != nil {
			log.Printf("Photo report delivery job failed: %v", err)
		}
	}); err != nil {
		log.Printf("Photo report delivery scheduling failed: %v", err)
	}

	// ---- SECURITY MIDDLEWARES ----
	// 1. Add security headers to all responses
	app.Use(middlewares.SecurityHeadersMiddleware())

	// 2. CORS configuration - restrict to allowed origins in production
	corsConfig := cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowOrigins:     "http://localhost:3001", // pour le pannel admin en dev
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, Content-Type",
	}

	// Use proper CORS origins in production, or * in development
	if os.Getenv("ENV") == "production" {
		corsConfig.AllowOrigins = cfg.AllowedOrigins
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
		Next:       func(c *fiber.Ctx) bool { return c.Path() == middlewares.HealthPath },
	}))

	// 8. Add statistics middleware to capture all requests
	app.Use(middlewares.StatisticsMiddleware(statisticsService))

	// Status Route
	app.Get("/status", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})

	// API Group --- NEW ROUTES
	routes.SetupAuthRoutes(app, db, jwtSecret, emailService, discordService)
	// Account deletion (on request and after 5 years without activity) removes or anonymises everything of a user.
	accountService := accounts.NewService(db, r2Service)
	userHandlers.Accounts = accountService
	if _, err := c.AddFunc("@daily", func() {
		report, err := accountService.PurgeInactive(context.Background(), emailService)
		if err != nil {
			log.Printf("Inactive accounts job failed: %v", err)
		} else if report.Warned > 0 || report.Deleted > 0 {
			log.Printf("Inactive accounts: %d warned, %d deleted", report.Warned, report.Deleted)
		}
	}); err != nil {
		log.Printf("Inactive accounts scheduling failed: %v", err)
	}
	routes.SetupUserRoutes(app, db, notificationService)
	routes.SetupDeviceRoutes(app, devicesHandler.NewHandler(devicesService.NewService(db, deviceConfig)))
	routes.SetupTraqRoutes(app, db)
	routes.SetupFileRoutes(app, db, r2Service)
	routes.SetupRestaurantRoutes(app, restHandler)
	routes.SetupClubRoutes(app, db, clubsHandler)
	routes.SetupPlanningRoutes(app, db)
	routes.SetupStatisticsRoutes(app, db, statisticsService)
	routes.SetupWashingMachineRoutes(app)
	routes.SetupWeatherRoutes(app, weatherHandler)
	routes.SetupBusDepartureRoutes(app, gtfsService)
	routes.SetupEventRoutes(app, db, eventHandler)
	routes.SetupReservationRoutes(app, db, reservationDiscordService, reservationService)
	routes.SetupBassineRoutes(app, db)
	routes.SetupAccessRoutes(app, db)
	routes.SetupPhotosRoutes(app, db, photosSync, photoReports)
	routes.SetupPhotoboothRoutes(app, photobooth)
	routes.SetupProfileRoutes(app, db, r2Service)
	adminGroup := routes.SetupAdminRoutes(app, db, notificationService, accountService)
	routes.SetupAdminReservationRoutes(adminGroup, db, reservationDiscordService, reservationService)

	app.Get(middlewares.HealthPath, func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Start Server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}

// firstEnv returns the first non-empty environment variable among names.
func firstEnv(names ...string) string {
	for _, n := range names {
		if v := strings.TrimSpace(os.Getenv(n)); v != "" {
			return v
		}
	}
	return ""
}
