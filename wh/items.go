package wh

import (
	jsoniter "github.com/json-iterator/go"
	"log"
	"os"
	"sort"
	"strconv"
)

const (
	NoviceWish         = 100
	PermanentWish      = 200
	CharacterEventWish = 301
	WeaponEventWish    = 302
)

var Wishes = []int{NoviceWish, PermanentWish, CharacterEventWish, WeaponEventWish}

type RawItem struct {
	UID      string `json:"uid"`
	WishType string `json:"gacha_type"`
	ItemID   string `json:"item_id"`
	Count    string `json:"count"`
	Time     string `json:"time"`
	Name     string `json:"name"`
	Lang     string `json:"lang"`
	ItemType string `json:"item_type"`
	RankType string `json:"rank_type"`
	ID       string `json:"id"`
}

type Item struct {
	*RawItem

	id *int64
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

func (item Item) String() string {
	return item.Name
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

func (items Items) Save(filename string) error {
	sort.Sort(sort.Reverse(items))

	data, err := jsoniter.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0666)
}
