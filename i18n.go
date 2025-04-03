package main

import (
	"embed"
	"log"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locale/*.toml
var LocaleFS embed.FS

var bundle *i18n.Bundle

func initI18n() error {
	bundle = i18n.NewBundle(language.French)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Charger tous les fichiers de traduction
	entries, err := LocaleFS.ReadDir("locale")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".toml") {
			path := filepath.Join("locale", entry.Name())
			_, err := bundle.LoadMessageFileFS(LocaleFS, path)
			if err != nil {
				log.Printf("Warning: Failed to load translation file %s: %v", path, err)
				continue
			}
			log.Printf("Loaded translation file: %s", path)
		}
	}

	return nil
}

func GetLocalizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, lang)
}
