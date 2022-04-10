package stat

import (
	"fmt"
	"github.com/fatih/color"
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
	return fmt.Sprintf("%s(%d)", i.ColoredString(), i.pulls)
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

func show(items wh.Items, title string, fourStarPity int, fiveStarPity int) {
	color.HiBlue(title)
	fmt.Println()
	infos, fourStar, fiveStar := stat(items)
	fmt.Printf("Next 4 star in %s pulls.\n", color.HiWhiteString("%2d", fourStarPity-fourStar))
	fmt.Printf("Next 5 star in %s pulls.\n", color.HiWhiteString("%2d", fiveStarPity-fiveStar))

	if items.FilterByRarity(wh.FiveStar).Count() != 0 {
		fmt.Println()
		show5stars(infos)
		fmt.Println()
	}
}

func drawLine(length int) {
	if color.NoColor {
		fmt.Println(strings.Repeat("─", length))
	} else {
		fmt.Printf("\x1B[38;5;239m%s\x1B[0m\n", strings.Repeat("─", length))
	}
}

func Stat(items wh.Items) {
	color.New()

	items = items.Unique()
	sort.Sort(items)

	drawLine(50)
	show(
		items.FilterByWishType(wh.CharacterEventWish, wh.CharacterEventWish2),
		"Character Event Wish and Character Event Wish-2", 10, 90,
	)

	drawLine(50)
	show(items.FilterByWishType(wh.WeaponEventWish), "Weapon Event Wish", 10, 80)

	drawLine(50)
	show(items.FilterByWishType(wh.StandardWish), "Standard Wish", 10, 90)

	drawLine(50)
}
