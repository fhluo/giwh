package cmd

import (
	"github.com/fhluo/giwh/stat"
	"github.com/fhluo/giwh/wh"
	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:  "stat",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := wh.LoadItems(args[0])
		if err != nil {
			return err
		}

		stat.Stat(items)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statCmd)
}
