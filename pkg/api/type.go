package api

import (
	"fmt"
	"strconv"
	"time"
)

type WishType int

const (
	BeginnersWish       WishType = 100 // Beginners' Wish (Novice Wish)
	StandardWish        WishType = 200 // Standard Wish (Permanent Wish)
	CharacterEventWish  WishType = 301 // Character Event Wish
	WeaponEventWish     WishType = 302 // Weapon Event Wish
	CharacterEventWish2 WishType = 400 // Character Event Wish-2
)

var WishTypes = []WishType{CharacterEventWish, CharacterEventWish2, WeaponEventWish, StandardWish, BeginnersWish}

func (w WishType) Str() string {
	return strconv.Itoa(int(w))
}

func (w WishType) Shared() SharedWishType {
	switch w {
	case CharacterEventWish, CharacterEventWish2:
		return SCharacterEventWish
	default:
		return SharedWishType(w)
	}
}

type SharedWishType int

const (
	SBeginnersWish      SharedWishType = 100 // Beginners' Wish (Novice Wish)
	SStandardWish       SharedWishType = 200 // Standard Wish (Permanent Wish)
	SCharacterEventWish SharedWishType = 301 // Character Event Wish and Character Event Wish-2
	SWeaponEventWish    SharedWishType = 302 // Weapon Event Wish
)

func (w SharedWishType) Str() string {
	return strconv.Itoa(int(w))
}

func (w SharedWishType) Pity(rarity Rarity) int {
	switch rarity {
	case Star5:
		switch w {
		case SCharacterEventWish:
			return 90
		case SWeaponEventWish:
			return 80
		case SStandardWish:
			return 90
		default:
			return 90
		}
	case Star4:
		return 10
	default:
		return 1
	}
}

var (
	SharedWishTypes = []SharedWishType{SCharacterEventWish, SWeaponEventWish, SStandardWish, SBeginnersWish}
)

type Rarity int

const (
	Star1 Rarity = iota + 1
	Star2
	Star3
	Star4
	Star5
)

func (r Rarity) Str() string {
	return strconv.Itoa(int(r))
}

type Time struct {
	time.Time
}

func (t *Time) String() string {
	return t.Time.Format("2006-01-02 15:04:05")
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Format(`"2006-01-02 15:04:05"`)), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var err error
	t.Time, err = time.Parse(`"2006-01-02 15:04:05"`, string(data))
	return err
}

type Item struct {
	Count    int      `json:"count,string"`
	WishType WishType `json:"gacha_type,string"`
	ID       int64    `json:"id,string"`
	ItemID   string   `json:"item_id"`
	ItemType string   `json:"item_type"`
	Lang     string   `json:"lang"`
	Name     string   `json:"name"`
	Rarity   Rarity   `json:"rank_type,string"`
	Time     Time     `json:"time"`
	UID      int      `json:"uid,string"`
}

func (item *Item) String() string {
	return fmt.Sprintf("Item{Name:%s Time:%v UID:%v ID:%v}", item.Name, item.Time.String(), item.UID, item.ID)
}

type Data struct {
	List   []*Item `json:"list"`
	Page   int     `json:"page,string"`
	Region string  `json:"region"`
	Size   int     `json:"size,string"`
	Total  int     `json:"total,string"`
}

type Result struct {
	Data    *Data  `json:"data"`
	Message string `json:"message"`
	RetCode int    `json:"retcode"`
}
