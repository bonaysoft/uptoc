package uploader

import (
	"fmt"
)

const (
	// Local added file's tag
	LocalObjectTypeAdded = "added"

	// Local changed file's tag
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

// Constructor defines the upload driver constructor used by some engine
type Constructor func(endpoint, accessKeyID, accessKeySecret, bucketName string) (Driver, error)

var supportDrivers = map[string]Constructor{
	"oss": AliOSSUploader,
}

// New is a instantiation function to find and init a upload driver.
func New(driver, endpoint, accessKeyID, accessKeySecret, bucketName string) (Driver, error) {
	if constructor, ok := supportDrivers[driver]; ok {
		return constructor(endpoint, accessKeyID, accessKeySecret, bucketName)
	}

	return nil, fmt.Errorf("driver[%s] not support", driver)
}
