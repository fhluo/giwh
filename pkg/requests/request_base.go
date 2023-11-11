package requests

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// RequestBase 是 GetGachaLog 的基础请求
type RequestBase struct {
	BaseURL    string `url:"-"`           // GetGachaLog 的基础 URL
	AuthKeyVer string `url:"authkey_ver"` // 授权密钥版本
	AuthKey    string `url:"authkey"`     // 授权密钥
	Lang       string `url:"lang"`        // 语言
}

var baseURLs = map[string]string{
	"mihoyo.com":    "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog",
	"hoyoverse.com": "https://hk4e-api-os.hoyoverse.com/event/gacha_info/api/getGachaLog",
	"mhyurl.cn":     "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog",
}

// getBaseURL 返回 hostname 对应的 GetGachaLog 基础 URL
func getBaseURL(hostname string) (string, bool) {
	for suffix, baseURL := range baseURLs {
		if strings.HasSuffix(hostname, suffix) {
			return baseURL, true
		}
	}

	return "", false
}

// decodeQuery 为 requestBase 的指定字段赋予 query 中对应的值
func decodeQuery(query url.Values, requestBase *RequestBase) error {
	// 反射 requestBase
	value := reflect.Indirect(reflect.ValueOf(requestBase))
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

// NewRequestBase 从 rawURL 解析出 RequestBase
func NewRequestBase(rawURL string) (requestBase RequestBase, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	// 获取 GetGachaLog 的基础 URL
	var ok bool
	requestBase.BaseURL, ok = getBaseURL(u.Hostname())
	if !ok {
		err = fmt.Errorf("unsupported hostname: %s", u.Hostname())
		return
	}

	// 为 requestBase 的字段赋值
	err = decodeQuery(u.Query(), &requestBase)
	return
}
