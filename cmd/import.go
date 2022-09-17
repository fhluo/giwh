package cmd

import (
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/repository"
	"github.com/spf13/cobra"
	"log"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import wish history",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filenames, err := ExpandPaths(args...)
		if err != nil {
			log.Fatalln(err)
		}

		for _, filename := range filenames {
			items, err := repository.Load(filename)
			if err != nil {
				log.Fatalln(err)
			}

			config.WishHistory = config.WishHistory.Append(items)
		}

		if err := config.SaveWishHistory(); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
