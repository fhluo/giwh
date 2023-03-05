package pipeline

import "github.com/fhluo/giwh/pkg/wish"

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

//func (p Pipeline) Update() Pipeline {
//
//}

//func (p Pipeline) Add(items_ ...wish.Item) Pipeline {
//	return Pipeline{}
//}
//
//func (p Pipeline) Reverse() {
//	lo.Reverse(p)
//}
//
//func (p Pipeline) IDAscending() bool {
//	for i := 1; i < len(p); i++ {
//		if p[i-1].ID > p[i].ID {
//			return false
//		}
//	}
//
//	return true
//}
//
//func (p Pipeline) IDDescending() bool {
//	for i := 1; i < len(p); i++ {
//		if p[i-1].ID < p[i].ID {
//			return false
//		}
//	}
//
//	return true
//}
//
//func (p Pipeline) GroupByUID() map[int]Pipeline {
//	result := lo.GroupBy(p, func(item wish.Item) int {
//		return item.UID
//	})
//
//	pipelines := make(map[int]Pipeline)
//	for k, v := range result {
//		pipelines[k] = v
//	}
//
//	return pipelines
//}
//
//func (p Pipeline) GroupBySharedWishType() map[api.SharedWishType]Pipeline {
//	result := lo.GroupBy(p, func(item wish.Item) api.SharedWishType {
//		return item.WishType.Shared()
//	})
//
//	pipelines := make(map[api.SharedWishType]Pipeline)
//	for k, v := range result {
//		pipelines[k] = v
//	}
//
//	return pipelines
//}
//
//func (p Pipeline) SortByIDAscending() {
//	switch {
//	case p.IDAscending():
//	case p.IDDescending():
//		p.Reverse()
//	default:
//		slices.SortFunc(p, func(a wish.Item, b wish.Item) bool {
//			return a.ID < b.ID
//		})
//	}
//}
//
//func (p Pipeline) SortByIDDescending() {
//	switch {
//	case p.IDAscending():
//		p.Reverse()
//	case p.IDDescending():
//	default:
//		slices.SortFunc(p, func(a wish.Item, b wish.Item) bool {
//			return a.ID > b.ID
//		})
//	}
//}
//
//func (p Pipeline) Unique() Pipeline {
//	return lo.UniqBy(p, func(item wish.Item) int64 {
//		return item.ID
//	})
//}
//
//func (p Pipeline) FilterByUID(uid int) Pipeline {
//	return lo.Filter(p, func(item wish.Item, _ int) bool {
//		return item.UID == uid
//	})
//}
//
//func (p Pipeline) FilterByWishType(wishTypes ...api.WishType) Pipeline {
//	switch len(wishTypes) {
//	case 0:
//		return nil
//	case 1:
//		return lo.Filter(p, func(item wish.Item, _ int) bool {
//			return item.WishType == wishTypes[0]
//		})
//	default:
//		return lo.Filter(p, func(item wish.Item, _ int) bool {
//			return lo.Contains(wishTypes, item.WishType)
//		})
//	}
//}
//
//func (p Pipeline) FilterBySharedWishType(wishTypes ...api.SharedWishType) Pipeline {
//	switch len(wishTypes) {
//	case 0:
//		return nil
//	case 1:
//		return lo.Filter(p, func(item wish.Item, _ int) bool {
//			return item.WishType.Shared() == wishTypes[0]
//		})
//	default:
//		return lo.Filter(p, func(item wish.Item, _ int) bool {
//			return lo.Contains(wishTypes, item.WishType.Shared())
//		})
//	}
//}
//
//func (p Pipeline) FilterByRarity(rarities ...api.Rarity) Pipeline {
//	switch len(rarities) {
//	case 0:
//		return nil
//	case 1:
//		return lo.Filter(p, func(item wish.Item, _ int) bool {
//			return item.Rarity == rarities[0]
//		})
//	default:
//		return lo.Filter(p, func(item wish.Item, _ int) bool {
//			return lo.Contains(rarities, item.Rarity)
//		})
//	}
//}
//
//func (p Pipeline) UIDs() []int {
//	return lo.Uniq(lo.Map(p, func(item wish.Item, _ int) int {
//		return item.UID
//	}))
//}
//
//func (p Pipeline) SharedWishTypes() []api.SharedWishType {
//	return lo.Uniq(lo.Map(p, func(item wish.Item, _ int) api.SharedWishType {
//		return item.WishType.Shared()
//	}))
//}
//
//func (p Pipeline) Progress5Star() int {
//	if len(p.UIDs()) != 1 || len(p.SharedWishTypes()) != 1 {
//		return -1
//	}
//
//	p.SortByIDDescending()
//	return slices.IndexFunc(p, func(item wish.Item) bool {
//		return item.Rarity == api.Star5
//	})
//}
//
//func (p Pipeline) Progress4Star() int {
//	if len(p.UIDs()) != 1 || len(p.SharedWishTypes()) != 1 {
//		return -1
//	}
//
//	p.SortByIDDescending()
//	return slices.IndexFunc(p, func(item wish.Item) bool {
//		return item.Rarity == api.Star4
//	})
//}
//
//func (p Pipeline) Pulls5Stars() map[int64]int {
//	if len(p.UIDs()) != 1 || len(p.SharedWishTypes()) != 1 {
//		return nil
//	}
//
//	p.SortByIDAscending()
//	pulls := make(map[int64]int)
//	prev := 0
//
//	for i, item := range p {
//		if item.Rarity == api.Star5 {
//			pulls[item.ID] = i - prev
//			prev = i
//		}
//	}
//
//	return pulls
//}
//
//func (p Pipeline) Pulls4Stars() map[int64]int {
//	if len(p.UIDs()) != 1 || len(p.SharedWishTypes()) != 1 {
//		return nil
//	}
//
//	p.SortByIDAscending()
//	pulls := make(map[int64]int)
//	prev := 0
//
//	for i, item := range p {
//		if item.Rarity == api.Star4 {
//			pulls[item.ID] = i - prev
//			prev = i
//		}
//	}
//
//	return pulls
//}
