package util

import (
	"fmt"
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
	names, err := SortFiles(names...)
	if err != nil {
		return "", err
	}
	return names[0], nil
}

// SortFiles sorts the files by modification time from newest to oldest. Files that fail to get file info are ignored.
func SortFiles(names ...string) ([]string, error) {
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
		if errs == nil {
			return nil, fmt.Errorf("not found")
		}
		return nil, fmt.Errorf("not found: %v", errs)
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
	var errs error

	return lo.FlatMap(paths, func(path string, _ int) []string {
		matches, err := filepath.Glob(path)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		return matches
	}), errs
}
