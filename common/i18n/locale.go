package i18n

import (
	"encoding/json"
	"github.com/samber/lo"
)

type Locale struct {
	Language          Language          `json:"language"`
	Characters        map[string]string `json:"characters"`
	CharactersInverse map[string]string `json:"charactersInverse"`
	Weapons           map[string]string `json:"weapons"`
	WeaponsInverse    map[string]string `json:"weaponsInverse"`
	Wishes            map[string]string `json:"wishes"`
	SharedWishes      map[string]string `json:"sharedWishes"`
}

func NewLocale(lang Language) Locale {
	return Locale{
		Language:          lang,
		Characters:        make(map[string]string),
		CharactersInverse: make(map[string]string),
		Weapons:           make(map[string]string),
		WeaponsInverse:    make(map[string]string),
		Wishes:            make(map[string]string),
		SharedWishes:      make(map[string]string),
	}
}

func (l Locale) BaseFilename() string {
	return l.Language.Tag().String() + ".json"
}

func (l Locale) JSON() []byte {
	return lo.Must(json.MarshalIndent(l, "", "  "))
}
