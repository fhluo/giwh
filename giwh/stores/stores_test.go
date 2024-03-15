package stores

import (
	"fmt"
	"github.com/fhluo/giwh/gacha-logs/gacha"
	"github.com/fhluo/giwh/hoyo-auth/clients"
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"testing"
	"time"
)

func init() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{Level: slog.LevelDebug, TimeFormat: time.TimeOnly}),
	))
}

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.SkipNow()
	}
}

func TestDefaultPath(t *testing.T) {
	skipCI(t)

	t.Log(DefaultPath())
}

func TestStore_Load(t *testing.T) {
	skipCI(t)

	store := New(nil)
	err := store.Load(DefaultPath())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Total pulls: %d\n.", len(store.Unique()))
}

func TestStore_Update(t *testing.T) {
	skipCI(t)

	store := New(nil)
	err := store.Load(DefaultPath())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Total pulls: %d\n.", len(store.Unique()))

	client := clients.Latest()
	if client == nil {
		t.Fatal("no client available")
	}

	auths := client.Auths()
	if len(auths) == 0 {
		t.Fatal("failed to find auth info")
	}

	count := 0
	err = store.Update(auths[len(auths)-1], func(log gacha.Log) {
		fmt.Printf("%s: pull %s\n", log.Time, log.Name)
		count++
	})
	if err != nil {
		t.Fatal(err)
	}

	if count == 0 {
		fmt.Println("No updates available, your wish history is up to date.")
		return
	}

	err = store.BackupAndSave(DefaultPath())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Updated %d pulls.\n", count)
}
