package db

import (
	"testing"
)

func TestLogs(t *testing.T) {
	err := Logs().ImportFromWishHistory()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(Logs().UIDList())
}
