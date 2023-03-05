package wish

import (
	"github.com/fhluo/giwh/pkg/auth"
	"github.com/fhluo/giwh/pkg/wiki"
	"github.com/google/go-querystring/query"
	"github.com/samber/lo"
	"net/http"
	"net/url"
	"time"
)

//go:generate go run ../../cmd/giwh-dev gen wishes

type GetWishHistoryResponseData struct {
	List   []Item `json:"list"`
	Page   int    `json:"page,string"`
	Region string `json:"region"`
	Size   int    `json:"size,string"`
	Total  int    `json:"total,string"`
}

func GetWishHistory(url string) ([]Item, error) {
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

	WishType Type  `url:"gacha_type"`
	Size     int   `url:"size"`
	BeginID  int64 `url:"begin_id,omitempty"`
	EndID    int64 `url:"end_id,omitempty"`
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

func New(base auth.Base) *Context {
	return &Context{
		url: lo.Must(url.Parse(base.GetGachaLogURL())),
		query: Query{
			AuthKeyVer: base.AuthKeyVer,
			AuthKey:    base.AuthKey,
			Lang:       base.Lang,
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

func (ctx *Context) WishType(wishType Type) *Context {
	ctx.query.WishType = wishType
	return ctx
}

func (ctx *Context) Size(size int) *Context {
	ctx.query.Size = size
	return ctx
}

func (ctx *Context) Begin(id int64) *Context {
	ctx.query.BeginID = id
	ctx.query.EndID = 0
	return ctx
}

func (ctx *Context) End(id int64) *Context {
	ctx.query.EndID = id
	ctx.query.BeginID = 0
	return ctx
}

func (ctx *Context) Interval(interval time.Duration) {
	ctx.interval = interval
}

func (ctx *Context) Fetch() (items []Item, err error) {
	items, err = GetWishHistory(ctx.URL())
	if err != nil || len(items) == 0 {
		return
	}

	if ctx.query.BeginID > 0 {
		ctx.query.BeginID = items[0].ID
	} else {
		ctx.query.EndID = items[len(items)-1].ID
	}

	time.Sleep(ctx.interval)
	return
}
