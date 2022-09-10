package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestFindDataPath(t *testing.T) {
	assert.Equal(t, "", findDataPath(nil))

	assert.Equal(t,
		`C:\Program Files\Genshin Impact\Genshin Impact game\GenshinImpact_Data`,
		findDataPath([]byte(`C:\Program Files\Genshin Impact\Genshin Impact game\GenshinImpact_Data\`)),
	)

	assert.Equal(t,
		`C:\Program Files\Genshin Impact\Genshin Impact game\GenshinImpact_Data`,
		findDataPath([]byte(`C:/Program Files\Genshin Impact\Genshin Impact game\GenshinImpact_Data`)),
	)

	assert.Equal(t,
		`C:\Program Files\Genshin Impact\Genshin Impact game\GenshinImpact_Data`,
		findDataPath([]byte(`C:/Program Files/Genshin Impact/Genshin Impact game/GenshinImpact_Data`)),
	)

	assert.Equal(t,
		`C:\Program Files\Genshin Impact\Genshin Impact game\GenshinImpact_Data`,
		findDataPath([]byte(`Warmup file C:/Program Files/Genshin Impact/Genshin Impact game/GenshinImpact_Data/`)),
	)
}

func TestFindAllURLs(t *testing.T) {
	assert.Equal(t,
		[]string{"https://example.com/index.html", "https://example.com/index.html?x=0&y=1"},
		findAllURLs([]byte("https://example.com/index.html\000???https://\000https://example.com/index.html?x=0&y=1")),
	)
}

func CreateText(path string, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = f.WriteString(content)
	if err != nil {
		_ = f.Close()
		return err
	}

	return f.Close()
}

type TempDir struct {
	path string
}

func NewTempDir() (*TempDir, error) {
	dir, err := os.MkdirTemp("", "giwh-test*")
	if err != nil {
		return nil, err
	}
	return &TempDir{path: dir}, nil
}

func (d *TempDir) DataPath() string {
	return filepath.Join(d.path, "GenshinImpact_Data")
}

func (d *TempDir) CreateDataDir() error {
	return os.Mkdir(d.DataPath(), 0666)
}

func (d *TempDir) OutputLogPath() string {
	return filepath.Join(d.path, "output_log.txt")
}

func (d *TempDir) CreateOutputLog(content string) error {
	err := d.CreateDataDir()
	if err != nil {
		return err
	}
	return CreateText(d.OutputLogPath(), content)
}

func (d *TempDir) Data2Path() string {
	return filepath.Join(d.DataPath(), `webCaches\Cache\Cache_Data\data_2`)
}

func (d *TempDir) CreateData2(content string) error {
	err := os.MkdirAll(filepath.Dir(d.Data2Path()), 0666)
	if err != nil {
		return err
	}

	return CreateText(d.Data2Path(), content)
}

func (d *TempDir) Clean() {
	_ = os.RemoveAll(d.path)
}

func (d *TempDir) CreateFiles(outputLog string, data2 string) error {
	err := d.CreateOutputLog(outputLog)
	if err != nil {
		return err
	}
	return d.CreateData2(data2)
}

func newTestRegion(outputLogPath string, apiBaseURL string) Region {
	return Region{
		Name:          "Test",
		OutputLogPath: outputLogPath,
		APIBaseURL:    apiBaseURL,
	}
}

func TestRegion_GetCacheDataPath(t *testing.T) {
	dir, err := NewTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Clean()

	err = dir.CreateOutputLog(fmt.Sprintf("Temp\nWarmup file %s\nTemp", dir.DataPath()))
	if err != nil {
		t.Fatal(err)
	}

	test := newTestRegion(dir.OutputLogPath(), "https://example.com/event/gacha_info/api/getGachaLog")
	path, err := test.GetCacheDataPath()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, dir.Data2Path(), path)
}

func TestRegion_GetURLsFromCacheData(t *testing.T) {
	dir, err := NewTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Clean()

	err = dir.CreateFiles(
		fmt.Sprintf("Temp\nWarmup file %s\nTemp", dir.DataPath()),
		"https://example.com/event/gacha_info/api/getGachaLog?authkey_ver=1&authkey=x&lang=zh-cn\000https://\000https://example.com/index.html https://example.com/event/gacha_info/api/getGachaLog?authkey_ver=1&authkey=y&lang=zh-cn",
	)
	if err != nil {
		t.Fatal(err)
	}

	test := newTestRegion(dir.OutputLogPath(), "https://example.com/event/gacha_info/api/getGachaLog")
	urls, err := test.GetURLsFromCacheData()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{"https://example.com/event/gacha_info/api/getGachaLog?authkey_ver=1&authkey=x&lang=zh-cn", "https://example.com/index.html", "https://example.com/event/gacha_info/api/getGachaLog?authkey_ver=1&authkey=y&lang=zh-cn"}, urls)
}

func TestRegion_GetAPIBase(t *testing.T) {
	dir, err := NewTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Clean()

	err = dir.CreateFiles(
		fmt.Sprintf("Temp\nWarmup file %s\nTemp", dir.DataPath()),
		"https://example.com/event/gacha_info/api/getGachaLog?authkey_ver=1&authkey=x&lang=zh-cn\000https://\000https://example.com/index.html https://example.com/event/gacha_info/api/getGachaLog?authkey_ver=1&authkey=y&lang=zh-cn",
	)
	if err != nil {
		t.Fatal(err)
	}

	test := newTestRegion(dir.OutputLogPath(), "https://example.com/event/gacha_info/api/getGachaLog")
	_, baseQuery, err := test.GetAPIBase()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, baseQuery, BaseQuery{
		AuthKeyVer: "1",
		AuthKey:    "y",
		Lang:       "zh-cn",
	})

	err = dir.CreateData2("https://example.com/index.html?authkey_ver=1&authkey=x&lang=zh-cn\000https://\000https://example.com/index.html https://example.com/index.html?authkey_ver=1&authkey=y&lang=zh-cn")
	if err != nil {
		t.Fatal(err)
	}

	_, baseQuery, err = test.GetAPIBase()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, baseQuery, BaseQuery{
		AuthKeyVer: "1",
		AuthKey:    "y",
		Lang:       "zh-cn",
	})
}
