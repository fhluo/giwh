package hywiki

import (
	"regexp"
	"strings"
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

func (m Menu) VarName() string {
	s := strings.ReplaceAll(m.Name, "'", "")
	return regexp.MustCompile(`\s`).ReplaceAllString(s, "")
}
