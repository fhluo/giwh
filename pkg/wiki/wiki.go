package wiki

import (
	"net/http"
)

type Wiki struct {
	Language Language
}

func (w Wiki) request(url string) Request {
	return Request{
		URL:      url,
		Language: w.Language,
		Client:   http.DefaultClient,
	}
}

func (w Wiki) GetMenus() ([]Menu, error) {
	resp, err := w.request("https://sg-wiki-api-static.hoyolab.com/hoyowiki/wapi/get_menus").Get()
	if err != nil {
		return nil, err
	}

	data, err := GetJSONResponseData[*GetMenusResponseData](resp)
	if err != nil {
		return nil, err
	}

	return data.Menus, nil
}

func (w Wiki) GetEntryPageList(payload GetEntryPageListPayload) ([]Entry, error) {
	resp, err := w.request("https://sg-wiki-api.hoyolab.com/hoyowiki/wapi/get_entry_page_list").JSONPost(payload)
	if err != nil {
		return nil, err
	}

	data, err := GetJSONResponseData[*GetEntryPageListResponseData](resp)
	if err != nil {
		return nil, err
	}

	return data.List, nil
}

func (w Wiki) GetMenuFilters(payload GetMenuFiltersPayload) ([]Filter, error) {
	resp, err := w.request("https://sg-wiki-api-static.hoyolab.com/hoyowiki/wapi/get_menu_filters").QueryGet(payload.Values())
	if err != nil {
		return nil, err
	}

	data, err := GetJSONResponseData[*GetMenuFiltersResponseData](resp)
	if err != nil {
		return nil, err
	}

	return data.Filters, nil
}
