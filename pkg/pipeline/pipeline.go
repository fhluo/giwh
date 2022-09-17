package pipeline

import (
	"github.com/fhluo/giwh/pkg/api"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type Pipeline struct {
	items []*api.Item

	_4star *int
	_5star *int
}

func New(items []*api.Item) Pipeline {
	return Pipeline{items: items}
}

func (p Pipeline) Count() int {
	return len(p.items)
}

func (p Pipeline) Count4Star() int {
	_4star, _ := p.Count4StarAnd5Star()
	return _4star
}

func (p Pipeline) Count5Star() int {
	_, _5star := p.Count4StarAnd5Star()
	return _5star
}

func (p Pipeline) Count4StarAnd5Star() (int, int) {
	if p._4star != nil && p._5star != nil {
		return *p._4star, *p._5star
	}

	p._4star = new(int)
	p._5star = new(int)
	for _, item := range p.items {
		switch item.Rarity {
		case api.Star4:
			*p._4star++
		case api.Star5:
			*p._5star++
		}
	}

	return *p._4star, *p._5star
}

func (p Pipeline) First() *api.Item {
	return p.items[0]
}

func (p Pipeline) Last() *api.Item {
	return p.items[len(p.items)-1]
}

func (p Pipeline) Traverse(f func(e *api.Item)) {
	for _, e := range p.items {
		f(e)
	}
}

func (p Pipeline) Append(items []*api.Item) Pipeline {
	return Pipeline{items: append(p.items, items...)}
}

func (p Pipeline) Items() []*api.Item {
	return p.items
}

func (p Pipeline) Copy() Pipeline {
	items := make([]*api.Item, len(p.items))
	copy(items, p.items)
	return Pipeline{items: items}
}

func (p Pipeline) Reverse() {
	lo.Reverse(p.items)
}

func (p Pipeline) IDAscending() bool {
	for i := 1; i < len(p.items); i++ {
		if p.items[i-1].ID > p.items[i].ID {
			return false
		}
	}

	return true
}

func (p Pipeline) IDDescending() bool {
	for i := 1; i < len(p.items); i++ {
		if p.items[i-1].ID < p.items[i].ID {
			return false
		}
	}

	return true
}

func (p Pipeline) GroupByUID() map[int]Pipeline {
	result := lo.GroupBy(p.items, func(item *api.Item) int {
		return item.UID
	})

	pipelines := make(map[int]Pipeline)
	for k, v := range result {
		pipelines[k] = Pipeline{items: v}
	}

	return pipelines
}

func (p Pipeline) GroupBySharedWishType() map[api.SharedWishType]Pipeline {
	result := lo.GroupBy(p.items, func(item *api.Item) api.SharedWishType {
		return item.WishType.Shared()
	})

	pipelines := make(map[api.SharedWishType]Pipeline)
	for k, v := range result {
		pipelines[k] = Pipeline{items: v}
	}

	return pipelines
}

func (p Pipeline) SortByIDAscending() {
	switch {
	case p.IDAscending():
	case p.IDDescending():
		p.Reverse()
	default:
		slices.SortFunc(p.items, func(a *api.Item, b *api.Item) bool {
			return a.ID < b.ID
		})
	}
}

func (p Pipeline) SortByIDDescending() {
	switch {
	case p.IDAscending():
		p.Reverse()
	case p.IDDescending():
	default:
		slices.SortFunc(p.items, func(a *api.Item, b *api.Item) bool {
			return a.ID > b.ID
		})
	}
}

func (p Pipeline) Unique() Pipeline {
	return Pipeline{
		items: lo.UniqBy(p.items, func(item *api.Item) int64 {
			return item.ID
		}),
	}
}

func (p Pipeline) FilterByUID(uid int) Pipeline {
	return Pipeline{
		items: lo.Filter(p.items, func(item *api.Item, _ int) bool {
			return item.UID == uid
		}),
	}
}

func (p Pipeline) FilterByWishType(wishTypes ...api.WishType) Pipeline {
	switch len(wishTypes) {
	case 0:
		return Pipeline{}
	case 1:
		return Pipeline{
			items: lo.Filter(p.items, func(item *api.Item, _ int) bool {
				return item.WishType == wishTypes[0]
			}),
		}
	default:
		return Pipeline{
			items: lo.Filter(p.items, func(item *api.Item, _ int) bool {
				return lo.Contains(wishTypes, item.WishType)
			}),
		}
	}
}

func (p Pipeline) FilterBySharedWishType(wishTypes ...api.SharedWishType) Pipeline {
	switch len(wishTypes) {
	case 0:
		return Pipeline{}
	case 1:
		return Pipeline{
			items: lo.Filter(p.items, func(item *api.Item, _ int) bool {
				return item.WishType.Shared() == wishTypes[0]
			}),
		}
	default:
		return Pipeline{
			items: lo.Filter(p.items, func(item *api.Item, _ int) bool {
				return lo.Contains(wishTypes, item.WishType.Shared())
			}),
		}
	}
}

func (p Pipeline) FilterByRarity(rarities ...api.Rarity) Pipeline {
	switch len(rarities) {
	case 0:
		return Pipeline{}
	case 1:
		return Pipeline{
			items: lo.Filter(p.items, func(item *api.Item, _ int) bool {
				return item.Rarity == rarities[0]
			}),
		}
	default:
		return Pipeline{
			items: lo.Filter(p.items, func(item *api.Item, _ int) bool {
				return lo.Contains(rarities, item.Rarity)
			}),
		}
	}
}

func (p Pipeline) GetIndex(f func(item *api.Item) bool) []int {
	var result []int
	for i, item := range p.items {
		if f(item) {
			result = append(result, i)
		}
	}
	return result
}

func (p Pipeline) Progress() map[api.SharedWishType]map[api.Rarity]int {
	p.SortByIDDescending()
	result := make(map[api.SharedWishType]map[api.Rarity]int)
	done := make(map[api.SharedWishType]map[api.Rarity]bool)

	for _, wishType := range api.SharedWishTypes {
		result[wishType] = make(map[api.Rarity]int)
		done[wishType] = make(map[api.Rarity]bool)
	}

	var wishType api.SharedWishType
	for _, item := range p.items {
		wishType = item.WishType.Shared()

		switch item.Rarity {
		case api.Star4:
			done[wishType][api.Star4] = true
			if !done[wishType][api.Star5] {
				result[wishType][api.Star5]++
			}
		case api.Star5:
			done[wishType][api.Star5] = true
			if !done[wishType][api.Star4] {
				result[wishType][api.Star4]++
			}
		default:
			if !done[wishType][api.Star4] {
				result[wishType][api.Star4]++
			}
			if !done[wishType][api.Star5] {
				result[wishType][api.Star5]++
			}
		}
	}

	return result
}

func (p Pipeline) Pulls() map[api.SharedWishType]map[int64]int {
	p.SortByIDAscending()
	progress := make(map[api.SharedWishType]map[int64]int)
	progress4Star := make(map[api.SharedWishType]int)
	progress5Star := make(map[api.SharedWishType]int)

	for _, wishType := range api.SharedWishTypes {
		progress[wishType] = make(map[int64]int)
	}

	var wishType api.SharedWishType
	for _, item := range p.items {
		wishType = item.WishType.Shared()

		switch item.Rarity {
		case api.Star4:
			progress[wishType][item.ID] = progress4Star[wishType] + 1
			progress4Star[wishType] = 0
			progress5Star[wishType]++
		case api.Star5:
			progress[wishType][item.ID] = progress5Star[wishType] + 1
			progress4Star[wishType]++
			progress5Star[wishType] = 0
		default:
			progress4Star[wishType]++
			progress5Star[wishType]++
		}
	}

	return progress
}
