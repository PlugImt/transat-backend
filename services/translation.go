package services

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/translate"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

// TranslationService provides an interface for text translation.
type TranslationService struct {
	client *translate.Client
	apiKey string
	ctx    context.Context
}

// NewTranslationService creates a new TranslationService.
func NewTranslationService(apiKey string) (*TranslationService, error) {
	ctx := context.Background()
	if apiKey == "" {
		utils.LogMessage(utils.LevelWarn, "GOOGLE_TRANSLATE_API_KEY is not set. Translation service will be disabled.")
		return &TranslationService{client: nil, apiKey: "", ctx: ctx}, nil
	}

	client, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Google Translate client: %w", err)
	}

	utils.LogMessage(utils.LevelInfo, "Google Translate client created successfully.")
	return &TranslationService{client: client, apiKey: apiKey, ctx: ctx}, nil
}

// TranslateText translates a single string to the target language.
func (ts *TranslationService) TranslateText(text string, targetLangCode string) (string, error) {
	if ts.client == nil {
		return text, nil // Service is disabled, return original text
	}

	targetLang, err := language.Parse(targetLangCode)
	if err != nil {
		log.Printf("Failed to parse target language '%s': %v", targetLangCode, err)
		return text, fmt.Errorf("failed to parse target language: %w", err)
	}

	if text == "" {
		return "", nil
	}

	translations, err := ts.client.Translate(ts.ctx, []string{text}, targetLang, nil)
	if err != nil {
		log.Printf("Failed to translate text '%s' to '%s': %v", text, targetLangCode, err)
		return text, fmt.Errorf("failed to translate text: %w", err)
	}

	if len(translations) == 0 {
		log.Printf("No translation returned for text '%s' to '%s'", text, targetLangCode)
		return text, fmt.Errorf("no translation returned")
	}

	return translations[0].Text, nil
}

// TranslateMenu translates all text fields within a MenuData struct.
func (ts *TranslationService) TranslateMenu(menu *models.MenuData, targetLangCode string) (*models.MenuData, error) {
	if ts.client == nil {
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

	translateSlice := func(source []string, target []string) error {
		if len(source) == 0 {
			return nil
		}

		targetLangParsed, err := language.Parse(targetLangCode)
		if err != nil {
			return fmt.Errorf("invalid target language for batch translation: %w", err)
		}

		translations, err := ts.client.Translate(ts.ctx, source, targetLangParsed, nil)
		if err != nil {
			log.Printf("Batch translation failed for %d items to '%s': %v. Falling back to individual translation.", len(source), targetLangCode, err)
			var firstErr error
			for i, item := range source {
				target[i], err = ts.TranslateText(item, targetLangCode)
				if err != nil && firstErr == nil {
					firstErr = err
				}
			}
			return firstErr
		}

		if len(translations) != len(source) {
			log.Printf("Warning: Batch translation returned %d results for %d inputs. Mismatch occurred.", len(translations), len(source))
			copyCount := len(source)
			if len(translations) < copyCount {
				copyCount = len(translations)
			}
			for i := 0; i < copyCount; i++ {
				target[i] = translations[i].Text
			}
			if len(source) > len(translations) {
				for i := len(translations); i < len(source); i++ {
					target[i] = source[i] // Fallback to original
				}
			}
			return fmt.Errorf("batch translation result count mismatch")
		}

		for i := 0; i < len(source); i++ {
			target[i] = translations[i].Text
		}
		return nil
	}

	log.Printf("Translating menu to %s...", targetLangCode)
	var firstErr error
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

	return translatedMenu, firstErr
}

// Close terminates the connection to the Google Cloud Translate API.
func (ts *TranslationService) Close() error {
	if ts.client != nil {
		log.Println("Closing Google Translate client.")
		return ts.client.Close()
	}
	return nil
}
