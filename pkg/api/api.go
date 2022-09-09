package api

import (
	"encoding/json"
	"fmt"
	"github.com/fhluo/giwh/pkg/wish"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

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

type API struct {
	url      *url.URL
	items    []Item
	interval time.Duration
	stopID   string

	wishType int    `url:"gacha_type"`
	page     int    `url:"page"`
	size     int    `url:"size"`
	endID    string `url:"end_id"`
}

func New(rawURL string, wishType wish.Type) (*API, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return &API{
		url:      u,
		wishType: int(wishType),
		page:     1,
		size:     5,
		endID:    "0",
		interval: DefaultInterval,
	}, nil
}

func (api *API) Init(wishType wish.Type) {
	api.wishType = int(wishType)
	api.page = 1
	api.endID = "0"
	api.items = nil
}

func (api *API) URL() string {
	api.url.Query().Set("gacha_type", strconv.Itoa(api.wishType))
	api.url.Query().Set("page", strconv.Itoa(api.page))
	api.url.Query().Set("size", strconv.Itoa(api.size))
	api.url.Query().Set("end_id", api.endID)
	return api.url.String()
}

func (api *API) wait() {
	time.Sleep(api.interval)
}

func (api *API) nextPage() {
	if len(api.items) == 0 {
		return
	}
	api.page++
	api.endID = api.items[len(api.items)-1].ID
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
	if item.ID == api.stopID {
		return item, Stop
	}

	api.items = api.items[1:]
	return item, nil
}

func (api *API) FetchALL() ([]Item, error) {
	return Collect[Item](api)
}

func (api *API) FetchUntil(id string) ([]Item, error) {
	api.stopID = id
	defer func() {
		api.stopID = ""
	}()
	return api.FetchALL()
}
