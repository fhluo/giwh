package cmd

import (
	"github.com/fhluo/giwh/config"
	"github.com/fhluo/giwh/util"
	"github.com/fhluo/giwh/wh"
	"github.com/spf13/cobra"
	"log"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import wish history",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filenames, err := util.ExpandPaths(args...)
		if err != nil {
			logger.Fatalln(err)
		}

		for _, filename := range filenames {
			items, err := wh.LoadItems(filename)
			if err != nil {
				logger.Fatalln(err)
			}
			config.WishHistory = append(config.WishHistory, items...)
		}

		if err := config.SaveWishHistory(); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
