package i18n

import (
	"fmt"
	"os"
	"testing"
)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.SkipNow()
	}
}

func TestReadLocaleFile(t *testing.T) {
	skipCI(t)
	data, err := ReadLocaleFile(English)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
}
