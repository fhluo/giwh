package i18n

import (
	"embed"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/samber/lo"
	"golang.org/x/sys/windows"
	"golang.org/x/text/language"
	"sync"
)

//go:generate go run ../cmd/giwh-dev gen lang
//go:generate go run ../cmd/giwh-dev gen locales

var (
	//go:embed locales/*.json
	localesFS embed.FS

	tagToLang = make(map[string]Language)
)

func init() {
	for _, lang := range Languages {
		tagToLang[lang.Tag().String()] = lang
	}
}

var (
	once    sync.Once
	locales = make(map[Language]Locale)
)

func Locales() map[Language]Locale {
	once.Do(func() {
		for _, lang := range Languages {
			locales[lang] = lo.Must(ReadLocale(lang))
		}
	})
	return locales
}

func ReadLocaleFile(lang Language) ([]byte, error) {
	return localesFS.ReadFile(fmt.Sprintf("locales/%s.json", lang.Tag().String()))
}

func ReadLocale(lang Language) (locale Locale, err error) {
	data, err := ReadLocaleFile(lang)
	if err != nil {
		return
	}
	err = sonic.Unmarshal(data, &locale)
	return
}

func Default() Language {
	languages, err := windows.GetUserPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
	if err != nil {
		languages, err = windows.GetSystemPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
		if err != nil {
			return English
		}
	}

	return Match(languages...)
}

func Match(languages ...string) Language {
	var tags []language.Tag

	for _, lang := range languages {
		tag, err := language.Parse(lang)
		if err != nil {
			return Default()
		}
		switch tag.String() {
		case "zh-CN":
			tags = append(tags, language.SimplifiedChinese)
		case "zh-TW":
			tags = append(tags, language.TraditionalChinese)
		default:
			tags = append(tags, tag)
		}
	}

	tag, _, _ := Matcher.Match(tags...)
	return tagToLang[tag.String()]
}
