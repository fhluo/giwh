package stat

import (
	_ "embed"
	"fmt"
	"github.com/fatih/color"
	"github.com/fhluo/giwh/common/config"
	"github.com/fhluo/giwh/common/i18n"
	"github.com/fhluo/giwh/gacha-logs/gacha"
	"github.com/fhluo/giwh/gacha-logs/pipeline"
	"github.com/fhluo/giwh/gacha-logs/store"
	"log/slog"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
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
			s := store.New(nil)
			err := s.LoadIfExists(args[0])
			if err != nil {
				return err
			}

			Stat(s.Unique())

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
	Pulls     map[string]int
	Progress  int
	Count     int
	Count5    int
	First     string
	Last      string
	Items5    []gacha.Log
	Progress4 int
	Progress5 int
}

func Stat(logs []gacha.Log) {
	p := pipeline.New(logs).Unique().SortByIDAscending()

	locale := lo.Must(i18n.ReadLocale(i18n.Match(config.Language.Get())))

	err := tmpl.Execute(os.Stdout, lo.FilterMap(gacha.SharedTypes, func(wishType gacha.Type, _ int) (Info, bool) {
		current := p.FilterBySharedWish(wishType)
		if current.Len() == 0 {
			return Info{}, false
		}

		fiveStars := current.FilterByRarity("5")

		pity5 := 90
		if wishType == gacha.WeaponEventWish {
			pity5 = 80
		}

		info := Info{
			Name:      locale.Wishes[wishType],
			Pity4:     10,
			Pity5:     pity5,
			Pulls:     current.Pulls5Stars(),
			Progress:  current.Progress5Star(),
			Count:     current.Len(),
			Count5:    fiveStars.Len(),
			Items5:    fiveStars.Logs(),
			Progress4: current.Progress4Star(),
			Progress5: current.Progress5Star(),
		}

		current.SortByIDAscending()
		if info.Count > 0 {
			r, err := current.First().ParseTime()
			if err != nil {
				slog.Error(err.Error())
				return Info{}, false
			}
			slog.Debug("stat", "first", current.First())
			info.First = r.Format("2006/01/02")

			r, err = current.Last().ParseTime()
			if err != nil {
				slog.Error(err.Error())
				return Info{}, false
			}
			slog.Debug("stat", "last", current.Last())
			info.Last = r.Format("2006/01/02")
		}
		slog.Debug("stat", "info", info)

		return info, true
	}))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
