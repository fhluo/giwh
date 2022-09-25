package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"io/fs"
	"log"
	"net/http"
)

var (
	//go:embed frontend/dist
	frontendAssets embed.FS

	//go:embed assets
	assets        embed.FS
	assetsHandler http.Handler
)

func init() {
	log.SetFlags(0)

	r, err := fs.Sub(assets, "assets")
	if err != nil {
		log.Fatalln(err)
	}

	assetsHandler = http.FileServer(http.FS(r))
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:            "GIWH",
		Width:            1024,
		Height:           768,
		Assets:           frontendAssets,
		AssetsHandler:    assetsHandler,
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
