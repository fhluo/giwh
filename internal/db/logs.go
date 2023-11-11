package db

import (
	"database/sql"
	_ "embed"
	"github.com/fhluo/giwh/internal/config"
	"github.com/goccy/go-json"
	"log/slog"
	"os"
	"sync"
)

var (
	once sync.Once
	logs LogsDB
	//go:embed logs.sql
	logsSQL string
)

func Logs() LogsDB {
	once.Do(func() {
		db, err := sql.Open("sqlite", config.LogsPath.Get())
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		logs.db = db
		logs.InitDB()
	})
	return logs
}

type LogsDB struct {
	db *sql.DB
}

func (logs LogsDB) InitDB() {
	_, err := logs.db.Exec(logsSQL)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (logs LogsDB) ImportFromJSON(filename string) error {
	tx, err := logs.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`insert into logs(id, uid, gacha_type, name, item_id, item_type, rank_type, count, time, lang)
values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			slog.Warn(err.Error())
		}
	}()

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var logs_ []Log
	err = json.Unmarshal(data, &logs_)
	if err != nil {
		return err
	}

	for _, log := range logs_ {
		_, err = stmt.Exec(log.ID, log.UID, log.GachaType, log.Name, log.ItemID, log.ItemType, log.RankType, log.Count, log.Time, log.Lang)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
