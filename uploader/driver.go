package uploader

import "fmt"

const (
	LocalObjectTypeAdded   = "added"
	LocalObjectTypeChanged = "changed"
)

type Object struct {
	Key      string // remote file path
	ETag     string // file md5
	FilePath string // local file path
	Type     string // local file type, added or changed
}

// walk and upload to the cloud storage
type Driver interface {
	ListObjects() ([]Object, error)
	Upload(object, rawPath string) error
	Delete(object string) error
}

type Constructor func(endpoint, accessKeyID, accessKeySecret, bucketName string) (Driver, error)

var supportDrivers = map[string]Constructor{
	"oss": AliOSSUploader,
}

func New(driver, endpoint, accessKeyID, accessKeySecret, bucketName string) (Driver, error) {
	if constructor, ok := supportDrivers[driver]; ok {
		return constructor(endpoint, accessKeyID, accessKeySecret, bucketName)
	}

	return nil, fmt.Errorf("driver[%s] not support", driver)
}
