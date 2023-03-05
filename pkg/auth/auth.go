package auth

import (
	"fmt"
	"golang.org/x/exp/slices"
	"net/url"
	"strings"
)

// TODO handle "mhyurl.cn"

type Domain string

const (
	MIHOYO    Domain = "mihoyo.com"
	HOYOVERSE Domain = "hoyoverse.com"
)

var (
	GetGachaLogURLs = map[Domain]string{
		MIHOYO:    "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog",
		HOYOVERSE: "https://hk4e-api-os.hoyoverse.com/event/gacha_info/api/getGachaLog",
	}
)

func (domain Domain) GetGachaLogURL() string {
	return GetGachaLogURLs[domain]
}

var Domains = []Domain{
	MIHOYO,
	HOYOVERSE,
}

type Base struct {
	Domain
	AuthKeyVer string
	AuthKey    string
	Lang       string
}

func FromURL(rawURL string) (base Base, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	hostname := u.Hostname()
	i := slices.IndexFunc(Domains, func(domain Domain) bool {
		return strings.HasSuffix(hostname, string(domain))
	})
	if i == -1 {
		err = fmt.Errorf("invalid hostname: %s", hostname)
		return
	}

	base.Domain = Domains[i]

	query := u.Query()

	base.AuthKeyVer = query.Get("authkey_ver")
	if base.AuthKeyVer == "" {
		err = fmt.Errorf("missing authkey_ver")
		return
	}

	base.AuthKey = query.Get("authkey")
	if base.AuthKey == "" {
		err = fmt.Errorf("missing authkey")
		return
	}

	base.Lang = query.Get("lang")
	if base.Lang == "" {
		err = fmt.Errorf("missing lang")
		return
	}

	return
}
