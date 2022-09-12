package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"log"
)

//go:embed frontend/dist
var assets embed.FS

func init() {
	log.SetFlags(0)
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:            "GIWH",
		Width:            1024,
		Height:           768,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup:        app.startup,
		Bind:             []interface{}{app},
	})

	if err != nil {
		log.Println(err)
	}
}
