package util

import (
	"github.com/hashicorp/go-multierror"
	"path/filepath"
)

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
