package auth

import (
	"errors"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

var (
	OutputLogPaths = []string{
		filepath.Join(os.Getenv("USERPROFILE"), `\AppData\LocalLow\miHoYo`, `原神`, `output_log.txt`),
		filepath.Join(os.Getenv("USERPROFILE"), `\AppData\LocalLow\miHoYo`, `Genshin Impact`, `output_log.txt`),
	}
	DataPathRE = regexp.MustCompile(`[a-zA-Z]:[/|\\].*?[/|\\](YuanShen_Data|GenshinImpact_Data)[/|\\]`)
)

func GetCacheDataPaths() (cacheDataPaths []string) {
	for _, outputLogPath := range OutputLogPaths {
		data, err := os.ReadFile(outputLogPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				slog.Debug("file does not exist", "path", outputLogPath)
			} else {
				slog.Error(err.Error())
			}
			continue
		}
		cacheDataPaths = append(cacheDataPaths, filepath.Clean(
			filepath.Join(string(DataPathRE.Find(data)), `webCaches\Cache\Cache_Data\data_2`),
		))
	}
	return
}

var urlRE = regexp.MustCompile(`https?://[-a-zA-Z0-9.:/=&?_%+]+`)

func ReadInfos(path string) (infos []Info, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	return lo.FilterMap(urlRE.FindAll(data, -1), func(url []byte, _ int) (info Info, ok bool) {
		info, err = FromURL(string(url))
		if err != nil {
			slog.Debug(err.Error(), "url", string(url))
			return
		}
		ok = true
		return
	}), nil
}

func GetAllInfos() (infos []Info) {
	cacheDataPaths := GetCacheDataPaths()
	if len(cacheDataPaths) == 0 {
		return
	}

	for _, cacheDataPath := range cacheDataPaths {
		result, err := ReadInfos(cacheDataPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				slog.Debug("file does not exist", "path", cacheDataPath)
			} else {
				slog.Error(err.Error())
			}
			continue
		}
		infos = append(infos, result...)
	}

	return lo.Uniq(infos)
}
