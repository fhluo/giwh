package assets

import (
	_ "embed"
)

var (
	//go:embed assets.json
	AssetsJSON string
)

//go:generate go run ../cmd/giwh-dev dl assets

type Assets struct {
	Characters map[string]string `json:"characters"`
	Weapons    map[string]string `json:"weapons"`
}

func New() Assets {
	return Assets{
		Characters: make(map[string]string),
		Weapons:    make(map[string]string),
	}
}
