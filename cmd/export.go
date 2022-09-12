package cmd

import (
	"fmt"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/repository"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export wish history",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		p := config.WishHistory.Unique()

		switch {
		case cmd.Flags().Changed("uid") && cmd.Flags().Changed("wish"):
			p = p.FilterByUID(uid).FilterByWishType(lo.Map(wishes, func(t int, _ int) string {
				return strconv.Itoa(t)
			})...)
		case cmd.Flags().Changed("uid"):
			p = p.FilterByUID(uid)
		case cmd.Flags().Changed("wish"):
			p = p.FilterByWishType(lo.Map(wishes, func(t int, _ int) string {
				return strconv.Itoa(t)
			})...)
		}

		if len(p.Elements()) == 0 {
			log.Fatalln("No such wish history.")
		}

		if err := repository.Save(args[0], p.Items()); err != nil {
			log.Fatalln(err)
		} else {
			fmt.Printf("%d items exported.\n", len(p.Elements()))
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
