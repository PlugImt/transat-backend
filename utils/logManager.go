package utils

import (
	"fmt"
	"log"
	"os"
)

type LogLevel string

const (
	LevelInfo  LogLevel = "✅ INFO"
	LevelError LogLevel = "💥 ERROR"
	LevelWarn  LogLevel = "⚠️ WARN"
	LevelDebug LogLevel = "🔍 DEBUG"
)

func _log(level LogLevel, msg string) {
	if level == LevelDebug && os.Getenv("ENV") == "production" {
		return
	}
	log.Printf("║ %s: %s\n", level, msg)
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
