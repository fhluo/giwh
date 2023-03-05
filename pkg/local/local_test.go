package local

import (
	"fmt"
	"golang.org/x/exp/slog"
	"os"
	"testing"
)

func init() {
	slog.SetDefault(slog.New(slog.HandlerOptions{Level: slog.LevelDebug}.NewTextHandler(os.Stderr)))
}

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.SkipNow()
	}
}

func TestSearchCacheDataPaths(t *testing.T) {
	skipCI(t)
	fmt.Println(GetCacheDataPaths())
}

func TestGetAuths(t *testing.T) {
	skipCI(t)
	for _, a := range GetAuths() {
		fmt.Println(a.Domain, a.AuthKeyVer, a.Lang)
		fmt.Println(a.AuthKey)
		fmt.Println(a.GetGachaLogURL())
	}
}
