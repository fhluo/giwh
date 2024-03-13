package db

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"github.com/fhluo/giwh/gacha-logs/gacha"
	"log/slog"
	_ "modernc.org/sqlite"
	"os"
)

//go:embed logs.sql
var logsSQL string

func NewLogsDB(path string) (LogsDB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return LogsDB{}, err
	}
	logsDB := LogsDB{
		db: db,
	}
	logsDB.InitDB()
	return logsDB, nil
}

type LogsDB struct {
	db *sql.DB
}

func (l LogsDB) InitDB() {
	_, err := l.db.Exec(logsSQL)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (l LogsDB) ImportFromJSON(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var logs []gacha.Log
	err = json.Unmarshal(data, &logs)
	if err != nil {
		return err
	}

	return l.Insert(logs...)
}

func (l LogsDB) ImportFromWishHistory() error {
	return l.ImportFromJSON(config.WishHistoryPath.Get())
}

func (l LogsDB) Insert(logs ...gacha.Log) error {
	tx, err := l.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`insert or ignore into log(id, uid, gacha_type, name, item_id, item_type, rank_type, count, time, lang)
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

	for _, log := range logs {
		_, err = stmt.Exec(log.ID, log.UID, log.GachaType, log.Name, log.ItemID, log.ItemType, log.RankType, log.Count, log.Time, log.Lang)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (l LogsDB) UIDList() ([]string, error) {
	rows, err := l.db.Query(`select uid from uid`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			slog.Warn(err.Error())
		}
	}()

	var uidList []string

	for rows.Next() {
		var uid string
		err = rows.Scan(&uid)
		if err != nil {
			return nil, err
		}
		uidList = append(uidList, uid)
	}

	return uidList, nil
}

// Log 表示抽卡记录
type Log struct {
	ID        string `json:"id"`         // 记录 ID
	UID       string `json:"uid"`        // 用户 ID
	GachaType string `json:"gacha_type"` // 卡池类型
	Name      string `json:"name"`       // 物品名称
	ItemID    string `json:"item_id"`    // 物品 ID
	ItemType  string `json:"item_type"`  // 物品类型
	RankType  string `json:"rank_type"`  // 稀有度
	Count     string `json:"count"`      // 物品数量
	Time      string `json:"time"`       // 时间
	Lang      string `json:"lang"`       // 语言
}

type Pity struct {
	UID       string `json:"uid"`
	GachaType string `json:"gacha_type"`
	Pity      int
}

func (l LogsDB) Pity() {

}

func (l LogsDB) UpdateFromURL() {

}
