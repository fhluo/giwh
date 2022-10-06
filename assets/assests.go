package assets

import (
	_ "embed"
	"encoding/json"
)

var (
	//go:embed  images/characters.json
	charactersData []byte
	//go:embed  images/weapons.json
	weaponsData []byte

	Characters map[string]string
	Weapons    map[string]string
)

func init() {
	if err := json.Unmarshal(charactersData, &Characters); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(weaponsData, &Weapons); err != nil {
		panic(err)
	}
}
