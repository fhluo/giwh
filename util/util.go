package util

import (
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
	"os"
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
