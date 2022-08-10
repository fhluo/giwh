package cmd

import (
	"github.com/fhluo/giwh/pkg/stat"
	"github.com/fhluo/giwh/pkg/wish"
	"github.com/spf13/cobra"
	"log"
)

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Show stats",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		items, err := wish.LoadWishHistory(args[0])
		if err != nil {
			log.Fatalln(err)
		}

		stat.Stat(items)
	},
}

func init() {
	rootCmd.AddCommand(statCmd)
}
