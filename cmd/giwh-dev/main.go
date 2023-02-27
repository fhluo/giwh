package main

import (
	"github.com/fhluo/giwh/pkg/cmd/gen"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"log"
	"os"
)

var rootCmd = cobra.Command{
	Use:   "giwh",
	Short: "Manage your Genshin Impact wish history",
}

func init() {
	log.SetFlags(0)

	rootCmd.AddCommand(gen.NewCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error(), nil)
		os.Exit(1)
	}
}
