package pipeline

import (
	"github.com/fhluo/giwh/pkg/api"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"strconv"
	"time"
)

type Element struct {
	*api.Item
	ID   int64
	Time time.Time
}

func NewElement(item *api.Item) (*Element, error) {
	id, err := strconv.ParseInt(item.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse("2006-01-02 15:04:05", item.Time)
	if err != nil {
		return nil, err
	}

	return &Element{
		Item: item,
		ID:   id,
		Time: t,
	}, nil
}

type Pipeline struct {
	elements []*Element
}

func New(items []*api.Item) (pipeline Pipeline, err error) {
	pipeline.elements, err = ItemsTo(items, NewElement)
	return
}

func (p Pipeline) Traverse(f func(e *Element) bool) {
	for _, e := range p.elements {
		if !f(e) {
			break
		}
	}
}

func (p Pipeline) Append(elements []*Element) Pipeline {
	return Pipeline{elements: append(p.elements, elements...)}
}

func (p Pipeline) Elements() []*Element {
	return p.elements
}

func (p Pipeline) Items() []*api.Item {
	return lo.Map(p.elements, func(e *Element, _ int) *api.Item {
		return e.Item
	})
}

func (p Pipeline) Copy() Pipeline {
	elements := make([]*Element, len(p.elements))
	copy(elements, p.elements)
	return Pipeline{elements: elements}
}

func (p Pipeline) SortedByIDAscending() Pipeline {
	r := p.Copy()
	slices.SortFunc(r.elements, func(a *Element, b *Element) bool {
		return a.ID < b.ID
	})
	return r
}

func (p Pipeline) SortedByIDDescending() Pipeline {
	r := p.Copy()
	slices.SortFunc(r.elements, func(a *Element, b *Element) bool {
		return a.ID > b.ID
	})
	return r
}

func (p Pipeline) Unique() Pipeline {
	return Pipeline{
		elements: lo.UniqBy(p.elements, func(e *Element) int64 {
			return e.ID
		}),
	}
}

func (p Pipeline) FilterByUID(uid string) Pipeline {
	return Pipeline{
		elements: lo.Filter(p.elements, func(e *Element, _ int) bool {
			return e.UID == uid
		}),
	}
}

func (p Pipeline) FilterByWishType(wishTypes ...string) Pipeline {
	switch len(wishTypes) {
	case 0:
		return Pipeline{}
	case 1:
		return Pipeline{
			elements: lo.Filter(p.elements, func(e *Element, _ int) bool {
				return e.WishType == wishTypes[0]
			}),
		}
	default:
		return Pipeline{
			elements: lo.Filter(p.elements, func(e *Element, _ int) bool {
				return lo.Contains(wishTypes, e.WishType)
			}),
		}
	}
}

func (p Pipeline) FilterByRarity(rarities ...string) Pipeline {
	switch len(rarities) {
	case 0:
		return Pipeline{}
	case 1:
		return Pipeline{
			elements: lo.Filter(p.elements, func(e *Element, _ int) bool {
				return e.Rarity == rarities[0]
			}),
		}
	default:
		return Pipeline{
			elements: lo.Filter(p.elements, func(e *Element, _ int) bool {
				return lo.Contains(rarities, e.Rarity)
			}),
		}
	}
}
