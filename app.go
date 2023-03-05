package main

import (
	assetsPkg "github.com/fhluo/giwh/assets"
	"github.com/fhluo/giwh/i18n"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/wish"
	"golang.org/x/net/context"
	"os"
	"strconv"
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
	return strconv.Itoa(int(wish.CharacterEventWishAndCharacterEventWish2))
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
	data, err := os.ReadFile(config.WishHistoryPath)
	if err != nil {
		return ""
	}

	return unsafe.String(unsafe.SliceData(data), len(data))
}
