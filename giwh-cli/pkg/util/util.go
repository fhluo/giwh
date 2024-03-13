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
