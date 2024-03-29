package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"regexp"
)

type Language struct {
	Key   string `mapstructure:"key" json:"key"`
	Name  string `mapstructure:"name" json:"name"`
	Short string `mapstructure:"short" json:"short"`
}

func (lang Language) MustParse() language.Tag {
	return language.MustParse(lang.Key)
}

func (lang Language) CanonicalKey() string {
	return language.MustParse(lang.Key).String()
}

func (lang Language) Tag() language.Tag {
	switch lang.CanonicalKey() {
	case "zh-CN":
		return language.SimplifiedChinese
	case "zh-TW":
		return language.TraditionalChinese
	default:
		return lang.MustParse()
	}
}

func (lang Language) Parent() language.Tag {
	switch lang.CanonicalKey() {
	case "zh-CN":
		return language.SimplifiedChinese
	case "zh-TW":
		return language.TraditionalChinese
	default:
		return lang.MustParse().Parent()
	}
}

func (lang Language) VarName() string {
	return regexp.MustCompile(`\s`).ReplaceAllString(display.English.Languages().Name(lang.Tag()), "")
}

func (lang Language) ParentVarName() string {
	return regexp.MustCompile(`\s`).ReplaceAllString(display.English.Languages().Name(lang.Parent()), "")
}
