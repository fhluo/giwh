package repository

import (
	"encoding/json"
	"errors"
	"github.com/fhluo/giwh/pkg/api"
	"io/fs"
	"os"
)

type Repository interface {
	GetUIDs() []string
	GetProgress(uid string, wishType string, rarity string) int
	GetPulls(uid string, wishType string, id int64) int
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
