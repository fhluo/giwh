package cmd

import (
	_ "embed"
	"fmt"
	"github.com/fatih/color"
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/i18n"
	"github.com/fhluo/giwh/pkg/pipeline"
	"github.com/fhluo/giwh/pkg/repository"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"text/template"
)

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Show stats",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		items, err := repository.Load(args[0])
		if err != nil {
			log.Fatalln(err)
		}

		Stat(items)
	},
}

func init() {
	rootCmd.AddCommand(statCmd)
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
	Progress  map[int64]int
	Progress4 int
	Progress5 int
	Count     int
	Count5    int
	First     string
	Last      string
	Items5    []*api.Item
}

func Stat(items []*api.Item) {
	p := pipeline.New(items)
	progress := p.Progress()
	pulls := p.Pulls()

	err := tmpl.Execute(os.Stdout, lo.Map(api.SharedWishTypes, func(wishType api.SharedWishType, _ int) Info {
		current := p.FilterBySharedWishType(wishType)
		info := Info{
			Name:      i18n.GetSharedWishName(wishType.Str()),
			Pity4:     wishType.Pity(api.Star4),
			Pity5:     wishType.Pity(api.Star5),
			Progress:  pulls[wishType],
			Progress4: progress[wishType][api.Star4],
			Progress5: progress[wishType][api.Star5],
			Count:     current.Count(),
			Count5:    current.Count5Star(),
			Items5:    current.FilterByRarity(api.Star5).Items(),
		}

		if info.Count > 0 {
			info.First = current.First().Time.Format("2006/01/02")
			info.Last = current.Last().Time.Format("2006/01/02")
		}
		return info
	}))
	if err != nil {
		log.Fatalln(err)
	}
}
