package utils

import (
	"log"
)

type LogLevel string

const (
	LevelInfo  LogLevel = "✅ INFO"
	LevelError LogLevel = "💥 ERROR"
	LevelWarn  LogLevel = "⚠️ WARN"
	LevelDebug LogLevel = "🔍 DEBUG"
)

func LogHeader(title string) {
	log.Printf("╔======== %s ========╗\n", title)
}

func LogMessage(level LogLevel, msg string) {
	log.Printf("║ %s: %s\n", level, msg)
}

func LogLineKeyValue(level LogLevel, key string, value any) {
	log.Printf("║ %s: %s: %v\n", level, key, value)
}

func LogFooter() {
	log.Println("╚=======================================╝")
}
