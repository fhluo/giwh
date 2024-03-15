package assets

import (
	_ "embed"
	"github.com/bytedance/sonic"
	"github.com/samber/lo"
	"sync"
)

//go:generate go run ../cmd/giwh-dev dl assets

//go:embed assets.json
var JSON []byte

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

var (
	once   sync.Once
	assets = New()
)

func Get() Assets {
	once.Do(func() {
		lo.Must0(sonic.Unmarshal(JSON, &assets))
	})
	return assets
}
