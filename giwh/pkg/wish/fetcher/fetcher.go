package fetcher

import (
	"github.com/fhluo/giwh/pkg/auth"
	"github.com/fhluo/giwh/pkg/wiki"
	"github.com/fhluo/giwh/pkg/wish"
	"github.com/google/go-querystring/query"
	"github.com/samber/lo"
	"net/http"
	"net/url"
	"time"
)

type GetWishHistoryResponseData struct {
	List   []wish.Item `json:"list"`
	Page   int         `json:"page,string"`
	Region string      `json:"region"`
	Size   int         `json:"size,string"`
	Total  int         `json:"total,string"`
}

func GetWishHistory(url string) ([]wish.Item, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	responseData, err := wiki.GetJSONResponseData[*GetWishHistoryResponseData](resp)
	return responseData.List, err
}

const (
	DefaultInterval = 500 * time.Millisecond
)

type Query struct {
	AuthKeyVer string `url:"authkey_ver"`
	AuthKey    string `url:"authkey"`
	Lang       string `url:"lang"`

	WishType wish.Type `url:"gacha_type"`
	Size     int       `url:"size"`
	BeginID  string    `url:"begin_id,omitempty"`
	EndID    string    `url:"end_id,omitempty"`
}

func (q Query) Encode() string {
	values, err := query.Values(q)
	if err != nil {
		panic(err)
	}
	return values.Encode()
}

type Context struct {
	url      *url.URL
	query    Query
	interval time.Duration
}

func New(info auth.Info) *Context {
	return &Context{
		url: lo.Must(url.Parse(info.BaseURL)),
		query: Query{
			AuthKeyVer: info.AuthKeyVer,
			AuthKey:    info.AuthKey,
			Lang:       info.Lang,
			Size:       5,
		},
		interval: DefaultInterval,
	}
}

func (ctx *Context) URL() string {
	ctx.url.RawQuery = ctx.query.Encode()
	return ctx.url.String()
}

func (ctx *Context) String() string {
	return ctx.URL()
}

func (ctx *Context) WishType(wishType wish.Type) *Context {
	ctx.query.WishType = wishType
	return ctx
}

func (ctx *Context) Size(size int) *Context {
	ctx.query.Size = size
	return ctx
}

func (ctx *Context) Begin(id string) *Context {
	ctx.query.BeginID = id
	ctx.query.EndID = ""
	return ctx
}

func (ctx *Context) End(id string) *Context {
	ctx.query.EndID = id
	ctx.query.BeginID = ""
	return ctx
}

func (ctx *Context) Interval(interval time.Duration) {
	ctx.interval = interval
}

func (ctx *Context) Fetch() (items []wish.Item, err error) {
	items, err = GetWishHistory(ctx.URL())
	if err != nil || len(items) == 0 {
		return
	}

	if ctx.query.BeginID != "" {
		ctx.query.BeginID = items[0].ID
	} else {
		ctx.query.EndID = items[len(items)-1].ID
	}

	time.Sleep(ctx.interval)
	return
}
