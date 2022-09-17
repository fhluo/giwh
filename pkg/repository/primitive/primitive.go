package primitive

import (
	"fmt"
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/pipeline"
	"github.com/fhluo/giwh/pkg/repository"
	"github.com/samber/lo"
)

type Repository struct {
	pipeline.Pipeline
	index map[int]map[api.SharedWishType]pipeline.Pipeline
}

func New(filename string) (repository.Repository, error) {
	items, err := repository.Load(filename)
	if err != nil {
		return nil, err
	}

	r := new(Repository)

	r.Pipeline = pipeline.New(items)

	r.index = make(map[int]map[api.SharedWishType]pipeline.Pipeline)
	for uid, p := range r.GroupByUID() {
		r.index[uid] = p.GroupBySharedWishType()
	}

	return r, nil
}

func (r *Repository) GetUIDs() []int {
	r.SortByIDDescending()
	return lo.Uniq(lo.Map(r.Items(), func(item *api.Item, _ int) int {
		return item.UID
	}))
}

func (r *Repository) GetProgress(uid int, wishType api.SharedWishType, rarity api.Rarity) int {
	p := r.index[uid][wishType]
	p.SortByIDDescending()

	_, i, _ := lo.FindIndexOf(p.Items(), func(item *api.Item) bool {
		return item.Rarity == rarity
	})
	return i
}

func (r *Repository) GetItems(uid int, wishType api.SharedWishType, rarity api.Rarity) []repository.Item {
	p := r.index[uid][wishType]
	p.SortByIDAscending()

	items := lo.Map(p.FilterByRarity(rarity).Items(), func(item *api.Item, _ int) repository.Item {
		return repository.Item{
			ID:   item.ID,
			Name: item.Name,
		}
	})

	fmt.Println(items)

	return items
}

func (r *Repository) GetPulls(uid int, wishType api.SharedWishType, id int64) int {
	p := r.index[uid][wishType]
	p.SortByIDAscending()

	target, i, ok := lo.FindIndexOf(p.Items(), func(item *api.Item) bool {
		return item.ID == id
	})
	if !ok {
		return -1
	}

	_, j, ok := lo.FindIndexOf(p.Items(), func(item *api.Item) bool {
		return item.ID < target.ID && item.Rarity == target.Rarity
	})
	if !ok {
		return -1
	}

	return i - j
}
