package locales

import (
	"github.com/fhluo/giwh/i18n"
	"github.com/goccy/go-json"
	"github.com/samber/lo"
)

type Locale struct {
	Language   i18n.Language     `json:"language"`
	Characters map[string]string `json:"characters"`
	Weapons    map[string]string `json:"weapons"`
}

func NewLocale(lang i18n.Language) Locale {
	return Locale{
		Language:   lang,
		Characters: make(map[string]string),
		Weapons:    make(map[string]string),
	}
}

func (l Locale) BaseFilename() string {
	return l.Language.Tag().String() + ".json"
}

func (l Locale) JSON() []byte {
	return lo.Must(json.MarshalIndent(l, "", "  "))
}
