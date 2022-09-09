package auth

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

func createText(t *testing.T, path string, content string) {
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	_, err = f.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
}

func createFiles(t *testing.T, dataPath, outputLogPath, data2Path string) {
	err := os.Mkdir(dataPath, 0666)
	if err != nil {
		t.Fatal(err)
	}

	createText(t, outputLogPath, fmt.Sprintf("Temp\nWarmup file %s\nTemp", dataPath))

	err = os.MkdirAll(filepath.Dir(data2Path), 0666)
	if err != nil {
		t.Fatal(err)
	}

	createText(t, data2Path, "https://example.com/api/getGachaLog?ver=1\000???https://\000https://example.com/index.html?x=0&y=1 https://example.com/api/getGachaLog?ver=2")
}

func TestRegion(t *testing.T) {
	dir, err := os.MkdirTemp("", "giwh-test*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			t.Log(err)
		}
	}()

	dataPath := filepath.Join(dir, "GenshinImpact_Data")
	outputLogPath := filepath.Join(dir, "output_log.txt")
	data2Path := filepath.Join(dataPath, `webCaches\Cache\Cache_Data\data_2`)
	createFiles(t, dataPath, outputLogPath, data2Path)

	test := Region{
		Name:          "Test",
		OutputLogPath: outputLogPath,
		APIBaseURL:    "https://example.com/api/getGachaLog",
	}

	// Test GetCacheDataPath
	path, err := test.GetCacheDataPath()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, data2Path, path)

	// Test GetURLsFromCacheData
	urls, err := test.GetURLsFromCacheData()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{"https://example.com/api/getGachaLog?ver=1", "https://example.com/index.html?x=0&y=1", "https://example.com/api/getGachaLog?ver=2"}, urls)

	// Test GetAPIURL
	url, err := test.GetAPIURL()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "https://example.com/api/getGachaLog?ver=2", url)

	createText(t, data2Path, "https://example.com/index.html?authkey=1\000???https://\000https://example.com/index.html?x=0&y=1 https://example.com/index.html?authkey=2")
	url, err = test.GetAPIURL()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "https://example.com/api/getGachaLog?authkey=2", url)
}
