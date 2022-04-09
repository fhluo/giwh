package util

import (
	"errors"
	"github.com/fhluo/giwh/wh"
	"github.com/hashicorp/go-multierror"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

const (
	APIBaseURLCN     = "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog"
	APIBaseURLGlobal = "https://hk4e-api-os.hoyoverse.com/event/gacha_info/api/getGachaLog"

	HostNameCN     = "webstatic.mihoyo.com"
	HostNameGlobal = "webstatic-sea.hoyoverse.com"
)

var (
	PathCN     = filepath.Join(os.Getenv("USERPROFILE"), `\AppData\LocalLow\miHoYo\原神`)
	PathGlobal = filepath.Join(os.Getenv("USERPROFILE"), `\AppData\LocalLow\miHoYo\Genshin Impact`)

	OutputLogCN     = filepath.Join(PathCN, `output_log.txt`)
	OutputLogGlobal = filepath.Join(PathGlobal, `output_log.txt`)

	UIDInfoCN     = filepath.Join(PathCN, `UidInfo.txt`)
	UIDInfoGlobal = filepath.Join(PathGlobal, `UidInfo.txt`)
)

var ErrNotFound = errors.New("not found")

type info struct {
	name string
	time time.Time
}

func FindLatest(names ...string) (string, error) {
	infos := make([]*info, 0, len(names))

	var errs error

	for _, name := range names {
		fi, err := os.Stat(name)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		infos = append(infos, &info{name: name, time: fi.ModTime()})
	}

	if len(infos) == 0 {
		return "", errs
	}

	latest := infos[0]
	for _, i := range infos[1:] {
		if i.time.After(latest.time) {
			latest = i
		}
	}

	return latest.name, nil
}

func SortExisting(names ...string) ([]string, error) {
	infos := make([]*info, 0, len(names))

	var errs error

	for _, name := range names {
		fi, err := os.Stat(name)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		infos = append(infos, &info{name: name, time: fi.ModTime()})
	}

	switch len(infos) {
	case 0:
		return nil, errs
	case 1:
		return []string{infos[0].name}, nil
	default:
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].time.After(infos[j].time)
		})
		return lo.Map(infos, func(i *info, _ int) string {
			return i.name
		}), nil
	}
}

func FindURLFromOutputLog(filename string, f func(u *url.URL) bool) (*url.URL, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	matches := regexp.MustCompile(`OnGetWebViewPageFinish:(.*?)\r?\n`).FindAllSubmatch(data, -1)

	var errs error
	for i := len(matches) - 1; i >= 0; i-- {
		u, err := url.Parse(string(matches[i][1]))
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		if f(u) {
			return u, nil
		}
	}

	if errs != nil {
		return nil, errs
	}

	return nil, ErrNotFound
}

func GetUID(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile(`\d{9}`).Find(data)
	if r == nil {
		return "", ErrNotFound
	}

	return string(r), nil
}

func GetUIDAndAPIBaseURL() (string, string, error) {
	result, err := FindLatest(OutputLogCN, OutputLogGlobal)
	if err != nil {
		return "", "", err
	}

	switch result {
	case OutputLogCN:
		u, err := FindURLFromOutputLog(OutputLogCN, func(u *url.URL) bool {
			return u.Query().Has("authkey") && u.Hostname() == HostNameCN
		})
		if err != nil {
			return "", "", err
		}

		uid, err := GetUID(UIDInfoCN)
		if err != nil {
			return "", "", err
		}

		return uid, APIBaseURLCN + "?" + u.RawQuery, nil

	default:
		u, err := FindURLFromOutputLog(OutputLogGlobal, func(u *url.URL) bool {
			return u.Query().Has("authkey") && u.Hostname() == HostNameGlobal
		})
		if err != nil {
			return "", "", err
		}

		uid, err := GetUID(UIDInfoGlobal)
		if err != nil {
			return "", "", err
		}

		return uid, APIBaseURLGlobal + "?" + u.RawQuery, nil
	}
}

func LoadItems(filename string) (wh.Items, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var items wh.Items
	return items, jsoniter.Unmarshal(data, &items)
}

func LoadItemsIfExits(filename string) (wh.Items, error) {
	_, err := os.Stat(filename)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		return nil, nil
	}

	return LoadItems(filename)
}
