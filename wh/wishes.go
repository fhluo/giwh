package wh

import (
	"github.com/fhluo/giwh/i18n"
	"strconv"
)

type WishType int

const (
	BeginnersWish       WishType = 100 // Beginners' Wish (Novice Wish)
	StandardWish        WishType = 200 // Standard Wish (Permanent Wish)
	CharacterEventWish  WishType = 301 // Character Event Wish
	WeaponEventWish     WishType = 302 // Weapon Event Wish
	CharacterEventWish2 WishType = 400 // Character Event Wish-2
)

var (
	Wishes       = []WishType{BeginnersWish, StandardWish, CharacterEventWish, WeaponEventWish, CharacterEventWish2}
	SharedWishes = []WishType{CharacterEventWish, WeaponEventWish, StandardWish, BeginnersWish}
)

func (t WishType) String() string {
	return i18n.GetWishName(strconv.Itoa(int(t)))
}

func (t WishType) GetSharedWishName() string {
	return i18n.GetSharedWishName(strconv.Itoa(int(t)))
}
