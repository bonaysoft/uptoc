package engine

import (
	"crypto/md5"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"uptoc/uploader"
)

func TestSync(t *testing.T) {
	// init test data
	files := map[string]string{
		"abc1.txt": "abcabc",
		"abc2.txt": "112233",
		"abc3.txt": "445566",
	}

	localObjects := make([]uploader.Object, 0)
	for name, content := range files {
		hash := md5.Sum([]byte(content))
		localObjects = append(localObjects, uploader.Object{
			Key:      name,
			ETag:     hex.EncodeToString(hash[:]),
			FilePath: name,
			Type:     "text/plain",
		})
	}

	// test
	syncer := NewSyncer(&mockUploader{})
	assert.NoError(t, syncer.Sync(localObjects, ""))
}
