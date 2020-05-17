package uploader

import (
	"fmt"
	"strings"
)

const (
	// LocalObjectTypeAdded tags the local added file's tag
	LocalObjectTypeAdded = "added"

	// LocalObjectTypeChanged tags the local changed file's tag
	LocalObjectTypeChanged = "changed"
)

// Object is the basic operation unit
type Object struct {
	Key      string // remote file path
	ETag     string // file md5
	FilePath string // local file path
	Type     string // local file type, added or changed
}

// Driver is the interface that must be implemented by a cloud
// storage driver.
type Driver interface {
	ListObjects() ([]Object, error)
	Upload(object, rawPath string) error
	Delete(object string) error
}

// driver => endpoint format template
var supportDrivers = map[string]string{
	"cos":    "cos.%s.myqcloud.com",
	"oss":    "oss-%s.aliyuncs.com",
	"qiniu":  "s3-%s.qiniucs.com",
	"google": "storage.googleapis.com",
	"aws":    "%s",
}

// New is a instantiation function to find and init a upload driver.
func New(driver, region, accessKey, secretKey, bucketName string) (Driver, error) {
	if _, exist := supportDrivers[driver]; !exist {
		return nil, fmt.Errorf("driver[%s] not support", driver)
	}

	endpoint := supportDrivers[driver]
	if strings.Contains(endpoint, "%s") {
		endpoint = fmt.Sprintf(endpoint, region)
	}
	if driver == "aws" {
		endpoint = ""
	}

	return NewS3Uploader(region, endpoint, accessKey, secretKey, bucketName)
}
