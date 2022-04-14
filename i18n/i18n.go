package i18n

import (
	_ "embed"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"golang.org/x/text/language"
	"log"
)

//go:embed data/items.json
var itemsData []byte

func init() {
	var result [][]Item
	err := jsoniter.Unmarshal(itemsData, &result)
	if err != nil {
		log.Fatalln(err)
	}

	items = make([]map[string]string, len(result))

	for i := range result {
		items[i] = make(map[string]string, len(result[i]))
		for _, item := range result[i] {
			items[i][item.Lang] = item.Name
		}
	}
}

var (
	Language string

	items         []map[string]string
	itemLanguages = []string{
		"en-us", "de-de", "es-es", "fr-fr", "id-id", "ja-jp", "ko-kr",
		"pt-pt", "ru-ru", "th-th", "vi-vn", "zh-cn", "zh-tw",
	}
	itemLangMatcher = language.NewMatcher(
		lo.Map(itemLanguages, func(lang string, _ int) language.Tag {
			return language.Make(lang)
		}),
	)
)

type Item struct {
	Name string `json:"name"`
	Lang string `json:"lang"`
}

func (item Item) GetName() string {
	for i := range items {
		for l, name := range items[i] {
			if l == item.Lang && name == item.Name {
				if r, ok := items[i][Language]; ok {
					return r
				} else {
					_, index := language.MatchStrings(itemLangMatcher, Language, l)
					return items[i][itemLanguages[index]]
				}
			}
		}
	}
	return item.Name
}
