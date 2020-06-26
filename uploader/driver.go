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
	ListObjects(prefix string) ([]Object, error)
	Upload(object, rawPath string) error
	Delete(object string) error
}

type Config struct {
	Name      string `yaml:"name"`
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
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
func New(driver Config) (Driver, error) {
	if _, exist := supportDrivers[driver.Name]; !exist {
		return nil, fmt.Errorf("driver[%s] not support", driver.Name)
	}

	endpoint := supportDrivers[driver.Name]
	if strings.Contains(endpoint, "%s") {
		endpoint = fmt.Sprintf(endpoint, driver.Region)
	}
	if driver.Name == "aws" {
		endpoint = ""
	}

	return NewS3Uploader(driver.Region, endpoint, driver.AccessKey, driver.SecretKey, driver.Bucket)
}
