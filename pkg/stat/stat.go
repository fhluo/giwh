package stat

import (
	_ "embed"
	"fmt"
	"github.com/fatih/color"
	"github.com/fhluo/giwh/pkg/wish"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
)

func gray(format string, a ...any) string {
	if color.NoColor {
		return fmt.Sprintf(format, a...)
	} else {
		return fmt.Sprintf("\x1B[38;5;239m"+format+"\x1B[0m", a...)
	}
}

//go:embed stat.tmpl
var tmpl string

type Info struct {
	Name      string
	Pity4     int
	Pity5     int
	Progress  map[int64]int
	Progress4 int
	Progress5 int
	Count     int
	Count5    int
	First     string
	Last      string
	Items5    wish.History
}

func stat(name string, pity4, pity5 int, items wish.History) Info {
	sort.Sort(items)

	var (
		progress  = make(map[int64]int)
		progress4 int
		progress5 int
	)
	for i := range items {
		switch items[i].Rarity() {
		case wish.FourStar:
			progress[items[i].ID()] = progress4 + 1
			progress4 = 0
			progress5++
		case wish.FiveStar:
			progress[items[i].ID()] = progress5 + 1
			progress4++
			progress5 = 0
		default:
			progress4++
			progress5++
		}
	}

	info := Info{
		Name:      name,
		Pity4:     pity4,
		Pity5:     pity5,
		Progress:  progress,
		Progress4: progress4,
		Progress5: progress5,
		Count:     items.Count(),
		Count5:    items.FilterByRarity(wish.FiveStar).Count(),
		Items5:    items.FilterByRarity(wish.FiveStar),
	}

	if info.Count > 0 {
		info.First = items[0].Time().Format("2006/01/02")
		info.Last = items[len(items)-1].Time().Format("2006/01/02")
	}

	return info
}

func Stat(items wish.History) {
	items = items.Unique()

	err := template.Must(template.New("").Funcs(template.FuncMap{
		"gray":    gray,
		"repeat":  strings.Repeat,
		"hiBlack": color.HiBlackString,
		"magenta": color.MagentaString,
		"white":   color.WhiteString,
		"yellow":  color.YellowString,
	}).Parse(tmpl)).Execute(os.Stdout, []Info{
		stat(wish.CharacterEventWish.GetSharedWishName(), 10, 90, items.FilterByWishType(wish.CharacterEventWish, wish.CharacterEventWish2)),
		stat(wish.WeaponEventWish.GetSharedWishName(), 10, 80, items.FilterByWishType(wish.WeaponEventWish)),
		stat(wish.StandardWish.GetSharedWishName(), 10, 90, items.FilterByWishType(wish.StandardWish)),
	})
	if err != nil {
		log.Fatalln(err)
	}
}
