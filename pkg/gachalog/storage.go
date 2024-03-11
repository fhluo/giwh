package gachalog

import (
	"encoding/json"
	"errors"
	"github.com/fhluo/giwh/pkg/gacha"
	"github.com/fhluo/giwh/pkg/requests"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// Storage 是抽卡记录存储
type Storage struct {
	GachaLogs []gacha.Log // 抽卡记录
	ids       map[string]struct{}
}

// NewStorage 创建抽卡记录存储
func NewStorage() *Storage {
	return &Storage{}
}

// Add 添加抽卡记录
func (s *Storage) Add(logs ...gacha.Log) {
	for _, log := range logs {
		s.ids[log.ID] = struct{}{}
	}
	s.GachaLogs = append(s.GachaLogs, logs...)
}

// ReadFromFile 从文件读取抽卡记录
func (s *Storage) ReadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	
	var logs []gacha.Log
	if err = json.Unmarshal(data, &logs); err != nil {
		return err
	}

	s.Add(logs...)

	return nil
}

// ignoreErrNotExist 忽略文件不存在的错误
func ignoreErrNotExist(err error) error {
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return err
}

// ReadFromFileIfExits 从文件读取抽卡记录，如果文件不存在则忽略
func (s *Storage) ReadFromFileIfExits(filename string) error {
	return ignoreErrNotExist(s.ReadFromFile(filename))
}

// WriteToFile 将抽卡记录写入文件
func (s *Storage) WriteToFile(filename string) error {
	data, err := json.Marshal(s.GachaLogs)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0666)
}

// BackupAndWriteToFile 备份并将抽卡记录写入文件
func (s *Storage) BackupAndWriteToFile(filename string) error {
	dir, base := filepath.Split(filename)
	ext := filepath.Ext(base)

	if err := os.Rename(filename, filepath.Join(dir, strings.TrimSuffix(base, ext)+"_backup"+ext)); err != nil {
		slog.Warn("backup failed", "err", err.Error())
	}

	return s.WriteToFile(filename)
}

// Update 更新抽卡记录
func (s *Storage) Update(url requests.RequestURL) error {
	logs, err := requests.FetchGachaLogs(url).Until(func(log gacha.Log) bool {
		_, ok := s.ids[log.ID]
		return ok
	})

	if err != nil {
		return err
	}

	s.Add(logs...)

	return nil
}

// UpdateAll 更新所有抽卡记录
func (s *Storage) UpdateAll(url requests.RequestURL) error {
	logs, err := requests.FetchAllUntil(url, func(log gacha.Log) bool {
		_, ok := s.ids[log.ID]
		return ok
	})

	if err != nil {
		return err
	}

	s.Add(logs...)

	return nil
}
