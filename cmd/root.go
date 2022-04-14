package cmd

import (
	"errors"
	"fmt"
	"github.com/fhluo/giwh/clients"
	"github.com/fhluo/giwh/config"
	"github.com/fhluo/giwh/i18n"
	"github.com/fhluo/giwh/stat"
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

		uid, err := client.GetUID()
		if err != nil {
			logger.Fatalln(err)
		}

		items := config.WishHistory.FilterByUID(uid)
		if len(items) == 0 {
			fmt.Printf("The wish history is empty. (UID: %s)\n", uid)
			fmt.Println("You can use the update subcommand to update the wish history.")
			fmt.Println()
			return
		}

		stat.Stat(items)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&i18n.Language, "lang", "l", "", "set language")
}

func Execute() {
	if config.GetLanguage() != "" {
		i18n.Language = config.GetLanguage()
	}
	if err := rootCmd.Execute(); err != nil {
		logger.Fatalln(err)
	}
	if rootCmd.PersistentFlags().Changed("lang") {
		config.SetLanguage(i18n.Language)
		if err := config.Save(); err != nil {
			logger.Fatalln(err)
		}
	}
}
