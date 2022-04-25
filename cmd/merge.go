package cmd

import (
	"github.com/fhluo/giwh/pkg/util"
	"github.com/fhluo/giwh/wh"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge wish histories",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filenames, err := util.ExpandPaths(args...)
		if err != nil {
			logger.Fatalln(err)
		}

		var result wh.WishHistory

		for _, filename := range filenames {
			items, err := wh.LoadWishHistory(filename)
			if err != nil {
				logger.Fatalln(err)
			}
			result = append(result, items...)
		}

		if err = result.Unique().Save(outputFilename); err != nil {
			logger.Fatalln(err)
		}
	},
}

var outputFilename string

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().StringVarP(&outputFilename, "output", "o", "", "specify output filename")
	if err := mergeCmd.MarkFlagRequired("output"); err != nil {
		logger.Fatalln(err)
	}
}
