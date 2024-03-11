package i18n

import (
	"github.com/fhluo/giwh/pkg/wish"
	"github.com/goccy/go-json"
	"github.com/samber/lo"
)

type Locale struct {
	Language          Language             `json:"language"`
	Characters        map[string]string    `json:"characters"`
	CharactersInverse map[string]string    `json:"charactersInverse"`
	Weapons           map[string]string    `json:"weapons"`
	WeaponsInverse    map[string]string    `json:"weaponsInverse"`
	Wishes            map[wish.Type]string `json:"wishes"`
	SharedWishes      map[wish.Type]string `json:"sharedWishes"`
}

func NewLocale(lang Language) Locale {
	return Locale{
		Language:          lang,
		Characters:        make(map[string]string),
		CharactersInverse: make(map[string]string),
		Weapons:           make(map[string]string),
		WeaponsInverse:    make(map[string]string),
		Wishes:            make(map[wish.Type]string),
		SharedWishes:      make(map[wish.Type]string),
	}
}

func (l Locale) BaseFilename() string {
	return l.Language.Tag().String() + ".json"
}

func (l Locale) JSON() []byte {
	return lo.Must(json.MarshalIndent(l, "", "  "))
}
