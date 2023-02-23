package clients

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	CN = 1 << iota
	OS
)

var (
	ErrURLNotFound    = errors.New("URL not found")
	ErrClientNotFound = errors.New("client not found")
	ErrUIDNotFound    = errors.New("UID not found")

	CNClient = Client{
		PersistentDataPath: filepath.Join(os.Getenv("USERPROFILE"), `\AppData\LocalLow\miHoYo\原神`),
		QueryLinkHostName:  "webstatic.mihoyo.com",
		APIGetWishHistory:  "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog",
	}
	OSClient = Client{
		PersistentDataPath: filepath.Join(os.Getenv("USERPROFILE"), `\AppData\LocalLow\miHoYo\Genshin Impact`),
		QueryLinkHostName:  "webstatic-sea.hoyoverse.com",
		APIGetWishHistory:  "https://hk4e-api-os.hoyoverse.com/event/gacha_info/api/getGachaLog",
	}
)

func DetectVersions() int {
	var versions int

	if _, err := os.Stat(CNClient.PersistentDataPath); err == nil {
		versions |= CN
	}

	if _, err := os.Stat(OSClient.PersistentDataPath); err == nil {
		versions |= OS
	}

	return versions
}

func Default() (client Client, err error) {
	versions := DetectVersions()
	switch versions {
	case CN:
		client = CNClient

	case OS:
		client = OSClient

	case CN | OS:
		log1, err := os.Stat(CNClient.OutputLogPath())
		if err != nil {
			return client, err
		}
		log2, err := os.Stat(OSClient.OutputLogPath())
		if err != nil {
			return client, err
		}

		if log1.ModTime().After(log2.ModTime()) {
			client = CNClient
		} else {
			client = OSClient
		}

	default:
		err = ErrClientNotFound
	}

	return
}

var DataPathRE = regexp.MustCompile(`.:/.*?/(YuanShen_Data|GenshinImpact_Data)/`)

func FindDataPath(outputLog []byte) string {
	return string(DataPathRE.Find(outputLog))
}

type Client struct {
	PersistentDataPath string
	QueryLinkHostName  string
	APIGetWishHistory  string
}

func (client Client) OutputLogPath() string {
	return filepath.Join(client.PersistentDataPath, `output_log.txt`)
}

func (client Client) UIDInfoPath() string {
	return filepath.Join(client.PersistentDataPath, `UidInfo.txt`)
}

func (client Client) GetUID() (string, error) {
	data, err := os.ReadFile(client.UIDInfoPath())
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile(`\d{9}`).Find(data)
	if r == nil {
		return "", ErrUIDNotFound
	}

	return string(r), nil
}

func (client Client) FindURLFromCacheData(f func(u *url.URL) bool) (*url.URL, error) {
	data, err := os.ReadFile(client.OutputLogPath())
	if err != nil {
		return nil, err
	}

	path := FindDataPath(data)
	if len(path) == 0 {
		return nil, fmt.Errorf("failed to find data path from output_log.txt")
	}

	data, err = os.ReadFile(filepath.Join(path, `webCaches\Cache\Cache_Data\data_2`))
	if err != nil {
		return nil, err
	}

	matches := regexp.MustCompile(`https://[^\s\0]+`).FindAll(data, -1)

	var errs error
	for i := len(matches) - 1; i >= 0; i-- {
		u, err := url.Parse(string(matches[i]))
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

	return nil, ErrURLNotFound
}

func (client Client) FindURLFromOutputLog(f func(u *url.URL) bool) (*url.URL, error) {
	data, err := os.ReadFile(client.OutputLogPath())
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

	return nil, ErrURLNotFound
}

func (client Client) FindQueryLinkWithAuthKey() (*url.URL, error) {
	return client.FindURLFromCacheData(func(u *url.URL) bool {
		return u.Query().Has("authkey") && u.Hostname() == client.QueryLinkHostName && strings.HasSuffix(u.Path, "index.html")
	})
}

func (client Client) GetBaseURL() (baseURL string, err error) {
	u, err := client.FindQueryLinkWithAuthKey()
	if err != nil {
		return
	}
	return client.APIGetWishHistory + "?" + u.RawQuery, nil
}
