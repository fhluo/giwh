package wish

type Item struct {
	Count    string `json:"count"`
	WishType Type   `json:"gacha_type"`
	ID       string `json:"id"`
	ItemID   string `json:"item_id"`
	ItemType string `json:"item_type"`
	Lang     string `json:"lang"`
	Name     string `json:"name"`
	Rarity   Rarity `json:"rank_type"`
	Time     Time   `json:"time"`
	UID      string `json:"uid"`
}

func (item *Item) SharedWishType() Type {
	switch item.WishType {
	case CharacterEventWish, CharacterEventWish2:
		return CharacterEventWishAndCharacterEventWish2
	default:
		return item.WishType
	}
}

type (
	Type   string
	Rarity string
)

const (
	OneStar   Rarity = "1"
	TwoStar   Rarity = "2"
	ThreeStar Rarity = "3"
	FourStar  Rarity = "4"
	FiveStar  Rarity = "5"
)
