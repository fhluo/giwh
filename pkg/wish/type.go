package wish

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
)

type Type struct {
	Key  int    `json:"key,string" mapstructure:"key"`
	Name string `json:"name" mapstructure:"name"`
}

func (t Type) VarName() string {
	return regexp.MustCompile(`-+|\s+`).ReplaceAllString(cases.Title(language.English).String(t.Name), "")
}
