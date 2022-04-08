package cmd

import (
	"github.com/fhluo/giwh/util"
	"github.com/fhluo/giwh/wh"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"path/filepath"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge wish histories",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		filenames := make([]string, 0, len(args)-1)

		for _, arg := range args[:len(args)-1] {
			matches, err := filepath.Glob(arg)
			if err != nil {
				return err
			}
			filenames = append(filenames, matches...)
		}

		var result wh.Items

		for _, filename := range filenames {
			items, err := util.LoadItems(filename)
			if err != nil {
				return err
			}
			result = append(result, items...)
		}

		result = lo.UniqBy(result, func(item wh.Item) int64 {
			return item.ID()
		})

		if cmd.PersistentFlags().Changed("uid") {
			result = lo.Filter(result, func(item wh.Item, _ int) bool {
				return item.UID == uid
			})
		}

		if cmd.PersistentFlags().Changed("wish") {
			result = lo.Filter(result, func(item wh.Item, _ int) bool {
				return item.WishType == wish
			})
		}

		return result.Save(args[len(args)-1])
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}
