package cmd

import (
	"errors"
	"fmt"
	"github.com/fhluo/giwh/config"
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
		uid, baseURL, err := util.GetUIDAndAPIBaseURL()
		if err != nil {
			return err
		}

		var (
			items wh.Items
		)

		var input string
		if len(args) != 0 {
			input = args[0]
		} else {
			input = uid + ".json"
		}

		items, err = util.LoadItemsIfExits(input)
		if err != nil {
			return err
		}

		items = append(items, lo.Filter(config.CachedItems, func(item wh.Item, _ int) bool {
			return item.UID == uid
		})...)

		visit := make(map[int64]bool)
		for _, item := range items {
			visit[item.ID()] = true
		}

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
			config.CachedItems = append(config.CachedItems, items_...)
		}
		_ = config.SaveCache()

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

		if len(args) == 2 && cmd.PersistentFlags().Changed("uid") {
			items = lo.Filter(items, func(item wh.Item, _ int) bool {
				return item.UID == uid
			})
		}

		if len(args) == 2 && cmd.PersistentFlags().Changed("wish") {
			items = lo.Filter(items, func(item wh.Item, _ int) bool {
				return item.WishType == wish
			})
		}

		items = lo.UniqBy(items, func(item wh.Item) int64 {
			return item.ID()
		})

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
