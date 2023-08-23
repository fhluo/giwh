package auth

import (
	"fmt"
	"log/slog"
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

func TestFindAllURL(t *testing.T) {
	skipCI(t)

	cacheDataPaths := GetCacheDataPaths()
	if len(cacheDataPaths) == 0 {
		return
	}

	data, err := os.ReadFile(cacheDataPaths[0])
	if err != nil {
		return
	}

	urls := FindAllURL(data)
	for _, url := range urls {
		fmt.Println(url)
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
