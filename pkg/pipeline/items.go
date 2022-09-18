package pipeline

import (
	"github.com/fhluo/giwh/pkg/api"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type Items []*api.Item

func (items Items) Len() int {
	return len(items)
}

func (items Items) Equal(items_ Items) bool {
	return slices.EqualFunc(items, items_, func(item1, item2 *api.Item) bool {
		return item1.ID == item2.ID
	})
}

func (items Items) First() *api.Item {
	return items[0]
}

func (items Items) Last() *api.Item {
	return items[len(items)-1]
}

func (items Items) Append(items_ ...*api.Item) Items {
	return append(items, Items(items_)...)
}

func (items Items) Copy() Items {
	r := make(Items, len(items))
	copy(r, items)
	return r
}

func (items Items) Reverse() {
	lo.Reverse(items)
}

func (items Items) IDAscending() bool {
	for i := 1; i < len(items); i++ {
		if items[i-1].ID > items[i].ID {
			return false
		}
	}

	return true
}

func (items Items) IDDescending() bool {
	for i := 1; i < len(items); i++ {
		if items[i-1].ID < items[i].ID {
			return false
		}
	}

	return true
}

func (items Items) GroupByUID() map[int]Items {
	result := lo.GroupBy(items, func(item *api.Item) int {
		return item.UID
	})

	pipelines := make(map[int]Items)
	for k, v := range result {
		pipelines[k] = v
	}

	return pipelines
}

func (items Items) GroupBySharedWishType() map[api.SharedWishType]Items {
	result := lo.GroupBy(items, func(item *api.Item) api.SharedWishType {
		return item.WishType.Shared()
	})

	pipelines := make(map[api.SharedWishType]Items)
	for k, v := range result {
		pipelines[k] = v
	}

	return pipelines
}

func (items Items) SortByIDAscending() {
	switch {
	case items.IDAscending():
	case items.IDDescending():
		items.Reverse()
	default:
		slices.SortFunc(items, func(a *api.Item, b *api.Item) bool {
			return a.ID < b.ID
		})
	}
}

func (items Items) SortByIDDescending() {
	switch {
	case items.IDAscending():
		items.Reverse()
	case items.IDDescending():
	default:
		slices.SortFunc(items, func(a *api.Item, b *api.Item) bool {
			return a.ID > b.ID
		})
	}
}

func (items Items) Unique() Items {
	return lo.UniqBy(items, func(item *api.Item) int64 {
		return item.ID
	})
}

func (items Items) FilterByUID(uid int) Items {
	return lo.Filter(items, func(item *api.Item, _ int) bool {
		return item.UID == uid
	})
}

func (items Items) FilterByWishType(wishTypes ...api.WishType) Items {
	switch len(wishTypes) {
	case 0:
		return nil
	case 1:
		return lo.Filter(items, func(item *api.Item, _ int) bool {
			return item.WishType == wishTypes[0]
		})
	default:
		return lo.Filter(items, func(item *api.Item, _ int) bool {
			return lo.Contains(wishTypes, item.WishType)
		})
	}
}

func (items Items) FilterBySharedWishType(wishTypes ...api.SharedWishType) Items {
	switch len(wishTypes) {
	case 0:
		return nil
	case 1:
		return lo.Filter(items, func(item *api.Item, _ int) bool {
			return item.WishType.Shared() == wishTypes[0]
		})
	default:
		return lo.Filter(items, func(item *api.Item, _ int) bool {
			return lo.Contains(wishTypes, item.WishType.Shared())
		})
	}
}

func (items Items) FilterByRarity(rarities ...api.Rarity) Items {
	switch len(rarities) {
	case 0:
		return nil
	case 1:
		return lo.Filter(items, func(item *api.Item, _ int) bool {
			return item.Rarity == rarities[0]
		})
	default:
		return lo.Filter(items, func(item *api.Item, _ int) bool {
			return lo.Contains(rarities, item.Rarity)
		})
	}
}

func (items Items) UIDs() []int {
	return lo.Uniq(lo.Map(items, func(item *api.Item, _ int) int {
		return item.UID
	}))
}

func (items Items) SharedWishTypes() []api.SharedWishType {
	return lo.Uniq(lo.Map(items, func(item *api.Item, _ int) api.SharedWishType {
		return item.WishType.Shared()
	}))
}

func (items Items) Progress5Star() int {
	if len(items.UIDs()) != 1 || len(items.SharedWishTypes()) != 1 {
		return -1
	}

	items.SortByIDDescending()
	return slices.IndexFunc(items, func(item *api.Item) bool {
		return item.Rarity == api.Star5
	})
}

func (items Items) Progress4Star() int {
	if len(items.UIDs()) != 1 || len(items.SharedWishTypes()) != 1 {
		return -1
	}

	items.SortByIDDescending()
	return slices.IndexFunc(items, func(item *api.Item) bool {
		return item.Rarity == api.Star4
	})
}

func (items Items) Pulls5Stars() map[int64]int {
	if len(items.UIDs()) != 1 || len(items.SharedWishTypes()) != 1 {
		return nil
	}

	items.SortByIDAscending()
	pulls := make(map[int64]int)
	prev := 0

	for i, item := range items {
		if item.Rarity == api.Star5 {
			pulls[item.ID] = i - prev
			prev = i
		}
	}

	return pulls
}

func (items Items) Pulls4Stars() map[int64]int {
	if len(items.UIDs()) != 1 || len(items.SharedWishTypes()) != 1 {
		return nil
	}

	items.SortByIDAscending()
	pulls := make(map[int64]int)
	prev := 0

	for i, item := range items {
		if item.Rarity == api.Star4 {
			pulls[item.ID] = i - prev
			prev = i
		}
	}

	return pulls
}
