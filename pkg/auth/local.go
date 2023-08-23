package auth

import (
	"errors"
	"github.com/samber/lo"
	"io/fs"
	"log/slog"
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

func FindAllURL(data []byte) []string {
	result := urlRE.FindAll(data, -1)
	urls := make([]string, len(result))
	for i, url := range result {
		urls[i] = string(url)
	}
	return urls
}

func ReadInfos(path string) (infos []Info, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	urls := FindAllURL(data)
	infos = make([]Info, 0, len(urls))
	for _, url := range urls {
		info, err := FromURL(url)
		if err != nil {
			slog.Debug(err.Error(), "url", url)
			continue
		}
		infos = append(infos, info)
	}

	return infos, nil
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
