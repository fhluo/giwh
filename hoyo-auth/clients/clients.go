//go:build windows

package clients

import (
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"time"

	"github.com/fhluo/giwh/hoyo-auth/auths"

	"github.com/samber/lo"
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

// Latest returns the latest auth info
func Latest() *auths.Auth {
	type Pair struct {
		Genshin
		time time.Time
	}

	pairs := lo.FilterMap([]Genshin{
		GenshinCN(),
		GenshinGlobal(),
	}, func(genshin Genshin, _ int) (Pair, bool) {
		info, err := os.Stat(genshin.outputLogPath())
		if err != nil {
			return Pair{}, false
		}
		return Pair{
			Genshin: genshin,
			time:    info.ModTime(),
		}, true
	})

	if len(pairs) == 0 {
		return nil
	}

	r := slices.MaxFunc(pairs, func(a Pair, b Pair) int {
		return a.time.Compare(b.time)
	}).Auths()

	if len(r) == 0 {
		return nil
	}

	return r[len(r)-1]
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
		path string
		time time.Time
	}

	pairs := lo.FilterMap(g.cacheDataPaths(), func(path string, _ int) (Pair, bool) {
		info, err := os.Stat(path)
		if err != nil {
			slog.Warn("failed to get file info", "path", path, "err", err)
			return Pair{}, false
		}

		return Pair{
			path: path,
			time: info.ModTime(),
		}, true
	})

	if len(pairs) == 0 {
		slog.Warn("failed to find latest cache data path")
		return ""
	}

	slices.SortFunc(pairs, func(a, b Pair) int {
		return a.time.Compare(b.time)
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

	i := re.SubexpIndex("url")
	return lo.Map(matches, func(match [][]byte, _ int) string {
		return string(match[i])
	})
}

func (g Genshin) Auths() []*auths.Auth {
	r := lo.FilterMap(g.urlsInCacheData(), func(url string, _ int) (*auths.Auth, bool) {
		auth, err := auths.New(url)
		if err != nil {
			slog.Debug("failed to get auth info from url", "url", url, "err", err)
			return nil, false
		}
		return auth, true
	})
	return lo.Uniq(r)
}
