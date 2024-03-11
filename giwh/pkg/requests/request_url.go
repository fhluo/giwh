package requests

import (
	"github.com/google/go-querystring/query"
	"net/url"
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

// RequestURL 是抽卡记录 API 请求 URL
type RequestURL struct {
	BaseURL string
	Query   Query
}

// NewRequestURL 返回一个初始化的 RequestURL
func NewRequestURL(r RequestBase) RequestURL {
	return RequestURL{
		BaseURL: r.BaseURL,
		Query: Query{
			AuthKeyVer: r.AuthKeyVer,
			AuthKey:    r.AuthKey,
			Lang:       r.Lang,
			Size:       DefaultSize,
		},
	}
}

// WithGachaType 设置卡池类型
func (r RequestURL) WithGachaType(gachaType string) RequestURL {
	r.Query.GachaType = gachaType
	return r
}

// WithGachaTypes 返回多个卡池类型的请求 URL
func (r RequestURL) WithGachaTypes(gachaTypes []string) []RequestURL {
	requestURLs := make([]RequestURL, len(gachaTypes))
	for i, gachaType := range gachaTypes {
		requestURLs[i] = r.WithGachaType(gachaType)
	}
	return requestURLs
}

// WithSize 设置每页数量
func (r RequestURL) WithSize(size int) RequestURL {
	r.Query.Size = size
	return r
}

// WithBeginID 设置开始 ID
func (r RequestURL) WithBeginID(beginID string) RequestURL {
	r.Query.BeginID = beginID
	return r
}

// WithEndID 设置结束 ID
func (r RequestURL) WithEndID(endID string) RequestURL {
	r.Query.EndID = endID
	return r
}

// FetchGachaLogs 返回一个抽卡记录获取器
func (r RequestURL) FetchGachaLogs() *Fetcher {
	return &Fetcher{
		URL: r,
	}
}

// String 返回请求 URL
func (r RequestURL) String() string {
	return r.BaseURL + "?" + r.Query.String()
}
