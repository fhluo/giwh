package cmd

import (
	"errors"
	"fmt"
	"github.com/fhluo/giwh/fetcher"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/clients"
	"github.com/spf13/cobra"
	"log"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update wish history",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := clients.RecentlyUsed()
		if err != nil {
			if errors.Is(err, clients.ErrURLNotFound) {
				log.Fatalln("Please open the wish history page in the game.")
			} else {
				log.Fatalln(err)
			}
		}

		uid, err := client.GetUID()
		if err != nil {
			log.Fatalln(err)
		}

		baseURL, err := client.GetBaseURL()
		if err != nil {
			if errors.Is(err, clients.ErrURLNotFound) {
				authInfo, ok := config.GetAuthInfo(uid)
				if ok {
					uid = authInfo.UID
					baseURL = authInfo.BaseURL
				} else {
					log.Fatalln(err)
				}
			} else {
				log.Fatalln(err)
			}
		}

		config.UpdateAuthInfo(fetcher.AuthInfo{UID: uid, BaseURL: baseURL})
		_ = config.Save()

		items := config.WishHistory.FilterByUID(uid)
		result, err := fetcher.FetchAllWishHistory(baseURL, items)
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
