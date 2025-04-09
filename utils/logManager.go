package utils

import (
	"fmt"
	"time"
)

type LogLevel string

const (
	LevelInfo  LogLevel = "✅ INFO"
	LevelError LogLevel = "💥 ERROR"
	LevelWarn  LogLevel = "⚠️ WARN"
	LevelDebug LogLevel = "🔍 DEBUG"
)

func LogHeader(title string) {
	fmt.Printf("\n╔======== %s ========╗\n", title)
	fmt.Printf("║ 🕒 %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func LogMessage(level LogLevel, msg string) {
	fmt.Printf("║ %s: %s\n", level, msg)
}

func LogLineKeyValue(level LogLevel, key string, value any) {
	fmt.Printf("║ %s: %s: %v\n", level, key, value)
}

func LogFooter() {
	fmt.Println("╚=======================================╝")
}
