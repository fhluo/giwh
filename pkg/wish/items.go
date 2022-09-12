package wish

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/i18n"
	"github.com/fhluo/giwh/pkg/util"
	"golang.org/x/exp/slices"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	*api.Item

	id       *int64
	rarity   *Rarity
	wishType *Type
	time     *time.Time
}

func (item Item) ID() int64 {
	if item.id == nil {
		i, err := strconv.ParseInt(item.Item.ID, 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		item.id = &i
	}

	return *item.id
}

func (item Item) Rarity() Rarity {
	if item.rarity == nil {
		i, err := strconv.Atoi(item.Item.Rarity)
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
		i, err := strconv.Atoi(item.Item.WishType)
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
		t, err := time.Parse("2006-01-02 15:04:05", item.Item.Time)
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
	return util.Unique(items, func(item Item) int64 {
		return item.ID()
	})
}

func (items Items) Count() int {
	return len(items)
}

func (items Items) FilterByUID(uid string) Items {
	return util.Filter(items, func(item Item) bool {
		return item.UID == uid
	})
}

func (items Items) FilterByWishType(wishTypes ...Type) Items {
	return util.Filter(items, func(item Item) bool {
		return slices.Contains(wishTypes, item.WishType())
	})
}

func (items Items) FilterByRarity(rarities ...Rarity) Items {
	return util.Filter(items, func(item Item) bool {
		return slices.Contains(rarities, item.Rarity())
	})
}

func (items Items) FilterByUIDAndWishType(uid string, wishTypes ...Type) Items {
	return util.Filter(items, func(item Item) bool {
		return item.UID == uid && slices.Contains(wishTypes, item.WishType())
	})
}

func (items Items) Save(filename string) error {
	sort.Sort(sort.Reverse(items))

	ext := strings.ToLower(filepath.Ext(filename))

	if e, ok := exporters[ext]; ok {
		return e.Export(items, filename)
	} else {
		return fmt.Errorf("format %s is not supported", ext)
	}
}
