package config

import (
	"github.com/fhluo/giwh/i18n"
	"github.com/fhluo/giwh/pkg/repository/primitive"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	Dir = filepath.Join(os.Getenv("LOCALAPPDATA"), "giwh")

	Path            = filepath.Join(Dir, "config.toml")
	WishHistoryPath = filepath.Join(Dir, "wish_history.json")

	Repository *primitive.Repository

	Language = NewItem("language", i18n.Default().Tag().String())
)

func init() {
	_ = os.MkdirAll(Dir, 0666)

	viper.AddConfigPath(Dir)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err = viper.WriteConfigAs(Path); err != nil {
				slog.Warn("failed to write config", "path", Path)
			}
		}
	}

	var err error

	Repository, err = primitive.Load(WishHistoryPath)
	if err != nil {
		log.Fatalln(err)
	}
}

var mutex sync.Mutex

type Item[T any] struct {
	Key          string
	DefaultValue T
}

func NewItem[T any](key string, defaultValue T) Item[T] {
	viper.SetDefault(key, defaultValue)
	return Item[T]{
		Key:          key,
		DefaultValue: defaultValue,
	}
}

func (item Item[T]) Get() T {
	mutex.Lock()
	defer mutex.Unlock()

	return viper.Get(item.Key).(T)
}

func (item Item[T]) Set(value T) {
	mutex.Lock()
	defer mutex.Unlock()

	viper.Set(item.Key, value)
}

func Save() {
	if err := viper.WriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			if err = viper.WriteConfigAs(Path); err != nil {
				slog.Warn("failed to write config", "path", Path)
			}
		} else {
			slog.Error(err.Error(), nil)
		}
	}
}
