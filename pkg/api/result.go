package api

type Item struct {
	UID      string `json:"uid"`
	WishType string `json:"gacha_type"`
	ItemID   string `json:"item_id"`
	Count    string `json:"count"`
	Time     string `json:"time"`
	Name     string `json:"name"`
	Lang     string `json:"lang"`
	ItemType string `json:"item_type"`
	Rarity   string `json:"rank_type"`
	ID       string `json:"id"`
}

type Data struct {
	Page   string `json:"page"`
	Size   string `json:"size"`
	Total  string `json:"total"`
	List   []Item `json:"list"`
	Region string `json:"region"`
}

type Result struct {
	RetCode int    `json:"retcode"`
	Message string `json:"message"`
	Data    `json:"data"`
}
