package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Item struct {
	UID      string `json:"uid"`
	WishType string `json:"gacha_type"`
	ItemID   string `json:"item_id"`
	Count    string `json:"count"`
	Time     string `json:"time"`
	Name     string `json:"name"`
	Lang     string `json:"lang"`
	ItemType string `json:"item_type"`
	Rarity   string `json:"rank_type"`
	ID       string `json:"id"`
}

type Data struct {
	Page   string `json:"page"`
	Size   string `json:"size"`
	Total  string `json:"total"`
	List   []Item `json:"list"`
	Region string `json:"region"`
}

type Result struct {
	*Data   `json:"data"`
	Message string `json:"message"`
	RetCode int    `json:"retcode"`
}

func GetWishHistory(url string) ([]Item, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(resp.Status)
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result Result

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if result.RetCode != 0 {
		return nil, fmt.Errorf(result.Message)
	}

	return result.List, nil
}

const (
	DefaultInterval = 500 * time.Millisecond
)

type BaseQuery struct {
	AuthKeyVer string `url:"authkey_ver"`
	AuthKey    string `url:"authkey"`
	Lang       string `url:"lang"`
}

type Query struct {
	BaseQuery

	WishType string `url:"gacha_type"`
	Size     string `url:"size"`
	BeginID  string `url:"begin_id,omitempty"`
	EndID    string `url:"end_id,omitempty"`
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
	items    []Item
	interval time.Duration

	query Query
}

func New(baseURL string, baseQuery BaseQuery) (*Context, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &Context{
		url: u,
		query: Query{
			BaseQuery: baseQuery,
			Size:      "5",
		},
		interval: DefaultInterval,
	}, nil
}

func (ctx *Context) WishType(wishType string) *Context {
	ctx.query.WishType = wishType
	ctx.query.BeginID = ""
	ctx.query.EndID = ""
	ctx.items = nil
	return ctx
}

func (ctx *Context) URL() string {
	ctx.url.RawQuery = ctx.query.Encode()
	return ctx.url.String()
}

func (ctx *Context) wait() {
	time.Sleep(ctx.interval)
}

func (ctx *Context) nextPage() {
	if len(ctx.items) == 0 {
		return
	}
	ctx.query.EndID = ctx.items[len(ctx.items)-1].ID
}

func (ctx *Context) fetch() error {
	items, err := GetWishHistory(ctx.URL())
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return Stop
	}
	ctx.items = append(ctx.items, items...)

	ctx.nextPage()
	ctx.wait()
	return nil
}

func (ctx *Context) Next() (Item, error) {
	if len(ctx.items) == 0 {
		err := ctx.fetch()
		if err != nil {
			return Item{}, err
		}
	}

	item := ctx.items[0]

	ctx.items = ctx.items[1:]
	return item, nil
}

func (ctx *Context) Begin(id string) *Context {
	ctx.query.BeginID = id
	return ctx
}

func (ctx *Context) End(id string) *Context {
	ctx.query.EndID = id
	return ctx
}

func (ctx *Context) FetchALL() ([]Item, error) {
	return Collect[Item](ctx)
}
