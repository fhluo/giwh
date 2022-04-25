package util

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createEmptyTemp(t *testing.T, dir string, pattern string) string {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			t.Log(err)
		}
	}()
	return f.Name()
}

func TestSortFiles(t *testing.T) {
	_, err := SortFiles()
	assert.NotNil(t, err)

	_, err = SortFiles("")
	assert.NotNil(t, err)

	dir, err := os.MkdirTemp("", "giwh-test*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			t.Log(err)
		}
	}()

	count := 3
	names := make([]string, count)
	for i := count - 1; i >= 0; i-- {
		names[i] = createEmptyTemp(t, dir, "file*")
		time.Sleep(time.Millisecond)
	}

	r, err := SortFiles(names[0])
	assert.Equal(t, []string{names[0]}, r)

	r, err = SortFiles(names[1], names[0])
	assert.Equal(t, names[:2], r)

	r, err = SortFiles(names[1], "", names[0])
	assert.Equal(t, names[:2], r)

	r, err = SortFiles(names[1], names[2], names[0])
	assert.Equal(t, names, r)
}

func TestExpandPaths(t *testing.T) {
	dir, err := os.MkdirTemp("", "giwh-test*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			t.Log(err)
		}
	}()

	count := 6
	names := make([]string, count)
	for i := 0; i < count/2; i++ {
		names[i] = createEmptyTemp(t, dir, "a*")
		names[i+count/2] = createEmptyTemp(t, dir, "b*")
	}

	result, err := ExpandPaths(filepath.Join(dir, "a*"), filepath.Join(dir, "b*"))
	if err != nil {
		t.Fatal(err)
	}

	slices.Sort(names)
	slices.Sort(result)

	assert.Equal(t, names, result)
}
