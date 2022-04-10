package cmd

import (
	"github.com/fhluo/giwh/config"
	"github.com/fhluo/giwh/wh"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import wish history",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		for _, filename := range args {
			items, err := wh.LoadItems(filename)
			if err != nil {
				return err
			}
			config.WishHistory = append(config.WishHistory, items...)
		}

		return config.SaveWishHistory()
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
