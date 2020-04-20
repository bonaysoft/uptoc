package uploader

import (
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSUploader
type OSSUploader struct {
	bucket *oss.Bucket
}

// OSSUploader Constructor
func AliOSSUploader(endpoint, accessKeyID, accessKeySecret, bucketName string) (Driver, error) {
	ossCli, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := ossCli.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return &OSSUploader{
		bucket: bucket,
	}, nil
}

// List the remote objects
func (u *OSSUploader) ListObjects() ([]Object, error) {
	objectsResult, err := u.bucket.ListObjects()
	if err != nil {
		return nil, err
	}

	objects := make([]Object, 0, len(objectsResult.Objects))
	for _, obj := range objectsResult.Objects {
		objects = append(objects, Object{Key: obj.Key, ETag: strings.ToLower(strings.Trim(obj.ETag, `"`))})
	}

	return objects, nil
}

// Upload the local file to the object
func (u *OSSUploader) Upload(objectKey, filePath string) error {
	return u.bucket.PutObjectFromFile(objectKey, filePath)
}

// Delete the object
func (u *OSSUploader) Delete(object string) error {
	return u.bucket.DeleteObject(object)
}
