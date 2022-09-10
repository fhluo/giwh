package api

import (
	"encoding/json"
	"fmt"
	"github.com/fhluo/giwh/pkg/wish"
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
	Page   string `json:"page"`
	Size   string `json:"size"`
	Total  string `json:"total"`
	List   []Item `json:"list"`
	Region string `json:"region"`
}

type Result struct {
	RetCode int    `json:"retcode"`
	Message string `json:"message"`
	Data    `json:"data"`
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

type API struct {
	url      *url.URL
	items    []Item
	interval time.Duration

	query Query
}

func New(baseURL string, baseQuery BaseQuery, wishType wish.Type) (*API, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &API{
		url: u,
		query: Query{
			BaseQuery: baseQuery,
			WishType:  strconv.Itoa(int(wishType)),
			Size:      "5",
			BeginID:   "",
			EndID:     "0",
		},
		interval: DefaultInterval,
	}, nil
}

func (api *API) Init(wishType wish.Type) {
	api.query.WishType = strconv.Itoa(int(wishType))
	api.query.EndID = "0"
	api.items = nil
}

func (api *API) URL() string {
	api.url.RawQuery = api.query.Encode()
	return api.url.String()
}

func (api *API) wait() {
	time.Sleep(api.interval)
}

func (api *API) nextPage() {
	if len(api.items) == 0 {
		return
	}
	api.query.EndID = api.items[len(api.items)-1].ID
}

func (api *API) fetch() error {
	items, err := GetWishHistory(api.URL())
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return Stop
	}
	api.items = append(api.items, items...)

	api.nextPage()
	api.wait()
	return nil
}

func (api *API) Next() (Item, error) {
	if len(api.items) == 0 {
		err := api.fetch()
		if err != nil {
			return Item{}, err
		}
	}

	item := api.items[0]

	api.items = api.items[1:]
	return item, nil
}

func (api *API) FetchALL() ([]Item, error) {
	return Collect[Item](api)
}
