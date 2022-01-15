package engine

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/saltbo/gopkg/fileutil"

	"uptoc/uploader"
)

// Config provides core configuration for the engine.
type Config struct {
	SaveRoot  string   `yaml:"save_root"`
	VisitHost string   `yaml:"visit_host"`
	ForceSync bool     `yaml:"force_sync"`
	Excludes  []string `yaml:"excludes"`
}

func (c *Config) buildRemoteExcludes() []string {
	return c.buildExcludes(c.SaveRoot)
}

func (c *Config) buildExcludes(root string) []string {
	rets := make([]string, 0)
	for _, exclude := range c.Excludes {
		rets = append(rets, filepath.Join(root, exclude))
	}

	return c.Excludes
}

// Engine provides the core logic to finish the feature
type Engine struct {
	echo
	conf Config

	uploader uploader.Driver
}

// New returns a new engine.
func New(conf Config, ud uploader.Driver) *Engine {
	return &Engine{
		conf:     conf,
		uploader: ud,
	}
}

// TailRun run the core logic with every path.
func (e *Engine) TailRun(paths ...string) {
	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatalln(err)
		}

		if stat.IsDir() && e.conf.ForceSync {
			e.syncTo(path)
			continue
		} else if stat.IsDir() {
			e.uploadDirectory(path)
			continue
		}

		e.uploadFile(path, filepath.Join(e.conf.SaveRoot, stat.Name()))
	}
}

func (e *Engine) syncTo(dirPath string) {
	localObjects, err := e.loadLocalObjects(dirPath)
	if err != nil {
		log.Fatalln(err)
	}

	remoteObjects, err := e.uploader.ListObjects(e.conf.SaveRoot)
	if err != nil {
		log.Fatalln(err)
	}

	s := NewSyncer(e.uploader, e.conf, dirPath)
	if err := s.Sync(localObjects, remoteObjects); err != nil {
		log.Fatalln(err)
	}
}

func (e *Engine) uploadDirectory(dirPath string) {
	localObjects, err := e.loadLocalObjects(dirPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, obj := range localObjects {
		e.uploadFile(obj.FilePath, obj.Key)
	}
}

func (e *Engine) uploadFile(filePath, object string) {
	if err := e.uploader.Upload(object, filePath); err != nil {
		e.Failed(filePath, err)
		return
	}

	e.Success(e.conf.VisitHost, object)
}

func (e *Engine) loadLocalObjects(dirPath string) ([]uploader.Object, error) {
	if !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	}

	localObjects := make([]uploader.Object, 0)
	visitor := func(filePath string, info os.FileInfo, err error) error {
		if os.IsNotExist(err) {
			return err
		}

		if info.IsDir() {
			return nil
		}

		localPath := strings.TrimPrefix(filePath, dirPath)
		localObjects = append(localObjects, uploader.Object{
			Key:      filepath.Join(e.conf.SaveRoot, localPath),
			ETag:     fileutil.MD5Hex(filePath),
			FilePath: filePath,
		})
		return nil
	}

	if err := filepath.Walk(dirPath, visitor); err != nil {
		return nil, err
	}

	return localObjects, nil
}

func shouldExclude(filepath string, excludes []string) bool {
	for _, ePath := range excludes {
		if strings.HasPrefix(filepath, ePath) {
			return true
		}
	}

	return false
}
