package translator

import (
	"embed"
	"log"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type Translations map[string]LangTranslations

type LangTranslations map[string]string

type Translator struct {
	defaultLang  string
	translations Translations
}

type TranslatorOptions func(*Translator)

//go:embed *.toml
var translationsFS embed.FS

func NewTranslator(options ...TranslatorOptions) *Translator {
	translator := &Translator{
		defaultLang:  "en-en",
		translations: make(map[string]LangTranslations),
	}
	for _, option := range options {
		option(&Translator{})
	}
	files, _ := translationsFS.ReadDir(".")
	for _, file := range files {
		fileName := file.Name()
		lang := strings.Split(fileName, ".")[0]
		content, _ := translationsFS.ReadFile(fileName)
		var parsedContent map[string]string
		toml.Unmarshal(content, &parsedContent)
		translator.translations[lang] = parsedContent
	}
	return translator
}

func WithDefaultLang(lang string) func(*Translator) {
	return func(t *Translator) {
		t.defaultLang = lang
	}
}

/* Returns a translation for a given lang and key*/
func (t *Translator) T(lang string, key string) string {
	isLangEmpty := lang == ""
	isLangNotSupported := t.translations[lang] == nil
	if isLangEmpty || isLangNotSupported {
		lang = t.defaultLang
	}
	translation := t.translations[lang][key]
	if translation == "" {
		log.Printf("Translation not found for key: %s", key)
		return key
	}
	return translation
}

/* Returns all the translations for a given lang*/
func (t *Translator) LangT(lang string) map[string]string {
	isLangEmpty := lang == ""
	isLangNotSupported := t.translations[lang] == nil
	if isLangEmpty || isLangNotSupported {
		lang = t.defaultLang
	}
	return t.translations[lang]
}
