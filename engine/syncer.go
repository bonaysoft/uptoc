package engine

import (
	"log"
	"strings"

	"uptoc/uploader"
)

// Syncer is the syncer to finish the logic
type Syncer struct {
	uploader uploader.Driver

	tobeUploadedObjects []uploader.Object
	tobeDeletedObjects  []uploader.Object
}

// New returns a new syncer.
func NewSyncer(uploadDriver uploader.Driver) *Syncer {
	return &Syncer{
		uploader: uploadDriver,

		tobeUploadedObjects: make([]uploader.Object, 0),
		tobeDeletedObjects:  make([]uploader.Object, 0),
	}
}

// Sync uploads the to be upload objects to the cloud
// and delete the not exist remote objects
func (s *Syncer) Sync(localObjects []uploader.Object, saveRoot string) error {
	remoteObjects, err := s.uploader.ListObjects(saveRoot)
	if err != nil {
		return err
	}
	log.Printf("find %d local objects", len(localObjects))
	log.Printf("find %d remote objects", len(remoteObjects))
	log.Printf("compare the local files and the remote objects...")
	s.compareObjects(localObjects, remoteObjects)

	log.Printf("found %d files to be uploaded, uploading...", len(s.tobeUploadedObjects))
	for _, obj := range s.tobeUploadedObjects {
		log.Printf("[%s] %s => %s", obj.Type, obj.FilePath, obj.Key)
		if err := s.uploader.Upload(obj.Key, obj.FilePath); err != nil {
			return err
		}
	}

	log.Printf("found %d files to be deleted, cleaning...", len(s.tobeDeletedObjects))
	for _, obj := range s.tobeDeletedObjects {
		log.Printf("[deleted] %s", obj.Key)
		if err := s.uploader.Delete(obj.Key); err != nil {
			return err
		}
	}

	log.Printf("files sync done.")
	return nil
}

// compareObjects compare local files with the remote objects
func (s *Syncer) compareObjects(localObjects, remoteObjects []uploader.Object) {
	for _, localObject := range localObjects {
		if !s.objectExist(localObject, remoteObjects) {
			localObject.Type = uploader.LocalObjectTypeAdded
			s.tobeUploadedObjects = append(s.tobeUploadedObjects, localObject) // the added objects
		} else if s.objectNotMatch(localObject, remoteObjects) {
			localObject.Type = uploader.LocalObjectTypeChanged
			s.tobeUploadedObjects = append(s.tobeUploadedObjects, localObject) // the changed objects
		}

		// there do nothing, skip the no change files.
	}

	// find the deleted objects
	for _, remoteObject := range remoteObjects {
		if s.objectExist(remoteObject, localObjects) {
			continue
		}

		s.tobeDeletedObjects = append(s.tobeDeletedObjects, remoteObject)
	}
}

func (s *Syncer) objectExist(object uploader.Object, objects []uploader.Object) bool {
	for _, obj := range objects {
		if obj.Key == object.Key {
			return true
		}
	}
	return false
}

func (s *Syncer) objectNotMatch(object uploader.Object, objects []uploader.Object) bool {
	for _, obj := range objects {
		if obj.Key == object.Key && strings.ToLower(obj.ETag) != object.ETag {
			return true
		}
	}
	return false
}
