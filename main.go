package main

import (
	"embed"
	"github.com/fhluo/giwh/internal/config"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"golang.org/x/exp/slog"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

var (
	//go:embed web/dist
	assets embed.FS
	//go:embed assets/images
	images embed.FS
)

func init() {
	log.SetFlags(0)
}

func main() {
	defer config.Save()

	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "GIWH",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: http.FileServer(http.FS(lo.Must(fs.Sub(images, "assets")))),
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Mica,
		},
	})

	if err != nil {
		slog.Error(err.Error(), nil)
		os.Exit(1)
	}
}
