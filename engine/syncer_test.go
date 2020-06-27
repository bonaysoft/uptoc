package engine

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

func TestSync2(t *testing.T) {
	objects := []uploader.Object{
		{
			Key:      "test",
			FilePath: "test",
			Type:     "text/plain",
		},
	}

	s := NewSyncer(&mockUploader{listErr: fmt.Errorf("list error")})
	assert.Error(t, s.Sync(objects, ""))

	s = NewSyncer(&mockUploader{uploadErr: fmt.Errorf("upload error")})
	assert.Error(t, s.Sync(objects, ""))

	s = NewSyncer(&mockUploader{deleteErr: fmt.Errorf("delete error")})
	assert.Error(t, s.Sync(objects, ""))
}
