//go:build windows

package db

import (
	"os"
	"path/filepath"
)

func DefaultPath() string {
	return filepath.Join(os.Getenv("LOCALAPPDATA"), "giwh", "logs.sqlite3")
}
