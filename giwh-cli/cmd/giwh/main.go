package main

import (
	"github.com/fhluo/giwh/common/config"
	"github.com/fhluo/giwh/gacha-logs/pipeline"
	"github.com/fhluo/giwh/gacha-logs/store"
	"github.com/fhluo/giwh/giwh-cli/pkg/export"
	_import "github.com/fhluo/giwh/giwh-cli/pkg/import"
	"github.com/fhluo/giwh/giwh-cli/pkg/merge"
	"github.com/fhluo/giwh/giwh-cli/pkg/stat"
	"github.com/fhluo/giwh/giwh-cli/pkg/update"
	"github.com/fhluo/giwh/giwh-cli/pkg/version"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"time"
)

var language string

var rootCmd = &cobra.Command{
	Use:   "giwh",
	Short: "Keep track of your Genshin Impact Wish History",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := store.New(nil)
		err := s.LoadIfExists(config.WishHistoryPath.Get())
		if err != nil {
			return err
		}

		p := pipeline.New(s.Unique())
		if p.Len() == 0 {
			return nil
		}

		if cmd.Flags().Changed("lang") {
			config.Language.Set(language)
		}

		stat.Stat(p.FilterByUID(p.Logs()[0].UID).Logs())

		return nil
	},
}

func init() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{Level: slog.LevelDebug, TimeFormat: time.TimeOnly}),
	))

	rootCmd.AddCommand(
		update.NewCmd(),
		stat.NewCmd(),
		_import.NewCmd(),
		export.NewCmd(),
		merge.NewCmd(),
		version.NewCmd(),
	)

	rootCmd.PersistentFlags().StringVarP(&language, "lang", "l", "", "set language")
}

func main() {
	defer config.Save()

	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
