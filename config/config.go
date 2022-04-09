package config

import (
	"github.com/fhluo/giwh/util"
	"github.com/fhluo/giwh/wh"
	"github.com/samber/lo"
	"log"
	"os"
	"path/filepath"
)

var (
	Path = filepath.Join(os.Getenv("LOCALAPPDATA"), "giwh")

	CachedItems wh.Items
)

func init() {
	_ = os.MkdirAll(Path, 0666)

	var err error

	CachedItems, err = util.LoadItemsIfExits(filepath.Join(Path, "cache.json"))
	if err != nil {
		log.Fatalln(err)
	}
	if len(CachedItems) != 0 {
		_ = CachedItems.Save(filepath.Join(Path, "cache_backup.json"))
	}
}

func SaveCache() error {
	CachedItems = lo.UniqBy(CachedItems, func(item wh.Item) int64 {
		return item.ID()
	})
	return CachedItems.Save(filepath.Join(Path, "cache.json"))
}
