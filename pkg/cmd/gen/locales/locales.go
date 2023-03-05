package locales

import (
	_ "embed"
	"github.com/fhluo/giwh/i18n"
	"github.com/fhluo/giwh/pkg/cmd/gen/wishes"
	"github.com/fhluo/giwh/pkg/wiki"
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"os"
	"path/filepath"
	"sync"
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

			results, err := GetAllEntries()
			if err != nil {
				slog.Error(err.Error(), nil)
				os.Exit(1)
			}

			en := results[i18n.English]

			var locales []Locale

			for lang, entries := range results {
				locale := i18n.NewLocale(lang)

				locale.Wishes = lo.SliceToMap(wishes.Wishes[lang.Key], func(wishType i18n.WishType) (int, string) {
					return wishType.Key, wishType.Name
				})
				locale.SharedWishes = lo.SliceToMap(wishes.SharedWishes[lang.Key], func(wishType i18n.WishType) (int, string) {
					return wishType.Key, wishType.Name
				})

				for index, entry := range entries {
					switch index.MenuID {
					case wiki.CharacterArchive.ID:
						locale.Characters[en[index].Name] = entry.Name
						locale.CharactersInverse[entry.Name] = en[index].Name
					case wiki.Weapons.ID:
						locale.Weapons[en[index].Name] = entry.Name
						locale.WeaponsInverse[entry.Name] = en[index].Name
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

func GetAllEntries() (map[i18n.Language]map[wiki.EntryIndex]wiki.Entry, error) {
	var (
		mutex  sync.Mutex
		wg     sync.WaitGroup
		errors error
	)
	results := make(map[i18n.Language]map[wiki.EntryIndex]wiki.Entry)

	wg.Add(len(i18n.Languages))

	for _, lang := range i18n.Languages {
		go func(lang i18n.Language) {
			defer wg.Done()

			w := wiki.New(lang)

			menusEntries, err := w.GetMenusEntries(wiki.CharacterArchive.ID, wiki.Weapons.ID)

			mutex.Lock()
			defer mutex.Unlock()

			results[lang] = menusEntries
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		}(lang)
	}

	wg.Wait()

	return results, errors
}
