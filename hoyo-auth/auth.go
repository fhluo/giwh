package hoyo_auth

import (
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strings"
)

type Auth struct {
	Hostname   string `url:"-" toml:"hostname"`
	AuthKeyVer string `url:"authkey_ver" toml:"authkey_ver"`
	AuthKey    string `url:"authkey" toml:"authkey"`
	Lang       string `url:"lang" toml:"lang"`
}

func New(rawURL string) (*Auth, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	auth := &Auth{Hostname: u.Hostname()}
	if !auth.valid() {
		return auth, fmt.Errorf("invalid hostname: %s", auth.Hostname)
	}

	return auth, auth.decode(u.Query())
}

var domains = []string{
	"mihoyo.com",
	"hoyoverse.com",
	"mhyurl.cn",
}

func (auth *Auth) valid() bool {
	return slices.ContainsFunc(domains, func(domain string) bool {
		return strings.HasSuffix(auth.Hostname, domain)
	})
}

// decode 为指定字段赋予 query 中对应的值
func (auth *Auth) decode(query url.Values) error {
	value := reflect.Indirect(reflect.ValueOf(auth))
	typ := value.Type()

	for i := 0; i < typ.NumField(); i++ {
		// 获取 url 标签的值
		key := typ.Field(i).Tag.Get("url")

		// 忽略不需要的字段
		if key == "-" || key == "" {
			continue
		}

		// 检查是否缺少必要的字段
		if !query.Has(key) {
			return fmt.Errorf("missing query key: %s", key)
		}

		// 为字段赋值
		value.Field(i).SetString(query.Get(key))
	}

	return nil
}
