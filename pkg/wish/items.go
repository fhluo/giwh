package wish

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/fhluo/giwh/pkg/i18n"
	"github.com/fhluo/giwh/pkg/util"
	"golang.org/x/exp/slices"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type RawItem struct {
	UID      string `json:"uid" toml:"uid" csv:"UID"`
	WishType string `json:"gacha_type" toml:"gacha_type" csv:"Wish Type"`
	ItemID   string `json:"item_id" toml:"item_id" csv:"Item ID"`
	Count    string `json:"count" toml:"count" csv:"Count"`
	Time     string `json:"time" toml:"time" csv:"Time"`
	Name     string `json:"name" toml:"name" csv:"Name"`
	Lang     string `json:"lang" toml:"lang" csv:"Language"`
	ItemType string `json:"item_type" toml:"item_type" csv:"Item Type"`
	Rarity   string `json:"rank_type" toml:"rank_type" csv:"Rarity"`
	ID       string `json:"id" toml:"id" csv:"ID"`
}

func (r *RawItem) ToCSVHeader() []string {
	t := reflect.Indirect(reflect.ValueOf(r)).Type()
	result := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		result[i] = t.Field(i).Tag.Get("csv")
	}
	return result
}

func (r *RawItem) ToCSVRecord() []string {
	v := reflect.Indirect(reflect.ValueOf(r))
	result := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		result[i] = v.Field(i).String()
	}
	return result
}

type Item struct {
	*RawItem

	id       *int64
	rarity   *Rarity
	wishType *Type
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

func (item Item) WishType() Type {
	if item.wishType == nil {
		i, err := strconv.Atoi(item.RawItem.WishType)
		if err != nil {
			log.Fatalln(err)
		}
		t := Type(i)
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
		return i18n.GetItemName(item.Name, item.Lang)
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

type History []Item

func (h History) Len() int {
	return len(h)
}

func (h History) Less(i, j int) bool {
	return h[i].ID() < h[j].ID()
}

func (h History) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h History) Equal(items2 History) bool {
	return slices.EqualFunc(h, items2, func(item1, item2 Item) bool {
		return item1.ID() == item2.ID()
	})
}

func (h History) Unique() History {
	return util.Unique(h, func(item Item) int64 {
		return item.ID()
	})
}

func (h History) Count() int {
	return len(h)
}

func (h History) FilterByUID(uid string) History {
	return util.Filter(h, func(item Item) bool {
		return item.UID == uid
	})
}

func (h History) FilterByWishType(wishTypes ...Type) History {
	return util.Filter(h, func(item Item) bool {
		return slices.Contains(wishTypes, item.WishType())
	})
}

func (h History) FilterByRarity(rarities ...Rarity) History {
	return util.Filter(h, func(item Item) bool {
		return slices.Contains(rarities, item.Rarity())
	})
}

func (h History) FilterByUIDAndWishType(uid string, wishTypes ...Type) History {
	return util.Filter(h, func(item Item) bool {
		return item.UID == uid && slices.Contains(wishTypes, item.WishType())
	})
}

func (h History) ToCSVRecords() [][]string {
	if len(h) == 0 {
		return nil
	}

	items := make([][]string, len(h))
	for i := range h {
		items[i] = h[i].ToCSVRecord()
	}
	return items
}

func (h History) Save(filename string) error {
	sort.Sort(sort.Reverse(h))

	ext := strings.ToLower(filepath.Ext(filename))

	if e, ok := exporters[ext]; ok {
		return e.Export(h, filename)
	} else {
		return fmt.Errorf("format %s is not supported", ext)
	}
}

func LoadWishHistory(filename string) (History, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	if i, ok := importers[ext]; ok {
		return i.Import(filename)
	} else {
		return nil, fmt.Errorf("format %s is not supported", ext)
	}
}

func LoadWishHistoryIfExits(filename string) (History, error) {
	_, err := os.Stat(filename)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		return nil, nil
	}

	return LoadWishHistory(filename)
}
