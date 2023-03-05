package pipeline

import (
	"github.com/fhluo/giwh/pkg/wish"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type Pipeline struct {
	items []wish.Item
	index map[int64]struct{}
}

func New(items []wish.Item) *Pipeline {
	index := make(map[int64]struct{})
	for _, item := range items {
		index[item.ID] = struct{}{}
	}
	return &Pipeline{
		items: items,
		index: index,
	}
}

func (p *Pipeline) Len() int {
	return len(p.items)
}

func (p *Pipeline) First() wish.Item {
	return p.items[0]
}

func (p *Pipeline) Last() wish.Item {
	return p.items[len(p.items)-1]
}

func (p *Pipeline) Items() []wish.Item {
	return p.items
}

func (p *Pipeline) Contains(item wish.Item) bool {
	_, ok := p.index[item.ID]
	return ok
}

func (p *Pipeline) ContainsAny(items ...wish.Item) bool {
	for _, item := range items {
		if _, ok := p.index[item.ID]; ok {
			return true
		}
	}
	return false
}

func (p *Pipeline) Append(items ...wish.Item) {
	for _, item := range items {
		if _, ok := p.index[item.ID]; ok {
			continue
		}
		p.index[item.ID] = struct{}{}
		p.items = append(p.items, item)
	}
}

func (p *Pipeline) Reverse() *Pipeline {
	lo.Reverse(p.items)
	return p
}

func (p *Pipeline) IDAscending() bool {
	for i := 1; i < len(p.items); i++ {
		if p.items[i-1].ID > p.items[i].ID {
			return false
		}
	}

	return true
}

func (p *Pipeline) IDDescending() bool {
	for i := 1; i < len(p.items); i++ {
		if p.items[i-1].ID < p.items[i].ID {
			return false
		}
	}

	return true
}

func (p *Pipeline) SortByIDAscending() *Pipeline {
	switch {
	case p.IDAscending():
	case p.IDDescending():
		p.Reverse()
	default:
		slices.SortFunc(p.items, func(a wish.Item, b wish.Item) bool {
			return a.ID < b.ID
		})
	}
	return p
}

func (p *Pipeline) SortByIDDescending() *Pipeline {
	switch {
	case p.IDAscending():
		p.Reverse()
	case p.IDDescending():
	default:
		slices.SortFunc(p.items, func(a wish.Item, b wish.Item) bool {
			return a.ID > b.ID
		})
	}
	return p
}

func (p *Pipeline) GroupByUID() map[int][]wish.Item {
	return lo.GroupBy(p.items, func(item wish.Item) int {
		return item.UID
	})
}

func (p *Pipeline) GroupBySharedWish() map[wish.Type][]wish.Item {
	return lo.GroupBy(p.items, func(item wish.Item) wish.Type {
		return item.SharedWishType()
	})
}

func (p *Pipeline) Unique() *Pipeline {
	return New(lo.UniqBy(p.items, func(item wish.Item) int64 {
		return item.ID
	}))
}

func (p *Pipeline) FilterByUID(uid int) *Pipeline {
	return New(lo.Filter(p.items, func(item wish.Item, _ int) bool {
		return item.UID == uid
	}))
}

func (p *Pipeline) FilterByWish(types ...wish.Type) *Pipeline {
	switch len(types) {
	case 0:
		return nil
	case 1:
		return New(lo.Filter(p.items, func(item wish.Item, _ int) bool {
			return item.WishType == types[0]
		}))
	default:
		return New(lo.Filter(p.items, func(item wish.Item, _ int) bool {
			return slices.Contains(types, item.WishType)
		}))
	}
}

func (p *Pipeline) FilterBySharedWish(types ...wish.Type) *Pipeline {
	if slices.Contains(types, wish.CharacterEventWishAndCharacterEventWish2) {
		types = append(types, wish.CharacterEventWish2)
	}
	return p.FilterByWish(types...)
}

func (p *Pipeline) FilterByRarity(rarities ...int) *Pipeline {
	switch len(rarities) {
	case 0:
		return nil
	case 1:
		return New(lo.Filter(p.items, func(item wish.Item, _ int) bool {
			return item.Rarity == rarities[0]
		}))
	default:
		return New(lo.Filter(p.items, func(item wish.Item, _ int) bool {
			return slices.Contains(rarities, item.Rarity)
		}))
	}
}

func (p *Pipeline) UIDs() []int {
	return lo.Uniq(lo.Map(p.items, func(item wish.Item, _ int) int {
		return item.UID
	}))
}

func (p *Pipeline) SharedWishes() []wish.Type {
	return lo.Uniq(lo.Map(p.items, func(item wish.Item, _ int) wish.Type {
		return item.SharedWishType()
	}))
}

func (p *Pipeline) Progress5Star() int {
	if len(p.UIDs()) != 1 || len(p.SharedWishes()) != 1 {
		return -1
	}

	p.SortByIDDescending()
	return slices.IndexFunc(p.items, func(item wish.Item) bool {
		return item.Rarity == wish.FiveStar
	})
}

func (p *Pipeline) Progress4Star() int {
	if len(p.UIDs()) != 1 || len(p.SharedWishes()) != 1 {
		return -1
	}

	p.SortByIDDescending()
	return slices.IndexFunc(p.items, func(item wish.Item) bool {
		return item.Rarity == wish.FourStar
	})
}

func (p *Pipeline) Pulls5Stars() map[int64]int {
	if len(p.UIDs()) != 1 || len(p.SharedWishes()) != 1 {
		return nil
	}

	p.SortByIDAscending()
	pulls := make(map[int64]int)
	prev := 0

	for i, item := range p.items {
		if item.Rarity == wish.FiveStar {
			pulls[item.ID] = i - prev
			prev = i
		}
	}

	return pulls
}

func (p *Pipeline) Pulls4Stars() map[int64]int {
	if len(p.UIDs()) != 1 || len(p.SharedWishes()) != 1 {
		return nil
	}

	p.SortByIDAscending()
	pulls := make(map[int64]int)
	prev := 0

	for i, item := range p.items {
		if item.Rarity == wish.FourStar {
			pulls[item.ID] = i - prev
			prev = i
		}
	}

	return pulls
}
