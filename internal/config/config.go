package config

import (
	"bytes"
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/fhluo/giwh/fetcher"
	"github.com/fhluo/giwh/wh"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
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

	WishHistory wh.WishHistory

	config *Config
)

func init() {
	_ = os.MkdirAll(Dir, 0666)

	var err error

	if config, err = Load(Path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			config = new(Config)
		} else {
			log.Fatalf("fail to open config file: %s\n", err)
		}
	}

	WishHistory, err = wh.LoadWishHistoryIfExits(WishHistoryPath)
	if err != nil {
		log.Fatalln(err)
	}
}

func SaveWishHistory() error {
	items, err := wh.LoadWishHistoryIfExits(WishHistoryPath)
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
	return lo.Find(config.AuthInfos, func(authInfo fetcher.AuthInfo) bool {
		return authInfo.UID == uid
	})
}

func GetAuthInfo(uid string) (fetcher.AuthInfo, bool) {
	return config.GetAuthInfo(uid)
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

func UpdateAuthInfo(authInfo fetcher.AuthInfo) {
	config.UpdateAuthInfo(authInfo)
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

func GetLanguage() string {
	return config.Language
}

func SetLanguage(lang string) {
	config.Language = lang
}

func Save() error {
	return config.Save()
}
