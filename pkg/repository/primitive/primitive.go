package primitive

import (
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/i18n"
	"github.com/fhluo/giwh/pkg/pipeline"
	"github.com/fhluo/giwh/pkg/repository"
	"github.com/samber/lo"
	"os"
	"path"
	"path/filepath"
)

type Repository struct {
	items    pipeline.Items
	index    map[int]map[api.SharedWishType]pipeline.Items
	modified bool
}

func Load(filename string) (*Repository, error) {
	items, err := repository.LoadIfExits(filename)
	if err != nil {
		return nil, err
	}

	r := new(Repository)
	r.items = items

	r.index = make(map[int]map[api.SharedWishType]pipeline.Items)
	for uid, p := range r.items.GroupByUID() {
		r.index[uid] = p.GroupBySharedWishType()
	}

	return r, nil
}

func (r *Repository) GetItems() pipeline.Items {
	return r.items
}

func (r *Repository) GetUIDs() []int {
	r.items.SortByIDDescending()
	return lo.Uniq(lo.Map(r.items, func(item *api.Item, _ int) int {
		return item.UID
	}))
}

func (r *Repository) Get5StarProgress(uid int, wishType api.SharedWishType) int {
	return r.index[uid][wishType].Progress5Star()
}

func (r *Repository) Get4StarProgress(uid int, wishType api.SharedWishType) int {
	return r.index[uid][wishType].Progress4Star()
}

func (r *Repository) Get5Stars(uid int, wishType api.SharedWishType) []repository.Item {
	items := r.index[uid][wishType]
	pulls := items.Pulls5Stars()

	return lo.Map(items.FilterByRarity(api.Star5), func(item *api.Item, _ int) repository.Item {
		name := i18n.Item{Name: item.Name, Lang: item.Lang}.GetNameWithLang("en")

		icon, ok := repository.Characters[name]
		if ok {
			icon = path.Join("/images/characters", icon)
		} else {
			icon = path.Join("/images/weapons", repository.Weapons[name])
		}

		return repository.Item{
			Item:  item,
			Pulls: pulls[item.ID],
			Icon:  icon,
		}
	})
}

func (r *Repository) Get4Stars(uid int, wishType api.SharedWishType) []repository.Item {
	items := r.index[uid][wishType]
	pulls := items.Pulls4Stars()

	return lo.Map(items.FilterByRarity(api.Star4), func(item *api.Item, _ int) repository.Item {
		name := i18n.Item{Name: item.Name, Lang: item.Lang}.GetNameWithLang("en")

		icon, ok := repository.Characters[name]
		if ok {
			icon = path.Join("/images/characters", icon)
		} else {
			icon = path.Join("/images/weapons", repository.Weapons[name])
		}

		return repository.Item{
			Item:  item,
			Pulls: pulls[item.ID],
			Icon:  icon,
		}
	})
}

func (r *Repository) AddItems(items []*api.Item) {
	r.modified = true
	r.items = r.items.Append(items...)
	
	for uid, p := range pipeline.Items(items).GroupByUID() {
		for k, v := range p.GroupBySharedWishType() {
			r.index[uid][k] = r.index[uid][k].Append(v...)
		}
	}
}

func (r *Repository) Save(filename string) error {
	if !r.modified {
		return nil
	}

	dir, base := filepath.Split(filename)
	ext := filepath.Ext(base)
	_ = os.Rename(filename, filepath.Join(dir, base[:len(base)-len(ext)]+"_backup"+ext))

	items := r.items.Unique()
	items.SortByIDDescending()
	return repository.Save(filename, items)
}
