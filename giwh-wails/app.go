package main

import (
	assetsPkg "github.com/fhluo/giwh/assets"
	"github.com/fhluo/giwh/common/config"
	"github.com/fhluo/giwh/common/i18n"
	"github.com/fhluo/giwh/gacha-logs/gacha"
	"golang.org/x/net/context"
	"os"
	"unsafe"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetDefaultSharedWish() string {
	return gacha.CharacterEventWish
}

func (a *App) GetAssets() string {
	return unsafe.String(unsafe.SliceData(assetsPkg.JSON), len(assetsPkg.JSON))
}

func (a *App) GetLanguage() string {
	if config.Language.Get() != "" {
		return i18n.Match(config.Language.Get()).Key
	} else {
		return i18n.Default().Key
	}
}

func (a *App) GetLocale(lang string) string {
	data, err := i18n.ReadLocaleFile(i18n.Match(lang))
	if err != nil {
		return ""
	}

	return unsafe.String(unsafe.SliceData(data), len(data))
}

func (a *App) GetItems() string {
	data, err := os.ReadFile(config.WishHistoryPath.Get())
	if err != nil {
		return ""
	}

	return unsafe.String(unsafe.SliceData(data), len(data))
}
