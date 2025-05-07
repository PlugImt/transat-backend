package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/translate"
	"github.com/plugimt/transat-backend/models"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

// TranslationService handles text translation using Google Cloud Translate API.
type TranslationService struct {
	client *translate.Client
	ctx    context.Context
}

// NewTranslationService creates a new TranslationService instance.
// It requires the GOOGLE_TRANSLATE_API_KEY environment variable to be set.
func NewTranslationService() (*TranslationService, error) {
	ctx := context.Background()

	// Get API key from environment variable
	apiKey := os.Getenv("GOOGLE_TRANSLATE_API_KEY")
	if apiKey == "" {
		log.Println("Warning: GOOGLE_TRANSLATE_API_KEY environment variable is not set. Translation service will be disabled.")
		// Return a service instance with a nil client to indicate disabled state
		return &TranslationService{client: nil, ctx: ctx}, nil // Not returning error, service is just disabled
	}

	// Create client with API key
	client, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("Failed to create Google Translate client: %v", err)
		return nil, fmt.Errorf("failed to create translation client: %w", err)
	}

	log.Println("Google Translate client created successfully.")
	return &TranslationService{
		client: client,
		ctx:    ctx,
	}, nil
}

// TranslateText translates a single string to the target language.
func (ts *TranslationService) TranslateText(text string, targetLangCode string) (string, error) {
	if ts.client == nil {
		// Service is disabled, return original text
		return text, nil
	}

	targetLang, err := language.Parse(targetLangCode)
	if err != nil {
		log.Printf("Failed to parse target language '%s': %v", targetLangCode, err)
		return text, fmt.Errorf("failed to parse target language: %w", err) // Return original text and error
	}

	// Handle empty input text
	if text == "" {
		return "", nil
	}

	translations, err := ts.client.Translate(ts.ctx, []string{text}, targetLang, nil)
	if err != nil {
		log.Printf("Failed to translate text '%s' to '%s': %v", text, targetLangCode, err)
		// Return original text on translation failure
		return text, fmt.Errorf("failed to translate text: %w", err)
	}

	if len(translations) == 0 {
		log.Printf("No translation returned for text '%s' to '%s'", text, targetLangCode)
		// Return original text if no translation is provided
		return text, fmt.Errorf("no translation returned")
	}

	// log.Printf("Translated '%s' to '%s' -> '%s'", text, targetLangCode, translations[0].Text)
	return translations[0].Text, nil
}

// TranslateMenu translates all text fields within a MenuData struct.
func (ts *TranslationService) TranslateMenu(menu *models.MenuData, targetLangCode string) (*models.MenuData, error) {
	if ts.client == nil {
		// Service is disabled, return original menu
		log.Println("Translation service is disabled, returning original menu.")
		return menu, nil
	}

	if targetLangCode == "fr" { // Assuming 'fr' is the source language
		log.Println("Target language is French (source), skipping translation.")
		return menu, nil
	}

	translatedMenu := &models.MenuData{
		GrilladesMidi: make([]string, len(menu.GrilladesMidi)),
		Migrateurs:    make([]string, len(menu.Migrateurs)),
		Cibo:          make([]string, len(menu.Cibo)),
		AccompMidi:    make([]string, len(menu.AccompMidi)),
		GrilladesSoir: make([]string, len(menu.GrilladesSoir)),
		AccompSoir:    make([]string, len(menu.AccompSoir)),
	}

	// Helper function to translate a slice of strings
	translateSlice := func(source []string, target []string) error {

		// If source is empty, nothing to do
		if len(source) == 0 {
			return nil
		}

		// Optimization: Translate the whole slice in one API call if possible
		// Note: Google Translate API v2 (which cloud.google.com/go/translate uses)
		// supports translating multiple texts in a single request.
		targetLangParsed, err := language.Parse(targetLangCode)
		if err != nil {
			return fmt.Errorf("invalid target language for batch translation: %w", err)
		}

		translations, err := ts.client.Translate(ts.ctx, source, targetLangParsed, nil)
		if err != nil {
			log.Printf("Batch translation failed for %d items to '%s': %v. Falling back to individual translation.", len(source), targetLangCode, err)
			// Fallback to individual translation if batch fails
			var firstErr error
			for i, item := range source {
				target[i], err = ts.TranslateText(item, targetLangCode)
				if err != nil && firstErr == nil {
					firstErr = err // Keep the first error encountered
				}
			}
			return firstErr
		}

		if len(translations) != len(source) {
			log.Printf("Warning: Batch translation returned %d results for %d inputs. Mismatch occurred.", len(translations), len(source))
			// Fallback or handle mismatch? For now, copy what we got.
			copyCount := len(source)
			if len(translations) < copyCount {
				copyCount = len(translations)
			}
			for i := 0; i < copyCount; i++ {
				target[i] = translations[i].Text
			}
			// Fill remaining with original text?
			if len(source) > len(translations) {
				for i := len(translations); i < len(source); i++ {
					target[i] = source[i] // Fallback to original
				}
			}
			return fmt.Errorf("batch translation result count mismatch")
		}

		// Copy results from batch translation
		for i := 0; i < len(source); i++ {
			target[i] = translations[i].Text
		}
		return nil
	}

	// Translate each category using the optimized helper
	log.Printf("Translating menu to %s...", targetLangCode)
	var firstErr error // Track the first error across all slices
	if err := translateSlice(menu.GrilladesMidi, translatedMenu.GrilladesMidi); err != nil {
		log.Printf("Error translating GrilladesMidi: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}
	if err := translateSlice(menu.Migrateurs, translatedMenu.Migrateurs); err != nil {
		log.Printf("Error translating Migrateurs: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}
	if err := translateSlice(menu.Cibo, translatedMenu.Cibo); err != nil {
		log.Printf("Error translating Cibo: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}
	if err := translateSlice(menu.AccompMidi, translatedMenu.AccompMidi); err != nil {
		log.Printf("Error translating AccompMidi: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}
	if err := translateSlice(menu.GrilladesSoir, translatedMenu.GrilladesSoir); err != nil {
		log.Printf("Error translating GrilladesSoir: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}
	if err := translateSlice(menu.AccompSoir, translatedMenu.AccompSoir); err != nil {
		log.Printf("Error translating AccompSoir: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}
	log.Printf("Finished translating menu to %s.", targetLangCode)

	// Return the translated menu. If errors occurred, some items might contain original text.
	// Return the first error encountered during the process.
	return translatedMenu, firstErr
}

// Close terminates the connection to the Google Cloud Translate API.
// Call this when the application is shutting down gracefully.
func (ts *TranslationService) Close() error {
	if ts.client != nil {
		log.Println("Closing Google Translate client.")
		return ts.client.Close()
	}
	return nil
}
