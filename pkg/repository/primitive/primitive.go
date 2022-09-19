package primitive

import (
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/pipeline"
	"github.com/fhluo/giwh/pkg/repository"
	"github.com/samber/lo"
)

type Repository struct {
	pipeline.Items
	index map[int]map[api.SharedWishType]pipeline.Items
}

func Load(filename string) (repository.Repository, error) {
	items, err := repository.Load(filename)
	if err != nil {
		return nil, err
	}

	r := new(Repository)
	r.Items = items

	r.index = make(map[int]map[api.SharedWishType]pipeline.Items)
	for uid, p := range r.GroupByUID() {
		r.index[uid] = p.GroupBySharedWishType()
	}

	return r, nil
}

func (r *Repository) GetUIDs() []int {
	r.SortByIDDescending()
	return lo.Uniq(lo.Map(r.Items, func(item *api.Item, _ int) int {
		return item.UID
	}))
}

func (r *Repository) Get5StarProgress(uid int, wishType api.SharedWishType) int {
	return r.index[uid][wishType].Progress5Star()
}

func (r *Repository) Get4StarProgress(uid int, wishType api.SharedWishType) int {
	return r.index[uid][wishType].Progress4Star()
}

func (r *Repository) Get5Stars(uid int, wishType api.SharedWishType) []repository.Item {
	items := r.index[uid][wishType]
	pulls := items.Pulls5Stars()

	return lo.Map(items.FilterByRarity(api.Star5), func(item *api.Item, _ int) repository.Item {
		return repository.Item{
			Item:  item,
			Pulls: pulls[item.ID],
		}
	})
}

func (r *Repository) Get4Stars(uid int, wishType api.SharedWishType) []repository.Item {
	items := r.index[uid][wishType]
	pulls := items.Pulls4Stars()

	return lo.Map(items.FilterByRarity(api.Star4), func(item *api.Item, _ int) repository.Item {
		return repository.Item{
			Item:  item,
			Pulls: pulls[item.ID],
		}
	})
}
