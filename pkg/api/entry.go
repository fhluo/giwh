package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
)

type FilterValue struct {
	Values []string `json:"values"`
}

type Entry struct {
	EntryPageID  int                     `json:"entry_page_id,string"`
	Name         string                  `json:"name"`
	IconURL      string                  `json:"icon_url"`
	DisplayField map[string]any          `json:"display_field"`
	FilterValues map[string]*FilterValue `json:"filter_values"`
}

var (
	NonAlphanumeric = regexp.MustCompile(`\W`)
	Special         = regexp.MustCompile(`['"]`)
)

func (entry *Entry) VarName() string {
	return NonAlphanumeric.ReplaceAllString(entry.Name, "")
}

func (entry *Entry) Filename() string {
	return Special.ReplaceAllString(entry.Name, "") + path.Ext(entry.IconURL)
}

type EntryList struct {
	List  []*Entry `json:"list"`
	Total int      `json:"total,string"`
}

type Parameters struct {
	Filters  []string `json:"filters"`
	MenuID   int      `json:"menu_id,string"`
	PageNum  int      `json:"page_num"`
	PageSize int      `json:"page_size"`
	UseES    bool     `json:"use_es"`
}

func NewParameters() *Parameters {
	return &Parameters{
		Filters:  []string{},
		MenuID:   CharactersMenu,
		PageNum:  1,
		PageSize: 30,
		UseES:    true,
	}
}

const GetEntryPageListURL = `https://sg-wiki-api.hoyolab.com/hoyowiki/wapi/get_entry_page_list`

func GetEntryPageList(parameters *Parameters) (*EntryList, error) {
	data, err := json.Marshal(parameters)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, GetEntryPageListURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RPC-Language", "en-us")
	req.Header.Set("Referer", "https://wiki.hoyolab.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r JSONResponse[*EntryList]
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	if r.RetCode != 0 {
		return nil, fmt.Errorf(r.Message)
	}

	return r.Data, nil
}

const (
	CharactersMenu = 2
	WeaponsMenu    = 4
)

func GetAllEntries(menu int) ([]*Entry, error) {
	p := NewParameters()
	p.MenuID = menu

	var entries []*Entry
	r, err := GetEntryPageList(p)
	if err != nil {
		return nil, err
	}
	entries = append(entries, r.List...)

	for p.PageNum*p.PageSize < r.Total {
		p.PageNum++
		r, err = GetEntryPageList(p)
		if err != nil {
			return nil, err
		}
		entries = append(entries, r.List...)
	}

	return entries, nil
}
