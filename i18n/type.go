package i18n

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
)

type WishType struct {
	Key  int    `json:"key,string" mapstructure:"key"`
	Name string `json:"name" mapstructure:"name"`
}

func (t WishType) VarName() string {
	return regexp.MustCompile(`-+|\s+`).ReplaceAllString(cases.Title(language.English).String(t.Name), "")
}
