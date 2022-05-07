package clients

import (
	"errors"
	"github.com/fhluo/giwh/pkg/util"
	"github.com/hashicorp/go-multierror"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ErrURLNotFound    = errors.New("URL not found")
	ErrClientNotFound = errors.New("client not found")
	ErrUIDNotFound    = errors.New("UID not found")

	CN = Client{
		PersistentDataPath: filepath.Join(os.Getenv("USERPROFILE"), `\AppData\LocalLow\miHoYo\原神`),
		QueryLinkHostName:  "webstatic.mihoyo.com",
		APIGetWishHistory:  "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog",
	}
	Global = Client{
		PersistentDataPath: filepath.Join(os.Getenv("USERPROFILE"), `\AppData\LocalLow\miHoYo\Genshin Impact`),
		QueryLinkHostName:  "webstatic-sea.hoyoverse.com",
		APIGetWishHistory:  "https://hk4e-api-os.hoyoverse.com/event/gacha_info/api/getGachaLog",
	}
)

func RecentlyUsed() (client Client, err error) {
	result, err := util.FindLatest(CN.OutputLogPath(), Global.OutputLogPath())
	if err != nil {
		return
	}

	switch result {
	case CN.OutputLogPath():
		client = CN
	case Global.OutputLogPath():
		client = Global
	default:
		err = ErrClientNotFound
	}
	return
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
	return client.FindURLFromOutputLog(func(u *url.URL) bool {
		return u.Query().Has("authkey") && u.Hostname() == client.QueryLinkHostName
	})
}

func (client Client) GetBaseURL() (baseURL string, err error) {
	u, err := client.FindQueryLinkWithAuthKey()
	if err != nil {
		return
	}
	return client.APIGetWishHistory + "?" + u.RawQuery, nil
}
