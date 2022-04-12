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
	return i.ColoredString() + color.HiBlackString("(%d)", i.pulls)
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
	sort.Sort(items)

	fmt.Println()
	color.HiBlack(title)
	fmt.Println()
	infos, fourStar, fiveStar := stat(items)
	fmt.Println(
		color.MagentaString("4-Star:"),
		color.WhiteString("%2d", fourStar), color.HiBlackString("/ %d", fourStarPity),
	)
	fmt.Println(
		color.YellowString("5-Star:"),
		color.WhiteString("%2d", fiveStar), color.HiBlackString("/ %d", fiveStarPity),
	)

	if items.FilterByRarity(wh.FiveStar).Count() != 0 {
		fmt.Println()
		show5stars(infos)
	}

	if len(items) >= 2 {
		fmt.Println()
		color.HiBlack("%s ~ %s", items[0].Time().Format("2006/01/02"), items[len(items)-1].Time().Format("2006/01/02"))
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
	items = items.Unique()

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
