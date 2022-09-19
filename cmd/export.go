package cmd

import (
	"fmt"
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
		p := config.WishHistory.Unique()

		switch {
		case cmd.Flags().Changed("uid") && cmd.Flags().Changed("wish"):
			p = p.FilterByUID(uid).FilterBySharedWishType(lo.Map(wishes, func(i int, _ int) api.SharedWishType {
				return api.SharedWishType(i)
			})...)
		case cmd.Flags().Changed("uid"):
			p = p.FilterByUID(uid)
		case cmd.Flags().Changed("wish"):
			p = p.FilterBySharedWishType(lo.Map(wishes, func(i int, _ int) api.SharedWishType {
				return api.SharedWishType(i)
			})...)
		}

		if p.Len() == 0 {
			log.Fatalln("No such wish history.")
		}

		if err := repository.Save(args[0], p); err != nil {
			log.Fatalln(err)
		} else {
			fmt.Printf("%d items exported.\n", p.Len())
		}
	},
}

var (
	uid    int
	wishes []int
)

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().IntVarP(&uid, "uid", "u", 0, "specify uid")
	exportCmd.Flags().IntSliceVarP(&wishes, "wishes", "w", nil, "specify wish types")
}
