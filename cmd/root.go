package cmd

import (
	"errors"
	"github.com/fhluo/giwh/clients"
	"github.com/fhluo/giwh/config"
	"github.com/fhluo/giwh/stat"
	"github.com/fhluo/giwh/util"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", 0)

var rootCmd = &cobra.Command{
	Use:   "giwh",
	Short: "Genshin Impact Wish History Exporter",
	Args:  cobra.NoArgs,
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
			logger.Fatalln("Wish history not found.")
		}

		stat.Stat(items)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatalln(err)
	}
	if err := config.Save(); err != nil {
		logger.Fatalf("fail to save config file: %s\n", err)
	}
}
