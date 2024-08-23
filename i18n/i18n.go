package i18n

import (
	"embed"
	"fmt"
	"log/slog"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var LocaleFS embed.FS

type (
	I18nBundle struct {
		bundle     *i18n.Bundle
		localizers map[string]*i18n.Localizer
	}
)

func New() *I18nBundle {
	localizers := make(map[string]*i18n.Localizer)
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	supportLangs := []string{"en", "ja"}
	for _, lang := range supportLangs {
		_, err := bundle.LoadMessageFileFS(LocaleFS, fmt.Sprintf("locales/%s.toml", lang))
		if err != nil {
			slog.Warn("[i18n] failed to load locale file", "error", err, "lang", lang)
			continue
		}
		localizers[lang] = i18n.NewLocalizer(bundle, lang, lang)
	}
	return &I18nBundle{
		bundle:     bundle,
		localizers: localizers,
	}
}

func (i *I18nBundle) Localizer(lang string) *i18n.Localizer {
	fmt.Printf("lang: %v\n", lang)
	if l, ok := i.localizers[lang]; ok {
		return l
	}
	return i.localizers["en"] // fallback to English
}

func (i *I18nBundle) T(lang, msgID string, data any) (string, error) {
	localizer := i.Localizer(lang)
	output, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    msgID,
		TemplateData: data,
	})
	if err != nil {
		return "", err
	}
	return output, nil
}

func (i *I18nBundle) MusT(lang, msgID string, data any) string {
	output, _ := i.T(lang, msgID, data)
	return output
}
