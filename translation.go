package main

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type TranslationService struct {
	client *translate.Client
}

func NewTranslationService() (*TranslationService, error) {
	ctx := context.Background()
	
	// Get API key from environment variable
	apiKey := os.Getenv("GOOGLE_TRANSLATE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_TRANSLATE_API_KEY environment variable is not set")
	}

	// Create client with API key
	client, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create translation client: %w", err)
	}

	return &TranslationService{
		client: client,
	}, nil
}

func (ts *TranslationService) TranslateText(text string, targetLangCode string) (string, error) {
	ctx := context.Background()

	targetLang, err := language.Parse(targetLangCode)
	if err != nil {
		return "", fmt.Errorf("failed to parse target language: %w", err)
	}

	translations, err := ts.client.Translate(ctx, []string{text}, targetLang, nil)
	if err != nil {
		return "", fmt.Errorf("failed to translate text: %w", err)
	}

	if len(translations) == 0 {
		return "", fmt.Errorf("no translation returned")
	}

	return translations[0].Text, nil
}

func (ts *TranslationService) TranslateMenu(menu *MenuData, targetLangCode string) (*MenuData, error) {
	translatedMenu := &MenuData{
		GrilladesMidi: make([]string, len(menu.GrilladesMidi)),
		Migrateurs:    make([]string, len(menu.Migrateurs)),
		Cibo:          make([]string, len(menu.Cibo)),
		AccompMidi:    make([]string, len(menu.AccompMidi)),
		GrilladesSoir: make([]string, len(menu.GrilladesSoir)),
		AccompSoir:    make([]string, len(menu.AccompSoir)),
	}

	// Helper function to translate a slice of strings
	translateSlice := func(source []string, target []string) error {
		for i, item := range source {
			translated, err := ts.TranslateText(item, targetLangCode)
			if err != nil {
				return fmt.Errorf("failed to translate menu item: %w", err)
			}
			target[i] = translated
		}
		return nil
	}

	// Translate each category
	if err := translateSlice(menu.GrilladesMidi, translatedMenu.GrilladesMidi); err != nil {
		return nil, err
	}
	if err := translateSlice(menu.Migrateurs, translatedMenu.Migrateurs); err != nil {
		return nil, err
	}
	if err := translateSlice(menu.Cibo, translatedMenu.Cibo); err != nil {
		return nil, err
	}
	if err := translateSlice(menu.AccompMidi, translatedMenu.AccompMidi); err != nil {
		return nil, err
	}
	if err := translateSlice(menu.GrilladesSoir, translatedMenu.GrilladesSoir); err != nil {
		return nil, err
	}
	if err := translateSlice(menu.AccompSoir, translatedMenu.AccompSoir); err != nil {
		return nil, err
	}

	return translatedMenu, nil
}
