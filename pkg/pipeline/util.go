package pipeline

import (
	"github.com/fhluo/giwh/pkg/api"
	"golang.org/x/exp/slices"
)

func ItemsEqual(items1 []*api.Item, items2 []*api.Item) bool {
	return slices.EqualFunc(items1, items2, func(item1, item2 *api.Item) bool {
		return item1.ID == item2.ID
	})
}

func ItemsTo[T any](items []*api.Item, f func(item *api.Item) (T, error)) ([]T, error) {
	elements := make([]T, 0, len(items))

	for _, item := range items {
		element, err := f(item)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}

	return elements, nil
}
