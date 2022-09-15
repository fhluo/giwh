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
	r, err := primitive.New(config.WishHistoryPath)
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

func (a *App) GetSharedWishName(wishType string) string {
	return i18n.GetSharedWishName(wishType)
}

func (a *App) GetPity(rarity string, wishType string) int {
	return api.Pity(rarity, wishType)
}

func (a *App) GetSharedWishTypes() []string {
	return api.SharedWishTypes
}
