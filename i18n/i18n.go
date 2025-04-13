package i18n

import (
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed active.*.toml
var localesFS embed.FS

var bundle *i18n.Bundle

// Init initializes the i18n bundle by loading translation files from the embedded FS.
func Init() error {
	bundle = i18n.NewBundle(language.French) // Default language
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Search for .toml files in the embedded FS
	matches, err := fs.Glob(localesFS, "active.*.toml")
	if err != nil {
		// An error here likely means a problem with the Glob pattern or the embed.FS
		return fmt.Errorf("failed to glob embedded locale files: %w", err)
	}

	// Check if any files were found before attempting to load them
	if len(matches) == 0 {
		log.Println("Warning: No embedded locale files found matching 'active.*.toml'.")
		// Decide if this is a fatal error or not. For now, we continue.
		return nil
	}

	loaded := false
	for _, match := range matches {
		// Load the message file from the embedded FS
		_, err := bundle.LoadMessageFileFS(localesFS, match)
		if err != nil {
			log.Printf("Warning: Failed to load embedded message file %s: %v", match, err)
		} else {
			log.Printf("Loaded embedded message file: %s", match)
			loaded = true
		}
	}

	if !loaded {
		log.Println("Warning: No translation files were successfully loaded from the embedded FS.")
		// Depends on requirements, this could be a fatal error.
	}

	log.Println("Embedded i18n bundle initialization process completed.")
	return nil
}

// GetLocalizer returns a localizer for the specified language code.
// Falls back to the default language (French) if the requested language is not available.
func GetLocalizer(langCode string) *i18n.Localizer {
	if bundle == nil {
		log.Println("Error: i18n bundle is not initialized. Attempting to initialize now.")
		// Attempt to initialize now as a fallback
		if err := Init(); err != nil { // Changed InitI18n to Init
			log.Printf("Fallback i18n initialization failed: %v. Returning default localizer.", err)
			// Return a minimal default localizer if fallback init fails
			return i18n.NewLocalizer(i18n.NewBundle(language.French), language.French.String())
		}
		log.Println("Fallback i18n initialization successful.")
	}
	// Try to find a localizer for the requested language, fallback to French
	return i18n.NewLocalizer(bundle, langCode, language.French.String())
}
