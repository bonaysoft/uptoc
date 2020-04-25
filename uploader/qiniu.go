package uploader

import (
	"context"
	"fmt"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

var zones = map[string]storage.Zone{
	"huadong":  storage.ZoneHuadong,
	"huabei":   storage.ZoneHuabei,
	"huanan":   storage.ZoneHuanan,
	"beimei":   storage.ZoneBeimei,
	"xinjiapo": storage.ZoneXinjiapo,
}

// Qiniu implements the Driver base on qiuniu.
type Qiniu struct {
	mac        *qbox.Mac
	cfg        *storage.Config
	bucketName string
}

// NewQiniuUploader returns a new Qiniu uploader
func NewQiniuUploader(endpoint, accessKey, accessSecret, bucketName string) (Driver, error) {
	zone, ok := zones[endpoint]
	if !ok {
		return nil, fmt.Errorf("endpoint %s not support", endpoint)
	}

	return &Qiniu{
		mac: qbox.NewMac(accessKey, accessSecret),
		cfg: &storage.Config{
			Zone: &zone,
		},
		bucketName: bucketName,
	}, nil
}

// ListObjects returns some remote objects
func (u *Qiniu) ListObjects() ([]Object, error) {
	limit := 1000
	prefix := ""
	delimiter := ""

	//初始列举marker为空
	marker := ""
	objects := make([]Object, 0)
	bucket := storage.NewBucketManager(u.mac, u.cfg)
	for {
		entries, _, nextMarker, hashNext, err := bucket.ListFiles(u.bucketName, prefix, delimiter, marker, limit)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			objects = append(objects, Object{
				Key:  entry.Key,
				ETag: entry.Hash,
			})
		}
		if hashNext {
			marker = nextMarker
		} else {
			//list end
			break
		}
	}

	return objects, nil
}

// Upload uploads the local file to the object
func (u *Qiniu) Upload(object, rawPath string) error {
	putPolicy := storage.PutPolicy{
		Scope: u.bucketName,
	}
	ctx := context.Background()
	upToken := putPolicy.UploadToken(u.mac)
	return storage.NewFormUploader(u.cfg).PutFile(ctx, &storage.PutRet{}, upToken, object, rawPath, nil)
}

// Delete deletes the object
func (u *Qiniu) Delete(object string) error {
	return storage.NewBucketManager(u.mac, u.cfg).Delete(u.bucketName, object)
}
