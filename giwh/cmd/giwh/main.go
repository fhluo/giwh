package main

import (
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/cmd/export"
	_import "github.com/fhluo/giwh/pkg/cmd/import"
	"github.com/fhluo/giwh/pkg/cmd/merge"
	"github.com/fhluo/giwh/pkg/cmd/stat"
	"github.com/fhluo/giwh/pkg/cmd/update"
	"github.com/fhluo/giwh/pkg/cmd/version"
	"github.com/fhluo/giwh/pkg/wish/pipeline"
	"github.com/fhluo/giwh/pkg/wish/repository"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
)

var language string

var rootCmd = &cobra.Command{
	Use:   "giwh",
	Short: "Genshin Impact Wish History Exporter",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := repository.LoadItems(config.WishHistoryPath.Get())
		if err != nil {
			return err
		}

		p := pipeline.New(items)
		if p.Len() == 0 {
			return nil
		}

		if cmd.Flags().Changed("lang") {
			config.Language.Set(language)
		}

		stat.Stat(p.FilterByUID(p.Items()[0].UID).Items())

		return nil
	},
}

func init() {
	log.SetFlags(0)

	rootCmd.AddCommand(update.NewCmd())
	rootCmd.AddCommand(stat.NewCmd())

	rootCmd.AddCommand(_import.NewCmd())
	rootCmd.AddCommand(export.NewCmd())
	rootCmd.AddCommand(merge.NewCmd())

	rootCmd.AddCommand(version.NewCmd())

	rootCmd.PersistentFlags().StringVarP(&language, "lang", "l", "", "set language")
}

func main() {
	defer config.Save()

	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
