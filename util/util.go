package util

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	jsoniter "github.com/json-iterator/go"
	"io/fs"
	"net/url"
	"os"
	"regexp"
	"wh/wh"
)

var ErrNotFound = errors.New("not found")

type info struct {
	name string
	time int64
}

func ReadLatestFile(names ...string) ([]byte, error) {
	infos := make([]*info, 0, len(names))

	var errs error

	for _, name := range names {
		fi, err := os.Stat(name)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		infos = append(infos, &info{name: name, time: fi.ModTime().UnixNano()})
	}

	if len(infos) == 0 {
		return nil, errs
	}

	latest := infos[0]
	for _, i := range infos[1:] {
		if i.time > latest.time {
			latest = i
		}
	}

	return os.ReadFile(latest.name)
}

func FindURLFromOutputLog(f func(u *url.URL) bool, filenames ...string) (*url.URL, error) {
	data, err := ReadLatestFile(filenames...)
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
