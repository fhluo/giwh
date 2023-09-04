package config

import (
	"bytes"
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/fhluo/giwh/pkg/fetcher"
	"github.com/fhluo/giwh/pkg/util"
	"github.com/fhluo/giwh/pkg/wish"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
)

var (
	Dir = filepath.Join(os.Getenv("LOCALAPPDATA"), "giwh")

	Path                  = filepath.Join(Dir, "config.toml")
	WishHistoryPath       = filepath.Join(Dir, "wish_history.json")
	WishHistoryBackupPath = filepath.Join(Dir, "wish_history_backup.json")

	WishHistory wish.Items

	config         = mustLoadConfig()
	GetAuthInfo    = config.GetAuthInfo
	UpdateAuthInfo = config.UpdateAuthInfo
	GetLanguage    = func() string { return config.Language }
	SetLanguage    = func(lang string) { config.Language = lang }
	Save           = func() error { return config.Save() }
)

func init() {
	_ = os.MkdirAll(Dir, 0666)

	var err error

	WishHistory, err = wish.LoadItemsIfExits(WishHistoryPath)
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
	items, err := wish.LoadItemsIfExits(WishHistoryPath)
	if err != nil {
		return err
	}

	if items.Equal(WishHistory) {
		return nil
	}

	_ = os.Rename(WishHistoryPath, WishHistoryBackupPath)
	return WishHistory.Unique().Save(WishHistoryPath)
}

type Config struct {
	Language  string             `toml:"language"`
	AuthInfos []fetcher.AuthInfo `toml:"auth_infos"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	return cfg, toml.Unmarshal(data, cfg)
}

func (config *Config) GetAuthInfo(uid string) (fetcher.AuthInfo, bool) {
	return util.Find(config.AuthInfos, func(authInfo fetcher.AuthInfo) bool {
		return authInfo.UID == uid
	})
}

func (config *Config) UpdateAuthInfo(authInfo fetcher.AuthInfo) {
	i := slices.IndexFunc(config.AuthInfos, func(info fetcher.AuthInfo) bool {
		return info.UID == authInfo.UID
	})
	if i < 0 {
		config.AuthInfos = append(config.AuthInfos, authInfo)
	} else {
		config.AuthInfos[i] = authInfo
	}
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
