package hoyo_wiki

import (
	"bytes"
	"github.com/fhluo/giwh/i18n"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	URL      string
	Language i18n.Language
	Client   *http.Client
}

func (r Request) New(method string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, r.URL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("referer", "https://wiki.hoyolab.com/")
	req.Header.Set("x-rpc-language", r.Language.Key)
	return req, nil
}

func (r Request) Do(req *http.Request) (*http.Response, error) {
	if r.Client == nil {
		r.Client = http.DefaultClient
	}

	return r.Client.Do(req)
}

func (r Request) Get() (resp *http.Response, err error) {
	req, err := r.New(http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	return r.Client.Do(req)
}

func (r Request) QueryGet(data url.Values) (resp *http.Response, err error) {
	req, err := r.New(http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = data.Encode()

	return r.Client.Do(req)
}

func (r Request) Post(contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := r.New(http.MethodPost, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return r.Client.Do(req)
}

func (r Request) JSONPost(value any) (resp *http.Response, err error) {
	data, err := sonic.Marshal(value)
	if err != nil {
		return nil, err
	}

	return r.Post("application/json", bytes.NewReader(data))
}
