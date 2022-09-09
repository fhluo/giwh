package auth

import (
	"fmt"
	"github.com/samber/lo"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Region struct {
	Name          string
	OutputLogPath string
	APIBaseURL    string
}

var (
	CN = Region{
		Name:          "CN",
		OutputLogPath: filepath.Join(os.Getenv("USERPROFILE"), `\AppData\LocalLow\miHoYo`, `原神`, `output_log.txt`),
		APIBaseURL:    "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog",
	}
	OS = Region{
		Name:          "OS",
		OutputLogPath: filepath.Join(os.Getenv("USERPROFILE"), `\AppData\LocalLow\miHoYo`, `Genshin Impact`, `output_log.txt`),
		APIBaseURL:    "https://hk4e-api-os.hoyoverse.com/event/gacha_info/api/getGachaLog",
	}
)

var dataPathRE = regexp.MustCompile(`.:.*?[/|\\](YuanShen_Data|GenshinImpact_Data)`)

func findDataPath(outputLog []byte) string {
	return filepath.FromSlash(string(dataPathRE.Find(outputLog)))
}

func (r Region) GetCacheDataPath() (string, error) {
	data, err := os.ReadFile(r.OutputLogPath)
	if err != nil {
		return "", err
	}

	path := findDataPath(data)
	if path == "" {
		return "", fmt.Errorf("failed to find data path from output_log.txt")
	}

	return filepath.Join(path, `webCaches\Cache\Cache_Data\data_2`), nil
}

var urlRE = regexp.MustCompile(`https?://[-a-zA-Z0-9./=&?_%]+`)

func findAllURLs(data []byte) []string {
	return lo.Map(urlRE.FindAll(data, -1), func(url []byte, _ int) string {
		return string(url)
	})
}

func (r Region) GetURLsFromCacheData() ([]string, error) {
	path, err := r.GetCacheDataPath()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return findAllURLs(data), nil
}

func (r Region) GetAPIURL() (string, error) {
	urls, err := r.GetURLsFromCacheData()
	if err != nil {
		return "", err
	}

	apiURL, _, ok := lo.FindLastIndexOf(urls, func(rawURL string) bool {
		u, err := url.Parse(rawURL)
		if err != nil {
			return false
		}

		return strings.HasSuffix(u.Path, "api/getGachaLog")
	})

	if ok {
		return apiURL, nil
	}

	rawURL, _, ok := lo.FindLastIndexOf(urls, func(rawURL string) bool {
		u, err := url.Parse(rawURL)
		if err != nil {
			return false
		}

		if u.Query().Has("authkey") {
			return true
		}

		return false
	})

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return r.APIBaseURL + "?" + u.RawQuery, nil
}
