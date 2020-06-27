package engine

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"uptoc/uploader"
)

type mockUploader struct {
	listErr   error
	uploadErr error
	deleteErr error
}

func (m *mockUploader) ListObjects(prefix string) ([]uploader.Object, error) {
	return []uploader.Object{
		{Key: "abc1.txt", ETag: "abc123"},
		{Key: "abc2.txt", ETag: "d0970714757783e6cf17b26fb8e2298f"},
		{Key: "abc4.txt", ETag: "aaa123"},
	}, m.listErr
}

func (m *mockUploader) Upload(object, rawPath string) error {
	if strings.HasSuffix(object, "failed.txt") {
		return fmt.Errorf("test error")
	}

	return m.uploadErr
}

func (m *mockUploader) Delete(object string) error {
	return m.deleteErr
}

func TestEngine_TailRun(t *testing.T) {
	// init test data
	tmp := "/tmp/uptoc-engine-ut/"
	assert.NoError(t, os.RemoveAll(tmp))
	assert.NoError(t, os.Mkdir(tmp, os.FileMode(0755)))
	files := map[string]string{
		"abc1.txt":         "abcabc",
		"abc2.txt":         "112233",
		"abc3.txt":         "445566",
		"exclude/test.txt": "abc123",
	}
	for name, content := range files {
		if strings.HasPrefix(name, "exclude") {
			os.Mkdir(tmp+"exclude", os.FileMode(0744))
		}
		assert.NoError(t, ioutil.WriteFile(tmp+name, []byte(content), os.FileMode(0644)))
	}

	conf := Config{
		ForceSync: true,
		Excludes:  []string{"exclude"},
	}
	e := New(conf, &mockUploader{})
	e.TailRun(tmp)
}

func TestEngine_TailRun2(t *testing.T) {
	// init test data
	tmp := "/tmp/uptoc-engine-ut/"
	assert.NoError(t, os.RemoveAll(tmp))
	assert.NoError(t, os.Mkdir(tmp, os.FileMode(0755)))
	files := map[string]string{
		"abc1.txt":     "abcabc",
		"abc2.txt":     "445566",
		"dir/test.txt": "abc123",
		"failed.txt":   "aaaaaa",
	}
	for name, content := range files {
		if strings.HasPrefix(name, "dir") {
			os.Mkdir(tmp+"dir", os.FileMode(0744))
		}

		assert.NoError(t, ioutil.WriteFile(tmp+name, []byte(content), os.FileMode(0644)))
	}

	e := New(Config{}, &mockUploader{})
	e.TailRun(tmp+"abc1.txt", tmp+"abc2.txt", tmp+"dir", tmp+"failed.txt")
}
