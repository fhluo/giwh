package wish

import (
	"github.com/fhluo/giwh/pkg/i18n"
	"strconv"
)

type Type int

const (
	BeginnersWish       Type = 100 // Beginners' Wish (Novice Wish)
	StandardWish        Type = 200 // Standard Wish (Permanent Wish)
	CharacterEventWish  Type = 301 // Character Event Wish
	WeaponEventWish     Type = 302 // Weapon Event Wish
	CharacterEventWish2 Type = 400 // Character Event Wish-2
)

var (
	Types       = []Type{BeginnersWish, StandardWish, CharacterEventWish, WeaponEventWish, CharacterEventWish2}
	SharedTypes = []Type{CharacterEventWish, WeaponEventWish, StandardWish, BeginnersWish}
)

func (t Type) String() string {
	return i18n.GetWishName(strconv.Itoa(int(t)))
}

func (t Type) GetSharedWishName() string {
	return i18n.GetSharedWishName(strconv.Itoa(int(t)))
}
