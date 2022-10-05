package main

import (
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/i18n"
	"github.com/fhluo/giwh/pkg/repository"
	"github.com/fhluo/giwh/pkg/repository/primitive"
	"golang.org/x/net/context"
	"log"
)

type App struct {
	ctx context.Context
	repository.Repository
}

func NewApp() *App {
	r, err := primitive.Load(config.WishHistoryPath)
	if err != nil {
		log.Fatalln(err)
	}
	return &App{Repository: r}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if config.GetLanguage() != "" {
		i18n.Language = config.GetLanguage()
	}
}

func (a *App) GetWishName(wishType api.WishType) string {
	return i18n.GetWishName(wishType.Str())
}

func (a *App) GetSharedWishName(wishType api.SharedWishType) string {
	return i18n.GetSharedWishName(wishType.Str())
}

func (a *App) GetPity(rarity api.Rarity, wishType api.SharedWishType) int {
	return wishType.Pity(rarity)
}

func (a *App) GetSharedWishTypes() []api.SharedWishType {
	return api.SharedWishTypes
}
