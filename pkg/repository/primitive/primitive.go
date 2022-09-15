package primitive

import (
	"github.com/fhluo/giwh/pkg/pipeline"
	"github.com/fhluo/giwh/pkg/repository"
	"github.com/samber/lo"
)

type Repository struct {
	pipeline.Pipeline
	index map[string]map[string]pipeline.Pipeline
}

func New(filename string) (repository.Repository, error) {
	items, err := repository.Load(filename)
	if err != nil {
		return nil, err
	}

	r := new(Repository)

	r.Pipeline, err = pipeline.New(items)
	if err != nil {
		return nil, err
	}

	r.index = make(map[string]map[string]pipeline.Pipeline)
	for uid, p := range r.GroupByUID() {
		r.index[uid] = p.GroupByWishType()
	}

	return r, nil
}

func (r *Repository) GetUIDs() []string {
	r.SortByIDDescending()
	return lo.Uniq(lo.Map(r.Elements(), func(element *pipeline.Element, _ int) string {
		return element.UID
	}))
}

func (r *Repository) GetProgress(uid string, wishType string, rarity string) int {
	p := r.index[uid][wishType]
	p.SortByIDDescending()

	_, i, _ := lo.FindIndexOf(p.Elements(), func(element *pipeline.Element) bool {
		return element.Rarity == rarity
	})
	return i
}

func (r *Repository) GetPulls(uid string, wishType string, id int64) int {
	p := r.index[uid][wishType]
	p.SortByIDAscending()

	target, i, ok := lo.FindIndexOf(p.Elements(), func(element *pipeline.Element) bool {
		return element.ID == id
	})
	if !ok {
		return -1
	}

	_, j, ok := lo.FindIndexOf(p.Elements(), func(element *pipeline.Element) bool {
		return element.ID < target.ID && element.Rarity == target.Rarity
	})
	if !ok {
		return -1
	}

	return i - j
}
