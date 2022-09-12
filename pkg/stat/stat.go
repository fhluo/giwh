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

func stat(name string, pity4, pity5 int, p pipeline.Pipeline) Info {
	var (
		progress  = make(map[int64]int)
		progress4 int
		progress5 int
	)
	p.Traverse(func(e *pipeline.Element) bool {
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
		return true
	})

	info := Info{
		Name:      name,
		Pity4:     pity4,
		Pity5:     pity5,
		Progress:  progress,
		Progress4: progress4,
		Progress5: progress5,
		Count:     len(p.Elements()),
		Count5:    len(p.FilterByRarity(api.FiveStar).Elements()),
		Items5:    p.FilterByRarity(api.FiveStar).Items(),
	}

	if info.Count > 0 {
		info.First = p.Elements()[0].Time.Format("2006/01/02")
		info.Last = p.Elements()[len(p.Elements())-1].Time.Format("2006/01/02")
	}

	return info
}

func Stat(items []*api.Item) {
	p, err := pipeline.New(items)
	if err != nil {
		log.Fatalln(err)
	}
	p = p.Unique()

	err = template.Must(template.New("").Funcs(template.FuncMap{
		"gray":    gray,
		"repeat":  strings.Repeat,
		"hiBlack": color.HiBlackString,
		"magenta": color.MagentaString,
		"white":   color.WhiteString,
		"yellow":  color.YellowString,
	}).Parse(tmpl)).Execute(os.Stdout, []Info{
		stat(api.CharacterEventWish, 10, 90, p.FilterByWishType(api.CharacterEventWish, api.CharacterEventWish2).SortedByIDAscending()),
		stat(api.WeaponEventWish, 10, 80, p.FilterByWishType(api.WeaponEventWish).SortedByIDAscending()),
		stat(api.StandardWish, 10, 90, p.FilterByWishType(api.StandardWish).SortedByIDAscending()),
	})
	if err != nil {
		log.Fatalln(err)
	}
}
