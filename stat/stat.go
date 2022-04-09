package stat

import (
	"fmt"
	"github.com/fhluo/giwh/wh"
	"github.com/samber/lo"
	"sort"
	"strings"
)

type info struct {
	*wh.Item

	pulls int
}

func (i info) String() string {
	return fmt.Sprintf("%s(%d)", i.Name, i.pulls)
}

func stat(items wh.Items) (infos []info, fourStar int, fiveStar int) {
	for i := range items {
		switch items[i].Rarity() {
		case wh.FourStar:
			infos = append(infos, info{Item: &items[i], pulls: fourStar + 1})
			fourStar = 0
			fiveStar++
		case wh.FiveStar:
			infos = append(infos, info{Item: &items[i], pulls: fiveStar + 1})
			fourStar++
			fiveStar = 0
		default:
			fourStar++
			fiveStar++
		}
	}

	return
}

func show5stars(infos []info) {
	fmt.Println(strings.Join(
		lo.Map(
			lo.Filter(infos, func(i info, _ int) bool {
				return i.Rarity() == wh.FiveStar
			}),
			func(i info, _ int) string {
				return i.String()
			},
		),
		" ",
	))
}

func Stat(items wh.Items) {
	items = items.Unique()
	sort.Sort(items)

	fmt.Println("Character Event Wish and Character Event Wish-2")
	infos, fourStar, fiveStar := stat(items.FilterByWishType(wh.CharacterEventWish, wh.CharacterEventWish2))
	fmt.Printf("Next 4 star in %d pulls\n", 10-fourStar)
	fmt.Printf("Next 5 star in %d pulls\n", 90-fiveStar)
	show5stars(infos)
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("Weapon Event Wish")
	infos, fourStar, fiveStar = stat(items.FilterByWishType(wh.WeaponEventWish))
	fmt.Printf("Next 4 star in %d pulls\n", 10-fourStar)
	fmt.Printf("Next 5 star in %d pulls\n", 80-fiveStar)
	show5stars(infos)
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("Standard Wish")
	infos, fourStar, fiveStar = stat(items.FilterByWishType(wh.StandardWish))
	fmt.Printf("Next 4 star in %d pulls\n", 10-fourStar)
	fmt.Printf("Next 5 star in %d pulls\n", 90-fiveStar)
	show5stars(infos)
	fmt.Println(strings.Repeat("-", 50))
}
