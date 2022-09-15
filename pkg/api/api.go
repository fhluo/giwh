package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
	"strconv"
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
	Page   string  `json:"page"`
	Size   string  `json:"size"`
	Total  string  `json:"total"`
	List   []*Item `json:"list"`
	Region string  `json:"region"`
}

type Result struct {
	*Data   `json:"data"`
	Message string `json:"message"`
	RetCode int    `json:"retcode"`
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

type Query struct {
	AuthKeyVer string `url:"authkey_ver"`
	AuthKey    string `url:"authkey"`
	Lang       string `url:"lang"`

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
			Size:       "5",
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

	if ctx.query.BeginID != "" {
		ctx.query.BeginID = ctx.items[0].ID
	} else {
		ctx.query.EndID = ctx.items[len(ctx.items)-1].ID
	}

	ctx.wait()
	return nil
}

func (ctx *Context) GetUID() (string, error) {
	item, err := ctx.Peek()
	if err != nil {
		return "", err
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

	if ctx.query.BeginID != "" {
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
	if ctx.query.BeginID != "" {
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

const (
	BeginnersWish       = "100" // Beginners' Wish (Novice Wish)
	StandardWish        = "200" // Standard Wish (Permanent Wish)
	CharacterEventWish  = "301" // Character Event Wish
	WeaponEventWish     = "302" // Weapon Event Wish
	CharacterEventWish2 = "400" // Character Event Wish-2

	OneStar   = "1"
	TwoStar   = "2"
	ThreeStar = "3"
	FourStar  = "4"
	FiveStar  = "5"
)

var (
	WishTypes = []string{
		CharacterEventWish,
		CharacterEventWish2,
		WeaponEventWish,
		StandardWish,
		BeginnersWish,
	}
	SharedWishTypes = []string{
		CharacterEventWish,
		WeaponEventWish,
		StandardWish,
		BeginnersWish,
	}
)

func Pity(rarity string, wishType string) int {
	switch rarity {
	case FiveStar:
		switch wishType {
		case CharacterEventWish:
			return 90
		case WeaponEventWish:
			return 80
		case StandardWish:
			return 90
		default:
			return 90
		}
	case FourStar:
		return 10
	default:
		return 1
	}
}

func Pity4Star(_ string) int {
	return 10
}

func Pity5Star(wishType string) int {
	switch wishType {
	case CharacterEventWish:
		return 90
	case WeaponEventWish:
		return 80
	case StandardWish:
		return 90
	default:
		return 90
	}
}

func (ctx *Context) WishType(wishType string) *Context {
	ctx.query.WishType = wishType
	return ctx
}

func (ctx *Context) Size(size int) *Context {
	ctx.query.Size = strconv.Itoa(size)
	return ctx
}

func (ctx *Context) reset() {
	ctx.query.BeginID = ""
	ctx.query.EndID = ""
	ctx.items = nil
}

func (ctx *Context) Begin(id string) *Context {
	ctx.reset()
	ctx.query.BeginID = id
	return ctx
}

func (ctx *Context) End(id string) *Context {
	ctx.reset()
	ctx.query.EndID = id
	return ctx
}

func (ctx *Context) FetchAll() ([]*Item, error) {
	return Collect[*Item](ctx)
}
