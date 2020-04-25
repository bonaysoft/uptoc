package uploader

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var driverConfigs = map[string]map[string]string{
	"oss": {
		"endpoint":      "oss-cn-hangzhou.aliyuncs.com",
		"access_key":    os.Getenv("UPLOADER_OSS_AK"),
		"access_secret": os.Getenv("UPLOADER_OSS_SK"),
	},
	"qiniu": {
		"endpoint":      "huadong",
		"access_key":    os.Getenv("UPLOADER_QINIU_AK"),
		"access_secret": os.Getenv("UPLOADER_QINIU_AK"),
	},
}

func TestUploader(t *testing.T) {
	tmp := "/tmp/uptoc/"
	files := map[string]string{
		"abc1.txt": "abcabcabc",
		"abc2.txt": "112233",
		"abc3.txt": "445566",
	}
	assert.NoError(t, os.Mkdir(tmp, os.FileMode(0755)))
	for name, content := range files {
		assert.NoError(t, ioutil.WriteFile(tmp+name, []byte(content), os.FileMode(0644)))
	}

	// test the all drivers
	for driver := range supportDrivers {
		config := driverConfigs[driver]
		uploader, err := New(driver, config["endpoint"], config["access_key"], config["access_secret"], "ut-uptoc")
		assert.NoError(t, err)

		// test object upload
		for object := range files {
			assert.NoError(t, uploader.Upload(object, tmp+object))
		}

		// test objects list
		objects, err := uploader.ListObjects()
		assert.NoError(t, err)
		assert.Equal(t, len(files), len(objects))

		// test object delete
		for object := range files {
			assert.NoError(t, uploader.Delete(object))
		}
	}

	// clean the test files.
	assert.NoError(t, os.RemoveAll(tmp))
}

func TestNotSupportDriver(t *testing.T) {
	_, err := New("abc", "test", "test", "test", "test")
	assert.Error(t, err)
}
