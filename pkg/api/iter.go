package api

import "errors"

type Iterator[T any] interface {
	Next() (T, error)
}

var Stop = errors.New("stop iteration")

func Collect[T any](iterator Iterator[T]) ([]T, error) {
	var items []T
	for {
		item, err := iterator.Next()
		if err != nil {
			if errors.Is(err, Stop) {
				return items, nil
			} else {
				return items, err
			}
		}
		items = append(items, item)
	}
}
