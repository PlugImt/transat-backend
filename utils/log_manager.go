package utils

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/getsentry/sentry-go"
)

type LogLevel string

const (
	LevelSilent LogLevel = "SILENT"
	LevelInfo   LogLevel = "✅ INFO"
	LevelError  LogLevel = "💥 ERROR"
	LevelWarn   LogLevel = "⚠️ WARN"
	LevelDebug  LogLevel = "🔍 DEBUG"
)

var logger sentry.Logger
var currentLogLevel LogLevel = LevelInfo // Default log level, will be configured by InitLogger
var isProductionEnv bool = false         // Default, will be configured by InitLogger

// InitLogger initializes the logger and sets the logging level and environment.
func InitLogger(env string, configuredLogLevel string) {
	ctx := context.Background()
	logger = sentry.NewLogger(ctx)

	isProductionEnv = strings.ToLower(env) == "production"

	empLogLevel := LevelInfo // Default to INFO
	switch strings.ToUpper(configuredLogLevel) {
	case "DEBUG":
		empLogLevel = LevelDebug
	case "INFO":
		empLogLevel = LevelInfo
	case "WARN":
		empLogLevel = LevelWarn
	case "ERROR":
		empLogLevel = LevelError
	case "SILENT":
		empLogLevel = LevelSilent
	default:
		if configuredLogLevel != "" && strings.ToUpper(configuredLogLevel) != "DEFAULT" { // Log warning only if a level was provided but not recognized, and not "DEFAULT"
			log.Printf("║ ⚠️ WARN: Invalid LOG_LEVEL '%s'. Defaulting to INFO.\n", configuredLogLevel)
		}
	}
	currentLogLevel = empLogLevel
	// Ensure "SILENT" is respected for the init message too
	if currentLogLevel != LevelSilent {
		log.Printf("║ ✅ INFO: Logger initialized. Environment: %s, Effective LogLevel: %s\n", env, currentLogLevel)
	}
}

func _log(level LogLevel, msg string) {
	if currentLogLevel == LevelSilent {
		return
	}

	// Determine if the message should be logged based on currentLogLevel
	shouldLog := false
	switch currentLogLevel {
	case LevelDebug:
		shouldLog = true // Log everything if current level is DEBUG
	case LevelInfo:
		shouldLog = (level == LevelInfo || level == LevelWarn || level == LevelError || level == LevelDebug) // INFO includes DEBUG in non-prod
		if isProductionEnv && level == LevelDebug {
			shouldLog = false
		} // Unless it is debug in prod
	case LevelWarn:
		shouldLog = (level == LevelWarn || level == LevelError)
	case LevelError:
		shouldLog = (level == LevelError)
	}

	if !shouldLog {
		return
	}

	log.Printf("║ %s: %s\n", level, msg)

	ctx := context.Background()
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
