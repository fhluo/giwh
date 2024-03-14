package hyauth

import (
	"fmt"
	"github.com/lmittmann/tint"
	"github.com/pelletier/go-toml/v2"
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

func TestGenshin(t *testing.T) {
	skipCI(t)

	data, err := toml.Marshal(map[string]any{
		"genshin_cn": map[string]any{
			"output_log_path":   GenshinCN().outputLogPath(),
			"program_data_path": GenshinCN().programDataPath(),
			"auths":             GenshinCN().Auths(),
		},
		"genshin_global": map[string]any{
			"output_log_path":   GenshinGlobal().outputLogPath(),
			"program_data_path": GenshinGlobal().programDataPath(),
			"auths":             GenshinGlobal().Auths(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(data))
}

func TestGenshin_urlsInCacheData(t *testing.T) {
	skipCI(t)

	if GenshinCN().Executed() {
		for _, url := range GenshinCN().urlsInCacheData() {
			fmt.Println(url)
		}
	}

	if GenshinGlobal().Executed() {
		for _, url := range GenshinGlobal().urlsInCacheData() {
			fmt.Println(url)
		}
	}
}
