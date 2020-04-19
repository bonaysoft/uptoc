package uploader

import "fmt"

const (
	TYPE_ADDED   = "added"
	TYPE_CHANGED = "changed"
)

type Object struct {
	Key      string // remote file path
	ETag     string // file md5
	FilePath string // local file path
	Type     string // local file type, added or changed
}

// walk and upload to the cloud storage
type Uploader interface {
	ListObjects() ([]Object, error)
	Upload(object, rawPath string) error
	Delete(object string) error
}

type Constructor func(endpoint, accessKeyID, accessKeySecret, bucketName string) (Uploader, error)

var supportDrivers = map[string]Constructor{
	"oss": AliOSSUploader,
}

func New(uploader, endpoint, accessKeyID, accessKeySecret, bucketName string) (Uploader, error) {
	if constructor, ok := supportDrivers[uploader]; ok {
		return constructor(endpoint, accessKeyID, accessKeySecret, bucketName)
	}

	return nil, fmt.Errorf("uploader[%s] not support", uploader)
}
