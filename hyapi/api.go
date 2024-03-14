package hyapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetBinary 获取二进制数据
func GetBinary(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// 如果响应状态码不是 200 OK，返回错误
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// GetJSON 获取 JSON 数据
func GetJSON[T any](url string) (r T, err error) {
	data, err := GetBinary(url)
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &r); err != nil {
		return
	}

	return
}

// Response 是通用的 API 响应
type Response[T any] struct {
	ReturnCode int    `json:"retcode"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

// OK 判断响应是否成功
func (r *Response[T]) OK() bool {
	return r.ReturnCode == 0
}

// Get 获取 API 响应
func Get[T any](url string) (resp Response[T], err error) {
	return GetJSON[Response[T]](url)
}

// GetData 获取 API 响应的数据
func GetData[T any](url string) (T, error) {
	resp, err := Get[T](url)
	if err != nil {
		var t T
		return t, err
	}

	if !resp.OK() {
		return resp.Data, fmt.Errorf(resp.Message)
	}

	return resp.Data, nil
}
