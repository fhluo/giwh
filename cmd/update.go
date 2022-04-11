package cmd

import (
	"errors"
	"fmt"
	"github.com/fhluo/giwh/clients"
	"github.com/fhluo/giwh/config"
	"github.com/fhluo/giwh/util"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update wish history",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := clients.RecentlyUsed()
		if err != nil {
			if errors.Is(err, clients.ErrNotFound) {
				logger.Fatalln("Please open the wish history page in the game.")
			} else {
				logger.Fatalln(err)
			}
		}

		authInfo, err := client.GetAuthInfo()
		if err != nil {
			logger.Fatalln(err)
		}

		items, err := util.FetchAllWishHistory(authInfo.BaseURL, config.WishHistory.FilterByUID(authInfo.UID))
		if err != nil {
			logger.Fatalln(err)
		}

		if len(items) == 0 {
			fmt.Println("Your wish history is up to date.")
			return
		}

		config.WishHistory = append(config.WishHistory, items...)
		if err := config.SaveWishHistory(); err != nil {
			logger.Fatalln(err)
		}

		fmt.Printf("%d items fetched.\n", len(items))
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
