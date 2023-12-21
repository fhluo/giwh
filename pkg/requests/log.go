package requests

type GachaType = string

const (
	NoviceWishes        GachaType = "100"
	PermanentWish       GachaType = "200"
	CharacterEventWish  GachaType = "301"
	WeaponEventWish     GachaType = "302"
	CharacterEventWish2 GachaType = "400"
)

var GachaTypes = []GachaType{
	NoviceWishes,        // 100
	PermanentWish,       // 200
	CharacterEventWish,  // 301
	WeaponEventWish,     // 302
	CharacterEventWish2, // 400
}

var SharedGachaTypes = []GachaType{
	NoviceWishes,       // 100
	PermanentWish,      // 200
	CharacterEventWish, // 301, 400
	WeaponEventWish,    // 302
}

// GachaLog 是抽卡记录
type GachaLog struct {
	Count     string `json:"count"`      // 物品数目
	GachaType string `json:"gacha_type"` // 卡池类型
	ID        string `json:"id"`         // 记录 ID
	ItemID    string `json:"item_id"`    // 物品 ID
	ItemType  string `json:"item_type"`  // 物品类型
	Lang      string `json:"lang"`       // 语言
	Name      string `json:"name"`       // 物品名称
	RankType  string `json:"rank_type"`  // 稀有度
	Time      string `json:"time"`       // 时间
	UID       string `json:"uid"`        // 用户 ID
}
