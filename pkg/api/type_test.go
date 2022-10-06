package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"testing"
)

func TestResult(t *testing.T) {
	data, err := os.ReadFile("testdata/response.json")
	if err != nil {
		t.Fatal(err)
	}
	data = regexp.MustCompile(`\r\n`).ReplaceAll(data, []byte{'\n'})

	var r JSONResponse[*ItemList]
	err = json.Unmarshal(data, &r)
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, data, b)
}
