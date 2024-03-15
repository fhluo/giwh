package stores

import (
	"cmp"
	"encoding/json"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/fhluo/giwh/gacha-logs/gacha"
	"github.com/fhluo/giwh/hoyo-auth/auths"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Store 是抽卡记录存储
type Store struct {
	logs  []gacha.Log // 抽卡记录
	index map[string]struct{}
}

// New 创建抽卡记录存储
func New(logs []gacha.Log) *Store {
	store := new(Store)
	store.Add(logs...)
	return store
}

func (store *Store) Unique() []gacha.Log {
	// descending order
	slices.SortFunc(store.logs, func(a, b gacha.Log) int {
		return -cmp.Compare(a.ID, b.ID)
	})

	store.logs = slices.CompactFunc(store.logs, func(a gacha.Log, b gacha.Log) bool {
		return a.ID == b.ID
	})

	return store.logs
}

// Add 添加抽卡记录
func (store *Store) Add(logs ...gacha.Log) {
	if store.index == nil {
		store.index = make(map[string]struct{})
	}

	for _, log := range logs {
		store.index[log.ID] = struct{}{}
	}
	store.logs = append(store.logs, logs...)
}

// Load 从文件读取抽卡记录
func (store *Store) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var logs []gacha.Log
	if err = sonic.Unmarshal(data, &logs); err != nil {
		return err
	}

	store.Add(logs...)
	return nil
}

// ignoreErrNotExist 忽略文件不存在的错误
func ignoreErrNotExist(err error) error {
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return err
}

// LoadIfExists 从文件读取抽卡记录，如果文件不存在则忽略
func (store *Store) LoadIfExists(filename string) error {
	return ignoreErrNotExist(store.Load(filename))
}

// Save 将抽卡记录写入文件
func (store *Store) Save(filename string) error {
	data, err := json.MarshalIndent(store.Unique(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0666)
}

// BackupAndSave 备份并将抽卡记录写入文件
func (store *Store) BackupAndSave(filename string) error {
	dir, base := filepath.Split(filename)
	ext := filepath.Ext(base)

	backupFilename := filepath.Join(dir, strings.TrimSuffix(base, ext)+"_backup"+ext)
	if err := os.Rename(filename, backupFilename); err != nil {
		slog.Warn("failed to backup", "filename", filename, "backupFilename", backupFilename, "err", err.Error())
	}

	return store.Save(filename)
}

// Update 更新抽卡记录
func (store *Store) Update(auth *auths.Auth, f func(log gacha.Log)) error {
	for _, typ := range gacha.SharedTypes {
		logs, err := gacha.NewClient(auth).Fetch(typ, func(logs []gacha.Log) bool {
			return slices.ContainsFunc(logs, func(log gacha.Log) bool {
				if _, ok := store.index[log.ID]; ok {
					return true
				}

				if f != nil {
					f(log)
				}
				return false
			})
		})

		if err != nil {
			return err
		}
		store.Add(logs...)
	}
	return nil
}
