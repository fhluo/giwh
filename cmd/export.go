package cmd

import (
	"fmt"
	"github.com/fhluo/giwh/config"
	"github.com/fhluo/giwh/wh"
	"github.com/spf13/cobra"
	"os"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export wish history",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defer fmt.Println()

		items := config.WishHistory.Unique()

		switch {
		case cmd.Flags().Changed("uid") && cmd.Flags().Changed("wish"):
			items = items.FilterByUIDAndWishType(uid, wh.WishType(wish))
		case cmd.Flags().Changed("uid"):
			items = items.FilterByUID(uid)
		case cmd.Flags().Changed("wish"):
			items = items.FilterByWishType(wh.WishType(wish))
		}

		if len(items) == 0 {
			fmt.Println("No such wish history.")
			return
		}

		if err := items.Save(args[0]); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Printf("%d items exported.\n", len(items))
		}
	},
}

var (
	uid  string
	wish int
)

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&uid, "uid", "u", "", "specify uid")
	exportCmd.Flags().IntVarP(&wish, "wish", "w", 0, "specify wish type")
}
