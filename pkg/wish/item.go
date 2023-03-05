package wish

type Item struct {
	Count    int    `json:"count,string"`
	WishType Type   `json:"gacha_type,string"`
	ID       int64  `json:"id,string"`
	ItemID   string `json:"item_id"`
	ItemType string `json:"item_type"`
	Lang     string `json:"lang"`
	Name     string `json:"name"`
	Rarity   int    `json:"rank_type,string"`
	Time     Time   `json:"time"`
	UID      int    `json:"uid,string"`
}

func (item *Item) SharedWishType() Type {
	switch item.WishType {
	case CharacterEventWish, CharacterEventWish2:
		return CharacterEventWishAndCharacterEventWish2
	default:
		return item.WishType
	}
}

const (
	OneStar int = iota + 1
	TwoStar
	ThreeStar
	FourStar
	FiveStar
)
