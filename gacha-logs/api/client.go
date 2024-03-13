package api

import (
	"github.com/fhluo/giwh/hoyo-api"
	"slices"
	"time"
)

const DefaultInterval = 500 * time.Millisecond // 默认请求间隔

type Data struct {
	List   []Log  `json:"list"`         // 抽卡记录列表
	Page   int    `json:"page,string"`  // 页码
	Region string `json:"region"`       // 地区
	Size   int    `json:"size,string"`  // 每页数量
	Total  int    `json:"total,string"` // 总数
}

// GetGachaLog 获取抽卡记录
func GetGachaLog(url string) (Data, error) {
	return hoyo_api.GetData[Data](url)
}

type Client struct {
	*URLBuilder // 构建请求 URL

	GachaLogs []Log // 抽卡记录
	Cache     []Log // 缓存

	Done bool // 是否完成

	LastBeginID     string    // 上次开始 ID
	LastEndID       string    // 上次结束 ID
	LastRequestTime time.Time // 上次请求时间
}

func NewClient(u *URLBuilder) *Client {
	return &Client{
		URLBuilder: u,
	}
}

// Wait 等待下次请求
func (c *Client) Wait() {
	// 第一次请求不等待
	if c.LastRequestTime.IsZero() {
		c.LastRequestTime = time.Now()
		return
	}

	now := time.Now()
	// 如果请求间隔小于默认间隔，则等待
	if now.Sub(c.LastRequestTime) < DefaultInterval {
		time.Sleep(DefaultInterval - now.Sub(c.LastRequestTime))
	}

	c.LastRequestTime = time.Now()
}

// PrevURL 获取上一页抽卡记录 URL
func (c *Client) PrevURL() URLBuilder {
	panic("not implemented")
}

// PrevPage 获取上一页抽卡记录
func (c *Client) PrevPage() ([]Log, error) {
	panic("not implemented")
}

// NextURL 获取下一页抽卡记录 URL
func (c *Client) NextURL() (string, error) {
	// 第一次请求不设置 EndID
	if c.LastEndID != "" {
		c.End(c.LastEndID)
	}
	return c.Build()
}

// NextPage 获取下一页抽卡记录
func (c *Client) NextPage() ([]Log, error) {
	// 等待下次请求
	c.Wait()

	url, err := c.NextURL()
	if err != nil {
		return nil, err
	}

	resp, err := GetGachaLog(url)
	if err != nil {
		return nil, err
	}

	if len(resp.List) == 0 {
		c.Done = true // 标记为完成
		return nil, nil
	}

	// 追加抽卡记录
	c.GachaLogs = append(c.GachaLogs, resp.List...)
	c.Cache = resp.List

	// 更新开始 ID 和结束 ID
	c.LastBeginID = resp.List[0].ID
	c.LastEndID = resp.List[len(resp.List)-1].ID

	return resp.List, nil
}

// Next 获取下一条抽卡记录
func (c *Client) Next() (Log, error) {
	// 如果缓存为空，则获取下一页
	if len(c.Cache) == 0 {
		_, err := c.NextPage()
		if err != nil {
			return Log{}, err
		}
	}

	// 获取缓存中的第一条记录并移除
	log := c.Cache[0]
	c.Cache = c.Cache[1:]
	return log, nil
}

// All 获取所有抽卡记录
func (c *Client) All() ([]Log, error) {
	for {
		list, err := c.NextPage()
		if err != nil {
			return nil, err
		}

		// 如果没有抽卡记录，则退出
		if len(list) == 0 {
			break
		}
	}

	return c.GachaLogs, nil
}

// Until 获取所有抽卡记录，遇到满足条件的记录则停止
func (c *Client) Until(condition func(Log) bool) ([]Log, error) {
	for {
		// 如果已经完成，则退出
		if c.Done {
			break
		}

		// 获取下一条抽卡记录
		log, err := c.Next()
		if err != nil {
			return nil, err
		}

		// 如果满足条件，则退出
		if condition(log) {
			i := slices.IndexFunc(c.GachaLogs, func(l Log) bool {
				return l.ID == log.ID
			})
			return c.GachaLogs[:i], nil
		}
	}

	return c.GachaLogs, nil
}

// FetchAll 获取所有卡池的所有抽卡记录
func FetchAll(u URLBuilder) ([]Log, error) {
	var logs []Log

	for _, u := range u.GachaTypes(SharedGachaTypes) {
		list, err := NewClient(u).All()
		if err != nil {
			return nil, err
		}
		logs = append(logs, list...)
	}

	return logs, nil
}

// FetchAllUntil 获取所有卡池的所有抽卡记录，直到满足条件
func FetchAllUntil(url URLBuilder, condition func(Log) bool) ([]Log, error) {
	var logs []Log

	for _, u := range url.GachaTypes(SharedGachaTypes) {
		list, err := NewClient(u).Until(condition)
		if err != nil {
			return nil, err
		}
		logs = append(logs, list...)
	}

	return logs, nil
}
