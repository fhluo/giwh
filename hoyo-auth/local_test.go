package hoyo_auth

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"log/slog"
	"os"
	"testing"
)

func init() {
	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}),
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
			"infos":             GenshinCN().AuthInfos(),
		},
		"genshin_global": map[string]any{
			"output_log_path":   GenshinGlobal().outputLogPath(),
			"program_data_path": GenshinGlobal().programDataPath(),
			"infos":             GenshinGlobal().AuthInfos(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(data))
}
