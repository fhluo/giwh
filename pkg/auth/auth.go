package auth

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type Info struct {
	BaseURL    string `url:"-"`
	AuthKeyVer string `url:"authkey_ver"`
	AuthKey    string `url:"authkey"`
	Lang       string `url:"lang"`
}

func GetBaseURL(hostname string) (string, error) {
	switch {
	case strings.HasSuffix(hostname, "mihoyo.com"):
		return "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog", nil
	case strings.HasSuffix(hostname, "hoyoverse.com"):
		return "https://hk4e-api-os.hoyoverse.com/event/gacha_info/api/getGachaLog", nil
	case strings.HasSuffix(hostname, "mhyurl.cn"):
		return "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog", nil
	default:
		return "", fmt.Errorf("invalid hostname: %s", hostname)
	}
}

func FromURL(rawURL string) (info Info, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	info.BaseURL, err = GetBaseURL(u.Hostname())
	if err != nil {
		return
	}

	query := u.Query()
	value := reflect.Indirect(reflect.ValueOf(&info))
	type_ := value.Type()

	for i := 0; i < value.NumField(); i++ {
		key := type_.Field(i).Tag.Get("url")
		if key != "" && key != "-" {
			if !query.Has(key) {
				err = fmt.Errorf("missing %s", key)
				return
			}
			value.Field(i).SetString(query.Get(key))
		}
	}

	return
}
