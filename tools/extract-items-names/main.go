package main

import (
	"flag"
	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	CharacterDataPath = `ExcelBinOutput/AvatarExcelConfigData.json`
	WeaponDataPath    = `ExcelBinOutput/WeaponExcelConfigData.json`
)

var Languages = map[string]string{
	"TextMapEN.json":  "en-us",
	"TextMapFR.json":  "fr-fr",
	"TextMapDE.json":  "de-de",
	"TextMapES.json":  "es-es",
	"TextMapPT.json":  "pt-pt",
	"TextMapRU.json":  "ru-ru",
	"TextMapJP.json":  "ja-jp",
	"TextMapKR.json":  "ko-kr",
	"TextMapTH.json":  "th-th",
	"TextMapVI.json":  "vi-vn",
	"TextMapID.json":  "id-id",
	"TextMapCHT.json": "zh-tw",
	"TextMapCHS.json": "zh-cn",
}

type Character struct {
	ID              int    `json:"Id"`
	NameTextMapHash int64  `json:"NameTextMapHash"`
	DescTextMapHash int64  `json:"DescTextMapHash"`
	UseType         string `json:"UseType,omitempty"`
}

type Weapon struct {
	ID              int   `json:"id"`
	NameTextMapHash int64 `json:"nameTextMapHash"`
	DescTextMapHash int64 `json:"descTextMapHash"`
}

func LoadSliceOf[T any](filename string) ([]T, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var items []T
	if err = json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	return items, nil
}

type Item struct {
	Hash int64  `json:"-"`
	Name string `json:"name"`
	Lang string `json:"lang"`
}

func LoadTextMap(filename string) (map[int64]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var textMap map[int64]string
	return textMap, json.Unmarshal(data, &textMap)
}

var src, dst string

func init() {
	flag.StringVar(&src, "i", "", "")
	flag.StringVar(&dst, "o", "", "")
}

func GetNamesHashes() ([]int64, error) {
	characters, err := LoadSliceOf[*Character](filepath.Join(src, CharacterDataPath))
	if err != nil {
		return nil, err
	}
	characters = lo.Filter(characters, func(character *Character, _ int) bool {
		return character.UseType == "AVATAR_FORMAL"
	})

	weapons, err := LoadSliceOf[*Weapon](filepath.Join(src, WeaponDataPath))
	if err != nil {
		return nil, err
	}

	hashes := make([]int64, 0, len(characters)+len(weapons))

	hashes = append(hashes, lo.Map(characters, func(character *Character, _ int) int64 {
		return character.NameTextMapHash
	})...)

	hashes = append(hashes, lo.Map(weapons, func(weapon *Weapon, _ int) int64 {
		return weapon.NameTextMapHash
	})...)

	return hashes, nil
}

func main() {
	flag.Parse()

	hashes, err := GetNamesHashes()
	if err != nil {
		log.Fatalln(err)
	}

	filenames, err := filepath.Glob(filepath.Join(src, "TextMap/TextMap*.json"))
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(filenames))

	var items []*Item
	var mutex sync.Mutex

	for _, filename := range filenames {
		go func(filename string) {
			defer wg.Done()

			textMap, err := LoadTextMap(filename)
			if err != nil {
				log.Println(err)
			}

			r := lo.Map(hashes, func(hash int64, _ int) *Item {
				return &Item{
					Hash: hash,
					Name: textMap[hash],
					Lang: Languages[filepath.Base(filename)],
				}
			})

			mutex.Lock()
			items = append(items, r...)
			mutex.Unlock()
		}(filename)
	}

	wg.Wait()

	result := maps.Values(lo.GroupBy(items, func(item *Item) int64 {
		return item.Hash
	}))

	result = lo.Filter(result, func(items []*Item, _ int) bool {
		for _, item := range items {
			if item.Name == "" {
				return false
			}
		}
		return true
	})

	for i := range result {
		slices.SortFunc(result[i], func(a *Item, b *Item) bool {
			return a.Lang < b.Lang
		})
	}

	index := make(map[int64]int)
	for i, hash := range hashes {
		index[hash] = i
	}

	slices.SortFunc(result, func(a []*Item, b []*Item) bool {
		return index[a[0].Hash] < index[b[0].Hash]
	})

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(dst, data, 0666)
	if err != nil {
		log.Fatalln(err)
	}
}
