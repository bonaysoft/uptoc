package uploader

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type COSUploader struct {
	*cos.Client
}

func NewCOSUploader(endpoint, accessKeyID, accessKeySecret, bucketName string) (Driver, error) {
	u, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucketName, endpoint))
	if err != nil {
		return nil, err
	}

	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  accessKeyID,
			SecretKey: accessKeySecret,
		},
	})

	return &COSUploader{
		Client: client,
	}, nil
}

func (u *COSUploader) ListObjects() ([]Object, error) {
	marker := ""
	objects := make([]Object, 0)
	for {
		objectsResult, _, err := u.Bucket.Get(context.Background(), &cos.BucketGetOptions{Marker: marker})
		if err != nil {
			return nil, err
		}

		for _, obj := range objectsResult.Contents {
			objects = append(objects, Object{Key: obj.Key, ETag: strings.ToLower(strings.Trim(obj.ETag, `"`))})
		}

		if objectsResult.IsTruncated {
			marker = objectsResult.NextMarker
		} else {
			break
		}
	}

	return objects, nil
}

func (u *COSUploader) Upload(object, rawPath string) (err error) {
	_, err = u.Object.PutFromFile(context.Background(), object, rawPath, nil)
	return
}

func (u *COSUploader) Delete(object string) (err error) {
	_, err = u.Object.Delete(context.Background(), object)
	return
}
