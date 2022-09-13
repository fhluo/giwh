package stat

import (
	_ "embed"
	"fmt"
	"github.com/fatih/color"
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/pipeline"
	"log"
	"os"
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
	Items5    []*api.Item
}

func stat(wishType string, p pipeline.Pipeline) Info {
	var (
		progress  = make(map[int64]int)
		progress4 int
		progress5 int
	)
	p.Traverse(func(e *pipeline.Element) {
		switch e.Rarity {
		case api.FourStar:
			progress[e.ID] = progress4 + 1
			progress4 = 0
			progress5++
		case api.FiveStar:
			progress[e.ID] = progress5 + 1
			progress4++
			progress5 = 0
		default:
			progress4++
			progress5++
		}
	})

	info := Info{
		Name:      wishType,
		Pity4:     api.Pity4Star(wishType),
		Pity5:     api.Pity5Star(wishType),
		Progress:  progress,
		Progress4: progress4,
		Progress5: progress5,
		Count:     p.Count(),
		Count5:    p.Count5Star(),
		Items5:    p.FilterByRarity(api.FiveStar).Items(),
	}

	if p.Count() > 0 {
		info.First = p.First().Time.Format("2006/01/02")
		info.Last = p.Last().Time.Format("2006/01/02")
	}

	return info
}

func Stat(items []*api.Item) {
	p, err := pipeline.New(items)
	if err != nil {
		log.Fatalln(err)
	}
	p.SortByIDAscending()

	err = template.Must(template.New("").Funcs(template.FuncMap{
		"gray":    gray,
		"repeat":  strings.Repeat,
		"hiBlack": color.HiBlackString,
		"magenta": color.MagentaString,
		"white":   color.WhiteString,
		"yellow":  color.YellowString,
	}).Parse(tmpl)).Execute(os.Stdout, []Info{
		stat(api.CharacterEventWish, p.FilterByWishType(api.CharacterEventWish, api.CharacterEventWish2)),
		stat(api.WeaponEventWish, p.FilterByWishType(api.WeaponEventWish)),
		stat(api.StandardWish, p.FilterByWishType(api.StandardWish)),
	})
	if err != nil {
		log.Fatalln(err)
	}
}
