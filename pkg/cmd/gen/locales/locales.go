package locales

import (
	_ "embed"
	"github.com/fhluo/giwh/i18n"
	"github.com/fhluo/giwh/pkg/wiki"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type Locale = i18n.Locale

func NewCmd() *cobra.Command {
	var (
		outputDir   string
		packageName string
	)

	cmd := &cobra.Command{
		Use:   "locales",
		Short: "Generate i18n/locales",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			_ = os.MkdirAll(outputDir, 0666)

			results := make(map[i18n.Language]map[wiki.EntryIndex]wiki.Entry)

			for _, lang := range i18n.Languages {
				w := wiki.New(lang)

				results[lang], err = w.GetMenusEntries(wiki.CharacterArchive.ID, wiki.Weapons.ID)
				if err != nil {
					return err
				}
			}

			en := results[i18n.English]

			var locales []Locale

			for lang, entries := range results {
				locale := i18n.NewLocale(lang)

				for index, entry := range entries {
					switch index.MenuID {
					case wiki.CharacterArchive.ID:
						locale.Characters[en[index].Name] = entry.Name
					case wiki.Weapons.ID:
						locale.Weapons[en[index].Name] = entry.Name
					}
				}

				locales = append(locales, locale)
			}

			// TODO Merge with existing file
			for _, local := range locales {
				if err = os.WriteFile(filepath.Join(outputDir, local.BaseFilename()), local.JSON(), 0666); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "locales", "output directory")
	cmd.Flags().StringVarP(&packageName, "pkg", "p", "i18n", "package name")

	return cmd
}
