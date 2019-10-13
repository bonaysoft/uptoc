package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Uploader struct {
	bucket *oss.Bucket
}

func NewUploader(endpoint, accessKeyID, accessKeySecret, bucketName string) (*Uploader, error) {
	ossCli, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := ossCli.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return &Uploader{
		bucket: bucket,
	}, nil
}

func (u *Uploader) ListObjects() ([]string, error) {
	objectsResult, err := u.bucket.ListObjects()
	if err != nil {
		return nil, err
	}

	objects := make([]string, 0, len(objectsResult.Objects))
	for _, obj := range objectsResult.Objects {
		objects = append(objects, obj.Key)
	}

	return objects, nil
}

func (u *Uploader) Upload(objectKey, filePath string) error {
	return u.bucket.PutObjectFromFile(objectKey, filePath)
}

func (u *Uploader) Delete(object string) error {
	return u.bucket.DeleteObject(object)
}
