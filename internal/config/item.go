package config

import (
	"github.com/spf13/viper"
	"sync"
)

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
