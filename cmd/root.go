package cmd

import (
	"errors"
	"fmt"
	"github.com/fhluo/giwh/config"
	"github.com/fhluo/giwh/util"
	"github.com/fhluo/giwh/wh"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "giwh",
	Short: "Genshin Impact Wish History Exporter",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		authInfo, err := util.GetAuthInfo()
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
			input = authInfo.UID + ".json"
		}

		items, err = wh.LoadItemsIfExits(input)
		if err != nil {
			return err
		}

		items = append(items, config.CachedItems.FilterByUID(authInfo.UID)...)

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
			items_, err := wh.NewFetcher(authInfo.BaseURL, wish, visit).FetchALL()
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
			items = items.FilterByUID(uid)
		}

		if len(args) == 2 && cmd.PersistentFlags().Changed("wish") {
			items = items.FilterByWishType(wh.WishType(wish))
		}

		return items.Unique().Save(filename)
	},
}

var (
	uid  string
	wish int
)

func init() {
	rootCmd.PersistentFlags().StringVar(&uid, "uid", "", "specify uid")
	rootCmd.PersistentFlags().IntVar(&wish, "wish", 0, "specify wish type")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
	if err := config.Save(); err != nil {
		log.Printf("fail to save config file: %s\n", err)
	}
}
