package core

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"uptoc/uploader"
)

type Engine struct {
	uploader uploader.Uploader

	tobeUploadedObjects []uploader.Object
	tobeDeletedObjects  []uploader.Object
}

func NewEngine(u uploader.Uploader) *Engine {
	return &Engine{
		uploader: u,

		tobeUploadedObjects: make([]uploader.Object, 0),
		tobeDeletedObjects:  make([]uploader.Object, 0),
	}
}

func (e *Engine) LoadAndCompareObjects(localDir string) error {
	localObjects, err := loadLocalObjects(localDir)
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
			localObject.Type = uploader.TYPE_ADDED
			e.tobeUploadedObjects = append(e.tobeUploadedObjects, localObject) // the added objects
		} else if objectNotMatch(localObject, remoteObjects) {
			localObject.Type = uploader.TYPE_CHANGED
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

func (e *Engine) Sync() error {
	if e.uploader == nil {
		return fmt.Errorf("empty uploader")
	}

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

func fileMD5(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		return ""
	}
	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return ""
	}

	return hex.EncodeToString(md5hash.Sum(nil)[:])
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
		if obj.Key == object.Key && obj.ETag != object.ETag {
			return true
		}
	}
	return false
}

func loadLocalObjects(dirPath string) ([]uploader.Object, error) {
	localObjects := make([]uploader.Object, 0)
	visitor := func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		localObjects = append(localObjects, uploader.Object{
			Key:      strings.TrimPrefix(filePath, dirPath),
			ETag:     fileMD5(filePath),
			FilePath: filePath,
		})
		return nil
	}

	if err := filepath.Walk(dirPath, visitor); err != nil {
		return nil, err
	}

	return localObjects, nil
}
