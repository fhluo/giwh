package api

import (
	"fmt"
	"github.com/fhluo/giwh/hoyo-auth"
	"github.com/google/go-querystring/query"
	"net/url"
	"strings"
)

// Query 是抽卡记录 API 查询参数
type Query struct {
	AuthKeyVer string `url:"authkey_ver"`        // 授权密钥版本
	AuthKey    string `url:"authkey"`            // 授权密钥
	Lang       string `url:"lang"`               // 语言
	GachaType  string `url:"gacha_type"`         // 卡池类型
	Size       int    `url:"size"`               // 每页数量
	BeginID    string `url:"begin_id,omitempty"` // 开始 ID
	EndID      string `url:"end_id,omitempty"`   // 结束 ID
}

// Values 使用 query.Values 编码 Query
func (q Query) Values() url.Values {
	// 忽略错误，因为 q 是结构体且无自定义 Encoder
	r, _ := query.Values(q)
	return r
}

// String 返回 query string
func (q Query) String() string {
	return q.Values().Encode()
}

// URLBuilder 用于构建抽卡记录 API 请求 URL
type URLBuilder struct {
	Hostname string
	Query    Query
}

const DefaultSize = 5 // 默认每页数量

// NewURLBuilder 返回一个初始化的 URLBuilder
func NewURLBuilder(auth *hoyo_auth.Auth) *URLBuilder {
	return &URLBuilder{
		Hostname: auth.Hostname,
		Query: Query{
			AuthKeyVer: auth.AuthKeyVer,
			AuthKey:    auth.AuthKey,
			Lang:       auth.Lang,
			Size:       DefaultSize,
		},
	}
}

func (u *URLBuilder) Copy() *URLBuilder {
	return &(*u)
}

func (u *URLBuilder) Reset() *URLBuilder {
	return u.Begin("").End("")
}

var baseURLs = map[string]string{
	"mihoyo.com":    "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog",
	"hoyoverse.com": "https://hk4e-api-os.hoyoverse.com/event/gacha_info/api/getGachaLog",
	"mhyurl.cn":     "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog",
}

// BaseURL 返回 hostname 对应的 GetGachaLog 基础 URL
func (u *URLBuilder) BaseURL() (string, error) {
	for suffix, baseURL := range baseURLs {
		if strings.HasSuffix(u.Hostname, suffix) {
			return baseURL, nil
		}
	}

	return "", fmt.Errorf("invalid hostname: %s", u.Hostname)
}

// GachaType 设置卡池类型
func (u *URLBuilder) GachaType(gachaType string) *URLBuilder {
	u.Query.GachaType = gachaType
	return u
}

// GachaTypes 返回多个卡池类型的 URL
func (u *URLBuilder) GachaTypes(gachaTypes []string) []*URLBuilder {
	builders := make([]*URLBuilder, len(gachaTypes))
	for i, gachaType := range gachaTypes {
		builders[i] = u.Copy().GachaType(gachaType)
	}
	return builders
}

// Size 设置每页数量
func (u *URLBuilder) Size(size int) *URLBuilder {
	u.Query.Size = size
	return u
}

// Begin 设置开始 ID
func (u *URLBuilder) Begin(id string) *URLBuilder {
	u.Query.BeginID = id
	return u
}

// End 设置结束 ID
func (u *URLBuilder) End(id string) *URLBuilder {
	u.Query.EndID = id
	return u
}

func (u *URLBuilder) Build() (string, error) {
	baseURL, err := u.BaseURL()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s?%s", baseURL, u.Query.String()), nil
}
