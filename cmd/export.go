package cmd

import (
	"fmt"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/util"
	"github.com/fhluo/giwh/wh"
	"github.com/spf13/cobra"
	"log"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export wish history",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		items := config.WishHistory.Unique()

		switch {
		case cmd.Flags().Changed("uid") && cmd.Flags().Changed("wish"):
			items = items.FilterByUIDAndWishType(uid, util.Map(wishes, func(t int) wh.WishType {
				return wh.WishType(t)
			})...)
		case cmd.Flags().Changed("uid"):
			items = items.FilterByUID(uid)
		case cmd.Flags().Changed("wish"):
			items = items.FilterByWishType(util.Map(wishes, func(t int) wh.WishType {
				return wh.WishType(t)
			})...)
		}

		if len(items) == 0 {
			log.Fatalln("No such wish history.")
		}

		if err := items.Save(args[0]); err != nil {
			log.Fatalln(err)
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
