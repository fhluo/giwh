package util

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
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

		return Map(infos, func(i *info) string {
			return i.name
		}), nil
	}
}

func ExpandPaths(paths ...string) (result []string, errs error) {
	result = make([]string, 0, len(paths))

	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		result = append(result, matches...)
	}

	return result, errs
}

func Map[T, T_ any](items []T, f func(T) T_) []T_ {
	result := make([]T_, 0, len(items))
	for _, item := range items {
		result = append(result, f(item))
	}
	return result
}

func Filter[T any](items []T, f func(T) bool) (result []T) {
	for _, item := range items {
		if f(item) {
			result = append(result, item)
		}
	}
	return result
}

func Find[T any](items []T, f func(T) bool) (r T, b bool) {
	for _, item := range items {
		if f(item) {
			return item, true
		}
	}
	return r, false
}

func Unique[T any, K comparable](items []T, f func(T) K) []T {
	exists := make(map[K]bool)
	result := make([]T, 0, len(items))

	for _, item := range items {
		if key := f(item); !exists[key] {
			exists[key] = true
			result = append(result, item)
		}
	}

	return result
}
