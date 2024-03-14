//go:build windows

package hyauth

import (
	"github.com/samber/lo"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"time"
)

func hoyoAppDataPath() string {
	return filepath.Join(os.Getenv("UserProfile"), `\AppData\LocalLow\miHoYo`)
}

type Genshin struct {
	dataPath string
}

func GenshinCN() Genshin {
	return Genshin{
		dataPath: "原神",
	}
}

func GenshinGlobal() Genshin {
	return Genshin{
		dataPath: "Genshin Impact",
	}
}

func (g Genshin) outputLogPath() string {
	return filepath.Join(hoyoAppDataPath(), g.dataPath, "output_log.txt")
}

func (g Genshin) Executed() bool {
	if _, err := os.Stat(g.outputLogPath()); err == nil {
		return true
	}
	return false
}

func (g Genshin) outputLog() []byte {
	data, err := os.ReadFile(g.outputLogPath())
	if err != nil {
		slog.Warn("failed to read file", "path", g.outputLogPath(), "err", err)
		return nil
	}
	return data
}

func (g Genshin) programDataPath() string {
	re := regexp.MustCompile(`[a-zA-Z]:[\\/].*?[\\/](YuanShen_Data|GenshinImpact_Data)[\\/]`)
	match := re.Find(g.outputLog())
	if match == nil {
		slog.Warn("failed to find program data path", "path", g.outputLogPath())
		return ""
	}
	return filepath.Clean(string(match))
}

func (g Genshin) cacheDataPaths() []string {
	paths, err := filepath.Glob(filepath.Join(g.programDataPath(), `webCaches\*\Cache\Cache_Data\data_2`))
	if err != nil {
		panic(err)
	}
	return paths
}

func (g Genshin) latestCacheDataPath() string {
	type Pair struct {
		path    string
		modTime time.Time
	}

	cacheDataPaths := g.cacheDataPaths()
	pairs := make([]Pair, 0, len(cacheDataPaths))

	for _, path := range cacheDataPaths {
		info, err := os.Stat(path)
		if err != nil {
			slog.Warn("failed to get file info", "path", path, "err", err)
			continue
		}

		pairs = append(pairs, Pair{
			path:    path,
			modTime: info.ModTime(),
		})
	}

	if len(pairs) == 0 {
		slog.Warn("failed to find latest cache data path")
		return ""
	}

	slices.SortFunc(pairs, func(a, b Pair) int {
		switch {
		case a.modTime.Before(b.modTime):
			return -1
		case a.modTime.After(b.modTime):
			return 1
		default:
			return 0
		}
	})

	return pairs[len(pairs)-1].path
}

func (g Genshin) latestCacheData() []byte {
	data, err := os.ReadFile(g.latestCacheDataPath())
	if err != nil {
		slog.Warn("failed to read file", "path", g.outputLogPath(), "err", err)
		return nil
	}
	return data
}

func (g Genshin) urlsInCacheData() []string {
	re := regexp.MustCompile(`\x001/0/(?P<url>https://.*?)\x00`)
	matches := re.FindAllSubmatch(g.latestCacheData(), -1)

	urls := make([]string, 0, len(matches))
	i := re.SubexpIndex("url")

	for _, match := range matches {
		urls = append(urls, string(match[i]))
	}

	return urls
}

func (g Genshin) Auths() []*Auth {
	urls := g.urlsInCacheData()
	auths := make([]*Auth, 0, len(urls))

	for _, url := range urls {
		auth, err := New(url)
		if err != nil {
			slog.Debug("failed to get auth info from url", "url", url, "err", err)
			continue
		}
		auths = append(auths, auth)
	}

	return lo.Uniq(auths)
}
