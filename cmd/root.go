package cmd

import (
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/i18n"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "giwh",
	Short: "Genshin Impact Wish History Exporter",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if config.WishHistory.Len() == 0 {
			return
		}

		Stat(config.WishHistory.FilterByUID(config.WishHistory.First().UID))
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
