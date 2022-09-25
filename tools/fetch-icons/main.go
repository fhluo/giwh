package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func init() {
	log.SetFlags(0)
}

type FilterValue struct {
	Values []string `json:"values"`
}
type Entry struct {
	EntryPageID  int                     `json:"entry_page_id,string"`
	Name         string                  `json:"name"`
	IconURL      string                  `json:"icon_url"`
	DisplayField map[string]any          `json:"display_field"`
	FilterValues map[string]*FilterValue `json:"filter_values"`
}

var NonAlphanumeric = regexp.MustCompile(`\W`)

func (entry *Entry) VarName() string {
	return NonAlphanumeric.ReplaceAllString(entry.Name, "")
}

func (entry *Entry) Filename() string {
	return strings.ReplaceAll(entry.Name, `"`, "") + path.Ext(entry.IconURL)
}

type Data struct {
	List  []*Entry `json:"list"`
	Total string   `json:"total"`
}

type Result struct {
	RetCode int    `json:"retcode"`
	Message string `json:"message"`
	Data    *Data  `json:"data"`
}

func LoadEntries(filename string) ([]*Entry, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result Result
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if result.RetCode != 0 {
		return nil, fmt.Errorf(result.Message)
	}

	return result.Data.List, nil
}

func LoadEntriesFromFiles(filenames ...string) ([]*Entry, error) {
	var entries []*Entry
	for _, filename := range filenames {
		r, err := LoadEntries(filename)
		if err != nil {
			return nil, err
		}
		entries = append(entries, r...)
	}

	slices.SortFunc(entries, func(a *Entry, b *Entry) bool {
		return a.EntryPageID < b.EntryPageID
	})

	return entries, nil
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

func DownloadAll(entries []*Entry, path string) {
	var wg sync.WaitGroup
	wg.Add(len(entries))

	for _, entry := range entries {
		go func(entry *Entry) {
			defer wg.Done()

			if entry.IconURL == "" {
				fmt.Printf("Icon not available: %s\n", entry.Name)
				return
			}

			dst := filepath.Join(path, entry.Filename())
			if _, err := os.Stat(dst); err == nil {
				fmt.Printf("File already exists: %s\n", dst)
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

func DownloadCharactersIcons(pattern string, path string) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalln(err)
	}

	characters, err := LoadEntriesFromFiles(filenames...)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d characters' icons:\n", len(characters))

	if err = os.MkdirAll(path, 0666); err != nil {
		log.Fatalln(err)
	}
	DownloadAll(characters, path)

	fmt.Println()
}

func DownloadWeaponsIcons(pattern string, path string) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalln(err)
	}

	weapons, err := LoadEntriesFromFiles(filenames...)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d weapons' icons:\n", len(weapons))

	if err = os.MkdirAll(path, 0666); err != nil {
		log.Fatalln(err)
	}
	DownloadAll(weapons, path)

	fmt.Println()
}

func main() {
	var c, w, o string
	flag.StringVar(&c, "c", "", "")
	flag.StringVar(&w, "w", "", "")
	flag.StringVar(&o, "o", "", "")
	flag.Parse()

	DownloadCharactersIcons(c, filepath.Join(o, "characters"))
	DownloadWeaponsIcons(w, filepath.Join(o, "weapons"))
}