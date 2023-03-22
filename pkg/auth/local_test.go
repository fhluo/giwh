package auth

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

	infos := GetAllInfos()
	for _, info := range infos {
		fmt.Println(info.BaseURL, info.AuthKeyVer, info.Lang)
		fmt.Println(info.AuthKey)
	}
}
