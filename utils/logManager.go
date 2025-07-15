package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/getsentry/sentry-go"
)

type LogLevel string

const (
	LevelInfo  LogLevel = "✅ INFO"
	LevelError LogLevel = "💥 ERROR"
	LevelWarn  LogLevel = "⚠️ WARN"
	LevelDebug LogLevel = "🔍 DEBUG"
)

var logger sentry.Logger

func init() {
	InitLogger()
}

func InitLogger() {
	if logger != nil {
		return // Already initialized
	}
	ctx := context.Background()
	logger = sentry.NewLogger(ctx)
}

func _log(level LogLevel, msg string) {
	ctx := context.Background()

	if level == LevelDebug && os.Getenv("ENV") == "production" {
		return
	}
	log.Printf("║ %s: %s\n", level, msg)

	// Safety check: only use sentry logger if it's initialized
	if logger != nil {
		switch level {
		case LevelInfo:
			logger.Info(ctx, msg)
		case LevelError:
			logger.Error(ctx, msg)
		case LevelWarn:
			logger.Warn(ctx, msg)
		case LevelDebug:
			logger.Debug(ctx, msg)
		}
	}
}

func LogHeader(title string) {
	log.Printf("╔======== %s ========╗\n", title)
}

func LogMessage(level LogLevel, msg string) {
	_log(level, msg)
}

func LogLineKeyValue(level LogLevel, key string, value any) {
	_log(level, fmt.Sprintf("%s: %v", key, value))
}

func LogFooter() {
	log.Println("╚=======================================╝")
}
