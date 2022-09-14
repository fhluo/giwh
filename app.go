package main

import (
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/i18n"
	"golang.org/x/net/context"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if config.GetLanguage() != "" {
		i18n.Language = config.GetLanguage()
	}
}

type Progress struct {
	Rarity string `json:"rarity"`
	Count  int    `json:"count"`
}

type StatResult struct {
	WishType   string     `json:"wishType"`
	Progresses []Progress `json:"progresses"`
}

func (a *App) GetSharedWishTypes() []string {
	return api.SharedWishTypes
}

func (a *App) Stat() []StatResult {
	progress := config.WishHistory.Progress()

	results := make([]StatResult, 0, len(api.SharedWishTypes))
	for _, wishType := range api.SharedWishTypes {
		results = append(results, StatResult{
			WishType: i18n.GetSharedWishName(wishType),
			Progresses: []Progress{
				{api.FiveStar, progress[wishType][api.FiveStar]},
				{api.FourStar, progress[wishType][api.FiveStar]},
			},
		})
	}

	return results
}
