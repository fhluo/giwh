package i18n

import (
	_ "embed"
	"github.com/fhluo/giwh/pkg/util"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/text/language"
	"log"
)

var (
	//go:embed data/items.json
	itemsData []byte
	//go:embed data/wishes.json
	wishesData []byte
	//go:embed data/shared_wishes.json
	sharedWishesData []byte

	Language string
)

var (
	items []map[string]string

	itemLanguages = []string{
		"en-us", "de-de", "es-es", "fr-fr", "id-id", "ja-jp", "ko-kr",
		"pt-pt", "ru-ru", "th-th", "vi-vn", "zh-cn", "zh-tw",
	}
	itemLangMatcher = language.NewMatcher(
		util.Map(itemLanguages, func(lang string) language.Tag {
			return language.Make(lang)
		}),
	)
)

var (
	wishes       map[string]map[string]string
	sharedWishes map[string]map[string]string
)

func init() {
	initItems()
	initWishes()
	initSharedWishes()
}

func initItems() {
	var result [][]Item

	if err := jsoniter.Unmarshal(itemsData, &result); err != nil {
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

func initWishes() {
	var result []Wish
	if err := jsoniter.Unmarshal(wishesData, &result); err != nil {
		log.Fatalln(err)
	}

	wishes = make(map[string]map[string]string)

	for _, wish := range result {
		if _, ok := wishes[wish.Lang]; !ok {
			wishes[wish.Lang] = make(map[string]string)
		}
		wishes[wish.Lang][wish.Type] = wish.Name
	}
}

func initSharedWishes() {
	var result []Wish
	if err := jsoniter.Unmarshal(sharedWishesData, &result); err != nil {
		log.Fatalln(err)
	}

	sharedWishes = make(map[string]map[string]string)

	for _, wish := range result {
		if _, ok := sharedWishes[wish.Lang]; !ok {
			sharedWishes[wish.Lang] = make(map[string]string)
		}
		sharedWishes[wish.Lang][wish.Type] = wish.Name
	}
}

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

func GetItemName(name string, lang string) string {
	return Item{Name: name, Lang: lang}.GetName()
}

type Wish struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Lang string `json:"lang"`
}

func GetWishName(wishType string) string {
	if _, ok := wishes[Language]; ok {
		return wishes[Language][wishType]
	} else {
		_, index := language.MatchStrings(itemLangMatcher, Language)
		return wishes[itemLanguages[index]][wishType]
	}
}

func GetSharedWishName(wishType string) string {
	if _, ok := sharedWishes[Language]; ok {
		return sharedWishes[Language][wishType]
	} else {
		_, index := language.MatchStrings(itemLangMatcher, Language)
		return sharedWishes[itemLanguages[index]][wishType]
	}
}
