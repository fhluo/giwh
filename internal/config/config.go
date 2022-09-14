package config

import (
	"bytes"
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/fhluo/giwh/pkg/pipeline"
	"github.com/fhluo/giwh/pkg/repository"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var (
	Dir = filepath.Join(os.Getenv("LOCALAPPDATA"), "giwh")

	Path                  = filepath.Join(Dir, "config.toml")
	WishHistoryPath       = filepath.Join(Dir, "wish_history.json")
	WishHistoryBackupPath = filepath.Join(Dir, "wish_history_backup.json")

	WishHistory pipeline.Pipeline

	config      = mustLoadConfig()
	GetLanguage = func() string { return config.Language }
	SetLanguage = func(lang string) { config.Language = lang }
	Save        = func() error { return config.Save() }
)

func init() {
	_ = os.MkdirAll(Dir, 0666)

	items, err := repository.LoadIfExits(WishHistoryPath)
	if err != nil {
		log.Fatalln(err)
	}

	WishHistory, err = pipeline.New(items)
	if err != nil {
		log.Fatalln(err)
	}
}

func mustLoadConfig() *Config {
	config, err := Load(Path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			config = new(Config)
		} else {
			log.Fatalf("fail to open config file: %s\n", err)
		}
	}
	return config
}

func SaveWishHistory() error {
	items, err := repository.LoadIfExits(WishHistoryPath)
	if err != nil {
		return err
	}

	if pipeline.ItemsEqual(items, WishHistory.Items()) {
		return nil
	}

	_ = os.Rename(WishHistoryPath, WishHistoryBackupPath)

	WishHistory.SortByIDDescending()
	return repository.Save(WishHistoryPath, WishHistory.Unique().Items())
}

type Config struct {
	Language string `toml:"language"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	return cfg, toml.Unmarshal(data, cfg)
}

func (config *Config) Save() error {
	buf := new(bytes.Buffer)
	e := toml.NewEncoder(buf)
	e.Indent = ""

	err := e.Encode(config)
	data := buf.Bytes()
	if err != nil {
		return err
	}

	return os.WriteFile(Path, data, 0666)
}
