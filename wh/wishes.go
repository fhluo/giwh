package wh

type WishType int

const (
	BeginnersWish       WishType = 100 // Beginners' Wish (Novice Wish)
	StandardWish        WishType = 200 // Standard Wish (Permanent Wish)
	CharacterEventWish  WishType = 301 // Character Event Wish
	WeaponEventWish     WishType = 302 // Weapon Event Wish
	CharacterEventWish2 WishType = 400 // Character Event Wish-2
)

var (
	Wishes = []WishType{BeginnersWish, StandardWish, CharacterEventWish, WeaponEventWish, CharacterEventWish2}

	wishes = map[WishType]string{
		BeginnersWish:       "Beginners' Wish",
		StandardWish:        "Standard Wish",
		CharacterEventWish:  "Character Event Wish",
		WeaponEventWish:     "Weapon Event Wish",
		CharacterEventWish2: "Character Event Wish-2",
	}
)

func (t WishType) String() string {
	return wishes[t]
}
