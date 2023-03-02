package icons

import (
	"fmt"
	"github.com/fhluo/giwh/i18n"
	"github.com/fhluo/giwh/pkg/wiki"
	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func NewCmd() *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "icons",
		Short: "Download icons",
		Run: func(cmd *cobra.Command, args []string) {
			DownloadCharactersIcons(outputDir)
			DownloadWeaponsIcons(outputDir)
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "images", "output directory")

	return cmd
}

func clean(entries []wiki.Entry) []wiki.Entry {
	entries = lo.Filter(entries, func(entry wiki.Entry, _ int) bool {
		return entry.IconURL != ""
	})

	entries = lo.UniqBy(entries, func(entry wiki.Entry) string {
		return entry.Name
	})

	slices.SortFunc(entries, func(a wiki.Entry, b wiki.Entry) bool {
		return a.EntryPageID < b.EntryPageID
	})

	return entries
}

func Download(url string, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0666)
}

func DownloadAll(entries []wiki.Entry, path string) {
	var wg sync.WaitGroup
	wg.Add(len(entries))

	for _, entry := range entries {
		go func(entry wiki.Entry) {
			defer wg.Done()

			if entry.IconURL == "" {
				slog.Warn("Icon not available", "name", entry.Name)
				return
			}

			dst := filepath.Join(path, entry.Filename())
			if _, err := os.Stat(dst); err == nil {
				return
			}

			err := Download(entry.IconURL, dst)
			if err != nil {
				log.Println(err)
				return
			}

			slog.Info("Downloaded", "file", dst)
		}(entry)
	}

	wg.Wait()
}

func DownloadCharactersIcons(path string) {
	enWiki := wiki.New(i18n.English)
	characters, err := enWiki.GetEntries(wiki.CharacterArchive.ID)
	if err != nil {
		slog.Error(err.Error(), nil)
		os.Exit(1)
	}
	characters = clean(characters)

	slog.Info("Download characters' icons", "length", len(characters))

	if err = os.MkdirAll(filepath.Join(path, "characters"), 0666); err != nil {
		slog.Error(err.Error(), nil)
		os.Exit(1)
	}
	DownloadAll(characters, filepath.Join(path, "characters"))

	fmt.Println()

	index := make(map[string]string)
	for _, character := range characters {
		index[character.Name] = character.Filename()
	}
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		slog.Error(err.Error(), nil)
		os.Exit(1)
	}

	err = os.WriteFile(filepath.Join(path, "characters.json"), data, 0666)
	if err != nil {
		slog.Error(err.Error(), nil)
		os.Exit(1)
	}
}

func DownloadWeaponsIcons(path string) {
	enWiki := wiki.New(i18n.English)
	weapons, err := enWiki.GetEntries(wiki.Weapons.ID)
	if err != nil {
		slog.Error(err.Error(), nil)
		os.Exit(1)
	}
	slog.Info("Download weapons' icons", "length", len(weapons))

	weapons = clean(weapons)

	if err = os.MkdirAll(filepath.Join(path, "weapons"), 0666); err != nil {
		slog.Error(err.Error(), nil)
		os.Exit(1)
	}
	DownloadAll(weapons, filepath.Join(path, "weapons"))

	fmt.Println()

	index := make(map[string]string)
	for _, weapon := range weapons {
		index[weapon.Name] = weapon.Filename()
	}
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		slog.Error(err.Error(), nil)
		os.Exit(1)
	}

	err = os.WriteFile(filepath.Join(path, "weapons.json"), data, 0666)
	if err != nil {
		slog.Error(err.Error(), nil)
		os.Exit(1)
	}
}
