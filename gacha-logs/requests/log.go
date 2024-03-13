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
