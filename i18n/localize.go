package i18n

import "github.com/nicksnyder/go-i18n/v2/i18n"

// Localize returns the message for id, or fallback when the language has no translation.
// data fills the message's {{.Template}} placeholders.
func Localize(localizer *i18n.Localizer, id, fallback string, data map[string]interface{}) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:      id,
		TemplateData:   data,
		DefaultMessage: &i18n.Message{ID: id, Other: fallback},
	})
}
