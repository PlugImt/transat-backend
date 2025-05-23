package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
)

var SentryHandler fiber.Handler

// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
func SentryInit() {
	env := GetEnvName()
	commitSHA := GetEnvCommitSHA()
	host := GetEnvHost()

	log.Printf("║ %s: %s\n", LevelInfo, "Initializing Sentry with env: "+env+", commitSHA: "+commitSHA+", host: "+host)

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         "https://121ef29a5b1be0fb38b2dc2eb116f255@o4509277512859648.ingest.de.sentry.io/4509277597597776",
		Environment: env,
		Tags: map[string]string{
			"host":    host,
			"version": commitSHA,
		},
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					// si besoin, on a accès à la requête ici
					LogMessage(LevelInfo, "Error on: "+req.URL.String())
				}
			}

			return event
		},
		EnableLogs:    true,
		EnableTracing: true,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	SentryHandler = sentryfiber.New(sentryfiber.Options{
		Repanic:         true,
		WaitForDelivery: true,
	})
}

func EnhanceSentryEventWithEmail(ctx *fiber.Ctx) error {
	email := ctx.Locals("email").(string)

	if hub := sentryfiber.GetHubFromContext(ctx); hub != nil {
		if email != "" {
			scope := hub.Scope()
			scope.SetTag("email", email)
		}
	}
	return ctx.Next()
}

// SentryTransactionMiddleware creates a new transaction for each incoming request
// and handles distributed tracing.
func SentryTransactionMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Clone the hub from the current scope. This picks up the hub
		// that might have been configured by utils.SentryHandler or is the global one.
		hub := sentry.CurrentHub().Clone()
		// Put the cloned hub onto the context of this request.
		ctxWithHub := sentry.SetHubOnContext(c.UserContext(), hub)

		traceId := string(c.Request().Header.Peek("x-sentry-trace"))
		if traceId != "" {
			// Log the traceId if present (as per user's addition)
			LogLineKeyValue(LevelDebug, "SentryMiddleware: traceId", traceId)
		}

		options := []sentry.SpanOption{
			sentry.WithOpName("http.server"),               // Standard operation name for HTTP server transactions
			sentry.WithTransactionSource(sentry.SourceURL), // Sets the transaction source to the URL path
		}
		if traceId != "" {
			options = append(options, sentry.ContinueFromTrace(traceId))
		}

		// Start the transaction with the request-specific context.
		transaction := sentry.StartTransaction(ctxWithHub,
			fmt.Sprintf("%s %s", c.Method(), c.Path()), // Transaction name: "GET /path"
			options...,
		)
		// Store the transaction in c.Locals so other handlers can create child spans.
		c.Locals("sentryTransaction", transaction)

		// Defer finishing the transaction until the request is done.
		defer func() {
			// Set HTTP status on the transaction.
			statusCode := c.Response().StatusCode()
			transaction.Status = sentry.HTTPtoSpanStatus(statusCode)
			transaction.SetTag("http.status_code", fmt.Sprintf("%d", statusCode))
			transaction.SetData("http.method", c.Method())
			// Use OriginalURL to get the full URL before any modifications by handlers
			transaction.SetData("http.url", c.OriginalURL())
			transaction.Finish()
		}()

		// Continue to the next middleware or handler.
		err := c.Next()

		// If c.Next() returns an error, capture it with the request-specific hub.
		// Sentry's default panic recovery (likely in utils.SentryHandler) would also catch panics.
		if err != nil {
			hub.CaptureException(err)
			// The transaction status will be set by the defer block using the response code,
			// which Fiber sets appropriately even on errors.
		}

		return err // Return the error to Fiber
	}
}
