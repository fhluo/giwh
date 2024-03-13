package hoyo_wiki

import (
	"encoding/json"
	"github.com/samber/lo"
	"net/url"
	"path"
	"regexp"
	"strconv"
)

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

func NewGetEntryPageListPayload(menuID int, filters ...FilterValue) GetEntryPageListPayload {
	return GetEntryPageListPayload{
		Filters:  filters,
		MenuID:   menuID,
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

var (
	NonAlphanumeric = regexp.MustCompile(`\W`)
	Special         = regexp.MustCompile(`['"]`)
)

func (entry Entry) VarName() string {
	return NonAlphanumeric.ReplaceAllString(entry.Name, "")
}

func (entry Entry) Filename() string {
	return Special.ReplaceAllString(entry.Name, "") + path.Ext(entry.IconURL)
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
