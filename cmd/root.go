package cmd

import (
	"errors"
	"fmt"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/clients"
	"github.com/fhluo/giwh/pkg/i18n"
	"github.com/fhluo/giwh/pkg/stat"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "giwh",
	Short: "Genshin Impact Wish History Exporter",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := clients.Default()
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
		log.Fatalln(err)
	}
	if rootCmd.PersistentFlags().Changed("lang") {
		config.SetLanguage(i18n.Language)
		if err := config.Save(); err != nil {
			log.Fatalln(err)
		}
	}
}
