package db

import (
	"github.com/fhluo/giwh/giwh/stores"
	"testing"
)

func TestLogs(t *testing.T) {
	logs, err := NewLogsDB(DefaultPath())
	if err != nil {
		t.Fatal(err)
	}

	if err = logs.ImportFromJSON(stores.DefaultPath()); err != nil {
		t.Fatal(err)
	}

	t.Log(logs.UIDList())
}
