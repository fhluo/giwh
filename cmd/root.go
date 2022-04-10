package cmd

import (
	"errors"
	"fmt"
	"github.com/fhluo/giwh/clients"
	"github.com/fhluo/giwh/config"
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
		client, err := clients.RecentlyUsed()
		if err != nil {
			return err
		}

		authInfo, err := client.GetAuthInfo()
		if err != nil {
			return err
		}

		fmt.Println("UID:", authInfo.UID)

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

		items = append(items, config.WishHistory.FilterByUID(authInfo.UID)...)

		visit := make(map[int64]bool)
		for _, item := range items {
			visit[item.ID()] = true
		}

		if err != nil {
			if errors.Is(err, clients.ErrNotFound) {
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
			config.WishHistory = append(config.WishHistory, items_...)
		}
		_ = config.SaveWishHistory()

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

		return items.Unique().Save(filename)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := config.Save(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fail to save config file: %s\n", err)
	}
}
