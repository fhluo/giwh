package wh

import (
	"errors"
	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"io/fs"
	"log"
	"os"
	"sort"
	"strconv"
)

type RawItem struct {
	UID      string `json:"uid"`
	WishType string `json:"gacha_type"`
	ItemID   string `json:"item_id"`
	Count    string `json:"count"`
	Time     string `json:"time"`
	Name     string `json:"name"`
	Lang     string `json:"lang"`
	ItemType string `json:"item_type"`
	Rarity   string `json:"rank_type"`
	ID       string `json:"id"`
}

type Item struct {
	*RawItem

	id       *int64
	rarity   *Rarity
	wishType *WishType
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

func (item Item) String() string {
	return item.Name
}

func (item Item) ColoredString() string {
	switch item.Rarity() {
	case FourStar:
		return color.MagentaString(item.Name)
	case FiveStar:
		return color.YellowString(item.Name)
	default:
		return color.CyanString(item.Name)
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

	data, err := jsoniter.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0666)
}

func LoadItems(filename string) (Items, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var items Items
	return items, jsoniter.Unmarshal(data, &items)
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
