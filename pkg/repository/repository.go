package repository

import (
	"encoding/json"
	"errors"
	"github.com/fhluo/giwh/pkg/api"
	"io/fs"
	"os"
)

type Item struct {
	*api.Item
	Pulls int    `json:"pulls"`
	Icon  string `json:"icon"`
}

type Repository interface {
	GetUIDs() []int
	Get5StarProgress(uid int, wishType api.SharedWishType) int
	Get4StarProgress(uid int, wishType api.SharedWishType) int
	Get5Stars(uid int, wishType api.SharedWishType) []Item
	Get4Stars(uid int, wishType api.SharedWishType) []Item
	AddItems(items []*api.Item)
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
