package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"log"
	"net/http"
)

var (
	//go:embed frontend/dist
	assets embed.FS
	//go:embed assets/images
	images embed.FS
)

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
		AssetsHandler:    http.FileServer(http.FS(images)),
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup:        app.startup,
		Bind:             []interface{}{app},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Mica,
		},
	})

	if err != nil {
		log.Println(err)
	}
}
