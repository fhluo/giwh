package cmd

import (
	"fmt"
	"github.com/fhluo/giwh/cmd/giwh"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/repository"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export wish history",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		items := config.Repository.GetItems()

		switch {
		case cmd.Flags().Changed("uid") && cmd.Flags().Changed("wish"):
			items = items.FilterByUID(uid).FilterBySharedWishType(lo.Map(wishes, func(i int, _ int) api.SharedWishType {
				return api.SharedWishType(i)
			})...)
		case cmd.Flags().Changed("uid"):
			items = items.FilterByUID(uid)
		case cmd.Flags().Changed("wish"):
			items = items.FilterBySharedWishType(lo.Map(wishes, func(i int, _ int) api.SharedWishType {
				return api.SharedWishType(i)
			})...)
		}

		if items.Len() == 0 {
			log.Fatalln("No such wish history.")
		}

		if err := repository.Save(args[0], items); err != nil {
			log.Fatalln(err)
		} else {
			fmt.Printf("%d items exported.\n", items.Len())
		}
	},
}

var (
	uid    int
	wishes []int
)

func init() {
	main.rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().IntVarP(&uid, "uid", "u", 0, "specify uid")
	exportCmd.Flags().IntSliceVarP(&wishes, "wishes", "w", nil, "specify wish types")
}
