package repository

import (
	"encoding/json"
	"errors"
	"github.com/fhluo/giwh/pkg/api"
	"io/fs"
	"os"
)

type Item struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Repository interface {
	GetUIDs() []int
	GetProgress(uid int, wishType api.SharedWishType, rarity api.Rarity) int
	GetPulls(uid int, wishType api.SharedWishType, id int64) int
	GetItems(uid int, wishType api.SharedWishType, rarity api.Rarity) []Item
}

func Load(filename string) ([]*api.Item, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var items []*api.Item
	return items, json.Unmarshal(data, &items)
}

func LoadIfExits(filename string) ([]*api.Item, error) {
	items, err := Load(filename)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		return nil, nil
	}

	return items, nil
}

func Save(filename string, items []*api.Item) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0666)
}
