package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Items []string

func (items *Items) Next() (string, error) {
	for len(*items) != 0 {
		item := (*items)[0]
		*items = (*items)[1:]
		return item, nil
	}
	return "", Stop
}

func TestCollect(t *testing.T) {
	var items []string
	temp := Items(items)
	result, err := Collect[string](&temp)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, items, result)

	items = []string{"1", "2", "3", "4", "5", "6", "7"}
	temp = Items(items)
	result, err = Collect[string](&temp)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, items, result)
}
