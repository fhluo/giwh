package wh

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/fhluo/giwh/i18n"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type RawItem struct {
	UID      string `json:"uid" toml:"uid"`
	WishType string `json:"gacha_type" toml:"gacha_type"`
	ItemID   string `json:"item_id" toml:"item_id"`
	Count    string `json:"count" toml:"count"`
	Time     string `json:"time" toml:"time"`
	Name     string `json:"name" toml:"name"`
	Lang     string `json:"lang" toml:"lang"`
	ItemType string `json:"item_type" toml:"item_type"`
	Rarity   string `json:"rank_type" toml:"rank_type"`
	ID       string `json:"id" toml:"id"`
}

type Item struct {
	*RawItem

	id       *int64
	rarity   *Rarity
	wishType *WishType
	time     *time.Time
}

func (item Item) ID() int64 {
	if item.id == nil {
		i, err := strconv.ParseInt(item.RawItem.ID, 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		item.id = &i
	}

	return *item.id
}

func (item Item) Rarity() Rarity {
	if item.rarity == nil {
		i, err := strconv.Atoi(item.RawItem.Rarity)
		if err != nil {
			log.Fatalln(err)
		}
		r := Rarity(i)
		item.rarity = &r
	}

	return *item.rarity
}

func (item Item) WishType() WishType {
	if item.wishType == nil {
		i, err := strconv.Atoi(item.RawItem.WishType)
		if err != nil {
			log.Fatalln(err)
		}
		t := WishType(i)
		item.wishType = &t
	}

	return *item.wishType
}

func (item Item) Time() time.Time {
	if item.time == nil {
		t, err := time.Parse("2006-01-02 15:04:05", item.RawItem.Time)
		if err != nil {
			log.Fatalln(err)
		}
		item.time = &t
	}

	return *item.time
}

func (item Item) String() string {
	if i18n.Language == "" {
		return item.Name
	} else {
		return i18n.Item{Name: item.Name, Lang: item.Lang}.GetName()
	}
}

func (item Item) ColoredString() string {
	switch item.Rarity() {
	case FourStar:
		return color.MagentaString(item.String())
	case FiveStar:
		return color.YellowString(item.String())
	default:
		return color.CyanString(item.String())
	}
}

type Items []Item

func (items Items) Len() int {
	return len(items)
}

func (items Items) Less(i, j int) bool {
	return items[i].ID() < items[j].ID()
}

func (items Items) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (items Items) Equal(items2 Items) bool {
	return slices.EqualFunc(items, items2, func(item1, item2 Item) bool {
		return item1.ID() == item2.ID()
	})
}

func (items Items) Unique() Items {
	return lo.UniqBy(items, func(item Item) int64 {
		return item.ID()
	})
}

func (items Items) Count() int {
	return len(items)
}

func (items Items) FilterByUID(uid string) Items {
	return lo.Filter(items, func(item Item, _ int) bool {
		return item.UID == uid
	})
}

func (items Items) FilterByWishType(wishTypes ...WishType) Items {
	return lo.Filter(items, func(item Item, _ int) bool {
		return lo.Contains(wishTypes, item.WishType())
	})
}

func (items Items) FilterByRarity(rarities ...Rarity) Items {
	return lo.Filter(items, func(item Item, _ int) bool {
		return lo.Contains(rarities, item.Rarity())
	})
}

func (items Items) FilterByUIDAndWishType(uid string, wishTypes ...WishType) Items {
	return lo.Filter(items, func(item Item, _ int) bool {
		return item.UID == uid && lo.Contains(wishTypes, item.WishType())
	})
}

func (items Items) Save(filename string) error {
	sort.Sort(sort.Reverse(items))

	var (
		data []byte
		err  error
	)

	switch filepath.Ext(filename) {
	case ".json":
		data, err = jsoniter.MarshalIndent(items, "", "  ")
		if err != nil {
			return err
		}
	case ".toml":
		buf := new(bytes.Buffer)
		e := toml.NewEncoder(buf)
		e.Indent = ""

		err = e.Encode(map[string]interface{}{"list": items})
		data = buf.Bytes()

		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("format %s is not supported", filepath.Ext(filename))
	}

	return os.WriteFile(filename, data, 0666)
}

func LoadItems(filename string) (Items, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	switch filepath.Ext(filename) {
	case ".json":
		var items Items
		return items, jsoniter.Unmarshal(data, &items)
	case ".toml":
		var result struct {
			List Items `toml:"list"`
		}
		return result.List, toml.Unmarshal(data, &result)
	default:
		return nil, fmt.Errorf("format %s is not supported", filepath.Ext(filename))
	}
}

func LoadItemsIfExits(filename string) (Items, error) {
	_, err := os.Stat(filename)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		return nil, nil
	}

	return LoadItems(filename)
}
