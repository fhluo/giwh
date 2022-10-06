package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fhluo/giwh/pkg/api"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func clean(entries []*api.Entry) []*api.Entry {
	entries = lo.Filter(entries, func(entry *api.Entry, _ int) bool {
		return entry.IconURL != ""
	})

	entries = lo.UniqBy(entries, func(entry *api.Entry) string {
		return entry.Name
	})

	slices.SortFunc(entries, func(a *api.Entry, b *api.Entry) bool {
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

func DownloadAll(entries []*api.Entry, path string) {
	var wg sync.WaitGroup
	wg.Add(len(entries))

	for _, entry := range entries {
		go func(entry *api.Entry) {
			defer wg.Done()

			if entry.IconURL == "" {
				fmt.Printf("Icon not available: %s\n", entry.Name)
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

			fmt.Printf("Downloaded: %s\n", dst)
		}(entry)
	}

	wg.Wait()
}

func DownloadCharactersIcons(path string) {
	characters, err := api.GetAllEntries(api.CharactersMenu)
	if err != nil {
		log.Fatalln(err)
	}
	characters = clean(characters)

	fmt.Printf("%d characters' icons.\n", len(characters))

	if err = os.MkdirAll(path, 0666); err != nil {
		log.Fatalln(err)
	}
	DownloadAll(characters, path)

	fmt.Println()

	index := make(map[string]string)
	for _, character := range characters {
		index[character.Name] = character.Filename()
	}
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(filepath.Join(o, "characters.json"), data, 0666)
	if err != nil {
		log.Fatalln(err)
	}
}

func DownloadWeaponsIcons(path string) {
	weapons, err := api.GetAllEntries(api.WeaponsMenu)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d weapons' icons.\n", len(weapons))
	weapons = clean(weapons)

	if err = os.MkdirAll(path, 0666); err != nil {
		log.Fatalln(err)
	}
	DownloadAll(weapons, path)

	fmt.Println()

	index := make(map[string]string)
	for _, weapon := range weapons {
		index[weapon.Name] = weapon.Filename()
	}
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(filepath.Join(o, "weapons.json"), data, 0666)
	if err != nil {
		log.Fatalln(err)
	}
}

var o string

func init() {
	log.SetFlags(0)

	flag.StringVar(&o, "o", "", "")
	flag.Parse()
}

//go:generate go run . -o ../../assets/images

func main() {
	DownloadCharactersIcons(filepath.Join(o, "characters"))
	DownloadWeaponsIcons(filepath.Join(o, "weapons"))
}
