//go:build windows

package stores

import (
	"os"
	"path/filepath"
)

func DefaultPath() string {
	return filepath.Join(os.Getenv("LocalAppData"), "giwh", "wish_history.json")
}
