package gacha

// Log 表示抽卡记录
type Log struct {
	ID        string `json:"id"`         // 记录 ID
	UID       string `json:"uid"`        // 用户 ID
	GachaType string `json:"gacha_type"` // 卡池类型
	Name      string `json:"name"`       // 物品名称
	ItemID    string `json:"item_id"`    // 物品 ID
	ItemType  string `json:"item_type"`  // 物品类型
	RankType  string `json:"rank_type"`  // 稀有度
	Count     string `json:"count"`      // 物品数量
	Time      string `json:"time"`       // 时间
	Lang      string `json:"lang"`       // 语言
}

type Type = string

const (
	NoviceWishes        Type = "100" // 新手祈愿
	PermanentWish       Type = "200" // 常驻祈愿
	CharacterEventWish  Type = "301" // 角色活动祈愿
	WeaponEventWish     Type = "302" // 武器活动祈愿
	CharacterEventWish2 Type = "400" // 角色活动祈愿-2
	ChronicledWish      Type = "500" // 集录祈愿
)

var Types = []Type{
	NoviceWishes,        // 100
	PermanentWish,       // 200
	CharacterEventWish,  // 301
	WeaponEventWish,     // 302
	CharacterEventWish2, // 400
	ChronicledWish,      // 500
}

var SharedTypes = []Type{
	NoviceWishes,       // 100 新手祈愿
	PermanentWish,      // 200 常驻祈愿
	CharacterEventWish, // 301, 400 角色活动祈愿与角色活动祈愿-2
	WeaponEventWish,    // 302 武器活动祈愿
	ChronicledWish,     // 500 集录祈愿
}
