package api

import (
	"github.com/fhluo/giwh/giwh/api/gacha"
	"github.com/fhluo/giwh/hoyo-api/requests"
	"github.com/fhluo/giwh/hoyo-auth/auths"
	"log/slog"
	"time"
)

const DefaultInterval = 500 * time.Millisecond // 默认请求间隔

type Data struct {
	List   []gacha.Log `json:"list"`         // 抽卡记录列表
	Page   int         `json:"page,string"`  // 页码
	Region string      `json:"region"`       // 地区
	Size   int         `json:"size,string"`  // 每页数量
	Total  int         `json:"total,string"` // 总数
}

// GetGachaLog 获取抽卡记录
func GetGachaLog(url string) (Data, error) {
	return requests.GetData[Data](url)
}

type Client struct {
	*auths.Auth
}

func NewClient(auth *auths.Auth) *Client {
	return &Client{
		Auth: auth,
	}
}

var LastRequestTime time.Time // 上次请求时间

// Wait 等待下次请求
func (c *Client) Wait() {
	// 第一次请求不等待
	if LastRequestTime.IsZero() {
		LastRequestTime = time.Now()
		slog.Debug("wait", "time", "0")
		return
	}

	now := time.Now()
	// 如果请求间隔小于默认间隔，则等待
	if now.Sub(LastRequestTime) < DefaultInterval {
		d := DefaultInterval - now.Sub(LastRequestTime)
		time.Sleep(d)
		slog.Debug("wait", "time", d.Milliseconds())
	}

	LastRequestTime = time.Now()
}

func (c *Client) GetGachaLog(url string) ([]gacha.Log, error) {
	c.Wait()

	resp, err := GetGachaLog(url)
	slog.Debug("GetGachaLog", "url", url)
	if err != nil {
		return nil, err
	}

	return resp.List, nil
}

type Fetcher struct {
	*Client
	*URLBuilder
}

func (c *Fetcher) NextURL() (string, error) {
	return c.Build()
}

func (c *Fetcher) NextPage() ([]gacha.Log, error) {
	u, err := c.NextURL()
	if err != nil {
		return nil, err
	}

	logs, err := c.GetGachaLog(u)
	if err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		return nil, nil
	}

	c.End(logs[len(logs)-1].ID)
	return logs, nil
}

func (c *Client) NewFetcher(gachaType gacha.Type) *Fetcher {
	return &Fetcher{
		Client:     c,
		URLBuilder: NewURLBuilder(c.Auth).GachaType(gachaType).Size(10),
	}
}

func (c *Client) Fetch(gachaType gacha.Type, f func(logs []gacha.Log) (stop bool)) ([]gacha.Log, error) {
	fetcher := c.NewFetcher(gachaType)

	var result []gacha.Log
	for {
		logs, err := fetcher.NextPage()
		if err != nil {
			return result, err
		}
		result = append(result, logs...)

		if len(logs) < fetcher.Query.Size || (f != nil && f(logs)) {
			break
		}
	}

	return result, nil
}
