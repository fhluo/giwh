package api

import (
	"errors"
	"github.com/samber/lo"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
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

var (
	ErrDataPathNotFound = errors.New("data path could not be found from output_log.txt")
	ErrURLNotFound      = errors.New("urls containing the following query parameters could not be found: authkey_ver, authkey and lang")
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
		return "", ErrDataPathNotFound
	}

	return filepath.Join(path, `webCaches\Cache\Cache_Data\data_2`), nil
}

var urlRE = regexp.MustCompile(`https?://[-a-zA-Z0-9.:/=&?_%+]+`)

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

func (r Region) GetValidURLs() ([]*url.URL, error) {
	urls, err := r.GetURLsFromCacheData()
	if err != nil {
		return nil, err
	}

	result := lo.FilterMap(urls, func(rawURL string, _ int) (*url.URL, bool) {
		u, err := url.Parse(rawURL)
		if err != nil {
			return nil, false
		}
		return u, u.Query().Has("authkey_ver") && u.Query().Has("authkey") && u.Query().Has("lang")
	})

	if len(result) == 0 {
		return nil, ErrURLNotFound
	}
	return result, nil
}

func (r Region) GetAPIBase() (baseURL string, baseQuery BaseQuery, err error) {
	urls, err := r.GetValidURLs()
	if err != nil {
		return
	}

	url_, _, ok := lo.FindLastIndexOf(urls, func(u *url.URL) bool {
		return lo.Contains([]string{
			"/event/gacha_info/api/getGachaLog", "/hk4e/event/e20190909gacha/index.html", "genshin/event/e20190909gacha/index.html",
		}, u.Path)
	})
	if !ok {
		url_ = urls[len(urls)-1]
	}

	baseURL = r.APIBaseURL
	baseQuery.AuthKeyVer = url_.Query().Get("authkey_ver")
	baseQuery.AuthKey = url_.Query().Get("authkey")
	baseQuery.Lang = url_.Query().Get("lang")

	return
}
