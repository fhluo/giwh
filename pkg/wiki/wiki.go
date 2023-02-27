package wiki

import (
	"github.com/fhluo/giwh/i18n"
	"net/http"
)

//go:generate go run ../../cmd/giwh-dev gen menus

type Wiki struct {
	Language i18n.Language
}

func New(lang i18n.Language) Wiki {
	return Wiki{Language: lang}
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

// GetLeafMenus returns all leaf menus.
func (w Wiki) GetLeafMenus() ([]Menu, error) {
	r, err := w.GetMenus()
	if err != nil {
		return nil, err
	}

	var menus []Menu

	for _, menu := range r {
		for _, m := range menu.LeafMenus() {
			menus = append(menus, m)
		}
	}

	return menus, nil
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

// GetEntries returns all the entries of a menu.
func (w Wiki) GetEntries(menuID int) ([]Entry, error) {
	payload := NewGetEntryPageListPayload(menuID)

	var entries []Entry

	for {
		result, err := w.GetEntryPageList(payload)
		if err != nil {
			return nil, err
		}

		if len(result) == 0 {
			break
		}

		entries = append(entries, result...)
		payload.PageNum++
	}

	return entries, nil
}

// GetMenusEntries GetEntries returns all the entries of menus.
func (w Wiki) GetMenusEntries(menuIDs ...int) (menusEntries map[int][]Entry, err error) {
	menusEntries = make(map[int][]Entry)
	for _, menuID := range menuIDs {
		menusEntries[menuID], err = w.GetEntries(menuID)
		if err != nil {
			return
		}
	}

	return menusEntries, nil
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
