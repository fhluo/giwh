package hywiki

import (
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
)

type JSONResponse[T any] struct {
	ReturnCode int    `json:"retcode"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

func (r *JSONResponse[T]) OK() bool {
	return r.ReturnCode == 0
}

func GetJSONResponseData[T any](resp *http.Response) (T, error) {
	defer func() {
		_ = resp.Body.Close()
	}()

	var jsonResponse JSONResponse[T]

	if resp.StatusCode != http.StatusOK {
		return jsonResponse.Data, fmt.Errorf(resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return jsonResponse.Data, err
	}

	if err = sonic.Unmarshal(data, &jsonResponse); err != nil {
		return jsonResponse.Data, err
	}

	if !jsonResponse.OK() {
		return jsonResponse.Data, fmt.Errorf(jsonResponse.Message)
	}

	return jsonResponse.Data, nil
}
