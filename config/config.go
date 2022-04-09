package config

import (
	"errors"
	"fmt"
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

	Path            = filepath.Join(Dir, "config.toml")
	CachePath       = filepath.Join(Dir, "cache.json")
	CacheBackupPath = filepath.Join(Dir, "cache_backup.json")

	CachedItems wh.Items

	config *Config
)

func init() {
	_ = os.MkdirAll(Dir, 0666)

	var err error

	if config, err = Load(Path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			config = new(Config)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "fail to open config file: %s\n", err)
			os.Exit(1)
		}
	}

	CachedItems, err = wh.LoadItemsIfExits(CachePath)
	if err != nil {
		log.Fatalln(err)
	}
	if len(CachedItems) != 0 {
		_ = CachedItems.Unique().Save(CacheBackupPath)
	}
}

func SaveCache() error {
	return CachedItems.Unique().Save(CachePath)
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
