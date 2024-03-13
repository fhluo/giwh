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
