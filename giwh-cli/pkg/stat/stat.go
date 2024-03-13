package stat

import (
	_ "embed"
	"fmt"
	"github.com/fatih/color"
	"github.com/fhluo/giwh/i18n"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/wish"
	"github.com/fhluo/giwh/pkg/wish/pipeline"
	"github.com/fhluo/giwh/pkg/wish/repository"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"text/template"
)

func NewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "stat",
		Short: "Show stats",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := repository.LoadItemsIfExits(args[0])
			if err != nil {
				return err
			}

			Stat(items)

			return nil
		},
	}

	return cmd

}

func gray(format string, a ...any) string {
	if color.NoColor {
		return fmt.Sprintf(format, a...)
	} else {
		return fmt.Sprintf("\x1B[38;5;239m"+format+"\x1B[0m", a...)
	}
}

var (
	//go:embed stat.tmpl
	tmplStr string
	tmpl    = template.Must(template.New("").Funcs(template.FuncMap{
		"gray":    gray,
		"repeat":  strings.Repeat,
		"hiBlack": color.HiBlackString,
		"magenta": color.MagentaString,
		"white":   color.WhiteString,
		"yellow":  color.YellowString,
	}).Parse(tmplStr))
)

type Info struct {
	Name      string
	Pity4     int
	Pity5     int
	Pulls     map[int64]int
	Progress  int
	Count     int
	Count5    int
	First     string
	Last      string
	Items5    []wish.Item
	Progress4 int
	Progress5 int
}

func Stat(items []wish.Item) {
	p := pipeline.New(items)

	locale := lo.Must(i18n.ReadLocale(i18n.Match(config.Language.Get())))

	err := tmpl.Execute(os.Stdout, lo.FilterMap(wish.SharedWishes, func(wishType wish.Type, _ int) (Info, bool) {
		current := p.FilterBySharedWish(wishType)
		if current.Len() == 0 {
			return Info{}, false
		}

		fiveStars := current.FilterByRarity(wish.FiveStar)

		pity5 := 90
		if wishType == wish.WeaponEventWish {
			pity5 = 80
		}

		info := Info{
			Name:      locale.Wishes[int(wishType)],
			Pity4:     10,
			Pity5:     pity5,
			Pulls:     current.Pulls5Stars(),
			Progress:  current.Progress5Star(),
			Count:     current.Len(),
			Count5:    fiveStars.Len(),
			Items5:    fiveStars.Items(),
			Progress4: current.Progress4Star(),
			Progress5: current.Progress5Star(),
		}

		if info.Count > 0 {
			info.First = current.First().Time.Format("2006/01/02")
			info.Last = current.Last().Time.Format("2006/01/02")
		}
		return info, true
	}))
	if err != nil {
		log.Fatalln(err)
	}
}
