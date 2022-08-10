package wish

type Rarity int

const (
	OneStar Rarity = iota + 1
	TwoStar
	ThreeStar
	FourStar
	FiveStar
)

var rarities = map[Rarity]string{
	OneStar:   "1-Star",
	TwoStar:   "2-Star",
	ThreeStar: "3-Star",
	FourStar:  "4-Star",
	FiveStar:  "5-Star",
}

func (r Rarity) String() string {
	return rarities[r]
}
