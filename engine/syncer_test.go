package engine

import (
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"uptoc/uploader"
)

func TestSync(t *testing.T) {
	// init test data
	cfg := Config{SaveRoot: "/test/", Excludes: []string{"./abc"}}
	e := New(cfg, &mockUploader{})
	localObjects, err := e.loadLocalObjects("./testdata")
	assert.NoError(t, err)

	// test
	syncer := NewSyncer(&mockUploader{}, cfg, "./testdata")
	assert.NoError(t, syncer.Sync(localObjects, mockRemoteObjects(cfg.SaveRoot)))
}

func mockRemoteObjects(remoteRoot string) []uploader.Object {
	remoteFiles := map[string]string{
		"./abc/abc1.txt":        "abcabc",
		"./aaa/abc2.txt":        "112233",
		"./bbb/abc3.txt":        "445566",
		"./testdata/abc/a1.txt": "445566",
		"./testdata/abc/a2.txt": "445566",
	}

	remoteObjects := make([]uploader.Object, 0)
	for name, content := range remoteFiles {
		hash := md5.Sum([]byte(content))
		remoteObjects = append(remoteObjects, uploader.Object{
			Key:      filepath.Join(remoteRoot, name),
			ETag:     hex.EncodeToString(hash[:]),
			FilePath: name,
			Type:     "text/plain",
		})
	}

	return remoteObjects
}
