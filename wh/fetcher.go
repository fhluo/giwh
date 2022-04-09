package wh

import (
	"fmt"
	"github.com/google/go-querystring/query"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultInterval = time.Second / 2
)

type AuthInfo struct {
	UID     string
	BaseURL string
}

// QP Query Parameters
type QP struct {
	WishType int   `url:"gacha_type"`
	Page     int   `url:"page"`
	Size     int   `url:"size"`
	EndID    int64 `url:"end_id"`
}

func NewQP(wishType int) *QP {
	return &QP{
		WishType: wishType,
		Page:     1,
		Size:     6,
		EndID:    0,
	}
}

func (qp *QP) Values() url.Values {
	values, err := query.Values(qp)
	if err != nil {
		log.Fatalln(err)
	}
	return values
}

type Fetcher struct {
	BaseURL  string
	Visit    map[int64]bool
	Interval time.Duration

	*QP
}

func NewFetcher(baseURL string, wishType WishType, visit map[int64]bool) *Fetcher {
	if visit == nil {
		visit = make(map[int64]bool)
	}
	return &Fetcher{
		BaseURL:  baseURL,
		Visit:    visit,
		Interval: DefaultInterval,
		QP:       NewQP(int(wishType)),
	}
}

func (f *Fetcher) URL() string {
	u, err := url.Parse(f.BaseURL)
	if err != nil {
		log.Fatalln(err)
	}
	u.RawQuery += "&" + f.Values().Encode()
	return u.String()
}

func (f *Fetcher) Fetch() (items Items, err error) {

	resp, err := http.Get(f.URL())
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(resp.Status)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if jsoniter.Get(data, "retcode").ToInt() != 0 {
		err = fmt.Errorf(jsoniter.Get(data, "message").ToString())
		return
	}

	jsoniter.Get(data, "data", "list").ToVal(&items)

	return
}

func (f *Fetcher) FetchNext() (Items, error) {
	items, err := f.Fetch()
	if err != nil {
		return nil, err
	}

	for i, item := range items {
		if f.Visit[item.ID()] {
			items = items[:i]
			break
		} else {
			f.Visit[item.ID()] = true
		}
	}

	if len(items) == 0 {
		return nil, nil
	}

	f.Page++
	f.EndID = items[len(items)-1].ID()

	return items, nil
}

func (f *Fetcher) FetchALL() (Items, error) {
	var result Items

	for {
		items, err := f.FetchNext()
		if err != nil {
			return nil, err
		}

		if len(items) == 0 {
			break
		}

		fmt.Printf("%s: %s\n", items[0].WishType(), strings.Join(
			lo.Map(items, func(item Item, _ int) string {
				return item.String()
			}),
			", ",
		))

		result = append(result, items...)
		time.Sleep(f.Interval)
	}

	return result, nil
}
