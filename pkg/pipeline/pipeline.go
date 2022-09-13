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

	_4star *int
	_5star *int
}

func New(items []*api.Item) (pipeline Pipeline, err error) {
	pipeline.elements, err = ItemsTo(items, NewElement)
	return
}

func (p Pipeline) Count() int {
	return len(p.elements)
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
	for _, element := range p.elements {
		switch element.Rarity {
		case api.FourStar:
			*p._4star++
		case api.FiveStar:
			*p._5star++
		}
	}

	return *p._4star, *p._5star
}

func (p Pipeline) First() *Element {
	return p.elements[0]
}

func (p Pipeline) Last() *Element {
	return p.elements[len(p.elements)-1]
}

func (p Pipeline) Traverse(f func(e *Element)) {
	for _, e := range p.elements {
		f(e)
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

func (p Pipeline) Reverse() {
	lo.Reverse(p.elements)
}

func (p Pipeline) IDAscending() bool {
	for i := 1; i < len(p.elements); i++ {
		if p.elements[i].ID < p.elements[i-1].ID {
			return false
		}
	}

	return true
}

func (p Pipeline) IDDescending() bool {
	for i := 1; i < len(p.elements); i++ {
		if p.elements[i].ID > p.elements[i-1].ID {
			return false
		}
	}

	return true
}

func (p Pipeline) SortByIDAscending() {
	switch {
	case p.IDAscending():
	case p.IDDescending():
		p.Reverse()
	default:
		slices.SortFunc(p.elements, func(a *Element, b *Element) bool {
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
		slices.SortFunc(p.elements, func(a *Element, b *Element) bool {
			return a.ID > b.ID
		})
	}
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

func (p Pipeline) FilterBySharedWishType(wishTypes ...string) Pipeline {
	switch len(wishTypes) {
	case 0:
		return Pipeline{}
	case 1:
		return Pipeline{
			elements: lo.Filter(p.elements, func(e *Element, _ int) bool {
				switch wishTypes[0] {
				case api.CharacterEventWish, api.CharacterEventWish2:
					return e.WishType == api.CharacterEventWish || e.WishType == api.CharacterEventWish2
				default:
					return e.WishType == wishTypes[0]
				}
			}),
		}
	default:
		return Pipeline{
			elements: lo.Filter(p.elements, func(e *Element, _ int) bool {
				switch e.WishType {
				case api.CharacterEventWish, api.CharacterEventWish2:
					return lo.Contains(wishTypes, api.CharacterEventWish) || lo.Contains(wishTypes, api.CharacterEventWish2)
				default:
					return lo.Contains(wishTypes, e.WishType)
				}
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

func (p Pipeline) Progress() map[string]map[string]int {
	p.SortByIDDescending()
	result := make(map[string]map[string]int)
	done := make(map[string]map[string]bool)

	for _, wishType := range api.SharedWishTypes {
		result[wishType] = make(map[string]int)
		done[wishType] = make(map[string]bool)
	}

	var wishType string
	for _, element := range p.elements {
		switch element.WishType {
		case api.CharacterEventWish, api.CharacterEventWish2:
			wishType = api.CharacterEventWish
		default:
			wishType = element.WishType
		}

		switch element.Rarity {
		case api.FourStar:
			done[wishType][api.FourStar] = true
			if !done[wishType][api.FiveStar] {
				result[wishType][api.FiveStar]++
			}
		case api.FiveStar:
			done[wishType][api.FiveStar] = true
			if !done[wishType][api.FourStar] {
				result[wishType][api.FourStar]++
			}
		default:
			if !done[wishType][api.FourStar] {
				result[wishType][api.FourStar]++
			}
			if !done[wishType][api.FiveStar] {
				result[wishType][api.FiveStar]++
			}
		}
	}

	return result
}

func (p Pipeline) Pulls() map[string]map[int64]int {
	p.SortByIDAscending()
	progress := make(map[string]map[int64]int)
	progress4Star := make(map[string]int)
	progress5Star := make(map[string]int)

	for _, wishType := range api.SharedWishTypes {
		progress[wishType] = make(map[int64]int)
	}

	var wishType string
	for _, element := range p.elements {
		switch element.WishType {
		case api.CharacterEventWish, api.CharacterEventWish2:
			wishType = api.CharacterEventWish
		default:
			wishType = element.WishType
		}

		switch element.Rarity {
		case api.FourStar:
			progress[wishType][element.ID] = progress4Star[wishType] + 1
			progress4Star[wishType] = 0
			progress5Star[wishType]++
		case api.FiveStar:
			progress[wishType][element.ID] = progress5Star[wishType] + 1
			progress4Star[wishType]++
			progress5Star[wishType] = 0
		default:
			progress4Star[wishType]++
			progress5Star[wishType]++
		}
	}

	return progress
}
