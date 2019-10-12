package oss

import (
	"fmt"
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

func (u *Uploader) Upload(objectKey, filePath string) error {
	fmt.Println(u.bucket.BucketName, objectKey, filePath)
	return u.bucket.PutObjectFromFile(objectKey, filePath)
}
