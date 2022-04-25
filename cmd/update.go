package cmd

import (
	"errors"
	"fmt"
	"github.com/fhluo/giwh/clients"
	"github.com/fhluo/giwh/fetcher"
	"github.com/fhluo/giwh/internal/config"
	"github.com/spf13/cobra"
	"log"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update wish history",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := clients.RecentlyUsed()
		if err != nil {
			if errors.Is(err, clients.ErrNotFound) {
				log.Fatalln("Please open the wish history page in the game.")
			} else {
				log.Fatalln(err)
			}
		}

		authInfo, err := client.GetAuthInfo()
		if err != nil {
			log.Fatalln(err)
		}

		items := config.WishHistory.FilterByUID(authInfo.UID)
		result, err := fetcher.FetchAllWishHistory(authInfo.BaseURL, items)
		if err != nil {
			log.Fatalln(err)
		}

		count := len(result) - len(items)
		if count == 0 {
			fmt.Println("No items fetched. Your wish history is up to date.")
			return
		}

		config.WishHistory = append(config.WishHistory, result...)
		if err := config.SaveWishHistory(); err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%d items fetched.\n", count)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
