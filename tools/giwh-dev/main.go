package main

import (
	"github.com/fhluo/giwh/common/config"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
	"tools/dl"
	"tools/gen"
)

var rootCmd = cobra.Command{
	Use:   "giwh",
	Short: "Manage your Genshin Impact wish history",
}

func init() {
	log.SetFlags(0)

	rootCmd.AddCommand(gen.NewCmd())
	rootCmd.AddCommand(dl.NewCmd())
}

func main() {
	defer config.Save()

	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
