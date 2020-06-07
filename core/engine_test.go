package core

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"uptoc/uploader"
)

type mockUploader struct {
}

func (m mockUploader) ListObjects() ([]uploader.Object, error) {
	return []uploader.Object{
		{Key: "abc1.txt", ETag: "abc123"},
		{Key: "abc4.txt", ETag: "aaa123"},
	}, nil
}

func (m mockUploader) Upload(object, rawPath string) error {
	return nil
}

func (m mockUploader) Delete(object string) error {
	return nil
}

func TestEngine(t *testing.T) {
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

	// test
	e := NewEngine(&mockUploader{})
	assert.NoError(t, e.LoadAndCompareObjects(tmp, "/exclude"))
	assert.NoError(t, e.Sync())
}

func TestNOTExistDir(t *testing.T) {
	e := NewEngine(&mockUploader{})
	assert.Error(t, e.LoadAndCompareObjects("tmp233", "/exclude"))
}
