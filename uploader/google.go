package uploader

import (
	"context"
	"encoding/base64"
	"io"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GoogleUploader implements the Driver base on the Storage of Google.
type GoogleUploader struct {
	client *storage.BucketHandle
}

// NewGoogleUploader returns a new google storage uploader
func NewGoogleUploader(endpoint, accessKeyID, accessKeySecret, bucketName string) (Driver, error) {
	j, err := base64.StdEncoding.DecodeString(accessKeySecret)
	if err != nil {
		return nil, err
	}

	opts := []option.ClientOption{
		//option.WithEndpoint(endpoint),
		option.WithCredentialsJSON(j),
	}
	client, err := storage.NewClient(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	return &GoogleUploader{
		client: client.Bucket(bucketName),
	}, nil
}

// ListObjects returns some remote objects
func (u *GoogleUploader) ListObjects() ([]Object, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	objects := make([]Object, 0)
	it := u.client.Objects(ctx, nil)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		objects = append(objects, Object{
			Key:  obj.Name,
			ETag: strings.ToLower(strings.Trim(obj.Etag, `"`)),
		})
	}

	return objects, nil
}

// Upload uploads the local file to the object
func (u *GoogleUploader) Upload(objectKey, filePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	bodyReader, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer bodyReader.Close()

	wc := u.client.Object(objectKey).NewWriter(ctx)
	if _, err = io.Copy(wc, bodyReader); err != nil {
		return err
	}

	return wc.Close()
}

// Delete deletes the object
func (u *GoogleUploader) Delete(object string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return u.client.Object(object).Delete(ctx)
}
