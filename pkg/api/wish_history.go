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
	Count    int      `json:"count,string"`
	WishType WishType `json:"gacha_type,string"`
	ID       int64    `json:"id,string"`
	ItemID   string   `json:"item_id"`
	ItemType string   `json:"item_type"`
	Lang     string   `json:"lang"`
	Name     string   `json:"name"`
	Rarity   Rarity   `json:"rank_type,string"`
	Time     Time     `json:"time"`
	UID      int      `json:"uid,string"`
}

func (item *Item) String() string {
	return fmt.Sprintf("Item{Name:%s Time:%v UID:%v ID:%v}", item.Name, item.Time.String(), item.UID, item.ID)
}

type ItemList struct {
	List   []*Item `json:"list"`
	Page   int     `json:"page,string"`
	Region string  `json:"region"`
	Size   int     `json:"size,string"`
	Total  int     `json:"total,string"`
}

func GetWishHistory(url string) ([]*Item, error) {
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
	var r JSONResponse[*ItemList]

	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	if r.RetCode != 0 {
		return nil, fmt.Errorf(r.Message)
	}

	return r.Data.List, nil
}

const (
	DefaultInterval = 500 * time.Millisecond
)

type Query struct {
	AuthKeyVer string `url:"authkey_ver"`
	AuthKey    string `url:"authkey"`
	Lang       string `url:"lang"`

	WishType SharedWishType `url:"gacha_type"`
	Size     int            `url:"size"`
	BeginID  int64          `url:"begin_id,omitempty"`
	EndID    int64          `url:"end_id,omitempty"`
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
	items    []*Item
	interval time.Duration

	query      Query
	handleNext []func(item *Item)
}

func New(base Base) (*Context, error) {
	u, err := url.Parse(base.URL)
	if err != nil {
		return nil, err
	}
	return &Context{
		url: u,
		query: Query{
			AuthKeyVer: base.Query.AuthKeyVer,
			AuthKey:    base.Query.AuthKey,
			Lang:       base.Query.Lang,
			Size:       5,
		},
		interval: DefaultInterval,
	}, nil
}

func (ctx *Context) SetInterval(interval time.Duration) {
	ctx.interval = interval
}

func (ctx *Context) URL() string {
	ctx.url.RawQuery = ctx.query.Encode()
	return ctx.url.String()
}

func (ctx *Context) wait() {
	time.Sleep(ctx.interval)
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

	if ctx.query.BeginID > 0 {
		ctx.query.BeginID = ctx.items[0].ID
	} else {
		ctx.query.EndID = ctx.items[len(ctx.items)-1].ID
	}

	ctx.wait()
	return nil
}

func (ctx *Context) GetUID() (int, error) {
	item, err := ctx.Peek()
	if err != nil {
		return 0, err
	}
	return item.UID, nil
}

func (ctx *Context) Peek() (*Item, error) {
	if len(ctx.items) == 0 {
		err := ctx.fetch()
		if err != nil {
			return nil, err
		}
	}

	if ctx.query.BeginID > 0 {
		return ctx.items[len(ctx.items)-1], nil
	} else {
		return ctx.items[0], nil
	}
}

func (ctx *Context) Use(handleNext ...func(item *Item)) {
	ctx.handleNext = append(ctx.handleNext, handleNext...)
}

func (ctx *Context) Next() (*Item, error) {
	if len(ctx.items) == 0 {
		err := ctx.fetch()
		if err != nil {
			return nil, err
		}
	}

	var item *Item
	if ctx.query.BeginID > 0 {
		item = ctx.items[len(ctx.items)-1]
		ctx.items = ctx.items[:len(ctx.items)-1]
	} else {
		item = ctx.items[0]
		ctx.items = ctx.items[1:]
	}

	for i := range ctx.handleNext {
		ctx.handleNext[i](item)
	}

	return item, nil
}

func (ctx *Context) WishType(wishType SharedWishType) *Context {
	ctx.query.WishType = wishType
	return ctx
}

func (ctx *Context) Size(size int) *Context {
	ctx.query.Size = size
	return ctx
}

func (ctx *Context) reset() {
	ctx.query.BeginID = 0
	ctx.query.EndID = 0
	ctx.items = nil
}

func (ctx *Context) Begin(id int64) *Context {
	ctx.reset()
	ctx.query.BeginID = id
	return ctx
}

func (ctx *Context) End(id int64) *Context {
	ctx.reset()
	ctx.query.EndID = id
	return ctx
}

func (ctx *Context) FetchAll() ([]*Item, error) {
	return Collect[*Item](ctx)
}
