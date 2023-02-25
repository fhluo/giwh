package wiki

import (
	"encoding/json"
	"github.com/samber/lo"
	"net/url"
	"strconv"
)

type Menu struct {
	ID       int    `json:"id,string"`
	Name     string `json:"name"`
	HasPage  bool   `json:"has_page"`
	Style    string `json:"style"`
	SubMenus []Menu `json:"sub_menus"`
}

func (m Menu) IsLeaf() bool {
	return len(m.SubMenus) == 0
}

func (m Menu) LeafMenus() []Menu {
	var leafMenus []Menu
	for _, menu := range m.SubMenus {
		if menu.IsLeaf() {
			leafMenus = append(leafMenus, menu)
		} else {
			leafMenus = append(leafMenus, menu.LeafMenus()...)
		}
	}
	return leafMenus
}

type GetMenusResponseData struct {
	Menus []Menu `json:"menus"`
}

type FilterValues []FilterValue

func (f FilterValues) MarshalJSON() ([]byte, error) {
	return json.Marshal(lo.Map(f, func(fv FilterValue, _ int) string { return strconv.Itoa(fv.ID) }))
}

type GetEntryPageListPayload struct {
	Filters  FilterValues `json:"filters"`
	MenuID   int          `json:"menu_id,string"`
	PageNum  int          `json:"page_num"`
	PageSize int          `json:"page_size"`
	UseES    bool         `json:"use_es"`
}

func NewGetEntryPageListPayload(menu Menu, filters ...FilterValue) GetEntryPageListPayload {
	return GetEntryPageListPayload{
		Filters:  filters,
		MenuID:   menu.ID,
		PageNum:  1,
		PageSize: 30,
		UseES:    true,
	}
}

type GetEntryPageListResponseData struct {
	List  []Entry `json:"list"`
	Total int     `json:"total,string"`
}

type Entry struct {
	EntryPageID  int            `json:"entry_page_id,string"`
	Name         string         `json:"name"`
	IconURL      string         `json:"icon_url"`
	DisplayField map[string]any `json:"display_field"`
	FilterValues map[string]any `json:"filter_values"`
}

type GetMenuFiltersPayload struct {
	MenuID int `json:"menu_id,string"`
}

func (payload GetMenuFiltersPayload) Values() url.Values {
	v := url.Values{}
	v.Set("menu_id", strconv.Itoa(payload.MenuID))
	return v
}

func NewGetMenuFiltersPayload(menu Menu) GetMenuFiltersPayload {
	return GetMenuFiltersPayload{
		MenuID: menu.ID,
	}
}

type GetMenuFiltersResponseData struct {
	Filters []Filter `json:"filters"`
}

type Filter struct {
	Key    string        `json:"key"`
	Text   string        `json:"text"`
	Values []FilterValue `json:"values"`
}

type FilterValue struct {
	ID    int    `json:"id,string"`
	Value string `json:"value"`
}
