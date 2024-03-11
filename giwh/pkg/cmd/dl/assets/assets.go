package assets

import (
	"fmt"
	"github.com/fhluo/giwh/app/assets"
	"github.com/fhluo/giwh/i18n"
	"github.com/fhluo/giwh/pkg/wiki"
	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sync"
)

func NewCmd() *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "assets",
		Short: "Download assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			downloader := NewDownloader(outputDir)

			downloader.DownloadCharactersIcons()
			downloader.DownloadWeaponsIcons()

			return downloader.SaveAssetsInfo()
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", ".", "output directory")

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

type Downloader struct {
	wiki.Wiki
	dst        string
	characters []wiki.Entry
	weapons    []wiki.Entry
}

func NewDownloader(dst string) *Downloader {
	return &Downloader{
		Wiki: wiki.New(i18n.English),
		dst:  dst,
	}
}

func (d Downloader) Characters() []wiki.Entry {
	var err error

	if d.characters == nil {
		d.characters, err = d.GetEntries(wiki.CharacterArchive.ID)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		d.characters = clean(d.characters)
	}

	return d.characters
}

func (d Downloader) Weapons() []wiki.Entry {
	var err error

	if d.weapons == nil {
		d.weapons, err = d.GetEntries(wiki.Weapons.ID)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		d.weapons = clean(d.weapons)
	}

	return d.weapons
}

func (d Downloader) DownloadCharactersIcons() {
	slog.Info("Download characters' icons", "length", len(d.Characters()))

	_ = os.MkdirAll(filepath.Join(d.dst, "images", "characters"), 0666)
	DownloadAll(d.Characters(), filepath.Join(d.dst, "images", "characters"))

	fmt.Println()
}

func (d Downloader) DownloadWeaponsIcons() {
	slog.Info("Download weapons' icons", "length", len(d.Weapons()))

	_ = os.MkdirAll(filepath.Join(d.dst, "images", "weapons"), 0666)
	DownloadAll(d.Weapons(), filepath.Join(d.dst, "images", "weapons"))

	fmt.Println()
}

func (d Downloader) SaveAssetsInfo() error {
	info := assets.New()

	for _, character := range d.Characters() {
		info.Characters[character.Name] = path.Join("images", "characters", character.Filename())
	}

	for _, weapon := range d.Weapons() {
		info.Weapons[weapon.Name] = path.Join("images", "weapons", weapon.Filename())
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(d.dst, "assets.json"), data, 0666)
}
