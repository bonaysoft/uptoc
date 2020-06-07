package core

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"uptoc/uploader"
	"uptoc/utils"
)

// Engine is the core to finish the logic
type Engine struct {
	uploader uploader.Driver

	tobeUploadedObjects []uploader.Object
	tobeDeletedObjects  []uploader.Object
}

// NewEngine returns a new engine.
func NewEngine(uploadDriver uploader.Driver) *Engine {
	return &Engine{
		uploader: uploadDriver,

		tobeUploadedObjects: make([]uploader.Object, 0),
		tobeDeletedObjects:  make([]uploader.Object, 0),
	}
}

// LoadAndCompareObjects loads local files and compare with the remote objects
func (e *Engine) LoadAndCompareObjects(localDir string, excludePaths ...string) error {
	localObjects, err := loadLocalObjects(localDir, excludePaths)
	if err != nil {
		return err
	}
	log.Printf("find %d local objects", len(localObjects))

	remoteObjects, err := e.uploader.ListObjects()
	if err != nil {
		return err
	}
	log.Printf("find %d remote objects", len(remoteObjects))

	log.Printf("compare the local files and the remote objects...")
	for _, localObject := range localObjects {
		if !objectExist(localObject, remoteObjects) {
			localObject.Type = uploader.LocalObjectTypeAdded
			e.tobeUploadedObjects = append(e.tobeUploadedObjects, localObject) // the added objects
		} else if objectNotMatch(localObject, remoteObjects) {
			localObject.Type = uploader.LocalObjectTypeChanged
			e.tobeUploadedObjects = append(e.tobeUploadedObjects, localObject) // the changed objects
		}

		// there do nothing, skip the no change files.
	}

	// find the deleted objects
	for _, remoteObject := range remoteObjects {
		if objectExist(remoteObject, localObjects) {
			continue
		}

		e.tobeDeletedObjects = append(e.tobeDeletedObjects, remoteObject)
	}

	return nil
}

// Sync uploads the to be upload objects to the cloud
// and delete the not exist remote objects
func (e *Engine) Sync() error {
	log.Printf("found %d files to be uploaded, uploading...", len(e.tobeUploadedObjects))
	for _, obj := range e.tobeUploadedObjects {
		log.Printf("[%s] %s", obj.Type, obj.FilePath)
		if err := e.uploader.Upload(obj.Key, obj.FilePath); err != nil {
			return err
		}
	}

	log.Printf("found %d files to be deleted, cleaning...", len(e.tobeDeletedObjects))
	for _, obj := range e.tobeDeletedObjects {
		log.Printf("[deleted] %s", obj.Key)
		if err := e.uploader.Delete(obj.Key); err != nil {
			return err
		}
	}

	log.Printf("files sync done.")
	return nil
}

func objectExist(object uploader.Object, objects []uploader.Object) bool {
	for _, obj := range objects {
		if obj.Key == object.Key {
			return true
		}
	}
	return false
}

func objectNotMatch(object uploader.Object, objects []uploader.Object) bool {
	for _, obj := range objects {
		if obj.Key == object.Key && strings.ToLower(obj.ETag) != object.ETag {
			return true
		}
	}
	return false
}

func shouldExclude(dirPath, filePath string, excludePaths []string) bool {
	for _, ePath := range excludePaths {
		if strings.HasPrefix(filePath, dirPath+strings.TrimPrefix(ePath, "/")) {
			return true
		}
	}

	return false
}

func loadLocalObjects(dirPath string, excludePaths []string) ([]uploader.Object, error) {
	if !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	}

	localObjects := make([]uploader.Object, 0)
	visitor := func(filePath string, info os.FileInfo, err error) error {
		if os.IsNotExist(err) {
			return err
		}

		if info.IsDir() || shouldExclude(dirPath, filePath, excludePaths) {
			return nil
		}

		localObjects = append(localObjects, uploader.Object{
			Key:      strings.TrimPrefix(filePath, dirPath),
			ETag:     utils.FileMD5(filePath),
			FilePath: filePath,
		})
		return nil
	}

	if err := filepath.Walk(dirPath, visitor); err != nil {
		return nil, err
	}

	return localObjects, nil
}
