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
