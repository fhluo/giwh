package config

import (
	"errors"
	"github.com/fhluo/giwh/wh"
	"github.com/pelletier/go-toml/v2"
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

	WishHistory wh.Items

	config *Config

	logger = log.New(os.Stderr, "", 0)
)

func init() {
	_ = os.MkdirAll(Dir, 0666)

	var err error

	if config, err = Load(Path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			config = new(Config)
		} else {
			logger.Fatalf("fail to open config file: %s\n", err)
		}
	}

	WishHistory, err = wh.LoadItemsIfExits(WishHistoryPath)
	if err != nil {
		logger.Fatalln(err)
	}
}

func SaveWishHistory() error {
	items, err := wh.LoadItemsIfExits(WishHistoryPath)
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
	AuthInfos []wh.AuthInfo `toml:"auth_infos"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	return cfg, toml.Unmarshal(data, cfg)
}

func (config *Config) GetAuthInfo(uid string) (wh.AuthInfo, bool) {
	return lo.Find(config.AuthInfos, func(authInfo wh.AuthInfo) bool {
		return authInfo.UID == uid
	})
}

func GetAuthInfo(uid string) (wh.AuthInfo, bool) {
	return config.GetAuthInfo(uid)
}

func (config *Config) UpdateAuthInfo(authInfo wh.AuthInfo) {
	i := slices.IndexFunc(config.AuthInfos, func(info wh.AuthInfo) bool {
		return info.UID == authInfo.UID
	})
	if i < 0 {
		config.AuthInfos = append(config.AuthInfos, authInfo)
	} else {
		config.AuthInfos[i] = authInfo
	}
}

func UpdateAuthInfo(authInfo wh.AuthInfo) {
	config.UpdateAuthInfo(authInfo)
}

func (config *Config) Save() error {
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(Path, data, 0666)
}

func Save() error {
	return config.Save()
}
