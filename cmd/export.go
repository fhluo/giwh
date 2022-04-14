package cmd

import (
	"fmt"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/wh"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export wish history",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defer fmt.Println()

		items := config.WishHistory.Unique()

		switch {
		case cmd.Flags().Changed("uid") && cmd.Flags().Changed("wish"):
			items = items.FilterByUIDAndWishType(uid, lo.Map(wishes, func(t int, _ int) wh.WishType {
				return wh.WishType(t)
			})...)
		case cmd.Flags().Changed("uid"):
			items = items.FilterByUID(uid)
		case cmd.Flags().Changed("wish"):
			items = items.FilterByWishType(lo.Map(wishes, func(t int, _ int) wh.WishType {
				return wh.WishType(t)
			})...)
		}

		if len(items) == 0 {
			logger.Fatalln("No such wish history.")
		}

		if err := items.Save(args[0]); err != nil {
			logger.Fatalln(err)
		} else {
			fmt.Printf("%d items exported.\n", len(items))
		}
	},
}

var (
	uid    string
	wishes []int
)

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&uid, "uid", "u", "", "specify uid")
	exportCmd.Flags().IntSliceVarP(&wishes, "wishes", "w", nil, "specify wish types")
}
