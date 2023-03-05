package repository

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/fhluo/giwh/pkg/local"
	"github.com/fhluo/giwh/pkg/wish"
	"github.com/fhluo/giwh/pkg/wish/pipeline"
	"golang.org/x/exp/slog"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func LoadItems(filename string) (items []wish.Item, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	return items, sonic.Unmarshal(data, &items)
}

func LoadItemsIfExits(filename string) ([]wish.Item, error) {
	items, err := LoadItems(filename)

	switch {
	case err == nil:
	case errors.Is(err, fs.ErrNotExist):
		return nil, nil
	default:
		return nil, err
	}

	return items, nil
}

func SaveItems(filename string, items []wish.Item) error {
	data, err := sonic.Marshal(items)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0666)
}

func BackupAndSaveItems(filename string, items []wish.Item) error {
	dir, base := filepath.Split(filename)
	ext := filepath.Ext(base)

	if err := os.Rename(filename, filepath.Join(dir, strings.TrimSuffix(base, ext)+"_backup"+ext)); err != nil {
		slog.Warn("backup failed", "err", err.Error())
	}

	return SaveItems(filename, items)
}

func UpdateItems(filename string, handlers ...func(item wish.Item)) ([]wish.Item, error) {
	items, err := LoadItemsIfExits(filename)
	if err != nil {
		return nil, err
	}

	auths := local.GetAuths()
	if len(auths) == 0 {
		return nil, fmt.Errorf("")
	}

	p := pipeline.New(items)
	length := p.Len()

	for _, sharedWish := range wish.SharedWishes {
		ctx := wish.New(auths[len(auths)-1]).WishType(sharedWish).Size(10)

		for {
			items, err = ctx.Fetch()
			if err != nil {
				return nil, err
			}

			if p.ContainsAny(items...) {
				for _, handle := range handlers {
					for _, item := range items {
						if !p.Contains(item) {
							handle(item)
						}
					}
				}
				p.Append(items...)
				break
			}

			for _, handle := range handlers {
				for _, item := range items {
					handle(item)
				}
			}
			p.Append(items...)
		}
	}

	if p.Len() != length {
		if err = BackupAndSaveItems(filename, p.Items()); err != nil {
			return p.Items(), err
		}
	}

	return p.Items(), nil
}
