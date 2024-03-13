package requests

import (
	"github.com/fhluo/giwh/gacha-logs/gacha"
	"slices"
	"time"
)

const (
	DefaultSize     = 5                      // 默认每页数量
	DefaultInterval = 500 * time.Millisecond // 默认请求间隔
)

// Response 是抽卡记录响应
type Response struct {
	List   []gacha.Log `json:"list"`         // 抽卡记录列表
	Page   int         `json:"page,string"`  // 页码
	Region string      `json:"region"`       // 地区
	Size   int         `json:"size,string"`  // 每页数量
	Total  int         `json:"total,string"` // 总数
}

// GetGachaLog 获取抽卡记录
func GetGachaLog(url string) (Response, error) {
	return GetDataFromAPI[Response](url)
}

// Fetcher 是抽卡记录获取器
type Fetcher struct {
	URL       RequestURL  // 请求 URL
	GachaLogs []gacha.Log // 抽卡记录
	Cache     []gacha.Log // 缓存

	Done bool // 是否完成

	LastBeginID     string    // 上次开始 ID
	LastEndID       string    // 上次结束 ID
	LastRequestTime time.Time // 上次请求时间
}

// FetchGachaLogs 返回一个抽卡记录获取器
func FetchGachaLogs(url RequestURL) *Fetcher {
	return &Fetcher{
		URL: url,
	}
}

// Wait 等待下次请求
func (f *Fetcher) Wait() {
	// 第一次请求不等待
	if f.LastRequestTime.IsZero() {
		f.LastRequestTime = time.Now()
		return
	}

	now := time.Now()
	// 如果请求间隔小于默认间隔，则等待
	if now.Sub(f.LastRequestTime) < DefaultInterval {
		time.Sleep(DefaultInterval - now.Sub(f.LastRequestTime))
	}

	f.LastRequestTime = time.Now()
}

// PrevURL 获取上一页抽卡记录 URL
func (f *Fetcher) PrevURL() RequestURL {
	panic("not implemented")
}

// PrevPage 获取上一页抽卡记录
func (f *Fetcher) PrevPage() ([]gacha.Log, error) {
	panic("not implemented")
}

// NextURL 获取下一页抽卡记录 URL
func (f *Fetcher) NextURL() RequestURL {
	// 第一次请求不设置 EndID
	if f.LastEndID != "" {
		f.URL = f.URL.WithEndID(f.LastEndID)
	}
	return f.URL
}

// NextPage 获取下一页抽卡记录
func (f *Fetcher) NextPage() ([]gacha.Log, error) {
	// 等待下次请求
	f.Wait()

	resp, err := GetGachaLog(f.NextURL().String())
	if err != nil {
		return nil, err
	}

	if len(resp.List) == 0 {
		f.Done = true // 标记为完成
		return nil, nil
	}

	// 追加抽卡记录
	f.GachaLogs = append(f.GachaLogs, resp.List...)
	f.Cache = resp.List

	// 更新开始 ID 和结束 ID
	f.LastBeginID = resp.List[0].ID
	f.LastEndID = resp.List[len(resp.List)-1].ID

	return resp.List, nil
}

// Next 获取下一条抽卡记录
func (f *Fetcher) Next() (gacha.Log, error) {
	// 如果缓存为空，则获取下一页
	if len(f.Cache) == 0 {
		_, err := f.NextPage()
		if err != nil {
			return gacha.Log{}, err
		}
	}

	// 获取缓存中的第一条记录并移除
	log := f.Cache[0]
	f.Cache = f.Cache[1:]
	return log, nil
}

// All 获取所有抽卡记录
func (f *Fetcher) All() ([]gacha.Log, error) {
	for {
		list, err := f.NextPage()
		if err != nil {
			return nil, err
		}

		// 如果没有抽卡记录，则退出
		if len(list) == 0 {
			break
		}
	}

	return f.GachaLogs, nil
}

// Until 获取所有抽卡记录，遇到满足条件的记录则停止
func (f *Fetcher) Until(condition func(gacha.Log) bool) ([]gacha.Log, error) {
	for {
		// 如果已经完成，则退出
		if f.Done {
			break
		}

		// 获取下一条抽卡记录
		log, err := f.Next()
		if err != nil {
			return nil, err
		}

		// 如果满足条件，则退出
		if condition(log) {
			i := slices.IndexFunc(f.GachaLogs, func(l gacha.Log) bool {
				return l.ID == log.ID
			})
			return f.GachaLogs[:i], nil
		}
	}

	return f.GachaLogs, nil
}

// FetchAll 获取所有卡池的所有抽卡记录
func FetchAll(url RequestURL) ([]gacha.Log, error) {
	var logs []gacha.Log

	for _, u := range url.WithGachaTypes(SharedGachaTypes) {
		list, err := u.FetchGachaLogs().All()
		if err != nil {
			return nil, err
		}
		logs = append(logs, list...)
	}

	return logs, nil
}

// FetchAllUntil 获取所有卡池的所有抽卡记录，直到满足条件
func FetchAllUntil(url RequestURL, condition func(gacha.Log) bool) ([]gacha.Log, error) {
	var logs []gacha.Log

	for _, u := range url.WithGachaTypes(SharedGachaTypes) {
		list, err := u.FetchGachaLogs().Until(condition)
		if err != nil {
			return nil, err
		}
		logs = append(logs, list...)
	}

	return logs, nil
}
