package pipeline

import (
	"cmp"
	"github.com/fhluo/giwh/giwh/api/gacha"
	"github.com/samber/lo"
	"log/slog"
	"slices"
)

type Pipeline struct {
	logs  []gacha.Log
	index map[string]struct{}
}

func New(logs []gacha.Log) *Pipeline {
	index := make(map[string]struct{})
	for _, log := range logs {
		index[log.ID] = struct{}{}
	}
	return &Pipeline{
		logs:  logs,
		index: index,
	}
}

func (p *Pipeline) Len() int {
	return len(p.logs)
}

func (p *Pipeline) First() gacha.Log {
	return p.logs[0]
}

func (p *Pipeline) Last() gacha.Log {
	return p.logs[len(p.logs)-1]
}

func (p *Pipeline) Logs() []gacha.Log {
	return p.logs
}

func (p *Pipeline) Contains(log gacha.Log) bool {
	_, ok := p.index[log.ID]
	return ok
}

func (p *Pipeline) ContainsAny(logs ...gacha.Log) bool {
	return slices.ContainsFunc(logs, func(log gacha.Log) bool {
		_, ok := p.index[log.ID]
		return ok
	})
}

func (p *Pipeline) Append(logs ...gacha.Log) {
	for _, log := range logs {
		if _, ok := p.index[log.ID]; ok {
			continue
		}

		p.index[log.ID] = struct{}{}
		p.logs = append(p.logs, log)
	}
}

func (p *Pipeline) Reverse() *Pipeline {
	slices.Reverse(p.logs)
	return p
}

func (p *Pipeline) IDAscending() bool {
	for i := 1; i < len(p.logs); i++ {
		if p.logs[i-1].ID > p.logs[i].ID {
			return false
		}
	}

	return true
}

func (p *Pipeline) IDDescending() bool {
	for i := 1; i < len(p.logs); i++ {
		if p.logs[i-1].ID < p.logs[i].ID {
			return false
		}
	}

	return true
}

func (p *Pipeline) SortByIDAscending() *Pipeline {
	switch {
	case p.IDAscending():
		slog.Debug("pipeline", "idAscending", p.IDAscending())
	case p.IDDescending():
		slog.Debug("pipeline", "idDescending", p.IDDescending())
		p.Reverse()
	default:
		slices.SortFunc(p.logs, func(a gacha.Log, b gacha.Log) int {
			return cmp.Compare(a.ID, b.ID)
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
		slices.SortFunc(p.logs, func(a gacha.Log, b gacha.Log) int {
			return -cmp.Compare(a.ID, b.ID)
		})
	}
	return p
}

func (p *Pipeline) GroupByUID() map[string][]gacha.Log {
	return lo.GroupBy(p.logs, func(log gacha.Log) string {
		return log.UID
	})
}

func (p *Pipeline) GroupBySharedWish() map[gacha.Type][]gacha.Log {
	return lo.GroupBy(p.logs, func(log gacha.Log) gacha.Type {
		return log.SharedWishType()
	})
}

func (p *Pipeline) Unique() *Pipeline {
	return New(lo.UniqBy(p.logs, func(log gacha.Log) string {
		return log.ID
	}))
}

func (p *Pipeline) FilterByUID(uid string) *Pipeline {
	return New(lo.Filter(p.logs, func(log gacha.Log, _ int) bool {
		return log.UID == uid
	}))
}

func (p *Pipeline) FilterByWish(types ...gacha.Type) *Pipeline {
	switch len(types) {
	case 0:
		return nil
	case 1:
		return New(lo.Filter(p.logs, func(log gacha.Log, _ int) bool {
			return log.GachaType == types[0]
		}))
	default:
		return New(lo.Filter(p.logs, func(log gacha.Log, _ int) bool {
			return slices.Contains(types, log.GachaType)
		}))
	}
}

func (p *Pipeline) FilterBySharedWish(types ...gacha.Type) *Pipeline {
	if slices.Contains(types, gacha.CharacterEventWish) {
		types = append(types, gacha.CharacterEventWish2)
	}
	return p.FilterByWish(types...)
}

func (p *Pipeline) FilterByRarity(rarities ...string) *Pipeline {
	switch len(rarities) {
	case 0:
		return nil
	case 1:
		return New(lo.Filter(p.logs, func(log gacha.Log, _ int) bool {
			return log.RankType == rarities[0]
		}))
	default:
		return New(lo.Filter(p.logs, func(log gacha.Log, _ int) bool {
			return slices.Contains(rarities, log.RankType)
		}))
	}
}

func (p *Pipeline) UIDs() []string {
	return lo.Uniq(lo.Map(p.logs, func(log gacha.Log, _ int) string {
		return log.UID
	}))
}

func (p *Pipeline) SharedWishes() []gacha.Type {
	return lo.Uniq(lo.Map(p.logs, func(log gacha.Log, _ int) gacha.Type {
		return log.SharedWishType()
	}))
}

func (p *Pipeline) Progress5Star() int {
	if len(p.UIDs()) != 1 || len(p.SharedWishes()) != 1 {
		return -1
	}

	p.SortByIDDescending()
	return slices.IndexFunc(p.logs, func(log gacha.Log) bool {
		return log.RankType == "5"
	})
}

func (p *Pipeline) Progress4Star() int {
	if len(p.UIDs()) != 1 || len(p.SharedWishes()) != 1 {
		return -1
	}

	p.SortByIDDescending()
	return slices.IndexFunc(p.logs, func(log gacha.Log) bool {
		return log.RankType == "4"
	})
}

func (p *Pipeline) Pulls5Stars() map[string]int {
	if len(p.UIDs()) != 1 || len(p.SharedWishes()) != 1 {
		return nil
	}

	p.SortByIDAscending()
	pulls := make(map[string]int)
	prev := 0

	for i, log := range p.logs {
		if log.RankType == "5" {
			pulls[log.ID] = i - prev
			prev = i
		}
	}

	return pulls
}

func (p *Pipeline) Pulls4Stars() map[string]int {
	if len(p.UIDs()) != 1 || len(p.SharedWishes()) != 1 {
		return nil
	}

	p.SortByIDAscending()
	pulls := make(map[string]int)
	prev := 0

	for i, log := range p.logs {
		if log.RankType == "4" {
			pulls[log.ID] = i - prev
			prev = i
		}
	}

	return pulls
}
