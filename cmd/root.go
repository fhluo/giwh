package cmd

import (
	"errors"
	"fmt"
	"github.com/fhluo/giwh/util"
	"github.com/fhluo/giwh/wh"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "giwh",
	Short: "Genshin Impact Wish History Exporter",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			items wh.Items
			err   error
		)

		if len(args) != 0 {
			items, err = util.LoadItemsIfExits(args[0])
			if err != nil {
				return err
			}
		}

		visit := make(map[int64]bool)
		for _, item := range items {
			visit[item.ID()] = true
		}

		baseURL, err := util.GetAPIBaseURL()

		if err != nil {
			if errors.Is(err, util.ErrNotFound) {
				_, _ = fmt.Fprintln(os.Stderr, "Please open the wish history page in the game.")
			}
			return err
		}

		for _, wish := range wh.Wishes {
			items_, err := wh.NewFetcher(baseURL, wish, visit).FetchALL()
			if err != nil {
				log.Fatalln(err)
			}

			items = append(items, items_...)
		}

		if len(items) == 0 {
			return fmt.Errorf("wish history not found")
		}

		var filename string
		switch len(args) {
		case 0:
			filename = items[0].UID + ".json"
		case 1:
			filename = args[0]
		default:
			filename = args[1]
		}

		if cmd.PersistentFlags().Changed("uid") {
			items = lo.Filter(items, func(item wh.Item, _ int) bool {
				return item.UID == uid
			})
		}

		if cmd.PersistentFlags().Changed("wish") {
			items = lo.Filter(items, func(item wh.Item, _ int) bool {
				return item.WishType == wish
			})
		}

		return items.Save(filename)
	},
}

var (
	uid  string
	wish string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&uid, "uid", "", "specify uid")
	rootCmd.PersistentFlags().StringVar(&wish, "wish", "", "specify wish type")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
