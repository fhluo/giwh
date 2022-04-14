package util

import (
	"fmt"
	"github.com/fhluo/giwh/wh"
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type info struct {
	name string
	time time.Time
}

func FindLatest(names ...string) (string, error) {
	infos := make([]*info, 0, len(names))

	var errs error

	for _, name := range names {
		fi, err := os.Stat(name)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		infos = append(infos, &info{name: name, time: fi.ModTime()})
	}

	if len(infos) == 0 {
		return "", errs
	}

	latest := infos[0]
	for _, i := range infos[1:] {
		if i.time.After(latest.time) {
			latest = i
		}
	}

	return latest.name, nil
}

func SortExisting(names ...string) ([]string, error) {
	infos := make([]*info, 0, len(names))

	var errs error

	for _, name := range names {
		fi, err := os.Stat(name)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		infos = append(infos, &info{name: name, time: fi.ModTime()})
	}

	switch len(infos) {
	case 0:
		return nil, errs
	case 1:
		return []string{infos[0].name}, nil
	default:
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].time.After(infos[j].time)
		})
		return lo.Map(infos, func(i *info, _ int) string {
			return i.name
		}), nil
	}
}

func ExpandPaths(paths ...string) ([]string, error) {
	result := make([]string, 0, len(paths))

	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}
		result = append(result, matches...)
	}

	return result, nil
}

func FetchAllWishHistory(baseURL string, items wh.Items) (wh.Items, error) {
	visit := make(map[int64]bool)
	for _, item := range items {
		visit[item.ID()] = true
	}

	wishes := []wh.WishType{wh.CharacterEventWish, wh.WeaponEventWish, wh.StandardWish, wh.BeginnersWish}
	descriptions := map[wh.WishType]string{
		wh.CharacterEventWish: wh.CharacterEventWish.String() + " and " + wh.CharacterEventWish2.String(),
		wh.WeaponEventWish:    wh.WeaponEventWish.String(),
		wh.StandardWish:       wh.StandardWish.String(),
		wh.BeginnersWish:      wh.BeginnersWish.String(),
	}

	for i, wish := range wishes {
		fmt.Printf("Fetching the wish history of %s.\n", descriptions[wish])
		r, err := wh.NewFetcher(baseURL, wish, visit).FetchALL()
		if err != nil {
			return nil, err
		}

		items = append(items, r...)
		if i != len(wishes)-1 {
			time.Sleep(wh.DefaultInterval)
		}
	}

	return items, nil
}
