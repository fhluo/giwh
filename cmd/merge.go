package cmd

import (
	"github.com/fhluo/giwh/wh"
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

		result, err := wh.LoadItemsIfExits(args[len(args)-1])
		if err != nil {
			return err
		}

		for _, filename := range filenames {
			items, err := wh.LoadItems(filename)
			if err != nil {
				return err
			}
			result = append(result, items...)
		}

		result = result.Unique()

		if cmd.PersistentFlags().Changed("uid") {
			result = result.FilterByUID(uid)
		}

		if cmd.PersistentFlags().Changed("wish") {
			result = result.FilterByWishType(wish)
		}

		return result.Save(args[len(args)-1])
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}
