package uploader

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Uploader implements the Driver base on the S3 of AWS.
type S3Uploader struct {
	client *s3.S3
	bucket string
}

// NewS3Uploader returns a new s3 uploader
func NewS3Uploader(endpoint, accessKeyID, accessKeySecret, bucketName string) (Driver, error) {
	cfg := aws.NewConfig().WithCredentials(credentials.NewStaticCredentials(accessKeyID, accessKeySecret, ""))
	s, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	return &S3Uploader{
		client: s3.New(s, cfg.WithRegion(endpoint)),
		bucket: bucketName,
	}, nil
}

// ListObjects returns some remote objects
func (u *S3Uploader) ListObjects() ([]Object, error) {
	marker := ""
	objects := make([]Object, 0)
	for {
		input := &s3.ListObjectsInput{
			Bucket: aws.String(u.bucket),
			Marker: aws.String(marker),
		}
		objectsResult, err := u.client.ListObjects(input)
		if err != nil {
			return nil, err
		}

		for _, obj := range objectsResult.Contents {
			objects = append(objects, Object{
				Key:  aws.StringValue(obj.Key),
				ETag: strings.ToLower(strings.Trim(aws.StringValue(obj.ETag), `"`))})
		}

		if aws.BoolValue(objectsResult.IsTruncated) {
			marker = aws.StringValue(objectsResult.NextMarker)
		} else {
			break
		}
	}

	return objects, nil
}

// Upload uploads the local file to the object
func (u *S3Uploader) Upload(objectKey, filePath string) (err error) {
	body, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer body.Close()

	_, err = u.client.PutObject(&s3.PutObjectInput{
		Body:   body,
		Bucket: aws.String(u.bucket),
		Key:    aws.String(objectKey),
	})
	return
}

// Delete deletes the object
func (u *S3Uploader) Delete(object string) (err error) {
	_, err = u.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(object),
	})
	return
}
